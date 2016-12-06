.PHONY: test build fmt vet errcheck lint install update

EXECUTABLE := drone-line
VERSION := $(shell git describe --tags --always || git rev-parse --short HEAD)

# for dockerhub
DEPLOY_ACCOUNT := "appleboy"
DEPLOY_IMAGE := $(EXECUTABLE)
DEPLOY_WEBHOOK_IMAGE := "$(EXECUTABLE)-webhook"

TARGETS ?= linux/*,darwin/*,windows/*
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?=

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

all: build

fmt:
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

test:
	go test -v -coverprofile=coverage.txt

html:
	go tool cover -html=coverage.txt

dep_install:
	glide install

dep_update:
	glide up

install: $(wildcard *.go)
	go install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS) -X main.Version=$(VERSION)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS) -X main.Version=$(VERSION)' -o $@

docker_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-X main.Version=$(VERSION)"

docker_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) .

docker_webhook_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_WEBHOOK_IMAGE) example

docker: docker_build docker_image docker_webhook_image

docker_deploy:
ifeq ($(tag),)
	@echo "Usage: make $@ tag=<tag>"
	@exit 1
endif
	# deploy line image
	docker tag $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):latest $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
	# deploy line webhook image
	docker tag $(DEPLOY_ACCOUNT)/$(DEPLOY_WEBHOOK_IMAGE):latest $(DEPLOY_ACCOUNT)/$(DEPLOY_WEBHOOK_IMAGE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(DEPLOY_WEBHOOK_IMAGE):$(tag)

clean:
	go clean -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST)

version:
	@echo $(VERSION)
