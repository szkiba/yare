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
	"encoding/json"
	"net/http"
)

// handle store internal handler configuration.
type handler struct {
	body bool
}

// EchoHander returns a handler that serves HTTP requests with Echo response.
func EchoHander(body bool) http.Handler {
	return &handler{body: body}
}

// ServeHTTP is a http handler method.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dict, err := MapRequest(r, h.body)
	status := http.StatusOK

	if err != nil {
		addError(w, err)

		status = http.StatusBadRequest
	}

	if dict == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if data, err := json.Marshal(dict); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, _ = w.Write(data)
	} else {
		addError(w, wrapError(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func addError(w http.ResponseWriter, err error) {
	w.Header().Add("X-Error", err.Error())
}
