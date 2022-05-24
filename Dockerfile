FROM docker.io/library/alpine:3.16 as runtime

RUN \
  apk add --update --no-cache \
    bash \
    curl \
    ca-certificates \
    tzdata

# TODO: Adjust binary file name
ENTRYPOINT ["go-bootstrap"]
COPY go-bootstrap /usr/bin/

USER 65536:0
