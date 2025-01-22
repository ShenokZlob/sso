package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: иницилизировать хранилище (storage)
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// TODO: init auth service (auth)
	authService := auth.New(log, storage, tokenTTL)

	// Создаем приложение и возвращаем его
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
