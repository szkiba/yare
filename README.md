# yare

[![codecov](https://codecov.io/gh/szkiba/yare/branch/master/graph/badge.svg)](https://codecov.io/gh/szkiba/yare)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/szkiba/yare)
[![Go Report Card](https://goreportcard.com/badge/github.com/szkiba/yare)](https://goreportcard.com/report/github.com/szkiba/yare)

**HTTP echo service with JSON output**

**yare** (*Yet Another REST Echo*) is a small HTTP echo server written in Go programming language.
Core features are also available as Go package.

## Features

- *JSON output* - Echoes back HTTP request data (version, method, headers, etc) as JSON object.
- *Parse body* - Supports JSON and JWT request body formats (other formats ignored silently).
- *Parse Authorization* - Supports `Bearer` authentication scheme with JWT tokens and `Basic` scheme.
The response will include parsed credentials.
- *Custom parsers* - The Go package supports custom body and authorization scheme parser registration.
- *Request method* - Any HTTP methods are supported (GET, POST, PUT, etc), the response will include the original request method.
- *Request path* - Accessible on any request path, the response will include the original path.
- *Query parameters* - Supports arbitrary query parameters, the response will include original parameters.
- *Form parameters* - Supports arbitrary form parameters, the response will include original parameters.
- *http.Request and http.Response mapping* - The Go package supports mapping request and response parameters
to `map[string]interface{}` for trace logging.

## Install

You can install the pre-compiled binary or use Docker.

### Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](releases) and
copy to the desired location.

### Running with Docker

You can also use it within a Docker container. To do that, you'll need to
execute something more-or-less like the following:

```sh
docker run --network=host szkiba/yare
```

### Verifying your installation

To verify your installation, use the `yare -v` command:

```sh
$ yare -v

yare/1.0.0 linux/amd64
```

You should see `yare/VERSION` in the output.

### Usage

To print usage information, use the `yare --help` command:

```sh
$ yare --help

Usage of yare:
  -port int
        port to listen on (default 8080)
  -v    prints version
```

## TODO

Document, document, document...
