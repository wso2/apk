/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/wso2/apk/test/integration/integration/utils/roundtripper"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
)

// ExpectedResponse defines the response expected for a given request.
type ExpectedResponse struct {
	// Request defines the request to make.
	Request Request

	// ExpectedRequest defines the request that
	// is expected to arrive at the backend. If
	// not specified, the backend request will be
	// expected to match Request.
	ExpectedRequest *ExpectedRequest

	RedirectRequest *roundtripper.RedirectRequest

	// BackendSetResponseHeaders is a set of headers
	// the echoserver should set in its response.
	BackendSetResponseHeaders map[string]string

	// Response defines what response the test case
	// should receive.
	Response Response

	Backend   string
	Namespace string

	// User Given TestCase name
	TestCaseName string
}

// Request can be used as both the request to make and a means to verify
// that echoserver received the expected request. Note that multiple header
// values can be provided, as a comma-separated value.
type Request struct {
	Host             string
	Method           string
	Path             string
	Headers          map[string]string
	Body             string
	UnfollowRedirect bool
}

// ExpectedRequest defines expected properties of a request that reaches a backend.
type ExpectedRequest struct {
	Request

	// AbsentHeaders are names of headers that are expected
	// *not* to be present on the request.
	AbsentHeaders []string
}

// Response defines expected properties of a response from a backend.
type Response struct {
	StatusCode    int
	Headers       map[string]string
	AbsentHeaders []string
}

const (
	backendJWTHeader = "X-Jwt-Assertion"
)

// MakeRequestAndExpectEventuallyConsistentResponse makes a request with the given parameters,
// understanding that the request may fail for some amount of time.
//
// Once the request succeeds consistently with the response having the expected status code, make
// additional assertions on the response body using the provided ExpectedResponse.
func MakeRequestAndExpectEventuallyConsistentResponse(t *testing.T, r roundtripper.RoundTripper, timeoutConfig config.TimeoutConfig, gwAddr string, expected ExpectedResponse) {
	t.Helper()

	req := MakeRequest(t, &expected, gwAddr, "HTTPS", "https")

	WaitForConsistentResponse(t, r, req, expected, timeoutConfig.RequiredConsecutiveSuccesses, timeoutConfig.MaxTimeToConsistency)
}

// MakeRequest make a request with the given parameters.
func MakeRequest(t *testing.T, expected *ExpectedResponse, gwAddr, protocol, scheme string) roundtripper.Request {
	t.Helper()

	if expected.Request.Method == "" {
		expected.Request.Method = "GET"
	}

	if expected.Response.StatusCode == 0 {
		expected.Response.StatusCode = 200
	}

	t.Logf("Making %s request to %s://%s%s", expected.Request.Method, scheme, gwAddr, expected.Request.Path)

	path, query, _ := strings.Cut(expected.Request.Path, "?")

	req := roundtripper.Request{
		Method:           expected.Request.Method,
		Host:             expected.Request.Host,
		Body:             expected.Request.Body,
		URL:              url.URL{Scheme: scheme, Host: gwAddr, Path: path, RawQuery: query},
		Protocol:         protocol,
		Headers:          map[string][]string{},
		UnfollowRedirect: expected.Request.UnfollowRedirect,
	}

	if expected.Request.Headers != nil {
		for name, value := range expected.Request.Headers {
			req.Headers[name] = []string{value}
		}
	}

	backendSetHeaders := []string{}
	for name, val := range expected.BackendSetResponseHeaders {
		backendSetHeaders = append(backendSetHeaders, name+":"+val)
	}

	req.Headers["X-Echo-Set-Header"] = []string{strings.Join(backendSetHeaders, ",")}

	return req
}

// AwaitConvergence runs the given function until it returns 'true' `threshold` times in a row.
// Each failed attempt has a 1s delay; successful attempts have no delay.
func AwaitConvergence(t *testing.T, threshold int, maxTimeToConsistency time.Duration, fn func(elapsed time.Duration) bool) {
	successes := 0
	attempts := 0
	start := time.Now()
	to := time.After(maxTimeToConsistency)
	delay := time.Second
	for {
		select {
		case <-to:
			t.Fatalf("timeout while waiting after %d attempts", attempts)
		default:
		}

		completed := fn(time.Now().Sub(start))
		attempts++
		if completed {
			successes++
			if successes >= threshold {
				return
			}
			// Skip delay if we have a success
			continue
		}

		successes = 0
		select {
		// Capture the overall timeout
		case <-to:
			t.Fatalf("timeout while waiting after %d attempts, %d/%d successes", attempts, successes, threshold)
			// And the per-try delay
		case <-time.After(delay):
		}
	}
}

