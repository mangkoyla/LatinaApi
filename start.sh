#!/bin/bash

# Update latina modules
go get -v github.com/mangkoyla/LatinaBot@gh-pages
go get -v github.com/mangkoyla/LatinaSub-go@gh-pages

# Tidy and verify all modules
go mod download && go mod tidy && go mod verify

# Compile software
go build -tags with_grpc -o ./latinaapi ./cmd/latinaapi/main.go

# Run software
./latinaapi
