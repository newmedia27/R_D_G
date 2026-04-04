package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	_ "project/docs"
	"project/internal/clients"
	"project/internal/clients/mongodb"
	"project/internal/config"
	"project/internal/handlers"
	"project/internal/logger"
	"project/internal/middlewares"
	"project/internal/repositories"
	"project/internal/services"
	"project/internal/ws"
)

// @title Project (chat api)
// @version 1.0
// @description Example API
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Default().Error("Error load config", slog.Any("err", err))
		os.Exit(1)
	}
	logger.InitLogger(cfg)

	ctx := context.Background()

	c, err := clients.NewClients(ctx, cfg)
	if err != nil {
		logger.LogFatalContext(ctx, err)
	}
	defer mongodb.DisconnectMongo(c.Mongo.Db)

	r := repositories.NewRepositories(c)
	s := services.NewServices(cfg, r)

	hub := ws.NewHub()
	go hub.Run()

	app := fiber.New()
	m := middlewares.NewMiddlewares(cfg, s)
	h := handlers.NewHandlers(m, s, cfg, hub)

	h.InitRouters(app)

	go func() {
		if err = app.Listen(":" + cfg.Port); err != nil {
			logger.LogFatalContext(ctx, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.LogFatalContextWithCode(ctx, err, 13)
	}
}
