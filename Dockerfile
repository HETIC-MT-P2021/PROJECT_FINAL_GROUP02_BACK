FROM golang:1.13-alpine

LABEL maintainer="Croto de le Four al Fakhouri"

# Set the working directory inside the container
WORKDIR /go/src
# Copy the full project to current directory
COPY . .
# Run command to install dependencies
RUN apk add git
RUN go mod download

EXPOSE 8080