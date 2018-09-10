FROM golang:1.7.3-alpine

ENV app /go/src/github.com/appleboy/drone-line
ENV CHANNEL_SECRET $CHANNEL_SECRET
ENV CHANNEL_TOKEN $CHANNEL_TOKEN
ENV PORT 8089

RUN apk update && apk upgrade && \
  apk add --no-cache git && \
  rm -rf /var/cache/apk/*

ADD server.go $app/
WORKDIR $app
RUN go get -t -v ./...

EXPOSE $PORT

CMD ["go", "run", "server.go"]
