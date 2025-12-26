# Stage 1: Build the Go binary
FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure proper static build for Alpine
RUN sh scripts/build_linux.sh



FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/bin/app-linux /app/app-linux

RUN apk add --no-cache chromium ttf-freefont nss
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' >> /etc/apk/repositories && \
  apk add --no-cache doppler
RUN chmod +x /app/app-linux


EXPOSE 80
RUN adduser -D -g '' nonroot
USER nonroot
CMD ["doppler", "run", "--", "/app/app-linux"]
