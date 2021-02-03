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

package yare

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_mapCookies(t *testing.T) {
	t.Parallel()

	type args struct {
		cookies []*http.Cookie
	}

	tests := []struct {
		name string
		args args
		want Dict
	}{
		{
			name: "normal",
			args: args{cookies: []*http.Cookie{{Name: "foo", Value: "bar"}}},
			want: Dict{"foo": "bar"},
		},
		{
			name: "empty",
			args: args{cookies: []*http.Cookie{}},
			want: Dict{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mapCookies(tt.args.cookies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapCookies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAuthorizationHeader(t *testing.T) {
	t.Parallel()

	name := t.Name()

	dummmyParser := func(in []byte) (Dict, error) {
		if len(in) == 0 {
			return nil, errDummy
		}

		return Dict{"dummy": string(in)}, nil
	}

	RegisterAuthScheme(name, dummmyParser)

	type args struct {
		hdr string
	}

	tests := []struct {
		name    string
		args    args
		want    Dict
		wantErr bool
	}{
		{
			name: "normal",
			args: args{hdr: name + " passed"},
			want: Dict{name: Dict{"dummy": "passed"}},
		},
		{
			name: "empty",
			args: args{hdr: ""},
			want: nil,
		},
		{
			name:    "error",
			args:    args{hdr: name + ""},
			wantErr: true,
		},
		{
			name: "unknown",
			args: args{hdr: "unknown foo"},
			want: Dict{"unknown": "foo"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseAuthorizationHeader(tt.args.hdr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAuthorizationHeader() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAuthorizationHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRequestBody(t *testing.T) {
	t.Parallel()

	name := t.Name()

	dummmyParser := func(in []byte) (Dict, error) {
		if len(in) == 0 {
			return nil, errDummy
		}

		return Dict{"dummy": string(in)}, nil
	}

	_ = RegisterContentType("test/"+name, dummmyParser)

	type args struct {
		r *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    Dict
		wantErr bool
	}{
		{
			name: "normal",
			args: args{r: httptest.NewRequest("POST", "http://localhost", bytes.NewBufferString("Ready"))},
			want: Dict{"dummy": "Ready"},
		},
		{
			name:    "error",
			args:    args{r: httptest.NewRequest("POST", "http://localhost", bytes.NewBufferString(""))},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.args.r.Header.Set("Content-Type", "test/"+name)
			got, err := parseRequestBody(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRequestBody() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseResponseBody(t *testing.T) {
	t.Parallel()

	name := t.Name()

	dummmyParser := func(in []byte) (Dict, error) {
		if len(in) == 0 {
			return nil, errDummy
		}

		return Dict{"dummy": string(in)}, nil
	}

	_ = RegisterContentType("test/"+name, dummmyParser)

	type args struct {
		body string
	}

	tests := []struct {
		name    string
		args    args
		cty     string
		want    Dict
		wantErr bool
	}{
		{
			name: "normal",
			args: args{body: "Hi"},
			cty:  "test/" + name,
			want: Dict{"dummy": "Hi"},
		},
		{
			name:    "error",
			args:    args{body: "Hi"},
			cty:     "?",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			rec.Body = bytes.NewBufferString(tt.args.body)
			resp := rec.Result()

			resp.Header.Set("Content-Type", tt.cty)
			defer resp.Body.Close()

			got, err := parseResponseBody(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponseBody() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseResponseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

type zeroErrReader struct {
	err error
}

func (r zeroErrReader) Read(p []byte) (int, error) {
	return copy(p, []byte{0}), r.err
}

func Test_wrapReader(t *testing.T) {
	t.Parallel()

	type args struct {
		r io.ReadCloser
	}

	tests := []struct {
		name    string
		args    args
		data    []byte
		wantErr bool
	}{
		{
			name: "normal",
			args: args{r: ioutil.NopCloser(bytes.NewBufferString("dummy"))},
			data: []byte("dummy"),
		},
		{
			name:    "error",
			args:    args{r: ioutil.NopCloser(zeroErrReader{errDummy})},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader, data, err := wrapReader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapReader() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if reader == nil {
				t.Error("wrapReader() returns nil reader")
			}

			if !reflect.DeepEqual(data, tt.data) {
				t.Errorf("wrapReader() data = %v, want %v", data, tt.data)
			}
		})
	}
}

func Test_omitEmpty(t *testing.T) {
	t.Parallel()

	type args struct {
		d Dict
	}

	tests := []struct {
		name string
		args args
		want Dict
	}{
		{
			name: "noempty",
			args: args{d: Dict{"foo": "bar"}},
			want: Dict{"foo": "bar"},
		},
		{
			name: "empty",
			args: args{d: Dict{}},
			want: nil,
		},
		{
			name: "nil",
			args: args{d: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := omitEmpty(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("omitEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
