#!/bin/bash
#rsync -vzal * 192.168.221.136:/usr/local/src/go-webrados/
#reflex -d none -r '.'  -s -- sh -c  'go mod tidy && go run ./ -config config.ini'
reflex -d none -r '.'  -s -- sh -c  'rsync -val * gate:/opt/go-webrados/ && sleep 128d'