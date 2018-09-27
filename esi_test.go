// Copyright (c) 2013 The go-github AUTHORS. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
//      notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
//      copyright notice, this list of conditions and the following disclaimer
//      in the documentation and/or other materials provided with the
//      distribution.
//    * Neither the name of Google Inc. nor the names of its
//      contributors may be used to endorse or promote products derived from
//      this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Modified 2018; Klaus Birkelund Abildgaard Jensen.

package esi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/esi"
)

func setup() (*Client, *http.ServeMux, string, func()) {
	mux := http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client := NewClient(nil)

	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("unexpected request method; got %q, want %q", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("unexpected header value; Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("error reading request body: %q", err)
	}
	if got := string(b); got != want {
		t.Errorf("unexpected request body; got %q, want %q", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}
	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected request parameters; got %q, want %q", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expected error")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("expected URL parse error; got: %+v", err)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), DefaultBaseURL; got != want {
		t.Errorf("unexpected BaseURL; got %q, want %q", got, want)
	}
	if got, want := c.UserAgent, DefaultUserAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", DefaultBaseURL+"foo"
	inBody, outBody := &CharacterPublicInfo{Name: String("pascal")}, `{"name":"pascal"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_invalidMethod(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("INVALID METHOD", "", nil)
	if err == nil {
		t.Error("Expected error to be returned.")
	}

	if !strings.Contains(err.Error(), "invalid method") {
		t.Error("Expected error to contain 'invalid method'")
	}
}
func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	client.Do(context.Background(), req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_requestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".", nil)
	req.URL = nil
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("expected error")
	}

	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestDo_jsonSyntaxError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `invalid json`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)

	if v, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError; got %q", v)
	}
}

func TestDo_ioEOF(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ``)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("expected nil error; got %q", err)
	}
}

func TestDo_ioWriter(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "some data")
	})

	req, _ := client.NewRequest("GET", ".", nil)
	var w bytes.Buffer
	client.Do(context.Background(), req, &w)

	if w.String() != "some data" {
		t.Fatalf("expected \"some data\"; got %q", w.String())
	}
}

func TestDo_addOptions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("language")
		if param == "" {
			t.Fatal("missing language query parameter")
		}

		if param != "en_US" {
			t.Fatalf("query parameter; want \"en_US\"; got %q", param)
		}
	})

	u, _ := addOptions(".", I18NOptions{"en_US"})

	req, _ := client.NewRequest("GET", u, nil)
	client.Do(context.Background(), req, nil)
}

func TestAddOptions(t *testing.T) {
	if s, _ := addOptions(".", I18NOptions{"en_US"}); s != ".?language=en_US" {
		t.Fatalf("want \".\"; got %q", s)
	}
}

func TestAddOptions_nilParam(t *testing.T) {
	var opt *I18NOptions
	if s, _ := addOptions(".", opt); s != "." {
		t.Fatalf("want \".\"; got %q", s)
	}
}

func TestAddOptions_invalidURL(t *testing.T) {
	if s, err := addOptions("invalid%", I18NOptions{"en_US"}); err == nil {
		t.Fatalf("expected error; got %q", s)
	}
}

func TestAddOptions_invalidOptions(t *testing.T) {
	if s, err := addOptions(".", 42); err == nil {
		t.Fatalf("expected error; got %q", s)
	}
}

func TestDo_warningHeadersAreLogged(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var out bytes.Buffer
	log.SetOutput(&out)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("warning", "299 - This route is deprecated.")
	})

	req, _ := client.NewRequest("GET", ".", nil)
	client.Do(context.Background(), req, nil)

	if !strings.Contains(out.String(), "deprecated") {
		t.Fatalf("deprecation warning not logged; got %q", out.String())
	}
}

func TestCustomLogger(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var out bytes.Buffer
	client.Logging.Error = log.New(&out, "", 0)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("warning", "299 - This route is deprecated.")
	})

	req, _ := client.NewRequest("GET", ".", nil)
	client.Do(context.Background(), req, nil)

	if !strings.Contains(out.String(), "deprecated") {
		t.Fatalf("deprecation warning not logged; got %q", out.String())
	}
}

func TestDo_errorRateLimit(t *testing.T) {
	client, mux, _, teardown := setup()

	staticNow := time.Date(2018, 1, 1, 18, 0, 0, 0, time.UTC)

	oldtime := now
	now = func() time.Time {
		return staticNow
	}

	defer func() {
		now = oldtime
		teardown()
	}()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-esi-error-limit-remain", "60")
		w.Header().Set("x-esi-error-limit-reset", "42")

		http.Error(w, `{"error":"some error"}`, 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(context.Background(), req, nil)
	v, ok := err.(*Error)
	if !ok {
		t.Fatal("expected an error of type *ErrorResponse")
	}

	if v.Error() != "some error" {
		t.Fatalf("expected error text \"some error\"; got %q", v.Err)
	}

	if !v.Rate.Reset.Equal(staticNow.Add(42 * time.Second)) {
		t.Fatalf("expected reset to happen in 42s; got %v", v.Rate.Reset.Sub(staticNow))
	}
}

func TestRate_string(t *testing.T) {
	r := Rate{Remaining: 60, Reset: time.Now().Add(5 * time.Second)}

	expected := "error rate limit: 60 remaining calls; reset in 5s"
	got := r.String()

	if got != expected {
		t.Fatalf("unexpected output from Rate.String(); got %q; want %q", got, expected)
	}
}
