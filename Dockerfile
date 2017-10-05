FROM plugins/base:multiarch

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/appleboy/drone-line.git"
LABEL org.label-schema.name="Line"
LABEL org.label-schema.vendor="Bo-Yi Wu"
LABEL org.label-schema.schema-version="1.0"
LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>"

ADD drone-line /

ENTRYPOINT ["/drone-line"]
