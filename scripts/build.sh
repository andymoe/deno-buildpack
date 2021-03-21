#!/usr/bin/env bash

GOOS=linux go build -ldflags="-s -w" -o ./bin/run ./cmd/run/main.go

pushd bin
ln -sf "run" detect
ln -sf "run" build
popd