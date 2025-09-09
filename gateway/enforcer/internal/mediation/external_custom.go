package mediation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	commonmediation "github.com/wso2/apk/common-go-libs/pkg/mediation"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

// ExternalCustom is a mediation that delegates processing to an external plugin runner subprocess.
type ExternalCustom struct {
	policy *dpv2alpha1.Mediation
	// configuration derived from parameters
	symbol          string
	runnerPath      string
	runnerURL       string
	timeout         time.Duration
	downloadTimeout time.Duration
	cfg             *config.Server
}

const (
	paramPluginPath        = "pluginPath"        // required: path to .so (or desired destination path if pluginURL provided)
	paramPluginURL         = "pluginURL"         // optional: URL to download the .so if pluginPath is missing
	paramSymbol            = "symbol"            // optional: function symbol, defaults to ProcessJSON
	paramRunnerPath        = "runnerPath"        // optional: path to runner binary, defaults to APK_PLUGIN_RUNNER or "apk-plugin-runner" in PATH
	paramRunnerURL         = "runnerURL"         // optional: URL to download the runner binary if not present
	paramTimeoutMs         = "timeoutMs"         // optional: int milliseconds
	paramDownloadTimeoutMs = "downloadTimeoutMs" // optional: int milliseconds for downloading pluginURL
)

// NewExternalCustom constructs an ExternalCustom mediation from the cluster policy.
func NewExternalCustom(m *dpv2alpha1.Mediation) *ExternalCustom {
	ec := &ExternalCustom{policy: m}
	ec.cfg = config.GetConfig()
	if v, ok := extractPolicyValue(m.Parameters, paramSymbol); ok && v != "" {
		ec.symbol = v
	} else {
		ec.symbol = "ProcessJSON"
	}
	if v, ok := extractPolicyValue(m.Parameters, paramRunnerPath); ok && v != "" {
		ec.runnerPath = v
	} else if env := os.Getenv("APK_PLUGIN_RUNNER"); env != "" {
		ec.runnerPath = env
	} else {
		ec.runnerPath = "apk-plugin-runner"
	}
	if v, ok := extractPolicyValue(m.Parameters, paramRunnerURL); ok {
		ec.runnerURL = v
	}
	if v, ok := extractPolicyValue(m.Parameters, paramTimeoutMs); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ec.timeout = time.Duration(n) * time.Millisecond
		}
	}
	if ec.timeout == 0 {
		ec.timeout = 2 * time.Second
	}
	if v, ok := extractPolicyValue(m.Parameters, paramDownloadTimeoutMs); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ec.downloadTimeout = time.Duration(n) * time.Millisecond
		}
	}
	if ec.downloadTimeout == 0 {
		ec.downloadTimeout = 15 * time.Second
	}

	// Best-effort: ensure runner binary is available locally as well.
	if localRunner, err := ec.ensureRunnerLocal(); err == nil && localRunner != "" {
		ec.runnerPath = localRunner
	}
	// Debug summary of constructed mediation
	ec.dbg("ExternalCustom initialized: policy=%s(id=%s,ver=%s), symbol=%s, runnerPath=%s, runnerURL=%s, timeout=%s, downloadTimeout=%s, GOOS=%s, GOARCH=%s",
		safeStr(ec.policy.PolicyName), safeStr(ec.policy.PolicyID), safeStr(ec.policy.PolicyVersion),
		safeStr(ec.symbol), safeStr(ec.runnerPath), safeStr(ec.runnerURL), ec.timeout, ec.downloadTimeout, runtime.GOOS, runtime.GOARCH)
	return ec
}

