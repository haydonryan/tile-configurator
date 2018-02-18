#!/bin/bash -e

GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT_DIR}/tile-configurator-osx
GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/tile-configurator-linux
