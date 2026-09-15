package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	data := []byte(`
log_level: "debug"
server:
  addr: "127.0.0.1:9090"
stt:
  model: "test.gguf"
  threads: 2
  max_upload_mb: 42
`)
	cfg, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug", cfg.LogLevel)
	}
	if cfg.Server.Addr != "127.0.0.1:9090" {
		t.Fatalf("Addr = %q, want 127.0.0.1:9090", cfg.Server.Addr)
	}
	if cfg.STT.Model != "test.gguf" || cfg.STT.Threads != 2 || cfg.STT.MaxUploadMB != 42 {
		t.Fatalf("STT = %+v, want test.gguf/2/42", cfg.STT)
	}
}

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse([]byte(""))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if cfg.LogLevel != DefaultLogLevel {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, DefaultLogLevel)
	}
	if cfg.Server.Addr != DefaultAddr {
		t.Fatalf("Addr = %q, want %q", cfg.Server.Addr, DefaultAddr)
	}
}

func TestLoad_PrefersEtc(t *testing.T) {
	dir := t.TempDir()
	etcPath := filepath.Join(dir, "etc", "dasha-standalone-complex")
	cwdPath := filepath.Join(dir, "cwd")
	for _, d := range []string{etcPath, cwdPath} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(filepath.Join(etcPath, ConfigName), []byte("log_level: error"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cwdPath, ConfigName), []byte("log_level: debug"), 0o644); err != nil {
		t.Fatal(err)
	}

	old := SearchPaths
	SearchPaths = func() []string {
		return []string{
			filepath.Join(etcPath, ConfigName),
			filepath.Join(cwdPath, ConfigName),
		}
	}
	defer func() { SearchPaths = old }()

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.LogLevel != "error" {
		t.Fatalf("LogLevel = %q, want error (из /etc)", cfg.LogLevel)
	}
	if path != filepath.Join(etcPath, ConfigName) {
		t.Fatalf("path = %q, want %q", path, filepath.Join(etcPath, ConfigName))
	}
}

func TestLoad_FallsBackToCwd(t *testing.T) {
	dir := t.TempDir()
	cwdPath := filepath.Join(dir, ConfigName)
	if err := os.WriteFile(cwdPath, []byte("server:\n  addr: \":9999\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	old := SearchPaths
	SearchPaths = func() []string {
		return []string{
			filepath.Join(dir, "etc", "dasha-standalone-complex", ConfigName),
			cwdPath,
		}
	}
	defer func() { SearchPaths = old }()

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Server.Addr != ":9999" {
		t.Fatalf("Addr = %q, want :9999", cfg.Server.Addr)
	}
	if path != cwdPath {
		t.Fatalf("path = %q, want %q", path, cwdPath)
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()

	old := SearchPaths
	SearchPaths = func() []string {
		return []string{
			filepath.Join(dir, "nope", ConfigName),
			filepath.Join(dir, "nope2", ConfigName),
		}
	}
	defer func() { SearchPaths = old }()

	_, _, err := Load()
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load() error = %v, want ErrNotFound", err)
	}
}

func TestLoad_BrokenConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ConfigName)
	if err := os.WriteFile(path, []byte("\t:\t: broken: [yaml"), 0o644); err != nil {
		t.Fatal(err)
	}

	old := SearchPaths
	SearchPaths = func() []string {
		return []string{path}
	}
	defer func() { SearchPaths = old }()

	_, _, err := Load()
	if err == nil {
		t.Fatal("Load() = nil, want error for broken yaml")
	}
	if errors.Is(err, ErrNotFound) {
		t.Fatalf("второй путь не должен пробоваться после ошибки парсинга, got %v", err)
	}
}
