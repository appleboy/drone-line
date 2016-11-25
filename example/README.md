# Start Line Bot

There are two ways to start line boot webhook service.

## Start webhook service

### Start with Golang

Download dependency packages.

```bash
$ go get -t -v ./...
```

```bash
$ export CHANNEL_SECRET=xxxxx
$ export CHANNEL_TOKEN=xxxxx
$ export PORT=8089
$ go run server.go
```

### Start with Docker

Build your own docker image.

```bash
$ docker build -t appleboy/drone-line-webhook .
```

or download image from [docker hub](https://hub.docker.com/r/appleboy/drone-line-webhook/).

```bash
$ docker pull appleboy/drone-line-webhook
```

then start service on host port `8089`.

```bash
$ docker run --rm \
  -e CHANNEL_SECRET=xxxx \
  -e CHANNEL_TOKEN=xxxx \
  -p 8089:8089 \
  appleboy/drone-line-webhook
```

## Use ngrok

Use ngrok to tunnel your locally runnning bot so that Line can reach the webhook.

```bash
$ ngrok http 8089
```
