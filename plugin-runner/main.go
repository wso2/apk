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
	flag.StringVar(&soPath, "so", "", "Path to the plugin .so file")
	flag.StringVar(&symbol, "symbol", "ProcessJSON", "Plugin function symbol to call")
	flag.IntVar(&timeoutMs, "timeout", 0, "Timeout in milliseconds (0 = no timeout)")
	flag.Parse()

	if soPath == "" {
		fail(2, errors.New("missing --so plugin path"))
	}

	// Read all stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail(3, fmt.Errorf("failed to read stdin: %w", err))
	}

	plg, err := plugin.Open(soPath)
	if err != nil {
		fail(4, fmt.Errorf("failed to open plugin: %w", err))
	}
	sym, err := plg.Lookup(symbol)
	if err != nil {
		fail(5, fmt.Errorf("failed to lookup symbol %q: %w", symbol, err))
	}

	fn, ok := sym.(func([]byte) ([]byte, error))
	if !ok {
		fail(6, fmt.Errorf("symbol %q has incompatible type; expected func([]byte) ([]byte, error)", symbol))
	}

	done := make(chan struct{})
	var out []byte
	var callErr error

	if timeoutMs > 0 {
		go func() {
			out, callErr = fn(input)
			close(done)
		}()
		select {
		case <-done:
			// proceed
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			fail(7, errors.New("plugin call timed out"))
		}
	} else {
		out, callErr = fn(input)
	}

	if callErr != nil {
		fail(8, fmt.Errorf("plugin returned error: %w", callErr))
	}

	// Validate output is JSON (nice to have; tolerate plain bytes otherwise)
	if len(out) == 0 {
		out = []byte("{}")
	} else if !json.Valid(out) {
		// Wrap non-JSON as {"result": "..."}
		wrapped, _ := json.Marshal(map[string]string{"result": string(out)})
		out = wrapped
	}

	// Write to stdout
	if _, err := os.Stdout.Write(out); err != nil {
		fail(9, fmt.Errorf("failed to write stdout: %w", err))
	}
}

func fail(code int, err error) {
	log.SetFlags(0)
	_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}
