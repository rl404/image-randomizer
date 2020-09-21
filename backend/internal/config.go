package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Config is config model for app.
type Config struct {
	// Key for encrypt and decrypt.
	Masterkey string

	// HTTP port.
	Port string

	// Database ip and port.
	Address string
	// Database name.
	DB string
	// Database schema name.
	Schema string
	// Database username.
	User string
	// Database password.
	Password string
}

const (
	// DefaultPort is default HTTP app port.
	DefaultPort = "31001"
	// EnvPath is .env file path.
	EnvPath = ".env"
	// EnvPrefix is environment
	EnvPrefix = "IR"
	// DefaultMaxIdleConn is default max db idle connection.
	DefaultMaxIdleConn = 10
	// DefaultMaxOpenConn is default max db open connection.
	DefaultMaxOpenConn = 10
	// DefaultConnMaxLifeTime is default db connection max life time.
	DefaultConnMaxLifeTime = 5 * time.Minute
)

// GetConfig to get config from env.
func GetConfig() (cfg Config) {
	cfg.Port = DefaultPort

	// Load .env file if exist.
	godotenv.Load(EnvPath)

	// Convert env to struct.
	envconfig.Process(EnvPrefix, &cfg)

	// Prepare the ":" for starting HTTP.
	cfg.Port = ":" + cfg.Port

	return cfg
}

// InitDB to intiate db connection.
func (c *Config) InitDB() (db *gorm.DB, err error) {
	if c.Address == "" {
		return nil, ErrRequiredDB
	}

	// Split address and port.
	split := strings.Split(c.Address, ":")
	if len(split) != 2 {
		return nil, ErrInvalidDB
	}

	// Open db connection.
	conn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", split[0], split[1], c.DB, c.User, c.Password)
	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 nil,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   c.Schema + ".",
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, err
	}

	// Set base connection setting.
	dbTmp, _ := db.DB()
	dbTmp.SetMaxIdleConns(DefaultMaxIdleConn)
	dbTmp.SetMaxOpenConns(DefaultMaxOpenConn)
	dbTmp.SetConnMaxLifetime(DefaultConnMaxLifeTime)

	// Set default schema.
	err = db.Exec(fmt.Sprintf("SET search_path TO %s", c.Schema)).Error
	if err != nil {
		return db, err
	}

	// Schema check.
	if !c.isSchemaExist(db) {
		err := c.createSchema(db)
		if err != nil {
			return db, err
		}
	}

	err = db.AutoMigrate(&User{}, &Image{})
	if err != nil {
		return db, err
	}

	return db, nil
}

// isSchemaExist to check if schema is exist.
func (c *Config) isSchemaExist(db *gorm.DB) (isExist bool) {
	db.Raw("SELECT EXISTS(SELECT 1 FROM pg_namespace WHERE nspname = ?)", c.Schema).Row().Scan(&isExist)
	return isExist
}

// createSchema to create new schema.
func (c *Config) createSchema(db *gorm.DB) error {
	query := fmt.Sprintf("CREATE SCHEMA \"%s\" AUTHORIZATION \"%s\"", c.Schema, c.User)
	return db.Exec(query).Error
}
