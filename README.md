# unidata_exporter
[![GoDoc](https://godoc.org/github.com/moorereason/unidata_exporter?status.svg)](https://godoc.org/github.com/moorereason/unidata_exporter)
[![Go Report Card](https://goreportcard.com/badge/moorereason/unidata_exporter)](https://goreportcard.com/report/moorereason/unidata_exporter)

Export [Rocket Unidata](http://www.rocketsoftware.com/) license statistics to [Prometheus](https://prometheus.io).

## Getting started

```bash
go build
./unidata_exporter -udtbin $UDTBIN
```
