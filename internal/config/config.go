// internal/config/config.go
package config

import "time"

// Config хранит метаданные и настройки приложения.
type Config struct {
	Name        string
	Description string
	Mission     string
	Author      string
	Year        int
	Port        string // Сетевой порт
	// DSN         string // Строка подключения к БД (добавим позже)
}

// Load загружает конфигурацию приложения.
func Load() Config {
	return Config{
		Name:        "shanraq",
		Description: "AI - ваш семейный доктор",
		Mission:     "AI на страже вашего здоровья",
		Author:      "Daulet Baimyrza & Companions",
		Year:        time.Now().Year(),
		Port:        ":8080", // Порт по умолчанию
	}
}