// WaitForConsistentResponse repeats the provided request until it completes with a response having
// the expected response consistently. The provided threshold determines how many times in
// a row this must occur to be considered "consistent".
func WaitForConsistentResponse(t *testing.T, r roundtripper.RoundTripper, req roundtripper.Request, expected ExpectedResponse, threshold int, maxTimeToConsistency time.Duration) {
	AwaitConvergence(t, threshold, maxTimeToConsistency, func(elapsed time.Duration) bool {
		cReq, cRes, err := r.CaptureRoundTrip(req)
		if err != nil {
			t.Logf("Request failed, not ready yet: %v (after %v)", err.Error(), elapsed)
			return false
		}

		if err := CompareRequest(&req, cReq, cRes, expected); err != nil {
			t.Logf("Response expectation failed for request: %v  not ready yet: %v (after %v)", req, err, elapsed)
			return false
		}

		return true
	})
	t.Logf("Request passed")
}

// CompareRequest compares the expected request and the captured request.
func CompareRequest(req *roundtripper.Request, cReq *roundtripper.CapturedRequest, cRes *roundtripper.CapturedResponse, expected ExpectedResponse) error {
	if expected.Response.StatusCode != cRes.StatusCode {
		return fmt.Errorf("expected status code to be %d, got %d", expected.Response.StatusCode, cRes.StatusCode)
	}
	if cRes.StatusCode == 200 {
		// The request expected to arrive at the backend is
		// the same as the request made, unless otherwise
		// specified.
		if expected.ExpectedRequest == nil {
			expected.ExpectedRequest = &ExpectedRequest{Request: expected.Request}
		}

		if expected.TestCaseName == "FetchAPIDefinition" {
			if len(cRes.APIDefinition) <= 0 {
				return fmt.Errorf("expected api definition should not be empty")
			}
			return nil
		}

		if expected.TestCaseName == "FetchNonExistingAPIDefinition" {
			if !cRes.IsError && !strings.Contains(cRes.ErrorMsg, "API Definition not found") {
				return fmt.Errorf("expected error response")
			}
			return nil
		}

		if expected.ExpectedRequest.Method == "" {
			expected.ExpectedRequest.Method = "GET"
		}

		if expected.ExpectedRequest.Host != "" && expected.ExpectedRequest.Host != cReq.Host {
			return fmt.Errorf("expected host to be %s, got %s", expected.ExpectedRequest.Host, cReq.Host)
		}

		if expected.ExpectedRequest.Path != cReq.Path {
			return fmt.Errorf("expected path to be %s, got %s", expected.ExpectedRequest.Path, cReq.Path)
		}
		if expected.ExpectedRequest.Method != "OPTIONS" && expected.ExpectedRequest.Method != cReq.Method {
			return fmt.Errorf("expected method to be %s, got %s", expected.ExpectedRequest.Method, cReq.Method)
		}
		if expected.Namespace != cReq.Namespace {
			return fmt.Errorf("expected namespace to be %s, got %s", expected.Namespace, cReq.Namespace)
		}
		if expected.ExpectedRequest.Headers != nil {
			if cReq.Headers == nil {
				return fmt.Errorf("no headers captured, expected %v", len(expected.ExpectedRequest.Headers))
			}
			for name, val := range cReq.Headers {
				cReq.Headers[strings.ToLower(name)] = val
			}
			for name, expectedVal := range expected.ExpectedRequest.Headers {
				actualVal, ok := cReq.Headers[strings.ToLower(name)]
				if strings.EqualFold(name, backendJWTHeader) {
					if !ok {
						return fmt.Errorf("expected %s header to be set by the enforcer", name)
					}
					if actualVal == nil || actualVal[0] == "" {
						return fmt.Errorf("expected %s header value should not be null", name)
					}
					continue
				}
				if !ok {
					return fmt.Errorf("expected %s header to be set, actual headers: %v", name, cReq.Headers)
				} else if strings.Join(actualVal, ",") != expectedVal {
					return fmt.Errorf("expected %s header to be set to %s, got %s", name, expectedVal, strings.Join(actualVal, ","))
				}

			}
		}

		if expected.Response.Headers != nil {
			if cRes.Headers == nil {
				return fmt.Errorf("no headers captured, expected %v", len(expected.ExpectedRequest.Headers))
			}
			for name, val := range cRes.Headers {
				cRes.Headers[strings.ToLower(name)] = val
			}

			for name, expectedVal := range expected.Response.Headers {
				actualVal, ok := cRes.Headers[strings.ToLower(name)]
				if !ok {
					return fmt.Errorf("expected %s header to be set, actual headers: %v", name, cRes.Headers)
				} else if strings.Join(actualVal, ",") != expectedVal {
					return fmt.Errorf("expected %s header to be set to %s, got %s", name, expectedVal, strings.Join(actualVal, ","))
				}
			}
		}

		if len(expected.Response.AbsentHeaders) > 0 {
			for name, val := range cRes.Headers {
				cRes.Headers[strings.ToLower(name)] = val
			}

			for _, name := range expected.Response.AbsentHeaders {
				val, ok := cRes.Headers[strings.ToLower(name)]
				if ok {
					return fmt.Errorf("expected %s header to not be set, got %s", name, val)
				}
			}

		}

		// Verify that headers expected *not* to be present on the
		// request are actually not present.
		if len(expected.ExpectedRequest.AbsentHeaders) > 0 {
			for name, val := range cReq.Headers {
				cReq.Headers[strings.ToLower(name)] = val
			}

			for _, name := range expected.ExpectedRequest.AbsentHeaders {
				val, ok := cReq.Headers[strings.ToLower(name)]
				if ok {
					return fmt.Errorf("expected %s header to not be set, got %s", name, val)
				}
			}
		}

		if !strings.HasPrefix(cReq.Pod, expected.Backend) {
			return fmt.Errorf("expected pod name to start with %s, got %s", expected.Backend, cReq.Pod)
		}
	} else if roundtripper.IsRedirect(cRes.StatusCode) {
		if expected.RedirectRequest == nil {
			return nil
		}

		setRedirectRequestDefaults(req, cRes, &expected)

		if expected.RedirectRequest.Host != cRes.RedirectRequest.Host {
			return fmt.Errorf("expected redirected hostname to be %s, got %s", expected.RedirectRequest.Host, cRes.RedirectRequest.Host)
		}

		if expected.RedirectRequest.Port != cRes.RedirectRequest.Port {
			return fmt.Errorf("expected redirected port to be %s, got %s", expected.RedirectRequest.Port, cRes.RedirectRequest.Port)
		}

		if expected.RedirectRequest.Scheme != cRes.RedirectRequest.Scheme {
			return fmt.Errorf("expected redirected scheme to be %s, got %s", expected.RedirectRequest.Scheme, cRes.RedirectRequest.Scheme)
		}

		if expected.RedirectRequest.Path != cRes.RedirectRequest.Path {
			return fmt.Errorf("expected redirected path to be %s, got %s", expected.RedirectRequest.Path, cRes.RedirectRequest.Path)
		}
	}
	return nil
}

