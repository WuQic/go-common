// The httpClient package provides lower level http control for handling
// http based messages.
package httpClient

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/goinggo/tracelog"
)

type (
	// ProxyFunc is the proxy callback used for the transport.
	ProxyFunc func(*http.Request) (*url.URL, error)

	// Transport provides a thin wrapper arounf http.Transport.
	Transport struct {
		// Proxy specifies a function to return a proxy for a given
		// *http.Request. If the function returns a non-nil error, the
		// request is aborted with the provided error.
		// If Proxy is nil or returns a nil *url.URL, no proxy is used.
		Proxy func(*http.Request) (*url.URL, error)

		// TLSClientConfig specifies the TLS configuration to use with
		// tls.Client. If nil, the default configuration is used.
		TLSClientConfig *tls.Config

		// DisableKeepAlives, if true, prevents re-use of TCP connections
		// between different HTTP requests.
		DisableKeepAlives bool

		// DisableCompression, if true, prevents the Transport from
		// requesting compression with an "Accept-Encoding: gzip"
		// request header when the Request contains no existing
		// Accept-Encoding value. If the Transport requests gzip on
		// its own and gets a gzipped response, it's transparently
		// decoded in the Response.Body. However, if the user
		// explicitly requested gzip it is not automatically
		// uncompressed.
		DisableCompression bool

		// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
		// (keep-alive) to keep per-host.  If zero,
		// http.DefaultMaxIdleConnsPerHost is used.
		MaxIdleConnsPerHost int

		// ConnectTimeout, if non-zero, is the maximum amount of time a dial will wait for
		// a connect to complete.
		ConnectTimeout time.Duration

		// ResponseHeaderTimeout, if non-zero, specifies the amount of
		// time to wait for a server's response headers after fully
		// writing the request (including its body, if any). This
		// time does not include the time to read the response body.
		ResponseHeaderTimeout time.Duration

		// RequestTimeout, if non-zero, specifies the amount of time for the entire
		// request to complete (including all of the above timeouts + entire response body).
		// This should never be less than the sum total of the above two timeouts.
		RequestTimeout time.Duration

		starter   sync.Once
		transport *http.Transport
	}
)

// bodyCloseInterceptor.
type bodyCloseInterceptor struct {
	io.ReadCloser
	timer *time.Timer
}

// Version returns the current version of the package.
func Version() string {
	return "0.4.1"
}

// NewTransport creates new transport instance, passing nil to proxy func will not use a proxy,
// otherwise it will use the callback passed in.
func NewTransport(proxy ProxyFunc) *Transport {
	return &Transport{
		ConnectTimeout:        25 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		RequestTimeout:        85 * time.Second,
		Proxy:                 proxy,
	}
}

// Get implements a get request.
func (t *Transport) Get(url string) ([]byte, error) {
	return t.GetWithHeaders(url, nil)
}

// Get performs a Get request with the specified headers.
func (t *Transport) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{Transport: t}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		tracelog.ERROR(err, "http_client", "GetWithHeaders")
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return loadResponse(resp)
}

// Post performs a post request.
func (t *Transport) Post(url string, postParams url.Values) ([]byte, error) {
	return t.PostWithHeaders(url, postParams, nil)
}

// Post performs a post request with the specified headers.
func (t *Transport) PostWithHeaders(url string, postParams url.Values, headers map[string]string) ([]byte, error) {
	client := &http.Client{Transport: t}
	req, err := http.NewRequest("POST", url, strings.NewReader(postParams.Encode()))
	if err != nil {
		tracelog.ERROR(err, "http_client", "PostWithHeaders")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		tracelog.ERROR(err, "http_client", "PostWithHeaders")
		return nil, err
	}

	return loadResponse(resp)
}

// loadResponse parse a response.
func loadResponse(resp *http.Response) ([]byte, error) {
	tracelog.STARTED("http_client", "loadResponse")

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tracelog.INFO("yodlee_api", "loadResponse", "Api Response => \n\n %s \n\n", contents)

	if resp.StatusCode != 200 {
		return nil, errors.New(string(contents))
	}

	tracelog.COMPLETED("http_client", "loadResponse")
	return contents, err
}

// lazyStart.
func (t *Transport) lazyStart() {
	dialer := &net.Dialer{Timeout: t.ConnectTimeout}
	t.transport = &http.Transport{
		Dial:                  dialer.Dial,
		Proxy:                 t.Proxy,
		TLSClientConfig:       t.TLSClientConfig,
		DisableKeepAlives:     t.DisableKeepAlives,
		DisableCompression:    t.DisableCompression,
		MaxIdleConnsPerHost:   t.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: t.ResponseHeaderTimeout,
	}
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.starter.Do(t.lazyStart)

	if t.RequestTimeout > 0 {
		timer := time.AfterFunc(t.RequestTimeout, func() {
			t.transport.CancelRequest(req)
		})

		resp, err = t.transport.RoundTrip(req)
		if err != nil {
			timer.Stop()
		} else {
			resp.Body = &bodyCloseInterceptor{ReadCloser: resp.Body, timer: timer}
		}
	} else {
		resp, err = t.transport.RoundTrip(req)
	}

	return
}

// Close.
func (bci *bodyCloseInterceptor) Close() error {
	bci.timer.Stop()
	return bci.ReadCloser.Close()
}
