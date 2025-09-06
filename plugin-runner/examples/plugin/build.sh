#!/usr/bin/env bash
set -euo pipefail

# Build the example Go plugin (.so) for linux/amd64 and linux/arm64 using glibc-based toolchains.
# Outputs:
#   plugin-runner/dist/example-linux-amd64.so
#   plugin-runner/dist/example-linux-arm64.so

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
RUNNER_DIR="${REPO_ROOT}/plugin-runner"
DIST_DIR="${RUNNER_DIR}/dist"
mkdir -p "${DIST_DIR}"

GO_VERSION=${GO_VERSION:-"1.22"}
GO_IMAGE="golang:${GO_VERSION}-bullseye"

build_arch() {
  local arch="$1"  # amd64 | arm64
  local out="example-linux-${arch}.so"
  echo "==> Building plugin for linux/${arch} -> dist/${out}"
  docker run --rm --platform="linux/${arch}" \
    -v "${RUNNER_DIR}":/src -w /src \
    -e CGO_ENABLED=1 \
    "${GO_IMAGE}" bash -lc '
      set -e
      apt-get update -qq >/dev/null && apt-get install -y -qq build-essential file >/dev/null
      /usr/local/go/bin/go env -w GOFLAGS=-trimpath
      /usr/local/go/bin/go version
      /usr/local/go/bin/go build -buildmode=plugin -ldflags "-s -w" -o dist/'"${out}"' ./examples/plugin
      if command -v file >/dev/null 2>&1; then file dist/'"${out}"'; fi
    '
}

build_arch amd64
build_arch arm64

echo "==> SHA-256 checksums"
if command -v shasum >/dev/null 2>&1; then
  (cd "${DIST_DIR}" && shasum -a 256 example-linux-*.so)
elif command -v sha256sum >/dev/null 2>&1; then
  (cd "${DIST_DIR}" && sha256sum example-linux-*.so)
fi

echo "Done. Artifacts are in ${DIST_DIR}" 