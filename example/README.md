# Start Line Bot

There are two ways to start line boot webhook service.

## Start by Golang

Download dependency packages.

```bash
$ go get -t -v ./...
```

```bash
$ export CHANNEL_SECRET=xxxxx
$ export CHANNEL_TOKEN=xxxxx
$ go run server.go
```

## Start by Docker

Build docker image

```bash
$ docker build -t line .
```

Start service.

```bash
$ docker run \
  -e CHANNEL_SECRET=xxxx \
  -e CHANNEL_TOKEN=xxxx \
  line
```
