ARG ARCH_ORG

# We must use a glibc based distro due to embedded python not supporting musl libc for aarch64
FROM $ARCH_ORG/debian:bullseye-slim

# We meed git for kustomize to support overlays from git
RUN apt update && apt install git -y && rm -rf /var/lib/apt/lists/*

COPY manager /manager
USER 65532:65532

ENTRYPOINT ["/manager"]
