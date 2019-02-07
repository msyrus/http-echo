# Defining App builder image
FROM golang:alpine AS builder

# Add git to determine build git version
RUN apk add --no-cache --update git

# Set GOPATH to build Go app
ENV GOPATH=/go

# Set apps source directory
ENV SRC_DIR=${GOPATH}/src/github.com/msyrus/http-echo

# Copy apps scource code to the image
COPY . ${SRC_DIR}

# Define current working directory
WORKDIR ${SRC_DIR}

# Build App
RUN ./build.sh

# Defining App image
FROM alpine:latest

# Copy App binary to image
COPY --from=builder /go/bin/http-echo /usr/local/bin/http-echo

EXPOSE 8000

ENTRYPOINT ["http-echo"]
