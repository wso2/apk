# apk-plugin-runner

A tiny subprocess to execute Go plugins (.so built with `-buildmode=plugin`) over a JSON contract.

Contract:
- Symbol signature: `func ProcessJSON(input []byte) ([]byte, error)` (symbol name configurable via `--symbol`).
- Input: JSON from stdin.
- Output: JSON to stdout. Non-JSON outputs are wrapped as `{ "result": "..." }`.

Flags:
- `--so`: path to plugin `.so` (optional, but required if `--symbol` is provided)
- `--symbol`: function symbol (optional, required if `--so` is provided)
- `--timeout`: timeout in milliseconds for the call (0 = no timeout)

Behavior:
- If both `--so` and `--symbol` are omitted, the runner writes a hardcoded valid JSON result to stdout: `{ "result": "ok" }` and exits 0.
- If only one of `--so` or `--symbol` is provided, the runner fails with an error.
- If both are provided, the plugin is loaded and the symbol invoked as described above.

Example:
```
go build -buildmode=plugin -o myplugin.so ./plugin

go build -o apk-plugin-runner ./plugin-runner

cat input.json | ./apk-plugin-runner --so ./myplugin.so --symbol ProcessJSON > out.json

# No-plugin mode example (returns hardcoded result):
echo '{}' | ./apk-plugin-runner
```

Environment integration:
- Enforcer mediation `ExternalCustom` will look for runner from param `runnerPath`, env `APK_PLUGIN_RUNNER`, or `$PATH` (`apk-plugin-runner`).



/usr/bin/env -S bash -lc 'set -euo pipefail
cd /Users/tharsanan/Documents/github/forked/apk/plugin-runner
mkdir -p dist
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/apk-plugin-runner-linux-amd64 .
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/apk-plugin-runner-linux-arm64 .
chmod +x dist/apk-plugin-runner-linux-amd64 dist/apk-plugin-runner-linux-arm64
file dist/apk-plugin-runner-linux-amd64 || true
file dist/apk-plugin-runner-linux-arm64 || true
shasum -a 256 dist/apk-plugin-runner-linux-* | cat'