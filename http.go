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
	"strings"
)

// MapRequest creates Dict from various request attributes.
func MapRequest(r *http.Request, body bool) (Dict, error) {
	var err error

	out := make(Dict)
	errs := []error{}

	// HTTP protocol
	out["version"] = r.Proto
	out["method"] = r.Method

	// HTTP headers
	if d := omitEmpty(MapValues(r.Header)); d != nil {
		out["headers"] = d
	}

	// Cookies
	if d := omitEmpty(mapCookies(r.Cookies())); d != nil {
		out["cookies"] = d
	}

	// query and form params
	u := r.URL
	out["path"] = u.Path

	if d := omitEmpty(MapValues(u.Query())); d != nil {
		out["query"] = d
	}

	if err = r.ParseForm(); err == nil {
		if d := omitEmpty(MapValues(r.PostForm)); d != nil {
			out["form"] = d
		}
	} else {
		errs = append(errs, err)
	}

	// request body
	if body {
		if v, err := parseRequestBody(r); err == nil {
			if v != nil {
				out["body"] = v
			}
		} else {
			errs = append(errs, err)
		}
	}

	// authorization
	if v, err := parseAuthorizationHeader(r.Header.Get("Authorization")); err == nil {
		if v != nil {
			out["authorization"] = v
		}
	} else {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return out, wrapError(errs...)
	}

	return out, nil
}

// MapResponse creates Dict from various response attributes.
func MapResponse(r *http.Response, body bool) (Dict, error) {
	out := make(Dict)
	errs := []error{}

	// HTTP protocol
	out["version"] = r.Proto
	out["status"] = r.StatusCode

	// HTTP headers
	if d := omitEmpty(MapValues(r.Header)); d != nil {
		out["headers"] = d
	}

	// Cookies
	if d := omitEmpty(mapCookies(r.Cookies())); d != nil {
		out["cookies"] = d
	}

	// request body
	if body {
		if v, err := parseResponseBody(r); err == nil {
			if v != nil {
				out["body"] = v
			}
		} else {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return out, wrapError(errs...)
	}

	return out, nil
}

func mapCookies(cookies []*http.Cookie) Dict {
	out := make(Dict)

	if len(cookies) == 0 {
		return out
	}

	for _, cookie := range cookies {
		out[cookie.Name] = cookie.Value
	}

	return out
}

func parseAuthorizationHeader(hdr string) (Dict, error) {
	if len(hdr) == 0 {
		return nil, nil
	}

	f := strings.SplitN(hdr, " ", 2)

	scheme := f[0]

	var credentials string
	if len(f) > 1 {
		credentials = f[1]
	}

	v, err := parseAuth(scheme, credentials)
	if err != nil {
		return nil, err
	}

	var val interface{}
	if v = omitEmpty(v); v == nil {
		val = credentials
	} else {
		val = v
	}

	return Dict{scheme: val}, nil
}

func parseRequestBody(r *http.Request) (Dict, error) {
	reader, body, err := wrapReader(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = reader

	if len(body) == 0 {
		return nil, nil
	}

	out, err := parseContent(r.Header.Get("Content-Type"), body)
	if err != nil {
		return nil, err
	}

	return omitEmpty(out), nil
}

func parseResponseBody(r *http.Response) (Dict, error) {
	reader, body, err := wrapReader(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = reader

	out, err := parseContent(r.Header.Get("Content-Type"), body)
	if err != nil {
		return nil, err
	}

	return omitEmpty(out), nil
}

func wrapReader(r io.ReadCloser) (io.ReadCloser, []byte, error) {
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, wrapError(err)
	}

	return ioutil.NopCloser(bytes.NewBuffer(data)), data, nil
}

func omitEmpty(d Dict) Dict {
	if d == nil || len(d) == 0 {
		return nil
	}

	return d
}
