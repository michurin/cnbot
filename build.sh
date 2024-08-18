#!/bin/sh

set -ex

go mod verify &&
go build -ldflags="-s -w -X 'main.Build=$(date -u '+Build_%F_%R')_$(git rev-parse HEAD)'" ./cmd/...
