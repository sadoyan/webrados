#!/bin/bash

go mod tidy

export APIKEY='a66cd04bb85a2daed5080fb41c3da6642f37f4390d76e37c2a57f4edd4c9324e'
export JWTSECRET='Super$ecter123765@'

#go run ./ -config config.yml

reflex -d none -r '.'  -s -- sh -c  'go mod tidy && go run ./ -config config.yml'