// Process implements the Mediation interface by delegating to the external runner.
func (e *ExternalCustom) Process(h *requestconfig.Holder) *Result {
	start := time.Now()
	if !e.cfg.ExternalCustomMediationEnabled {
		return NewResult()
	}
	// Log high-level context only (avoid sensitive data)
	e.dbg("Process start: policy=%s(id=%s,ver=%s), phase=%s, symbol=%s, runner=%s",
		safeStr(e.policy.PolicyName), safeStr(e.policy.PolicyID), safeStr(e.policy.PolicyVersion),
		string(h.ProcessingPhase), safeStr(e.symbol), safeStr(e.runnerPath))
	// Build minimal, stable JSON input for the plugin.
	in := e.buildInput(h)
	e.dbg("Built input: params=%d, attrs=%d, reqHdrs=%d, respHdrs=%d, reqBodyLen=%d, respBodyLen=%d",
		len(in.Parameters), len(in.Attributes), lenOrZero(in.RequestHeaders), lenOrZero(in.ResponseHeaders), len(in.RequestBody), len(in.ResponseBody))
	payload, err := json.Marshal(in)
	if err != nil {
		e.cfg.Logger.Sugar().Errorf("ExternalCustom JSON marshal error: %v", err)
		// On marshalling error, fail open (no change) to avoid breaking traffic.
		return NewResult()
	}
	e.dbg("Input JSON size=%d bytes, sample=%q", len(payload), truncate(string(payload), 256))

	outPayload, err := e.invokeRunner(payload)
	if err != nil {
		e.cfg.Logger.Sugar().Errorf("ExternalCustom runner invocation error: %v", err)
		// On runner errors, fail open.
		return NewResult()
	}

	// Decode output and map to Result
	var extOut commonmediation.ExternalOutput
	if err := json.Unmarshal(outPayload, &extOut); err != nil {
		return NewResult()
	}
	e.dbg("Runner output size=%d bytes, flags: modifyBody=%t, immediate=%t(code=%d, hdrs=%d, bodyLen=%d), addHdrs=%d, rmHdrs=%d, stopFurther=%t, metadata=%d",
		len(outPayload), extOut.ModifyBody, extOut.ImmediateResponse, extOut.ImmediateResponseCode,
		len(extOut.ImmediateResponseHeaders), len(extOut.ImmediateResponseBody),
		len(extOut.AddHeaders), len(extOut.RemoveHeaders), extOut.StopFurtherProcessing, len(extOut.Metadata))
	e.dbg("Process end in %s", time.Since(start))
	return mapExternalOutputToResult(&extOut)
}

// externalInput defines the JSON contract sent to the plugin runner.
// externalInput now lives in common-go-libs as mediation.ExternalInput

// externalOutput defines the JSON contract returned by the plugin runner.
// externalOutput now lives in common-go-libs as mediation.ExternalOutput

