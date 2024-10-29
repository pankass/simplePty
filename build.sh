#!/bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o server ./server/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o server ./client/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -o server ./client/
