#!/bin/bash

# builds and installs firemodel with a git-based version

go build \
  -ldflags "-X github.com/mickeyreiss/firemodel/version.Version=$(git describe --tags --always --dirty)" \
  -o "./.build/firemodel" ./firemodel/main.go

