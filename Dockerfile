FROM golang:1.20-alpine AS builder

WORKDIR /app

ADD go.mod .
ADD go.sum .
ADD main.go .
ADD getobject.go .

RUN go build -o /usr/bin/s3-ssec-get

FROM alpine:latest

COPY --from=builder /usr/bin/s3-ssec-get /usr/bin/s3-ssec-get