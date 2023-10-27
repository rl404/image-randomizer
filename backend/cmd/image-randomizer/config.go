package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/fairy/monitoring/newrelic/database"
	"github.com/rl404/image-randomizer/internal/errors"
	"github.com/rl404/image-randomizer/internal/utils"
	"github.com/rl404/image-randomizer/pkg/cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type config struct {
	App      appConfig      `envconfig:"APP"`
	Cache    cacheConfig    `envconfig:"CACHE"`
	DB       dbConfig       `envconfig:"DB"`
	JWT      jwtConfig      `envconfig:"JWT"`
	Log      logConfig      `envconfig:"LOG"`
	Newrelic newrelicConfig `envconfig:"NEWRELIC"`
}

type appConfig struct {
	Env             string        `envconfig:"ENV" validate:"required,oneof=dev prod" mod:"default=dev,no_space,lcase"`
	Port            string        `envconfig:"PORT" validate:"required" mod:"default=31001,no_space"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" validate:"required,gt=0" mod:"default=10s"`
}

type cacheConfig struct {
	Dialect  string        `envconfig:"DIALECT" validate:"required,oneof=nocache redis inmemory" mod:"default=inmemory,no_space,lcase"`
	Address  string        `envconfig:"ADDRESS"`
	Password string        `envconfig:"PASSWORD"`
	Time     time.Duration `envconfig:"TIME" default:"60m" validate:"required,gt=0"`
}

type dbConfig struct {
	Address         string        `envconfig:"ADDRESS" validate:"required" mod:"default=localhost:5432,no_space"`
	Name            string        `envconfig:"NAME" validate:"required" mod:"default=image-randomizer"`
	User            string        `envconfig:"USER" validate:"required" mod:"default=postgres"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN" validate:"required,gt=0" mod:"default=10"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE" validate:"required,gt=0" mod:"default=10"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" validate:"required,gt=0" mod:"default=1m"`
}

type jwtConfig struct {
	AccessSecret   string        `envconfig:"ACCESS_SECRET" validate:"required"`
	AccessExpired  time.Duration `envconfig:"ACCESS_EXPIRED" default:"15m" validate:"required,gt=0"`
	RefreshSecret  string        `envconfig:"REFRESH_SECRET" validate:"required"`
	RefreshExpired time.Duration `envconfig:"REFRESH_EXPIRED" default:"168h" validate:"required,gt=0"`
}

type logConfig struct {
	Level utils.LogLevel `envconfig:"LEVEL" default:"-1"`
	JSON  bool           `envconfig:"JSON" default:"false"`
	Color bool           `envconfig:"COLOR" default:"true"`
}

type newrelicConfig struct {
	Name       string `envconfig:"NAME" default:"image-randomizer"`
	LicenseKey string `envconfig:"LICENSE_KEY"`
}

const envPath = "../../.env"
const envPrefix = "IR"

var cacheType = map[string]cache.CacheType{
	"nocache":  cache.NOP,
	"redis":    cache.Redis,
	"inmemory": cache.InMemory,
}

func getConfig() (*config, error) {
	var cfg config

	// Load .env file.
	_ = godotenv.Load(envPath)

	// Convert env to struct.
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	// Override PORT env.
	if port := os.Getenv("PORT"); port != "" {
		cfg.App.Port = port
	}

	// Validate.
	if err := utils.Validate(&cfg); err != nil {
		return nil, err
	}

	// Init global log.
	utils.InitLog(cfg.Log.Level, cfg.Log.JSON, cfg.Log.Color)

	return &cfg, nil
}

func newDB(cfg dbConfig) (*gorm.DB, error) {
	// Split host and port.
	split := strings.Split(cfg.Address, ":")
	if len(split) != 2 {
		return nil, errors.ErrInvalidDBFormat
	}

	dialector := postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", split[0], split[1], cfg.User, cfg.Password, cfg.Name))

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	tmp, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set basic config.
	tmp.SetMaxIdleConns(cfg.MaxConnIdle)
	tmp.SetMaxOpenConns(cfg.MaxConnOpen)
	tmp.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifetime) * time.Second)

	db.Use(database.NewGORM(cfg.Address, cfg.Name))

	return db, nil
}
