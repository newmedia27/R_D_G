package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mongodb "db/internal/clients/mongo-db"
	"db/internal/handlers"
	"db/internal/repositories"
	"db/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	client, err := mongodb.ConnectToMongo()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer func(c *mongo.Database) {
		mongodb.DisconnectMongo(c)
	}(client)

	r := repositories.NewRepositories(client)
	s := services.NewServices(r)
	h := handlers.NewHandler(s)

	mux := h.InitRoutes()

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("Server started", "port", 8080)
		if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			slog.Error("Server error:", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")

}
