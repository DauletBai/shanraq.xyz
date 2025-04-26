// cmd/web/main.go
package main

import (
	"database/sql" 
	"flag"
	"html/template" 
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"shanraq.xyz/internal/config"
	"shanraq.xyz/internal/handler"
	"shanraq.xyz/internal/view"
)
type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
	config        config.Config
	db            *sql.DB 
}

func main() {
	cfg := config.Load() 
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	addr := flag.String("addr", cfg.Port, "Адрес HTTP сервера")
	flag.Parse()

	logger.Info("Подключение к базе данных SQLite...")
	db, err := openDB(cfg.DBFile) 
	if err != nil {
		logger.Error("Не удалось подключиться к БД SQLite", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("База данных SQLite успешно подключена!")

	templateCache, err := view.NewTemplateCache()
	if err != nil {
		logger.Error("Не удалось создать кэш шаблонов", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		templateCache: templateCache,
		config:        cfg,
		db:            db, 
	}

	appHandler := handler.NewHandler(logger, templateCache, cfg, app.db) 

	srv := &http.Server{
		Addr:         *addr,
		Handler:      appHandler.Routes(http.Dir("./static/assets")),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Запуск веб-сервера", slog.String("addr", srv.Addr))
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Ошибка при запуске сервера", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	// Проверяем реальное соединение
	err = db.Ping()
	if err != nil {
		db.Close() 
		return nil, err
	}

	// В SQLite часто включают PRAGMA для улучшения работы с внешними ключами
	// _, err = db.Exec("PRAGMA foreign_keys = ON;")
    // if err != nil {
    //    db.Close()
    //    return nil, err
    // }

	// Настройки пула для SQLite менее критичны, но можно оставить
	// db.SetMaxOpenConns(1) // SQLite обычно лучше работает с одним пишущим соединением

	return db, nil
}