package esi // import "corpus.space/esi"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
)

// for testing
var now = time.Now

type endpoint struct {
	api *Client
}

const (
	// DefaultBaseURL is the URL of the default public ESI API.
	DefaultBaseURL = "https://esi.evetech.net/"

	// DefaultUserAgent is the value to use in the User-Agent header if none
	// has been explicitly configured.
	DefaultUserAgent = "corpus.space/esi"

	headerErrorRateRemaining = "X-ESI-Error-Limit-Remain"
	headerErrorRateReset     = "X-ESI-Error-Limit-Reset"
)

// A Client handles communication with the EVE Online Swagger Interface (ESI)
// API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public EVE Online ESI API.
	BaseURL *url.URL

	// User agent used when communicating with ESI. You should set this.
	UserAgent string

	// Logging holds optional loggers. If any are nil, logging is done via the
	// log package's standard logger.
	Logging struct {
		Info, Error, Debug *log.Logger
	}

	mu struct {
		sync.Mutex
		Rate
	}

	common endpoint // reuse a single struct for all endpoints

	// Endpoints for talking to different parts of ESI.
	Fleets *FleetsEndpoint
}

// NewClient returns a new ESI API client. If a nil httpClient is provided,
// http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
//
// It would be prudent to use a caching transport such as the excellent
// httpcache by gregjones.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

	api := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: DefaultUserAgent,
	}

	api.common.api = api

	// endpoints
	api.Fleets = (*FleetsEndpoint)(&api.common)

	return api
}

// NewRequest creates an API request. An url relative to the BaseURL of the
// client is provided. If body is non-nil, it is encoded to JSON and included
// in the request body.
func (api *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	u, err := api.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)

		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	return req, nil
}

// Response is an EVE Online ESI API response. This wraps the standard
// http.Response returned from ESI and provides convenient access to deprecation
// warnings and rate limit information.
type Response struct {
	*http.Response
}

func makeResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

// Error represents an ESI API error.
type Error struct {
	Response       *http.Response
	HTTPStatusCode int
	Err            string `json:"error"`

	Rate
}

func (e Error) Error() string {
	return e.Err
}

func makeError(r *http.Response) *Error {
	e := &Error{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, e)
	}

	e.HTTPStatusCode = r.StatusCode
	e.Rate = parseRate(r)

	return e
}

// Rate represents rate limit information.
type Rate struct {
	Remaining int
	Reset     time.Time
}

func (r Rate) String() string {
	return fmt.Sprintf("error rate limit: %d remaining calls; reset in %.fs", r.Remaining, r.Reset.Sub(now()).Seconds())
}

func parseRate(r *http.Response) Rate {
	var rate Rate
	if remaining := r.Header.Get(headerErrorRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}

	if reset := r.Header.Get(headerErrorRateReset); reset != "" {
		if v, _ := strconv.Atoi(reset); v != 0 {
			rate.Reset = now().Add(time.Duration(v) * time.Second)
		}
	}

	return rate
}

// Do carries out a request and stores the result in v.
func (api *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	// send request
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	// deferred closing of response body
	defer resp.Body.Close()

	response := makeResponse(resp)

	if err := api.check(resp); err != nil {
		api.mu.Lock()
		api.mu.Rate = err.(*Error).Rate
		api.mu.Unlock()

		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				if err == io.EOF {
					err = nil
				}

				return response, err
			}
		}
	}

	return response, nil
}

func (api *Client) check(resp *http.Response) error {
	if rc := resp.StatusCode; 200 <= rc && rc <= 299 {
		// check for any waning headers and log them
		if v := resp.Header.Get("warning"); v != "" {
			logf(api.Logging.Error, "warning header received (%s %v): %s",
				resp.Request.Method, resp.Request.URL.Path, v,
			)
		}

		return nil
	}

	return makeError(resp)
}

func logf(logger *log.Logger, format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
		return
	}

	log.Printf(format, args...)
}

// I18NOptions specifies optional parameters to various methods that support
// internationalization.
type I18NOptions struct {
	Language string `url:"language"`
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Bool is a helper function that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper function that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper function that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Float64 is a helper function that allocates a new float64 value
// to store v and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// String is a helper function that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