func (e *ExternalCustom) buildInput(h *requestconfig.Holder) *commonmediation.ExternalInput {
	params := map[string]string{}
	if e.policy != nil {
		for _, p := range e.policy.Parameters {
			if p != nil && p.Key != "" && p.Value != "" {
				params[p.Key] = p.Value
			}
		}
	}

	attrs := map[string]string{}
	if h != nil && h.RequestAttributes != nil {
		if h.RequestAttributes.RouteName != "" {
			attrs["routeName"] = h.RequestAttributes.RouteName
		}
		if h.RequestAttributes.RequestID != "" {
			attrs["requestID"] = h.RequestAttributes.RequestID
		}
	}

	// Only include simple header maps if available; ignore complex Envoy options.
	reqHeaders := map[string]string{}
	if h != nil && h.RequestHeaders != nil && h.RequestHeaders.Headers != nil {
		for _, hv := range h.RequestHeaders.Headers.Headers {
			if hv.GetKey() != "" {
				reqHeaders[hv.GetKey()] = string(hv.GetRawValue())
			}
		}
	}
	respHeaders := map[string]string{}
	if h != nil && h.ResponseHeaders != nil && h.ResponseHeaders.Headers != nil {
		for _, hv := range h.ResponseHeaders.Headers.Headers {
			if hv.GetKey() != "" {
				respHeaders[hv.GetKey()] = string(hv.GetRawValue())
			}
		}
	}

	var reqBody, respBody string
	if h != nil && h.RequestBody != nil {
		reqBody = string(h.RequestBody.Body)
	}
	if h != nil && h.ResponseBody != nil {
		respBody = string(h.ResponseBody.Body)
	}

	phase := string(h.ProcessingPhase)

	in := &commonmediation.ExternalInput{
		Phase:           phase,
		PolicyName:      e.policy.PolicyName,
		PolicyID:        e.policy.PolicyID,
		PolicyVersion:   e.policy.PolicyVersion,
		Parameters:      params,
		Attributes:      attrs,
		RequestHeaders:  nilIfEmptyMap(reqHeaders),
		ResponseHeaders: nilIfEmptyMap(respHeaders),
		RequestBody:     reqBody,
		ResponseBody:    respBody,
	}
	// Debug summary for input (with redactions)
	e.dbg("buildInput: phase=%s, params=%d, attrs=%d, reqHdrs=%v, respHdrs=%v, reqBody=%q, respBody=%q",
		phase, len(params), len(attrs), redactAndSampleHeaders(reqHeaders, 5), redactAndSampleHeaders(respHeaders, 5),
		truncate(reqBody, 128), truncate(respBody, 128))
	return in
}

func nilIfEmptyMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	return m
}

func (e *ExternalCustom) invokeRunner(input []byte) ([]byte, error) {
	// Ensure the runner binary is available locally; if not, try to download/resolve.
	localRunner, err := e.ensureRunnerLocal()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, localRunner)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	e.dbg("Invoking runner: %s (GOOS=%s, GOARCH=%s, timeout=%s, input=%d bytes)", localRunner, runtime.GOOS, runtime.GOARCH, e.timeout, len(input))
	t0 := time.Now()
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("runner timeout after %s: %w, runner=%s, GOOS=%s, GOARCH=%s, stderr: %s", e.timeout, err, localRunner, runtime.GOOS, runtime.GOARCH, truncate(stderr.String(), 2048))
		}
		// Provide helpful context for exec format errors.
		if errors.Is(err, syscall.ENOEXEC) || strings.Contains(err.Error(), "exec format error") {
			return nil, fmt.Errorf("runner binary is not executable for this platform: %v (runner=%s, GOOS=%s, GOARCH=%s). Ensure the runner matches the container host architecture. Stderr: %s", err, localRunner, runtime.GOOS, runtime.GOARCH, truncate(stderr.String(), 2048))
		}
		return nil, fmt.Errorf("runner error: %v (runner=%s, GOOS=%s, GOARCH=%s), stderr: %s", err, localRunner, runtime.GOOS, runtime.GOARCH, stderr.String())
	}
	e.dbg("Runner finished in %s, stdout=%d bytes, sample=%q, stderrSample=%q", time.Since(t0), stdout.Len(), truncate(stdout.String(), 256), truncate(stderr.String(), 256))
	return stdout.Bytes(), nil
}

func isHTTPURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// truncate returns s limited to limit bytes, appending "..." if truncated.
func truncate(s string, limit int) string {
	if limit <= 0 || len(s) <= limit {
		return s
	}
	if limit <= 3 {
		return s[:limit]
	}
	return s[:limit-3] + "..."
}

