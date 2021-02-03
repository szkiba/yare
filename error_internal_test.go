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
	"errors"
	"testing"
)

var (
	errDummy     = errors.New("dummy")
	errFoo       = errors.New("foo")
	errWithComma = errors.New("foo,bar")
)

func Test_wrapError(t *testing.T) {
	t.Parallel()

	type args struct {
		err []error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name:    "empty",
			args:    args{err: []error{}},
			wantErr: true,
			want:    "parse error,",
		},
		{
			name:    "single",
			args:    args{err: []error{errDummy}},
			wantErr: true,
			want:    "parse error,dummy",
		},
		{
			name:    "multiple",
			args:    args{err: []error{errDummy, errFoo}},
			wantErr: true,
			want:    "parse error,dummy,foo",
		},
		{
			name:    "wrapped",
			args:    args{err: []error{wrapError(errFoo)}},
			wantErr: true,
			want:    "parse error,foo",
		},
		{
			name:    "comma",
			args:    args{err: []error{errWithComma}},
			wantErr: true,
			want:    "parse error,\"foo,bar\"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := wrapError(tt.args.err...); (err != nil) != tt.wantErr {
				t.Errorf("wrapError() error = %v, wantErr %v", err, tt.wantErr)
			} else if err.Error() != tt.want {
				t.Errorf("wrapError() = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}
