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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/iszkiba/yare/v1"
)

func TestParseJWT(t *testing.T) {
	t.Parallel()

	type args struct {
		in string
	}

	tests := []struct {
		name    string
		args    args
		want    yare.Dict
		wantErr bool
	}{
		{
			name:    "unsigned",
			args:    args{in: "eyJhbGciOiJub25lIn0.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ."},
			wantErr: false,
			want: yare.Dict{
				"header":   yare.Dict{"alg": "none"},
				"payload":  yare.Dict{"iss": "joe", "exp": json.Number("1300819380"), "http://example.com/is_root": true},
				"verified": false,
			},
		},
		{
			name:    "invalid",
			args:    args{in: "invalid"},
			wantErr: true,
		},
		{
			name:    "empty",
			args:    args{in: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := yare.ParseJWT([]byte(tt.args.in))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJWT() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		in string
	}

	tests := []struct {
		name    string
		args    args
		want    yare.Dict
		wantErr bool
	}{
		{
			name: "normal",
			args: args{in: `{"iss":"joe","exp":1300819380,"http://example.com/is_root":true}`},
			want: yare.Dict{"iss": "joe", "exp": json.Number("1300819380"), "http://example.com/is_root": true},
		},
		{
			name:    "invalid",
			args:    args{in: "not a json"},
			wantErr: true,
		},
		{
			name:    "empty",
			args:    args{in: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := yare.ParseJSON([]byte(tt.args.in))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
