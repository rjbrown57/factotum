# Build the manager binary
FROM gcr.io/distroless/static:nonroot
ARG TARGETOS
ARG TARGETARCH
COPY factotum ./
USER 65532:65532
ENTRYPOINT ["/factotum"]
