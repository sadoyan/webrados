#!/bin/bash

export GOROOT=/usr/local/go
go mod tidy
go build webrados.go
