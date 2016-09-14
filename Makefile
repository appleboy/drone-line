.PHONY: test

install:
	glide install

update:
	glide up

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo

image:
	docker build --rm -t appleboy/drone-line .

docker: build image
