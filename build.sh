#!/bin/bash

export GOROOT=/usr/local/go
export GOPATH=`pwd`

rm -rf pkg/*
rm -rf src/{github.com,golang.org,gopkg.in}
go get github.com/ceph/go-ceph
go get gopkg.in/ini.v1
go get golang.org/x/sys/unix
go get go.opentelemetry.io/otel
#go get github.com/go-redis/redis
go get github.com/gomodule/redigo/redis
#go build src/webrados.go
