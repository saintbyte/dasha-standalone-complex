#!/usr/bin/env bash
# Скачивает библиотеку ONNX Runtime (libonnxruntime.so) в third_party/onnxruntime/lib.
# Библиотека не хранится в git (см. .gitignore), её нужно загрузить перед сборкой Go-движка.
set -euo pipefail

VERSION="${ONNXRUNTIME_VERSION:-1.28.0}"
DEST_DIR="third_party/onnxruntime/lib"
TARBALL="onnxruntime-linux-x64-${VERSION}.tgz"
URL="https://github.com/microsoft/onnxruntime/releases/download/v${VERSION}/${TARBALL}"
TMP="$(mktemp -d)"

if ls "${DEST_DIR}/libonnxruntime.so."* >/dev/null 2>&1; then
  echo "onnxruntime already present in ${DEST_DIR}"
  exit 0
fi

echo "downloading ${URL}"
curl -sSL -o "${TMP}/${TARBALL}" "${URL}"
tar -xzf "${TMP}/${TARBALL}" -C "${TMP}"

mkdir -p "${DEST_DIR}"
cp -a "${TMP}/onnxruntime-linux-x64-${VERSION}/lib/libonnxruntime.so"* "${DEST_DIR}/"
rm -rf "${TMP}"
echo "installed onnxruntime ${VERSION} into ${DEST_DIR}"