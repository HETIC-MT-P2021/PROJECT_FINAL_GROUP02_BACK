FROM golang:1.15-alpine

RUN apk add --no-cache git
RUN mkdir /app

# Add Maintainer Info
LABEL maintainer="Team Lambert"

ADD . /app
WORKDIR /app

RUN go get
RUN go mod vendor
RUN go build -o main .

CMD ["./main"]