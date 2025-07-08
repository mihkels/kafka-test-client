#!/usr/bin/env bash

export CGO_ENABLED=0
export GOOS=linux

go build -a -installsuffix cgo -o build/kafka_app ./cmd
