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
	"reflect"
	"testing"
)

func TestRegisterContentType(t *testing.T) {
	t.Parallel()

	dummmyParser := func(in []byte) (Dict, error) {
		return Dict{"dummy": string(in)}, nil
	}

	type args struct {
		cty    string
		parser ParserFunc
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid",
			args:    args{cty: "foo bar"},
			wantErr: true,
		},
		{
			name:    "normal",
			args:    args{cty: "application/dummy", parser: dummmyParser},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := RegisterContentType(tt.args.cty, tt.args.parser); (err != nil) != tt.wantErr {
				t.Errorf("RegisterContentType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			got, err := parseContent(tt.args.cty, []byte("hello"))
			if err != nil || !reflect.DeepEqual(got, Dict{"dummy": "hello"}) {
				t.Errorf("Registration failed for %v", tt.args.cty)
			}
		})
	}
}

func TestRegisterAuthScheme(t *testing.T) {
	t.Parallel()

	dummmyParser := func(in []byte) (Dict, error) {
		return Dict{"dummy": string(in)}, nil
	}

	type args struct {
		scheme string
		parser ParserFunc
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{scheme: "dummy", parser: dummmyParser},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			RegisterAuthScheme(tt.args.scheme, tt.args.parser)

			got, err := parseAuth(tt.args.scheme, "hello")
			if err != nil || !reflect.DeepEqual(got, Dict{"dummy": "hello"}) {
				t.Errorf("Registration failed for %v", tt.args.scheme)
			}
		})
	}
}

func registerTestContentParsers(t *testing.T) {
	t.Helper()
	name := t.Name()

	dummmyParser := func(in []byte) (Dict, error) {
		return Dict{"dummy": string(in)}, nil
	}

	nopParser := func(in []byte) (Dict, error) {
		return nil, nil
	}

	errParser := func(in []byte) (Dict, error) {
		return nil, errDummy
	}

	_ = RegisterContentType(name+"/Content", dummmyParser)
	_ = RegisterContentType(name+"/WithNil", nopParser)
	_ = RegisterContentType(name+"/WithNil", dummmyParser)
	_ = RegisterContentType(name+"/WithErr", errParser)
}

func Test_parseContent(t *testing.T) {
	registerTestContentParsers(t)

	main := t.Name()

	t.Parallel()

	type args struct {
		cty     string
		content string
	}

	tests := []struct {
		name    string
		args    args
		want    Dict
		wantErr bool
	}{
		{
			name: "normal", args: args{cty: main + "/Content", content: "Hi"},
			want: Dict{"dummy": "Hi"},
		},
		{
			name: "suffix", args: args{cty: main + "/foo+Content", content: "Hi"},
			want: Dict{"dummy": "Hi"},
		},
		{
			name: "prefix", args: args{cty: main + "/Content+foo", content: "Hi"},
			want: Dict{"dummy": "Hi"},
		},
		{
			name: "unknown", args: args{cty: main + "/Unknown", content: "Hi"}, want: nil,
		},
		{
			name: "badtype", args: args{cty: "?", content: "Hi"}, wantErr: true,
		},
		{
			name: "empty", args: args{cty: main + "/WithNil", content: "Hi"},
			want: Dict{"dummy": "Hi"},
		},
		{
			name: "error", args: args{cty: main + "/WithErr", content: "Hi"}, wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseContent(tt.args.cty, []byte(tt.args.content))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContent() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func registerTestAuthParsers(t *testing.T) {
	t.Helper()
	name := t.Name()

	dummmyParser := func(in []byte) (Dict, error) {
		return Dict{"dummy": string(in)}, nil
	}

	nopParser := func(in []byte) (Dict, error) {
		return nil, nil
	}

	errParser := func(in []byte) (Dict, error) {
		return nil, errDummy
	}

	RegisterAuthScheme(name+"Auth", dummmyParser)
	RegisterAuthScheme(name+"AuthWithNil", nopParser)
	RegisterAuthScheme(name+"AuthWithNil", dummmyParser)
	RegisterAuthScheme(name+"AuthWithErr", errParser)
}

func Test_parseAuth(t *testing.T) {
	registerTestAuthParsers(t)

	scheme := t.Name() + "Auth"

	t.Parallel()

	type args struct {
		scheme      string
		credentials string
	}

	tests := []struct {
		name    string
		args    args
		want    Dict
		wantErr bool
	}{
		{
			name: "normal",
			args: args{scheme: scheme, credentials: "hello"},
			want: Dict{"dummy": "hello"},
		},
		{
			name: "unknown",
			args: args{scheme: "unknown", credentials: "hello"},
			want: nil,
		},
		{
			name: "empty",
			args: args{scheme: scheme + "WithNil", credentials: "hello"},
			want: Dict{"dummy": "hello"},
		},
		{
			name:    "error",
			args:    args{scheme: scheme + "WithErr", credentials: "hello"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseAuth(tt.args.scheme, tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAuth() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseContentType(t *testing.T) {
	t.Parallel()

	type args struct {
		cty string
	}

	tests := []struct {
		name    string
		args    args
		main    string
		sub     string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{cty: "application/json"},
			main: "application",
			sub:  "json",
		},
		{
			name: "partial",
			args: args{cty: "application"},
			main: "",
			sub:  "",
		},
		{
			name:    "invalid",
			args:    args{cty: "?"},
			wantErr: true,
		},
		{
			name:    "empty",
			args:    args{cty: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, got1, err := parseContentType(tt.args.cty)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContentType() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.main {
				t.Errorf("parseContentType() got = %v, want %v", got, tt.main)
			}
			if got1 != tt.sub {
				t.Errorf("parseContentType() got1 = %v, want %v", got1, tt.sub)
			}
		})
	}
}
