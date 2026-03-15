FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /shlink-cli .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /shlink-cli /usr/local/bin/shlink-cli

ENTRYPOINT ["shlink-cli"]
