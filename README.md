# stringslice

[![Build Status](https://cloud.drone.io/api/badges/pfmt/stringslice/status.svg)](https://cloud.drone.io/pfmt/stringslice)
[![Go Reference](https://pkg.go.dev/badge/github.com/pfmt/stringslice.svg)](https://pkg.go.dev/github.com/pfmt/stringslice)

String slice utils for Go.  
Source files are distributed under the BSD-style license.

## About

The software is considered to be at a alpha level of readiness,
its extremely slow and allocates a lots of memory.

## Benchmark

```sh
$ go test -count=1 -race -bench ./...
goos: linux
goarch: amd64
pkg: github.com/pfmt/stringslice
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkUniqueCopy/slice_test.go:30-8         	  555273	      1970 ns/op
PASS
ok  	github.com/pfmt/stringslice	1.142s
```
