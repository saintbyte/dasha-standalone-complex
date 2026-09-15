# TeraTTS — Makefile для сборки единого бинарника и скачивания зависимостей.
#
# Основные цели:
#   make deps    — скачать все зависимости (ONNX Runtime + веса модели)
#   make build   — скачать зависимости (если нужно) и собрать бинарник ./teratts
#   make test    — запустить тесты
#   make check   — тесты + go vet
#   make format  — отформатировать Go-код
#   make clean   — удалить собранный бинарник
SHELL := /bin/bash

GO          ?= go
PYTHON      ?= python3
BIN         ?= bin/dasha
THIRD_PARTY := third_party/llama-go

export LIBRARY_PATH  := $(abspath $(THIRD_PARTY))
export C_INCLUDE_PATH := $(abspath $(THIRD_PARTY))


.PHONY: all deps deps-onnx deps-models build test check format clean

all: build

# Скачивание всех требуемых библиотек: ONNX Runtime и весов модели TeraTTSv2.
# Скрипты идемпотентны: повторный запуск пропускает уже загруженные файлы.
deps: deps-onnx deps-models

# Скачивание ONNX Runtime (libonnxruntime.so) в third_party/onnxruntime/lib.
deps-onnx:
	./scripts/download_onnxruntime.sh

# Скачивание весов модели TeraTTSv2 (ONNX) в models/teratts.
deps-models:
	$(PYTHON) scripts/download_teratts.py

# Сборка единого бинарника (CLI + HTTP-сервис). Перед сборкой гарантированно
# скачиваются зависимости (если их ещё нет).
build: deps
	mkdir -p $(dir $(BIN))
	$(GO) build -o $(BIN) ./cmd/dasha

test:
	$(GO) test ./...

check: test
	$(GO) vet ./...

format:
	$(GO) fmt ./...

clean:
	rm -rf $(BIN)
	rm -rf $(THIRD_PARTY)
	rm -rf build

## Клонирует и собирает библиотеку llama-go (нужна один раз)
lib:
	@if [ ! -f "$(THIRD_PARTY)/libbinding.a" ]; then \
		echo "Сборка llama-go..." && \
		git clone --depth 1 --branch llama.cpp-b10069 \
			--recurse-submodules --shallow-submodules \
			https://github.com/tcpipuk/llama-go "$(THIRD_PARTY)" && \
		$(MAKE) -C "$(THIRD_PARTY)" libbinding.a; \
	fi
