FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk update && apk upgrade
RUN apk --update add --no-cache git tzdata
ADD . .
RUN GOPROXY=direct  go build -o api

# final stage
FROM alpine:3.18
WORKDIR /app
RUN apk update && apk upgrade && apk --no-cache add curl
COPY --from=builder /app/api /app/
EXPOSE 3000
ENTRYPOINT ./api