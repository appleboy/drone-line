# drone-line

[![Build Status](https://travis-ci.org/appleboy/drone-line.svg?branch=master)](https://travis-ci.org/appleboy/drone-line) [![codecov](https://codecov.io/gh/appleboy/drone-line/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-line) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-line)](https://goreportcard.com/report/github.com/appleboy/drone-line)

Drone plugin for sending line notifications.

## Register Line BOT API Trial

Please refer to [Getting started with BOT API Trial](https://developers.line.me/bot-api/getting-started-with-bot-api-trial).

## Build

Build the binary with the following commands:

```
$ make build
```

## Testing

Test the package with the following command:

```
$ make test
```

## Docker

Build the docker image with the following commands:

```
$ make docker
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
  -e PLUGIN_IMAGE=https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-128.png \
  -e PLUGIN_VIDEO=http://www.sample-videos.com/video/mp4/480/big_buck_bunny_480p_5mb.mp4 \
  -e PLUGIN_DELIMITER=:: \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  appleboy/drone-line
```
