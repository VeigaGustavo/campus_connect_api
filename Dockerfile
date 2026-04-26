FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/campus_connect_api ./main.go

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /bin/campus_connect_api /usr/local/bin/campus_connect_api

EXPOSE 8080

USER app
ENTRYPOINT ["campus_connect_api"]
