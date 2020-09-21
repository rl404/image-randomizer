# Golang base image
FROM golang:1.15 as go_builder
WORKDIR /go/src/github.com/rl404/image-randomizer/backend
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o randomizer

# New stage from scratch
FROM alpine:3.10
LABEL maintainer="axel.rl.404@gmail.com"
RUN apk add --no-cache ca-certificates
WORKDIR /docker/bin
COPY --from=go_builder /go/src/github.com/rl404/image-randomizer/backend/randomizer randomizer
CMD ["/docker/bin/randomizer"]