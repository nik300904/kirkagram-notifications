package main

import (
	"context"
	"kirkagram-notification/internal/config"
	"kirkagram-notification/internal/kafka"
	"kirkagram-notification/internal/lib/logger/handlers/slogpretty"
	"kirkagram-notification/internal/service"
	"kirkagram-notification/internal/storage"
	"kirkagram-notification/internal/storage/psgr"
	"kirkagram-notification/internal/transport"
	"kirkagram-notification/internal/transport/ws"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal     = "local"
	envProd      = "prod"
	envDev       = "dev"
	numConsumers = 3
)

func main() {
	cfg := config.New()
	log := setupLogger(cfg.Env)
	db := storage.New(cfg)

	postRepo := psgr.NewPostStorage(db)
	likeRepo := psgr.NewLikeStorage(db)

	postService := service.NewPostService(postRepo, log)
	likeService := service.NewLikeService(likeRepo, log)

	wsManager := ws.NewWebSocketManager(log, likeService, postService)

	likeConsumers := kafka.NewConsumer("like-group", log)
	defer likeConsumers.Close()

	postConsumers := kafka.NewConsumer("like-group", log)
	defer postConsumers.Close()

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < numConsumers; i++ {
		go func(consumerNum int) {
			for {
				err := likeConsumers.ConsumeMessages(ctx, wsManager, "like")
				if err != nil {
					log.Error("Error consuming like messages",
						slog.String("error", err.Error()),
						slog.Int("consumer", consumerNum),
						slog.String("consumer", "post"))
					time.Sleep(5 * time.Second)
				}
				if ctx.Err() != nil {
					return
				}
			}
		}(i)
		go func(consumerNum int) {
			for {
				err := postConsumers.ConsumeMessages(ctx, wsManager, "post")
				if err != nil {
					log.Error("Error consuming like messages",
						slog.String("error", err.Error()),
						slog.Int("consumer", consumerNum),
						slog.String("consumer", "post"))
					time.Sleep(5 * time.Second)
				}
				if ctx.Err() != nil {
					return
				}
			}
		}(i)
	}

	log.Info("Consumer started. Press Ctrl+C to stop")

	router := transport.NewHandler(wsManager)

	srv := &http.Server{
		Addr:    "localhost:8083",
		Handler: router.InitRouter(),
	}

	log.Info("SERVER STARTED AT", slog.String("address", cfg.HttpServe.Address))

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	<-sigChan
	log.Info("Shutting down...")

	cancel()
	time.Sleep(time.Second)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug}}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
