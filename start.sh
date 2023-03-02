#!/bin/bash

go get -u ./... && go mod tidy && go mod verify
go run cmd/latinaapi/main.go