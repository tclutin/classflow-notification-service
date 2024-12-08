package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-notification-service/internal/config"
	"github.com/tclutin/classflow-notification-service/internal/handler"
	"github.com/tclutin/classflow-notification-service/internal/repository"
	"github.com/tclutin/classflow-notification-service/internal/service"
	"github.com/tclutin/classflow-notification-service/pkg/client/postgresql"
	"github.com/tclutin/classflow-notification-service/pkg/client/telegram"
	"github.com/tclutin/classflow-notification-service/pkg/logger"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	server              *http.Server
	logger              *slog.Logger
	db                  *pgxpool.Pool
	notificationService *service.NotificationService
}

func New() *App {
	cfg := config.MustLoad()

	appLogger := logger.New(cfg.Environment, "logs/app.log")

	dsn := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName)

	postgres := postgresql.NewPool(context.Background(), dsn)

	repo := repository.NewScheduleRepository(postgres, appLogger)

	tgClient := telegram.NewTGClient(cfg.Telegram.Token)

	notifyService := service.NewNotificationService(appLogger, tgClient, repo)

	router := handler.New().Init()

	return &App{
		server: &http.Server{
			Addr:    net.JoinHostPort(cfg.HTTPServer.Address, cfg.HTTPServer.Port),
			Handler: router,
		},
		notificationService: notifyService,
		db:                  postgres,
		logger:              appLogger,
	}
}

func (app *App) Run() {
	app.logger.Info("Starting application...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		app.logger.Info("Server is starting...")
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Server stopped with error", "error", err)
			os.Exit(1)
		}
	}()

	app.logger.Info("Server started successfully")

	app.logger.Info("Notification service is starting...")

	go app.notificationService.Start(context.Background())

	app.logger.Info("Notification service started successfully")

	<-stop

	app.logger.Info("Shutting down app...")

	err := app.server.Shutdown(context.Background())
	if err != nil {
		app.logger.Error("Error during server shutdown", "error", err)
		os.Exit(1)
	}

	defer app.db.Close()

	app.logger.Info("Application shutdown completed")
}
