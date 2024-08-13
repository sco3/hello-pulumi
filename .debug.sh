#!/bin/bash

#go build -o cortex -gcflags "all=-N -l" .
#exec dlv --listen=:2345 --headless=true --api-version=2 exec ./goapp

dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient
