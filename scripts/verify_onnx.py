#!/usr/bin/env python3
"""Проверяет, что ONNX-графы TeraTTSv2 загружаются через ONNX Runtime без ошибок.

Запуск (из корня проекта):
    python3 scripts/verify_onnx.py
"""

import os
import sys

import onnxruntime as ort

ROOT = os.path.join("models", "teratts")
MODELS = [
    os.path.join("models", "text_encoder.onnx"),
    os.path.join("models", "duration_predictor.onnx"),
    os.path.join("models", "sampler_distilled_cfg3_8step.onnx"),
    os.path.join("models", "sampler_teacher_8step.onnx"),
    os.path.join("models", "vocoder.onnx"),
]


def main():
    providers = ["CPUExecutionProvider"]
    available = ort.get_available_providers()
    print("ONNX Runtime:", ort.__version__, "| providers:", available)
    failures = []
    for rel in MODELS:
        path = os.path.join(ROOT, rel)
        if not os.path.exists(path):
            print(f"FAIL {rel}: file not found")
            failures.append(rel)
            continue
        try:
            so = ort.SessionOptions()
            so.log_severity_level = 3
            sess = ort.InferenceSession(path, so, providers=providers)
            print(f"OK   {rel}: {len(sess.get_inputs())} inputs, "
                  f"{len(sess.get_outputs())} outputs")
        except Exception as exc:  # noqa: BLE001
            print(f"FAIL {rel}: {exc}")
            failures.append(rel)
    if failures:
        print(f"\nERROR: {len(failures)} model(s) failed to load: {failures}")
        return 1
    print("\nOK: all ONNX models loaded without errors")
    return 0


if __name__ == "__main__":
    sys.exit(main())