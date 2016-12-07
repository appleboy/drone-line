FROM centurylink/ca-certs

ENV PORT 8089

ADD drone-line-webhook /

EXPOSE $PORT

ENTRYPOINT ["/drone-line-webhook"]
