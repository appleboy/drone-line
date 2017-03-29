# drone-facebook

[![Build Status](https://travis-ci.org/appleboy/drone-facebook.svg?branch=master)](https://travis-ci.org/appleboy/drone-facebook) [![codecov](https://codecov.io/gh/appleboy/drone-facebook/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-facebook) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-facebook)](https://goreportcard.com/report/github.com/appleboy/drone-facebook)

[Drone](https://github.com/drone/drone) plugin for sending [Facebook Messages](https://developers.facebook.com/docs/messenger-platform).

## Feature

* [x] Send with Text Message (Support [handlebar](https://github.com/aymerick/raymond) template).
* [x] Send with New Image.
* [x] Send with New Audio.
* [x] Send with New Video.
* [x] Send with New File.

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
'/bin/drone-facebook' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_FB_PAGE_TOKEN=xxxxxxx \
  -e PLUGIN_FB_VERIFY_TOKEN=xxxxxxx \
  -e PLUGIN_TO=xxxxxxx \
  -e PLUGIN_MESSAGE=test \
  -e PLUGIN_IMAGES=http://example.com/test.png \
  -e PLUGIN_AUDIOS=http://example.com/test.mp3 \
  -e PLUGIN_VIDEOS=http://example.com/test.mp4 \
  -e PLUGIN_FILES=http://example.com/test.pdf \
  -e PLUGIN_ONLY_MATCH_EMAIL=false \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_COMMIT_AUTHOR_EMAIL=appleboy@gmail.com \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  -e DRONE_JOB_STARTED=1477550550 \
  -e DRONE_JOB_FINISHED=1477550750 \
  -e DRONE_TAG=1.0.0 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-facebook
```

You can get more [information](DOCS.md) about how to use facebook plugin in drone.
