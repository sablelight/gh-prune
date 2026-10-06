FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gh-prune ./cmd/gh-prune

FROM alpine:3.24

RUN apk add --no-cache ca-certificates

COPY --from=builder /gh-prune /usr/local/bin/gh-prune

ENTRYPOINT ["gh-prune"]