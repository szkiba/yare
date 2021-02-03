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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/iszkiba/yare/v1"
)

func TestEchoHandler(t *testing.T) {
	t.Parallel()

	cty := registerJSON(t)

	tests := []struct {
		name string
		par  par
		want yare.Dict
	}{
		{
			name: "parameters", par: par{
				method: http.MethodPost, url: "http://localhost/?dummy=yes",
				header: kv{"Content-Type": "application/x-www-form-urlencoded"}, body: "foo=bar",
			},
			want: yare.Dict{
				"method": http.MethodPost, "version": "HTTP/1.1", "path": "/",
				"headers": map[string]interface{}{"Content-Type": "application/x-www-form-urlencoded"},
				"form":    map[string]interface{}{"foo": "bar"}, "query": map[string]interface{}{"dummy": "yes"},
			},
		},
		{
			name: "normal", par: par{
				method: http.MethodPost, body: "{\"foo\":\"bar\"}",
				header: kv{"Content-Type": cty, "Authorization": "Bearer dummy"},
			},
			want: yare.Dict{
				"method": http.MethodPost, "version": "HTTP/1.1", "path": "/",
				"headers":       map[string]interface{}{"Content-Type": cty, "Authorization": "Bearer dummy"},
				"authorization": map[string]interface{}{"Bearer": "dummy"}, "body": map[string]interface{}{"foo": "bar"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRequest(tt.par)
			w := httptest.NewRecorder()
			handler := yare.EchoHander(tt.par.method == http.MethodPost)

			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			data, _ := ioutil.ReadAll(resp.Body)
			got, _ := yare.ParseJSON(data)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EchoHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEchoHandlerError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		par    par
		status int
	}{
		{
			name: "parameters", par: par{
				method: http.MethodPost,
				header: kv{"Content-Type": "application/x-www-form-urlencoded"}, body: "%1",
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newRequest(tt.par)
			w := httptest.NewRecorder()
			handler := yare.EchoHander(tt.par.method == http.MethodPost)

			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.status {
				t.Errorf("EchoHandler() status = %v, want %v", resp.StatusCode, tt.status)
			}
		})
	}
}
