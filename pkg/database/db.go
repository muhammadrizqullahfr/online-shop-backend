package database

import (
	"context"
	"fmt"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Timezone string
	Retry    config.RetryConfig
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
		c.Timezone,
	)
}

func (c *DBConfig) InitConnectionDB(ctx context.Context, l *logger.SlogAdapter, appConf *config.Config) (*gorm.DB, error) {
	// var level logger.LogLevel

	// switch appConf.LogLevel {
	// case "debug":
	// 	level = logger.Info
	// case "info":
	// 	level = logger.Info
	// case "warn":
	// 	level = logger.Warn
	// case "error":
	// 	level = logger.Error
	// }
	gcf := gorm.Config{
		SkipDefaultTransaction: true,
		// Logger:                 log.LogMode(level),
	}

	var db *gorm.DB
	var err error

	for i := 0; i < c.Retry.Max || c.Retry.Max == -1; i++ {
		db, err = gorm.Open(postgres.Open(c.GetDSN()), &gcf)
		if err == nil {

			sqlDB, err := db.DB()
			if err != nil {
				return nil, fmt.Errorf("failed to get sql.DB: %v", err)
			}

			l.Info(ctx, "Database is connected")
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(10)

			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(100)

			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDB.SetConnMaxLifetime(5 * time.Minute)

			return db, nil
		}

		l.Info(ctx, fmt.Sprintf("Postgres connection failed: %v", err))
		l.Info(ctx, fmt.Sprintf("Retrying in %v", c.Retry.Delay))
		time.Sleep(c.Retry.Delay)
	}

	return nil, err

}
