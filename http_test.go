// MIT License
//
// Copyright (c) 2021 Iván Szkiba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package yare_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/iszkiba/yare/v1"
)

type kv map[string]string

type par struct {
	header kv
	body   string
	method string
	url    string
}

func newRequest(in par) *http.Request {
	if in.url == "" {
		in.url = "http://localhost/"
	}

	r := httptest.NewRequest(in.method, in.url, bytes.NewBufferString(in.body))

	if in.header != nil {
		for k, v := range in.header {
			r.Header.Set(k, v)
		}
	}

	return r
}

func newResponse(in par) *http.Response {
	if in.url == "" {
		in.url = "http://localhost/"
	}

	rec := httptest.NewRecorder()
	r := rec.Result()

	if in.header != nil {
		for k, v := range in.header {
			r.Header.Set(k, v)
		}
	}

	if in.body != "" {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(in.body))
	}

	return r
}

func registerJSON(t *testing.T) string {
	t.Helper()

	cty := "test/" + t.Name()
	_ = yare.RegisterContentType(cty, yare.ParseJSON)

	return cty
}

func registerAuth(t *testing.T) string {
	t.Helper()

	scheme := t.Name()
	yare.RegisterAuthScheme(scheme, yare.ParseJWT)

	return scheme
}

func TestMapRequest(t *testing.T) {
	t.Parallel()

	cty := registerJSON(t)

	tests := []struct {
		name string
		par  par
		want yare.Dict
	}{
		{
			name: "minimal", par: par{method: "GET"}, want: yare.Dict{"method": "GET", "version": "HTTP/1.1", "path": "/"},
		},
		{
			name: "cookies", par: par{method: "GET", header: kv{"Cookie": "foo=bar"}},
			want: yare.Dict{
				"method": "GET", "version": "HTTP/1.1", "path": "/",
				"headers": yare.Dict{"Cookie": "foo=bar"},
				"cookies": yare.Dict{"foo": "bar"},
			},
		},
		{
			name: "parameters", par: par{
				method: "PUT", url: "http://localhost/?dummy=yes",
				header: kv{"Content-Type": "application/x-www-form-urlencoded"}, body: "foo=bar",
			},
			want: yare.Dict{
				"method": "PUT", "version": "HTTP/1.1", "path": "/",
				"headers": yare.Dict{"Content-Type": "application/x-www-form-urlencoded"},
				"form":    yare.Dict{"foo": "bar"}, "query": yare.Dict{"dummy": "yes"},
			},
		},
		{
			name: "normal", par: par{
				method: http.MethodPost, body: "{\"foo\":\"bar\"}",
				header: kv{"Content-Type": cty, "Authorization": "Bearer dummy"},
			},
			want: yare.Dict{
				"method": http.MethodPost, "version": "HTTP/1.1", "path": "/",
				"headers":       yare.Dict{"Content-Type": cty, "Authorization": "Bearer dummy"},
				"authorization": yare.Dict{"Bearer": "dummy"}, "body": yare.Dict{"foo": "bar"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := yare.MapRequest(newRequest(tt.par), tt.par.method == http.MethodPost)
			if err != nil {
				t.Errorf("MapRequest() error = %v", err)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapRequestError(t *testing.T) {
	t.Parallel()

	cty := registerJSON(t)

	scheme := registerAuth(t)

	tests := []struct {
		name string
		par  par
	}{
		{
			name: "parameters", par: par{
				method: "PUT",
				header: kv{"Content-Type": "application/x-www-form-urlencoded"},
				body:   "%1",
			},
		},
		{
			name: "body", par: par{
				method: http.MethodPost,
				header: kv{"Content-Type": cty},
				body:   "{\"foo\"",
			},
		},
		{
			name: "authorization", par: par{
				method: "GET",
				header: kv{"Authorization": scheme + " dummy"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := yare.MapRequest(newRequest(tt.par), tt.par.method == http.MethodPost)
			if err == nil {
				t.Error("MapRequest() error is nil")
			}
		})
	}
}

func TestMapResponse(t *testing.T) {
	t.Parallel()

	cty := registerJSON(t)

	tests := []struct {
		name string
		par  par
		want yare.Dict
	}{
		{
			name: "minimal", par: par{}, want: yare.Dict{"status": 200, "version": "HTTP/1.1"},
		},
		{
			name: "normal", par: par{
				body:   "{\"foo\":\"bar\"}",
				header: kv{"Content-Type": cty},
				method: http.MethodPost,
			},
			want: yare.Dict{
				"status":  200,
				"version": "HTTP/1.1",
				"headers": yare.Dict{"Content-Type": cty},
				"body":    yare.Dict{"foo": "bar"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := newResponse(tt.par)
			defer resp.Body.Close()
			got, err := yare.MapResponse(resp, tt.par.method == http.MethodPost)
			if err != nil {
				t.Errorf("MapResponse() error = %v", err)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapResposeError(t *testing.T) {
	t.Parallel()

	cty := registerJSON(t)

	tests := []struct {
		name string
		par  par
	}{
		{
			name: "body", par: par{
				method: http.MethodPost,
				header: kv{"Content-Type": cty},
				body:   "{\"foo\"",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := newResponse(tt.par)
			defer resp.Body.Close()
			_, err := yare.MapResponse(resp, tt.par.method == http.MethodPost)
			if err == nil {
				t.Error("MapResponse() error is nil")
			}
		})
	}
}
