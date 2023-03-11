#!/bin/bash

# Update latinasub module
go get -v github.com/LalatinaHub/LatinaSub-go@main

# Tidy and verify all modules
go mod download && go mod tidy && go mod verify

# Compile software
go build -tags with_grpc -o ./latinaapi ./cmd/latinaapi/main.go

# Run software
./latinaapi