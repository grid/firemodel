.PHONY: default
default: setup check build

.PHONY: setup
setup:
	go get -v -t -d ./...

.PHONY: check
check:
	go vet ./...
	go test -v ./...

build: .build/firemodel

.build: .build
	mkdir -p .build

.build/firemodel: .build
	go build -o .build/firemodel .

