package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
// Содержит все настраиваемые параметры сервера.
type Config struct {
	Addr                    string
	DatabaseURL             string
	PendingBatch            int
	RetainDelivered         time.Duration // Время, через которое нужно удалить доставленные сообщения
	RetainUndelivered       time.Duration // Время, сколько хранить недоставленные сообщения
	AllowedOrigins          []string
	JanitorInterval         time.Duration
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	HandshakeTimeout        time.Duration
}

// Читает конфигурацию из ENV и возвращает готовую структуру.
func Load() (Config, error) {
	get := func(k, def string) string {
		if v := os.Getenv(k); v != "" { return v }
		return def
	}
	atoi := func(k string, def int) int {
		if v := os.Getenv(k); v != "" {
			if n, err := strconv.Atoi(v); err == nil { return n }
		}
		return def
	}

	cfg := Config{
		Addr:             get("ADDR", ":8080"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		PendingBatch:     atoi("PENDING_BATCH", 1000),
		RetainDelivered:  time.Duration(atoi("RETAIN_DELIVERED_MINUTES", 60)) * time.Minute,
		RetainUndelivered: time.Duration(atoi("RETAIN_UNDELIVERED_DAYS", 14)) * 24 * time.Hour,
		AllowedOrigins:   split(get("ALLOWED_ORIGINS", "")),
		JanitorInterval:  time.Duration(atoi("JANITOR_INTERVAL_SECONDS", 60)) * time.Second,
		ReadTimeout:      time.Duration(atoi("READ_TIMEOUT_SECONDS", 60)) * time.Second,
		WriteTimeout:     time.Duration(atoi("WRITE_TIMEOUT_SECONDS", 20)) * time.Second,
		HandshakeTimeout: time.Duration(atoi("HANDSHAKE_TIMEOUT_SECONDS", 5)) * time.Second,
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL required")
	}
	return cfg, nil
}

func split(s string) []string {
	if s == "" { return nil }
	parts := strings.Split(s, ",")
	for i := range parts { parts[i] = strings.TrimSpace(parts[i]) }
	return parts
}