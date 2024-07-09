# We must use a glibc based distro due to embedded python not supporting musl libc for aarch64 (only amd64+musl is supported)
# see https://github.com/indygreg/python-build-standalone/issues/87
# use `docker buildx imagetools inspect cgr.dev/chainguard/wolfi-base:latest` to find latest sha256 of multiarch image
FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/wolfi-base@sha256:5cc4736f5a03e9ba618760e5b21c1f1bc0ce275fc6ba86589c8dd26374c2647b

# See https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG TARGETPLATFORM

RUN apk update && apk add tzdata

ENV KLUCTL_CACHE_DIR=/tmp/kluctl-cache

COPY bin/template-controller /usr/bin/
USER 65532:65532

ENTRYPOINT ["/usr/bin/template-controller"]
