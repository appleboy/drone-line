.PHONY: test build fmt vet errcheck lint install update

DIST := dist
EXECUTABLE := drone-line
VERSION := $(shell git describe --tags --always || git rev-parse --short HEAD)
IMPORT := github.com/appleboy/$(EXECUTABLE)

# for dockerhub
DEPLOY_ACCOUNT := "appleboy"
DEPLOY_IMAGE := $(EXECUTABLE)
DEPLOY_WEBHOOK_IMAGE := "$(EXECUTABLE)-webhook"

TARGETS ?= linux/*,darwin/*,windows/*
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?=
LDFLAGS += -X "main.Version=$(VERSION)"

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
	for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.txt $$PKG || exit 1; done;

html:
	go tool cover -html=coverage.txt

dep_install:
	glide install

dep_update:
	glide up

install: $(wildcard *.go)
	go install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $@

.PHONY: release
release: release-dirs release-build release-copy release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-build
release-build:
	@which xgo > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/karalabe/xgo; \
	fi
	xgo -dest $(DIST)/binaries -tags '$(TAGS)' -ldflags '-s -w $(LDFLAGS)' -targets '$(TARGETS)' -out $(EXECUTABLE)-$(VERSION) $(IMPORT)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

docker_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '$(LDFLAGS)'

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
