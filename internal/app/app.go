package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"medods-test-task/internal/config"
	"medods-test-task/internal/service"
	"medods-test-task/internal/storage"
	"medods-test-task/internal/transport/rest"
	"medods-test-task/internal/transport/rest/handler"
	"os/signal"
	"syscall"
)

type App struct {
	conf config.Config
}

func New() *App {
	conf := config.New()

	return &App{
		conf: conf,
	}
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("starting app")

	slog.Info("connecting to mongodb")
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", app.conf.Mongo.Host, app.conf.Mongo.Port)))
	if err != nil {
		slog.Error("failed to connect to mongodb", slog.Any("error", err))
		return
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		slog.Error("failed to ping mongodb", slog.Any("error", err))
		return
	}

	sessionStorage := storage.NewSessionStorage(mongoClient)
	sessionService := service.NewSessionService(sessionStorage, app.conf.Session)
	sessionHandler := handler.NewSessionHandler(sessionService)

	server := rest.New(app.conf.Server).Handle(sessionHandler)

	slog.Info("starting web service")
	go func() {
		err = server.Listen()
		if err != nil {
			slog.Error("failed to start web service", slog.Any("error", err))
			return
		}
	}()

	<-ctx.Done()

	slog.Info("stopping web service")
	err = server.Shutdown()
	if err != nil {
		slog.Error("failed to stop web service", slog.Any("error", err))
		return
	}
}
