package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/pkg/bootstrapper"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/RakaMurdiarta/online-shop-system/pkg/logger"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer/mailslurp"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer/stub"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func main() {

	v := validator.New()
	ctx := context.Background()

	log := logger.NewSlogAdapter(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	//load config
	cfg := config.LoadConfig()

	//validate .env config
	cfg.Validate(v)

	DB_PORT, err := strconv.Atoi(cfg.DBPort)

	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	//connection db
	dbConfig := &database.DBConfig{
		Host:     cfg.DBHost,
		Port:     DB_PORT,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
		Database: cfg.DBDatabase,
		Timezone: cfg.DBTimezone,
		Retry: config.RetryConfig{
			Max:   cfg.RetryConfig.Max,
			Delay: cfg.RetryConfig.Delay,
		},
	}

	db, err := dbConfig.InitConnectionDB(ctx, log, cfg)

	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("[DB] Database Connection Failed ,%v", err))

	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("[DB] failed to get sql.DB ,%v", err))
	}

	migrationDir := "./migrations/"
	if err := database.AutoMigrate(ctx, sqlDB, cfg, log, migrationDir); err != nil {
		log.Error(ctx, fmt.Sprintf("Database Migration %v", err))
	}

	//supabase client storage
	supabaseStorageClient := shared.NewSupabaseStorageClient(cfg)

	//mailer transport (provider-agnostic; swap subpackage to change provider)
	var mailTransport mailer.Transport
	switch cfg.MailerDriver {
	case "stub":
		log.Info(ctx, "[Mailer] using stub transport (no external sends)")
		mailTransport = stub.New()
	default:
		mailTransport = mailslurp.New(cfg.MailerAPIKey, cfg.MailerInboxID)
	}

	echo := echo.New()

	apiServer := bootstrapper.NewServer(echo, cfg, db, supabaseStorageClient, mailTransport)

	apiServer.InitAPI()

	listen := net.JoinHostPort(cfg.AppHost, cfg.AppPort)

	srv := &http.Server{
		Addr:    listen,
		Handler: echo,
	}

	log.Info(ctx, fmt.Sprintf("[Server] Running on : http://%v:%v", cfg.AppHost, cfg.AppPort))

	//blocking operation
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(ctx, fmt.Sprintf("[Server] Error running server : %v", err))
	}

}
