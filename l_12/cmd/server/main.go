package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"tcp/cmd/server/handler"
	"tcp/internal/documentstore"
	"tcp/internal/users"
)

func main() {
	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(l))

	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("Failed to listen", slog.Any("error", err))
		return
	}

	slog.Info("Server started", slog.Any("address", server.Addr()))
	defer func() {
		err = server.Close()
		if err != nil {
			slog.Error("Failed to close server", slog.Any("error", err))
		}
	}()

	store := documentstore.NewStore()
	userService := users.NewUserService(store)

	for {
		conn, err := server.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", slog.Any("error", err))
			continue
		}
		fmt.Println("New connection accepted: ", conn.RemoteAddr())

		go func(c net.Conn) {
			h := handler.NewHandleConnection(c, userService)
			h.Handle()
		}(conn)
	}
}
