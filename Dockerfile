FROM centurylink/ca-certs

ADD drone-line /

ENTRYPOINT ["/drone-line"]