// GetTestCaseName gets the user-defined test case name or generates one from expected response to a given request.
func (er *ExpectedResponse) GetTestCaseName(i int) string {

	// If TestCase name is provided then use that or else generate one.
	if er.TestCaseName != "" {
		return er.TestCaseName
	}

	headerStr := ""
	reqStr := ""

	if er.Request.Headers != nil {
		headerStr = " with headers"
	}

	reqStr = fmt.Sprintf("%d request to '%s%s'%s", i, er.Request.Host, er.Request.Path, headerStr)

	if er.Backend != "" {
		return fmt.Sprintf("%s should go to %s", reqStr, er.Backend)
	}
	return fmt.Sprintf("%s should receive a %d", reqStr, er.Response.StatusCode)
}

func setRedirectRequestDefaults(req *roundtripper.Request, cRes *roundtripper.CapturedResponse, expected *ExpectedResponse) {
	// If the expected host is nil it means we do not test host redirect.
	// In that case we are setting it to the one we got from the response because we do not know the ip/host of the gateway.
	if expected.RedirectRequest.Host == "" {
		expected.RedirectRequest.Host = cRes.RedirectRequest.Host
	}

	if expected.RedirectRequest.Port == "" {
		expected.RedirectRequest.Port = req.URL.Port()
	}

	if expected.RedirectRequest.Scheme == "" {
		expected.RedirectRequest.Scheme = req.URL.Scheme
	}

	if expected.RedirectRequest.Path == "" {
		expected.RedirectRequest.Path = req.URL.Path
	}
}

// GetTestToken get test token from test token endpoint call
func GetTestToken(t *testing.T, scopes ...string) string {
	t.Helper()
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("POST", "https://localhost:9095/testkey",
		strings.NewReader(fmt.Sprintf("scope=%s", strings.Join(scopes, " "))))

	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4=")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = "localhost"
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Logf("failed to get token: %v retrying after 100s ...", err)
		time.Sleep(100 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("failed to get token: %v", err)
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Logf("Status is: %v, failed to get token: %v retrying after 100s ...", resp.StatusCode, err)
		time.Sleep(100 * time.Second)
		resp, err = client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status is: %v, failed to get token: %v", resp.StatusCode, err)
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read token response: %v", err)
	}
	return string(body)
}
