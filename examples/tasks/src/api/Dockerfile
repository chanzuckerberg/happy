FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
RUN apk update && apk upgrade
RUN apk --update add --no-cache git tzdata
ADD . .
RUN GOPROXY=direct  go build -o api

# final stage
FROM alpine:3.16
WORKDIR /app
RUN apk update && apk upgrade && apk --no-cache add curl
COPY --from=builder /app/api /app/
ADD delete.sh .
ADD migrate.sh .
EXPOSE 3000
ENTRYPOINT ./api