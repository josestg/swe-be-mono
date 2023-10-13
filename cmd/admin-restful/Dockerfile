# Build the Go Binary.
FROM golang:1.21.3 AS builder

ARG GITHUB_TOKEN
ENV GOOS linux
ENV CGO_ENABLED 0

RUN mkdir /service
COPY go.* /service/
WORKDIR /service

# Configure github auth
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go mod download

COPY . /service

# prevent make command from being cached.
ARG CACHEBUST=1
RUN make -C . bin/admin-restful

FROM alpine:3.18

ARG BUILD_DATE
ARG BUILD_VERSION
ARG IMAGE_NAME

COPY --from=builder /service/resources/ /service/resources/
COPY --from=builder /service/bin/admin-restful /service/
COPY --from=builder /service/cmd/admin-restful/config/ /service/config/

WORKDIR /service

CMD ["./admin-restful"]

LABEL \
    org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="${IMAGE_NAME}" \
    org.opencontainers.image.revision="${BUILD_VERSION}" \
    org.opencontainers.image.authors="Jose Alfredo Sitanggang <josealredositanggang@gmail.com>" \
    org.opencontainers.image.source="https://github.com/josestg/swe-be-mono/cmd/admin-restful" \
    org.opencontainers.image.vendor="josestg.com"
