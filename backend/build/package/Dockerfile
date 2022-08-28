# Golang base image
FROM golang:1.18 as go_builder
LABEL stage="image_randomizer_builder"
WORKDIR /image-randomizer
COPY . .
WORKDIR /image-randomizer/cmd/image-randomizer
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o image-randomizer

# New stage from scratch
FROM alpine:3.16
LABEL stage="image_randomizer_binary"
RUN apk add --no-cache ca-certificates
COPY --from=go_builder /image-randomizer/cmd/image-randomizer/image-randomizer /cmd/image-randomizer/image-randomizer
WORKDIR /cmd/image-randomizer
