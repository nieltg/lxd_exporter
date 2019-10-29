# lxd_exporter

[![Build Status](https://travis-ci.org/nieltg/lxd_exporter.svg?branch=master)](https://travis-ci.org/nieltg/lxd_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/nieltg/lxd_exporter)](https://goreportcard.com/report/github.com/nieltg/lxd_exporter)
[![Coverage Status](https://coveralls.io/repos/github/nieltg/lxd_exporter/badge.svg?branch=master)](https://coveralls.io/github/nieltg/lxd_exporter?branch=master)

LXD metrics exporter for Prometheus.

## Usage

Download latest precompiled binary of this exporter from [the release page](https://github.com/nieltg/lxd_exporter/releases).

Extract archive, then run the exporter:
```
./lxd_exporter
```

The exporter must have access to LXD socket which can be guided by specifying `LXD_SOCKET` or `LXD_DIR` environment variable.
For more information, you can see the documentation from [Go LXD client library](https://godoc.org/github.com/lxc/lxd/client#ConnectLXDUnix).

## Hacking

Install [Go](https://golang.org/dl) before hacking this library.

To run all tests:
```
go test ./...
```

To build exporter binary:
```
mkdir build
go build -o build ./...
```

Binary will be available on `build/` directory.

## License

[MIT](LICENSE).
