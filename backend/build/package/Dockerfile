# Golang base image
FROM golang:1.27.3 AS go_builder
LABEL stage="image_randomizer_be_builder"
WORKDIR /image-randomizer
COPY . .
WORKDIR /image-randomizer/cmd/image-randomizer
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o image-randomizer

# New stage from scratch
FROM alpine:3.24.3
LABEL stage="image_randomizer_be_binary"
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 image-randomizer
COPY --from=go_builder /image-randomizer/cmd/image-randomizer/image-randomizer /usr/local/bin/image-randomizer
WORKDIR /home/image-randomizer
USER image-randomizer
ENTRYPOINT ["/usr/local/bin/image-randomizer"]
