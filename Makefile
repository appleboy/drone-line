.PHONY: test

DEPLOY_ACCOUNT := "appleboy"
DEPLOY_IMAGE := "drone-line"

install:
	glide install

build:
	go build

test:
	go test -v -coverprofile=coverage.txt

html:
	go tool cover -html=coverage.txt

update:
	glide up

docker_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo

docker_image:
	docker build --rm -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) .

docker: docker_build docker_image

docker_deploy:
ifeq ($(tag),)
	@echo "Usage: make $@ tag=<tag>"
	@exit 1
endif
	docker tag $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):latest $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
