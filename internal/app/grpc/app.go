package grpcapp

import (
	"fmt"
	authgrpc "sso/internal/grpc/auth"

	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	// Метка для функции, чтобы было проще ее отслеживать
	// Оставляем такие метки в функциях, которые что-то логируют "наружу"(?)
	const op = "grpcapp.Run"

	// Обертка для логов этой функции с добавлением ее названия и порта
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	// Слушаем TCP, так как gRPC работает на его базе
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		// Удобная запись для дебага. Название функции: ошибка
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	// Тут мы говорим нашему серверу обрабатывать запросы, которые слушает l
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
