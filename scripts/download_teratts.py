#!/usr/bin/env python3
"""Скачивает веса модели TeraTTSv2 (ONNX) в каталог models/teratts.

Модель распространяется авторами уже в формате ONNX (графы для ONNX Runtime),
поэтому "конвертация" выполняется на стороне издателя, а здесь происходит
загрузка готовых ONNX-графов и вспомогательных ресурсов в проект.

Запуск:
    python3 scripts/download_teratts.py
"""

import json
import os
import sys
import urllib.request

MODEL_ID = "TeraSpace/TeraTTSv2"
DEST = os.path.join("models", "teratts")
API_URL = f"https://huggingface.co/api/models/{MODEL_ID}"
RESOLVE_URL = f"https://huggingface.co/{MODEL_ID}/resolve/main"

SKIP_PREFIXES = ("__pycache__",)


def list_files():
    with urllib.request.urlopen(API_URL) as resp:
        payload = json.load(resp)
    return [s["rfilename"] for s in payload.get("siblings", [])]


def download(url, dest_path):
    os.makedirs(os.path.dirname(dest_path), exist_ok=True)
    tmp = dest_path + ".part"
    req = urllib.request.Request(url, headers={"User-Agent": "teratts-download"})
    with urllib.request.urlopen(req) as resp, open(tmp, "wb") as out:
        while True:
            chunk = resp.read(1024 * 1024)
            if not chunk:
                break
            out.write(chunk)
    os.replace(tmp, dest_path)


def main():
    files = list_files()
    skipped = 0
    for rel in files:
        if rel.startswith(SKIP_PREFIXES):
            skipped += 1
            continue
        dest = os.path.join(DEST, rel)
        if os.path.exists(dest):
            print(f"skip  {rel}")
            continue
        url = f"{RESOLVE_URL}/{urllib.request.quote(rel)}"
        print(f"get   {rel}")
        download(url, dest)
    print(f"done: {len(files) - skipped} files downloaded into {DEST}")


if __name__ == "__main__":
    sys.exit(main())