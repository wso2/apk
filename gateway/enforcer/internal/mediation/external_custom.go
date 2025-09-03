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
	"strconv"
	"strings"
	"time"

	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	commonmediation "github.com/wso2/apk/common-go-libs/pkg/mediation"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"google.golang.org/protobuf/types/known/structpb"
)

// ExternalCustom is a mediation that delegates processing to an external plugin runner subprocess.
type ExternalCustom struct {
	policy *dpv2alpha1.Mediation
	// configuration derived from parameters
	pluginPath string
	pluginURL  string
	symbol     string
	runnerPath string
	runnerURL  string
	timeout    time.Duration
	downloadTimeout time.Duration
}

const (
	paramPluginPath = "pluginPath" // required: path to .so (or desired destination path if pluginURL provided)
	paramPluginURL  = "pluginURL"  // optional: URL to download the .so if pluginPath is missing
	paramSymbol     = "symbol"     // optional: function symbol, defaults to ProcessJSON
	paramRunnerPath = "runnerPath" // optional: path to runner binary, defaults to APK_PLUGIN_RUNNER or "apk-plugin-runner" in PATH
	paramRunnerURL  = "runnerURL"  // optional: URL to download the runner binary if not present
	paramTimeoutMs  = "timeoutMs"  // optional: int milliseconds
	paramDownloadTimeoutMs = "downloadTimeoutMs" // optional: int milliseconds for downloading pluginURL
)

// NewExternalCustom constructs an ExternalCustom mediation from the cluster policy.
func NewExternalCustom(m *dpv2alpha1.Mediation) *ExternalCustom {
	ec := &ExternalCustom{policy: m}

	if v, ok := extractPolicyValue(m.Parameters, paramPluginPath); ok {
		ec.pluginPath = v
	}
	if v, ok := extractPolicyValue(m.Parameters, paramPluginURL); ok {
		ec.pluginURL = v
	}
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

	// Best-effort: ensure the plugin is available locally at construction time
	// to avoid repeated downloads/checks on every request. We intentionally
	// ignore errors here and fall back to runtime ensure in invokeRunner.
	if localPath, err := ec.ensurePluginLocal(); err == nil && localPath != "" {
		ec.pluginPath = localPath
	}

	// Best-effort: ensure runner binary is available locally as well.
	if localRunner, err := ec.ensureRunnerLocal(); err == nil && localRunner != "" {
		ec.runnerPath = localRunner
	}
	return ec
}

// Process implements the Mediation interface by delegating to the external runner.
func (e *ExternalCustom) Process(h *requestconfig.Holder) *Result {
	// Build minimal, stable JSON input for the plugin.
	in := e.buildInput(h)
	payload, err := json.Marshal(in)
	if err != nil {
		// On marshalling error, fail open (no change) to avoid breaking traffic.
		return NewResult()
	}

	outPayload, err := e.invokeRunner(payload)
	if err != nil {
		// On runner errors, fail open.
		return NewResult()
	}

	// Decode output and map to Result
	var extOut commonmediation.ExternalOutput
	if err := json.Unmarshal(outPayload, &extOut); err != nil {
		return NewResult()
	}
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
	return in
}

func nilIfEmptyMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	return m
}

func (e *ExternalCustom) invokeRunner(input []byte) ([]byte, error) {
	if e.pluginPath == "" && e.pluginURL == "" {
		return nil, errors.New("pluginPath parameter is required (or provide pluginURL)")
	}

	// Ensure the plugin .so is available locally; if not, try to download.
	localPath, err := e.ensurePluginLocal()
	if err != nil {
		return nil, err
	}

	// Ensure the runner binary is available locally; if not, try to download/resolve.
	localRunner, err := e.ensureRunnerLocal()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, localRunner, "--so", localPath, "--symbol", e.symbol)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("runner timeout: %w", err)
		}
		return nil, fmt.Errorf("runner error: %v, stderr: %s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}

// ensurePluginLocal ensures there's a readable local file for the plugin and
// returns the path to use with the runner. Behaviors:
// - If pluginPath is an existing file, return it.
// - If pluginURL is provided and pluginPath does not exist but is non-empty,
//   download pluginURL to pluginPath (creating parent dirs) and return pluginPath.
// - If pluginPath itself is an HTTP/HTTPS URL, download it into a cache dir and
//   return the cached file path.
func (e *ExternalCustom) ensurePluginLocal() (string, error) {
	// Case 1: pluginPath points to existing local file
	if e.pluginPath != "" {
		if fi, err := os.Stat(e.pluginPath); err == nil && !fi.IsDir() {
			return e.pluginPath, nil
		}
	}

	// Case 2: pluginURL provided + pluginPath provided as destination
	if e.pluginURL != "" && e.pluginPath != "" {
		if err := downloadToFile(e.pluginURL, e.pluginPath, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download plugin to %s: %w", e.pluginPath, err)
		}
		return e.pluginPath, nil
	}

	// Case 3: pluginPath is a URL; download to cache directory
	if isHTTPURL(e.pluginPath) {
		destDir := filepath.Join(os.TempDir(), "apk-plugins-cache")
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create cache dir: %w", err)
		}
		base := filepath.Base(e.pluginPath)
		if base == "." || base == "/" || base == "" {
			base = "plugin.so"
		}
		dest := filepath.Join(destDir, base)
		if err := downloadToFile(e.pluginPath, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download plugin: %w", err)
		}
		return dest, nil
	}

	// Case 4: Only pluginURL provided; download to cache using URL basename
	if e.pluginURL != "" {
		destDir := filepath.Join(os.TempDir(), "apk-plugins-cache")
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create cache dir: %w", err)
		}
		base := filepath.Base(e.pluginURL)
		if base == "." || base == "/" || base == "" {
			base = "plugin.so"
		}
		dest := filepath.Join(destDir, base)
		if err := downloadToFile(e.pluginURL, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download plugin: %w", err)
		}
		return dest, nil
	}

	return "", errors.New("plugin file not found and no valid URL to download")
}

func isHTTPURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func downloadToFile(url, dest string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
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
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
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
	// Try PATH resolution first for convenience (handles bare names like "apk-plugin-runner").
	if rp != "" && !isHTTPURL(rp) {
		if abs, err := exec.LookPath(rp); err == nil {
			return abs, nil
		}
		// If file exists at the given path, use it.
		if fi, err := os.Stat(rp); err == nil && !fi.IsDir() {
			return rp, nil
		}
	}

	// Download to a specified destination if both URL and path are provided.
	if e.runnerURL != "" && rp != "" && !isHTTPURL(rp) {
		if err := downloadToFile(e.runnerURL, rp, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner to %s: %w", rp, err)
		}
		// Ensure executable
		_ = os.Chmod(rp, 0o755)
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
		if err := downloadToFile(rp, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner: %w", err)
		}
		_ = os.Chmod(dest, 0o755)
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
		if err := downloadToFile(e.runnerURL, dest, e.downloadTimeout); err != nil {
			return "", fmt.Errorf("failed to download runner: %w", err)
		}
		_ = os.Chmod(dest, 0o755)
		return dest, nil
	}

	// Last resort: try PATH again for default name
	if abs, err := exec.LookPath("apk-plugin-runner"); err == nil {
		return abs, nil
	}
	return "", errors.New("runner binary not found and no valid URL to download")
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
