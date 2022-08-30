package main

import (
	imageDB "github.com/rl404/image-randomizer/internal/domain/image/repository/db"
	userDB "github.com/rl404/image-randomizer/internal/domain/user/repository/db"
	"github.com/rl404/image-randomizer/internal/utils"
)

func migrate() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Migrate.
	utils.Info("migrating...")
	if err := db.AutoMigrate(
		&userDB.User{},
		&imageDB.Image{},
	); err != nil {
		return err
	}

	utils.Info("done")
	return nil
}
