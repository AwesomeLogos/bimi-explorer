# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as builder
RUN apk update && \
    apk upgrade && \
    apk --no-cache add git
RUN mkdir /build
ADD . /build/
WORKDIR /build
ARG COMMIT
ARG LASTMOD
RUN echo "INFO: installing sqlc"
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN echo "INFO: generating sqlc"
RUN sqlc generate
RUN echo "INFO: building for $COMMIT on $LASTMOD"
RUN \
    CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags "-X main.COMMIT=$COMMIT -X main.LASTMOD=$LASTMOD -extldflags '-static'" \
    -o bimi-explorer *.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bimi-explorer /app/
WORKDIR /app
ENV PORT 4000
ENTRYPOINT ["./bimi-explorer"]
