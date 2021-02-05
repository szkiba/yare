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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/szkiba/yare"
)

var version = "dev"

type options struct {
	port    int
	version bool
}

const defaultPort = 8080

func getopt(args []string) *options {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	o := options{port: defaultPort}

	if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		o.port = p
	}

	flags.IntVar(&o.port, "port", o.port, "port to listen on")

	ver := flags.Bool("v", false, "prints version")

	_ = flags.Parse(args[1:])

	o.version = *ver

	return &o
}

func main() {
	o := getopt(os.Args)

	if o.version {
		fmt.Fprintf(os.Stderr, "yare/%s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	http.Handle("/", yare.EchoHander(true))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", o.port), nil))
}

func init() {
	_ = yare.RegisterContentType("application/json", yare.ParseJSON)
	_ = yare.RegisterContentType("application/jwt", yare.ParseJWT)

	yare.RegisterAuthScheme("Bearer", yare.ParserFunc(yare.ParseJWT).Optional())
}
