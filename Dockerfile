FROM docker.io/library/alpine:3.18 as runtime

RUN \
  apk add --update --no-cache \
  bash \
  curl \
  ca-certificates \
  tzdata

ENTRYPOINT ["provider-minio"]
COPY provider-minio /usr/bin/

USER 65536:0
