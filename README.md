# drone-line

Drone plugin for sending line notifications.

## Register Line BOT API Trial

Please refer to [Getting started with BOT API Trial](https://developers.line.me/bot-api/getting-started-with-bot-api-trial).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm -t appleboy/drone-line .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-line' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_CHANNEL_ID=xxxxxxx \
  -e PLUGIN_CHANNEL_SECRET=xxxxxxx \
  -e PLUGIN_MID=xxxxxxx \
  -e PLUGIN_TO=xxxxxxx \
  -e PLUGIN_MESSAGE=test \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=hello-world \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  appleboy/drone-line
```