func downloadToFile(url, dest string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	// Debug: begin download
	// (No auth headers used; safe to log URL and destination path)
	// Destination directory will be created if needed.
	// The file will be atomically moved to final path after download completes.
	// Truncate long paths for logs.
	// Note: this function might run during hot path only on first use or cache miss.
	// Keep debug level to avoid noisy logs in normal operations.
	//nolint:staticcheck // debug-only logs
	//
	// Log start
	// Use a lightweight call to print
	// The actual logger is on ExternalCustom, not available here – print via fmt only if needed.
	// We'll avoid printing; rely on caller-side logs.
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	// ensure parent dir
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	// Download to a temporary file in the same directory, then atomically rename.
	tmpFile, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer func() { _ = os.Remove(tmpName) }()
	// Ensure close on all paths
	defer func() { _ = tmpFile.Close() }()
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}
	// Flush contents
	if err := tmpFile.Sync(); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	// Make it readable by others by default; callers can chmod further if needed.
	_ = os.Chmod(tmpName, 0o644)
	// Atomic replace
	if err := os.Rename(tmpName, dest); err != nil {
		return err
	}
	// Best-effort chmod +x if looks like a binary name
	_ = os.Chmod(dest, 0o755)
	return nil
}

// ensureRunnerLocal ensures there's a runnable local binary for the runner and
// returns the path to use. Behaviors:
// - If runnerPath resolves in PATH or exists locally, return its absolute path.
// - If runnerURL is provided and runnerPath is provided as destination, download there and chmod +x.
// - If runnerPath itself is a URL, download to cache and chmod +x.
// - If only runnerURL is provided, download to cache and chmod +x.
func (e *ExternalCustom) ensureRunnerLocal() (string, error) {
	rp := e.runnerPath
	e.dbg("ensureRunnerLocal: start, runnerPath=%s, runnerURL=%s", safeStr(rp), safeStr(e.runnerURL))
	// Try PATH resolution first for convenience (handles bare names like "apk-plugin-runner").
	if rp != "" && !isHTTPURL(rp) {
		if abs, err := exec.LookPath(rp); err == nil {
			e.dbg("ensureRunnerLocal: found in PATH => %s", abs)
			return abs, nil
		}
		// If file exists at the given path, use it.
		if fi, err := os.Stat(rp); err == nil && !fi.IsDir() {
			e.dbg("ensureRunnerLocal: using existing file => %s", rp)
			return rp, nil
		}
	}

	// Download to a specified destination if both URL and path are provided.
	if e.runnerURL != "" && rp != "" && !isHTTPURL(rp) {
		e.dbg("ensureRunnerLocal: downloading runnerURL to specified path => url=%s dest=%s", e.runnerURL, rp)
		if err := downloadToFile(e.runnerURL, rp, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner to %s: %w", rp, err)
		}
		// Ensure executable
		_ = os.Chmod(rp, 0o755)
		e.dbg("ensureRunnerLocal: download complete => %s", rp)
		return rp, nil
	}

	// If runnerPath is a URL, download to cache directory.
	if isHTTPURL(rp) {
		destDir := filepath.Join(os.TempDir(), "apk-runners-cache")
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create runner cache dir: %w", err)
		}
		base := filepath.Base(rp)
		if base == "." || base == "/" || base == "" {
			base = "apk-plugin-runner"
		}
		dest := filepath.Join(destDir, base)
		if fi, err := os.Stat(dest); err == nil && !fi.IsDir() {
			_ = os.Chmod(dest, 0o755)
			// Quick sanity: try invoking with --version to catch exec format errors early.
			_ = quickCheckRunner(dest, e.timeout/4)
			e.dbg("ensureRunnerLocal: cached runner found => %s", dest)
			return dest, nil
		}
		e.dbg("ensureRunnerLocal: downloading runner from path-URL => %s to %s", rp, dest)
		if err := downloadToFile(rp, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner: %w", err)
		}
		_ = os.Chmod(dest, 0o755)
		e.dbg("ensureRunnerLocal: download complete => %s", dest)
		return dest, nil
	}

	// If only runnerURL is provided, download to cache directory.
	if e.runnerURL != "" {
		destDir := filepath.Join(os.TempDir(), "apk-runners-cache")
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create runner cache dir: %w", err)
		}
		base := filepath.Base(e.runnerURL)
		if base == "." || base == "/" || base == "" {
			base = "apk-plugin-runner"
		}
		dest := filepath.Join(destDir, base)
		if fi, err := os.Stat(dest); err == nil && !fi.IsDir() {
			_ = os.Chmod(dest, 0o755)
			// Quick sanity: try invoking with --version to catch exec format errors early.
			_ = quickCheckRunner(dest, e.timeout/4)
			e.dbg("ensureRunnerLocal: cached runner (from runnerURL) found => %s", dest)
			return dest, nil
		}
		e.dbg("ensureRunnerLocal: downloading runner from runnerURL => %s to %s", e.runnerURL, dest)
		if err := downloadToFile(e.runnerURL, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner: %w", err)
		}
		_ = os.Chmod(dest, 0o755)
		e.dbg("ensureRunnerLocal: download complete => %s", dest)
		return dest, nil
	}

	// Last resort: try PATH again for default name
	if abs, err := exec.LookPath("apk-plugin-runner"); err == nil {
		_ = quickCheckRunner(abs, e.timeout/4)
		e.dbg("ensureRunnerLocal: fallback PATH found => %s", abs)
		return abs, nil
	}
	e.dbg("ensureRunnerLocal: runner not found")
	return "", errors.New("runner binary not found and no valid URL to download")
}

