# Golang base image
FROM golang:1.5 as go_builder
LABEL stage="image_randomizer_be_builder"
WORKDIR /image-randomizer
COPY . .
WORKDIR /image-randomizer/cmd/image-randomizer
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o image-randomizer

# New stage from scratch
FROM alpine:3.20
LABEL stage="image_randomizer_be_binary"
RUN apk add --no-cache ca-certificates
COPY --from=go_builder /image-randomizer/cmd/image-randomizer/image-randomizer /cmd/image-randomizer/image-randomizer
WORKDIR /cmd/image-randomizer
