.PHONY: test

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
	docker build --rm -t appleboy/drone-line .

docker: build image