// quickCheckRunner tries to run the runner with --version within a short timeout.
// It's best-effort; errors are ignored but help surface format/permission issues earlier.
func quickCheckRunner(path string, d time.Duration) error {
	if d <= 0 {
		d = 500 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	return err
}

func mapExternalOutputToResult(o *commonmediation.ExternalOutput) *Result {
	r := NewResult()
	if o == nil {
		return r
	}
	for k, v := range o.AddHeaders {
		r.AddHeaders[k] = v
	}
	if len(o.RemoveHeaders) > 0 {
		r.RemoveHeaders = append(r.RemoveHeaders, o.RemoveHeaders...)
	}
	r.ModifyBody = o.ModifyBody
	r.Body = o.Body
	r.ImmediateResponse = o.ImmediateResponse
	r.ImmediateResponseCode = v32.StatusCode(o.ImmediateResponseCode)
	r.ImmediateResponseBody = o.ImmediateResponseBody
	r.ImmediateResponseDetail = o.ImmediateResponseDetail
	for k, v := range o.ImmediateResponseHeaders {
		r.ImmediateResponseHeaders[k] = v
	}
	r.ImmediateResponseContentType = o.ImmediateResponseContentType
	r.StopFurtherProcessing = o.StopFurtherProcessing
	if len(o.Metadata) > 0 {
		for k, v := range o.Metadata {
			// Convert to structpb.Value where possible
			if spb, err := structpb.NewValue(v); err == nil {
				r.Metadata[k] = spb
			}
		}
	}
	return r
}

// -----------------
// Logging helpers
// -----------------

// dbg logs at debug level if a logger is available.
func (e *ExternalCustom) dbg(format string, args ...any) {
	if e == nil || e.cfg == nil {
		return
	}
	e.cfg.Logger.Sugar().Debugf(format, args...)
}

func safeStr(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func lenOrZero(m map[string]string) int {
	if m == nil {
		return 0
	}
	return len(m)
}

// redactAndSampleHeaders returns a small, redacted subset of headers for debug logs.
// It redacts Authorization and Cookie values and returns up to maxKeys entries.
func redactAndSampleHeaders(h map[string]string, maxKeys int) map[string]string {
	if h == nil || maxKeys <= 0 {
		return nil
	}
	out := make(map[string]string, 0)
	n := 0
	for k, v := range h {
		kk := strings.ToLower(k)
		vv := v
		if kk == "authorization" || kk == "proxy-authorization" || kk == "cookie" || kk == "set-cookie" {
			vv = "<redacted>"
		} else {
			vv = truncate(v, 64)
		}
		out[k] = vv
		n++
		if n >= maxKeys {
			break
		}
	}
	return out
}
