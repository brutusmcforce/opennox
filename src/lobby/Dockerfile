# Build from the src directory: docker build -f lobby/Dockerfile .
FROM golang:1.25-alpine AS builder

WORKDIR /app

ADD go.* ./
RUN go mod download

ADD . .
RUN CGO_ENABLED=0 go build -o nox-lobby ./lobby/cmd/nox-lobby

FROM alpine:3.22

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/nox-lobby /usr/bin/nox-lobby

EXPOSE 80
EXPOSE 6060

ENTRYPOINT ["nox-lobby", "serve", "--host=:80", "--monitor=:6060"]
