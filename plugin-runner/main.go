package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"plugin"
	"time"
)

// Contract: plugin symbol is a function with signature:
//   func ProcessJSON(input []byte) ([]byte, error)
// Symbol name is configurable via --symbol; default is ProcessJSON.

func main() {
	var soPath string
	var symbol string
	var timeoutMs int
	var verbose bool
	flag.StringVar(&soPath, "so", "", "Path to the plugin .so file")
	// Symbol is optional only when --so is also omitted. If one is provided without the other, we fail.
	flag.StringVar(&symbol, "symbol", "", "Plugin function symbol to call (required if --so is provided)")
	flag.IntVar(&timeoutMs, "timeout", 0, "Timeout in milliseconds (0 = no timeout)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging to stderr")
	flag.Parse()

	// Configure logging; ensure timestamps for easier tracing.
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	vprintf := func(format string, args ...any) {
		if verbose {
			log.Printf(format, args...)
		}
	}

	vprintf("runner start: so=%s symbol=%s timeoutMs=%d pid=%d", soPath, symbol, timeoutMs, os.Getpid())

	// Validate pairing of --so and --symbol. Both required together, or neither.
	if (soPath == "" && symbol != "") || (soPath != "" && symbol == "") {
		fail(2, errors.New("must provide both --so and --symbol, or neither"))
	}

	// If neither is provided, return a hardcoded valid JSON result and exit 0.
	if soPath == "" && symbol == "" {
		// Emit a hardcoded ExternalOutput-shaped JSON with explicit fields.
		out := []byte(`{
			"addHeaders": {},
			"removeHeaders": [],
			"modifyBody": false,
			"body": "",
			"immediateResponse": true,
			"immediateResponseCode": 211,
			"immediateResponseBody": "",
			"immediateResponseDetail": "",
			"immediateResponseHeaders": {},
			"immediateResponseContentType": "",
			"stopFurtherProcessing": false,
			"metadata": {}
		}`)
		if _, err := os.Stdout.Write(out); err != nil {
			fail(9, fmt.Errorf("failed to write stdout: %w", err))
		}
		return
	}

	// Read all stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail(3, fmt.Errorf("failed to read stdin: %w", err))
	}
	vprintf("stdin read: %d bytes", len(input))

	vprintf("opening plugin: %s", soPath)
	plg, err := plugin.Open(soPath)
	if err != nil {
		fail(4, fmt.Errorf("failed to open plugin: %w", err))
	}
	vprintf("plugin opened")
	sym, err := plg.Lookup(symbol)
	if err != nil {
		fail(5, fmt.Errorf("failed to lookup symbol %q: %w", symbol, err))
	}
	vprintf("symbol %q resolved", symbol)

	fn, ok := sym.(func([]byte) ([]byte, error))
	if !ok {
		fail(6, fmt.Errorf("symbol %q has incompatible type; expected func([]byte) ([]byte, error)", symbol))
	}

	done := make(chan struct{})
	var out []byte
	var callErr error
	start := time.Now()
	vprintf("invoking symbol %q", symbol)

	if timeoutMs > 0 {
		go func() {
			out, callErr = fn(input)
			close(done)
		}()
		select {
		case <-done:
			// proceed
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			vprintf("plugin call timed out after %dms", timeoutMs)
			fail(7, errors.New("plugin call timed out"))
		}
	} else {
		out, callErr = fn(input)
	}
	vprintf("invocation finished in %s (err=%v)", time.Since(start), callErr)

	if callErr != nil {
		fail(8, fmt.Errorf("plugin returned error: %w", callErr))
	}

	// Validate output is JSON (nice to have; tolerate plain bytes otherwise)
	if len(out) == 0 {
		vprintf("plugin returned empty output; defaulting to {}")
		out = []byte("{}")
	} else if !json.Valid(out) {
		// Wrap non-JSON as {"result": "..."}
		vprintf("plugin output not valid JSON; wrapping as {result: <string>}")
		wrapped, _ := json.Marshal(map[string]string{"result": string(out)})
		out = wrapped
	}
	vprintf("writing output: %d bytes", len(out))

	// Write to stdout
	if _, err := os.Stdout.Write(out); err != nil {
		fail(9, fmt.Errorf("failed to write stdout: %w", err))
	}
	vprintf("output write complete")
}

func fail(code int, err error) {
	log.SetFlags(0)
	_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}
