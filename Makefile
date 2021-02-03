# MIT License
#
# Copyright (c) 2021 Iván Szkiba
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

#: Build everything (default target)
all: build test

#: Install all the build and lint dependencies
setup:
	@go mod download
	@go get github.com/google/addlicense
	@go mod tidy
	@addlicense -f LICENSE . cmd
.PHONY: setup

#: Run all the tests
test:
	@CGO_ENABLED=0 go test ./... -coverprofile=coverage.txt
.PHONY: test

#: Build server
build:
	@goreleaser  build --snapshot --rm-dist
.PHONY: build

#: Genereate test coverage report
cover: test
	@go tool cover -html=coverage.txt
.PHONY: cover

#: Run all the linters
lint:
	golangci-lint run ./...
.PHONY: lint

#: Run all the tests and code checks
ci: lint test build
.PHONY: ci

#: Execute server
run:
	@go run ./cmd/yare/main.go
.PHONY: run

#: Clean up working directory
clean:
	@rm -rf dist
.PHONY: clean

#: Print this help
help:
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
		| grep -v -- -- \
		| sed 'N;s/\n/###/' \
		| sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
		| column -t  -s '###'
.PHONY: help
