export GO111MODULE=on
export CGO_ENABLED=0

MAIN_DIR := ./cmd
OUTPUT_DIR := ./dist
CONFIG_DIR := ./config
BIN_OUTPUT := $(if $(filter $(shell go env GOOS), windows), ${OUTPUT_DIR}/plc-gw.exe, ${OUTPUT_DIR}/plc-gw)

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

default: clean generate-dns docs build
.PHONY: all clean generate-dns docs build

clean:
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	rm -rf dist/ builds/ cover.out

generate-dns:
	go generate ./...

docs:
	swag init -o docs/ -g api/*

test: clean build
	go test -v -cover ./...

build: clean docs
	@echo Version: $(VERSION)
	go build -v -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIR}
	if test `go env GOOS` = "linux"; then \
		GOOS=windows GOARCH=386 \
		go build -v -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT}-i386.exe ${MAIN_DIR} \
		&& \
		GOOS=windows GOARCH=amd64 \
		go build -v -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT}.exe ${MAIN_DIR} \
		; fi
	cp ${CONFIG_DIR}/app.yaml ${OUTPUT_DIR}

install:
	@echo Installing
	${BIN_OUTPUT} install