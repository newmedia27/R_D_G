package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"docker/internal/documentstore"
	"github.com/k0kubun/pp"
)

func main() {
	l := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(l))
	store := documentstore.NewStore()
	pp.Println(store)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong\n"))
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
