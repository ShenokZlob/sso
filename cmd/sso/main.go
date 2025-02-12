package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: иницилизировать объект конфига
	cfg := config.MustLoad()

	// TODO: иницилизировать логер
	log := setupLogger(cfg.Env)
	// В реальных прогах, скорее всего, опасно выводить содержимое конфига!
	log.Info("starting application", slog.Any("cfg", cfg))

	// TODO: иницилизировать приложение
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// TODO: запустить gRPC-сервер приложения
	go application.GRPCSrv.MustRun()

	// Graceful shutdown
	// Теперь программа перед выключением будет пытаться что-то сделать (завершить задачи?)
	// Ждем конкретные сигналы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Ждем когда канал получит сигнал от ОС
	s := <-stop
	log.Info("stopping application", slog.String("signal", s.String()))
	application.GRPCSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	// Используем TextHandler для вывода в консоль во время разработки
	case envLocal:
		log = slog.New(
			// В опциях выбираем уровень вывода ошибок
			// Выбираем Debug так как хотим видеть все ошибки в логах
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	// Используем формат JSON так как с этими логами будет работать машина, а не чел
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
