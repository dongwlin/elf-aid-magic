build:
	go run ./scripts/build

install:
	go run ./scripts/install

install-clear:
	go run ./scripts/install --clear

.PHONY: build install install-clear