package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/pkg/logger"
	"github.com/pressly/goose/v3"
)

func AutoMigrate(ctx context.Context, db *sql.DB, config *config.Config, log *logger.SlogAdapter, migrationDir string) error {

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect %v", err)
	}

	current, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current migration version %v", err)
	}

	log.Info(ctx, fmt.Sprintf("Current migration version : %v", current))

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("failed to apply migration: %v", err)
	}

	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}

	if newVersion > current {
		log.Info(ctx, "Successfully migrated")

	} else {
		log.Info(ctx, "Database Already Up to date")
	}

	return nil
}
