FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk update && apk upgrade
RUN apk --update add --no-cache git tzdata
ADD . .
RUN GOPROXY=direct  go build -o api

# final stage
FROM alpine:3.9
WORKDIR /app
RUN apk update && apk upgrade && apk --no-cache add curl
RUN apk add --no-cache git make gcc g++ libc-dev pkgconfig \
    libxml2-dev libxslt-dev postgresql-dev coreutils curl wget bash \
    gnupg tar linux-headers bison readline-dev readline zlib-dev \
    zlib yaml-dev autoconf ncurses-dev curl-dev apache2-dev \
    libx11-dev libffi-dev tcl-dev tk-dev openjdk8
COPY --from=builder /app/api /app/
EXPOSE 3000
ENTRYPOINT ./api