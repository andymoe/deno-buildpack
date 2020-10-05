#!/usr/bin/env bash

GOOS=linux go build -ldflags="-s -w" -o ./bin/detect ./cmd/detect/main.go
GOOS=linux go build -ldflags="-s -w" -o ./bin/build ./cmd/build/main.go