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
	"reflect"
	"testing"

	"github.com/iszkiba/yare/v1"
)

func TestMapValues(t *testing.T) {
	t.Parallel()

	type args struct {
		in map[string][]string
	}

	tests := []struct {
		name string
		args args
		want yare.Dict
	}{
		{
			name: "nil",
			args: args{in: nil},
			want: yare.Dict{},
		},
		{
			name: "single",
			args: args{in: map[string][]string{"foo": {"foo"}}},
			want: yare.Dict{"foo": "foo"},
		},
		{
			name: "multi",
			args: args{in: map[string][]string{"foo": {"foo", "bar"}}},
			want: yare.Dict{"foo": []string{"foo", "bar"}},
		},
		{
			name: "mixed",
			args: args{in: map[string][]string{"foo": {"foo", "bar"}, "bar": {"bar"}}},
			want: yare.Dict{"foo": []string{"foo", "bar"}, "bar": "bar"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := yare.MapValues(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
