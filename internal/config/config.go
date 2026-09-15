// Package config загружает и разбирает YAML-конфигурацию приложения.
//
// Файл конфигурации (config.yml) ищется в следующем порядке:
//  1. /etc/dasha-standalone-complex/config.yml
//  2. config.yml в текущей рабочей директории
//
// Если файл не найден ни в одном месте — возвращается ошибка ErrNotFound.
// Путь загруженного файла логируется через slog.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigName — имя файла конфигурации.
const ConfigName = "config.yml"

// Defaults для полей конфига на случай отсутствия значений в файле.
const (
	DefaultLogLevel = "info"
	DefaultAddr     = ":8080"
)

// Config — корневая структура конфигурации.
type Config struct {
	LogLevel string `yaml:"log_level"`
	Server   Server `yaml:"server"`
	STT      STT    `yaml:"stt"`
}

// Server — настройки HTTP-сервера.
type Server struct {
	Addr string `yaml:"addr"`
}

// STT — настройки распознавания речи.
type STT struct {
	Model       string `yaml:"model"`
	Threads     int    `yaml:"threads"`
	MaxUploadMB int64  `yaml:"max_upload_mb"`
}

// ErrNotFound возвращается, если config.yml не найден ни в одном из путей поиска.
var ErrNotFound = errors.New("config.yml не найден")

// Default возвращает конфиг со значениями по умолчанию.
func Default() Config {
	return Config{
		LogLevel: DefaultLogLevel,
		Server:   Server{Addr: DefaultAddr},
	}
}

// SearchPaths возвращает пути, в которых ищется конфигурация, в порядке приоритета.
// Объявлена как переменная, чтобы тесты могли её переопределить.
var SearchPaths = func() []string {
	return []string{
		filepath.Join("/etc", "dasha-standalone-complex", ConfigName),
		filepath.Join(".", ConfigName),
	}
}

// Load ищет config.yml в путях из SearchPaths, парсит его и применяет
// значения по умолчанию для отсутствующих полей. Возвращает путь,
// из которого конфиг был загружен.
func Load() (*Config, string, error) {
	for _, path := range SearchPaths() {
		cfg, err := loadFile(path)
		if err == nil {
			slog.Info("конфиг загружен", "path", path, "log_level", cfg.LogLevel, "addr", cfg.Server.Addr)
			return cfg, path, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, "", err
		}
	}
	return nil, "", ErrNotFound
}

// loadFile читает и разбирает конфиг из конкретного файла.
func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось разобрать %s: %w", path, err)
	}

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = DefaultAddr
	}
	return &cfg, nil
}

// Parse разбирает конфигурацию из YAML-данных (для тестов).
func Parse(data []byte) (*Config, error) {
	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = DefaultAddr
	}
	return &cfg, nil
}
