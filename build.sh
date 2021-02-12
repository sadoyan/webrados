#!/bin/bash

export GOROOT=/usr/local/go
export export GOPATH=`pwd`

rm -rf pkg/*
rm -rf src/{github.com,golang.org,gopkg.in}
go get github.com/ceph/go-ceph
go get gopkg.in/ini.v1
go get golang.org/x/sys/unix
go build src/webrados.go