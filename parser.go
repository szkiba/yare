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
	"encoding/json"
	"mime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dgrijalva/jwt-go"
)

// ParserFunc used to register custom body and authorization parsers.
type ParserFunc func([]byte) (Dict, error)

// Optional returns a new ParserFunc which silently ignore errors.
func (p ParserFunc) Optional() ParserFunc {
	return func(in []byte) (Dict, error) {
		dict, _ := p(in)

		return dict, nil
	}
}

type contentType struct {
	main   string
	sub    string
	parser ParserFunc
}

type authScheme struct {
	scheme string
	parser ParserFunc
}

var (
	contentTypeMu      sync.Mutex
	atomicContentTypes atomic.Value
	authSchemeMu       sync.Mutex
	atomicAutScheme    atomic.Value
)

// RegisterContentType registers custom content parser for a given Content-Type.
func RegisterContentType(cty string, parser ParserFunc) error {
	main, sub, err := parseContentType(cty)
	if err != nil {
		return err
	}

	contentTypeMu.Lock()
	values, _ := atomicContentTypes.Load().([]contentType)
	atomicContentTypes.Store(append(values, contentType{main, sub, parser}))
	contentTypeMu.Unlock()

	return nil
}

// RegisterAuthScheme registers custom authentication credentials parser for a given scheme.
func RegisterAuthScheme(scheme string, parser ParserFunc) {
	authSchemeMu.Lock()
	values, _ := atomicAutScheme.Load().([]authScheme)
	atomicAutScheme.Store(append(values, authScheme{scheme, parser}))
	authSchemeMu.Unlock()
}

// ParseJWT is a ParserFunc for parsing JWT to Dict.
//
// Can use as RegisterContentType or RegisterAuthScheme parser argument.
func ParseJWT(in []byte) (Dict, error) {
	p := jwt.Parser{}
	p.UseJSONNumber = true
	t, err := p.Parse(string(in), func(tt *jwt.Token) (interface{}, error) {
		return tt, nil
	})

	if t == nil || t.Claims == nil {
		return nil, wrapError(err)
	}

	// must be ok because it is just marshall/unmarshall jwt.Claims
	data, _ := json.Marshal(t.Claims)
	claims, _ := ParseJSON(data)

	return Dict{"header": Dict(t.Header), "payload": claims, "verified": false}, nil
}

// ParseJSON is a ParserFunc for parsing JSON string to Dict.
//
// Can use as RegisterContentType parser argument.
func ParseJSON(in []byte) (Dict, error) {
	out := make(Dict)

	d := json.NewDecoder(bytes.NewBuffer(in))
	d.UseNumber()

	if err := d.Decode(&out); err != nil {
		return nil, wrapError(err)
	}

	return out, nil
}

func parseContent(cty string, content []byte) (Dict, error) {
	main, sub, err := parseContentType(cty)
	if err != nil {
		return nil, err
	}

	contentTypeMu.Lock()
	contentTypes, _ := atomicContentTypes.Load().([]contentType)
	contentTypeMu.Unlock()

	for _, c := range contentTypes {
		if c.main == main && (strings.HasPrefix(sub, c.sub) || strings.HasSuffix(sub, c.sub)) {
			if dict, err := c.parser(content); err != nil || dict != nil {
				return dict, err
			}
		}
	}

	return nil, nil
}

func parseAuth(scheme, credentials string) (Dict, error) {
	authSchemeMu.Lock()
	authorizations, _ := atomicAutScheme.Load().([]authScheme)
	authSchemeMu.Unlock()

	for _, a := range authorizations {
		if a.scheme == scheme {
			if dict, err := a.parser([]byte(credentials)); err != nil || dict != nil {
				return dict, err
			}
		}
	}

	return nil, nil
}

func parseContentType(cty string) (string, string, error) {
	mt, _, err := mime.ParseMediaType(cty)
	if err != nil {
		return "", "", wrapError(err)
	}

	idx := strings.Index(mt, "/")
	if idx <= 0 {
		return "", "", nil
	}

	return mt[:idx], mt[idx+1:], nil
}
