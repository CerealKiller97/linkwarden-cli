# 🐳 linkwarden-cli image.
#
#   docker build -t linkwarden-cli .                      # builds from source (production stage)
#   docker run --rm -e LINKWARDEN_TOKEN linkwarden-cli links list
#
# Releases use the goreleaser-production stage, which only packages the binary GoReleaser built.

FROM golang:1.27 AS build

ENV CGO_ENABLED=0

WORKDIR /linkwarden-cli

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN go build -trimpath \
    -ldflags "-s -w -X github.com/CerealKiller97/linkwarden-cli/cmd.version=${VERSION}" \
    -o /out/linkwarden .

# distroless/static ships CA certificates for HTTPS to Linkwarden, and runs as a non-root user.
FROM gcr.io/distroless/static-debian13:nonroot AS prepare-production

LABEL org.opencontainers.image.source="https://github.com/CerealKiller97/linkwarden-cli"
LABEL org.opencontainers.image.title="linkwarden-cli"
LABEL org.opencontainers.image.description="CLI for the Linkwarden bookmark manager"
LABEL org.opencontainers.image.licenses="AGPL-3.0-only"

ENTRYPOINT ["/usr/local/bin/linkwarden"]

FROM prepare-production AS goreleaser-production

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/linkwarden /usr/local/bin/linkwarden

FROM prepare-production AS production

COPY --from=build /out/linkwarden /usr/local/bin/linkwarden
