# apk-plugin-runner

A tiny subprocess to execute Go plugins (.so built with `-buildmode=plugin`) over a JSON contract.

Contract:
- Symbol signature: `func ProcessJSON(input []byte) ([]byte, error)` (symbol name configurable via `--symbol`).
- Input: JSON from stdin.
- Output: JSON to stdout. Non-JSON outputs are wrapped as `{ "result": "..." }`.

Flags:
- `--so`: path to plugin `.so` (required)
- `--symbol`: function symbol (default `ProcessJSON`)
- `--timeout`: timeout in milliseconds for the call (0 = no timeout)

Example:
```
go build -buildmode=plugin -o myplugin.so ./plugin

go build -o apk-plugin-runner ./plugin-runner

cat input.json | ./apk-plugin-runner --so ./myplugin.so --symbol ProcessJSON > out.json
```

Environment integration:
- Enforcer mediation `ExternalCustom` will look for runner from param `runnerPath`, env `APK_PLUGIN_RUNNER`, or `$PATH` (`apk-plugin-runner`).