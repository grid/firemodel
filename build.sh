#!/bin/bash

# builds and installs firemodel with a git-based version

go build \
  -ldflags "-X github.com/visor-tax/firemodel/version.Version=$(git describe --tags --always --dirty)" \
  -o "./.build/firemodel" ./firemodel/main.go

