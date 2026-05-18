FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath -ldflags="-s -w" \
    -o /bin/campus_connect_api \
    .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app && adduser -S app -G app \
    && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /bin/campus_connect_api /usr/local/bin/campus_connect_api

EXPOSE 8080
USER app

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["campus_connect_api"]
