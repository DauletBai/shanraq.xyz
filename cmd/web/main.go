// cmd/web/main.go
package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"shanraq.xyz/internal/config"
	"shanraq.xyz/internal/handler"
	"shanraq.xyz/internal/view"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	addr := flag.String("addr", cfg.Port, "Адрес HTTP сервера")
	flag.Parse()

	templateCache, err := view.NewTemplateCache()
	if err != nil {
		logger.Error("Не удалось создать кэш шаблонов", slog.String("error", err.Error()))
		os.Exit(1)
	}

	appHandler := handler.NewHandler(logger, templateCache, cfg)

	staticAssetFS := http.Dir("./static/assets")
	srv := &http.Server{
		Addr: *addr,
		Handler:      appHandler.Routes(staticAssetFS), 
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Запуск веб-сервера (чтение файлов с диска)", slog.String("addr", srv.Addr))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Ошибка при запуске сервера", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
