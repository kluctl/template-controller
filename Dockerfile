# We must use a glibc based distro due to embedded python not supporting musl libc for aarch64 (only amd64+musl is supported)
# see https://github.com/indygreg/python-build-standalone/issues/87
# use `docker buildx imagetools inspect cgr.dev/chainguard/wolfi-base:latest` to find latest sha256 of multiarch image
FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/wolfi-base@sha256:1c7fcbc7d67145b82f8d393eabf43cf12a14e3759f2e42c8cb61799550909a08

# See https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG TARGETPLATFORM

RUN apk update && apk add tzdata

ENV KLUCTL_CACHE_DIR=/tmp/kluctl-cache

COPY bin/template-controller /usr/bin/
USER 65532:65532

ENTRYPOINT ["/usr/bin/template-controller"]
