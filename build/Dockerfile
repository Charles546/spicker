FROM golang:1.18.1-alpine AS build

ARG GOINSTALLOPTS
ARG GOGC

ENV GOGC=$GOGC

# Enable module mode (see https://github.com/golang/go/wiki/Modules)
ENV GO111MODULE=auto

WORKDIR /build

RUN apk add --no-cache git gcc libc-dev
COPY ./ ./
RUN go install -v ./cmd/...

FROM alpine:3.15.4

LABEL description="Stock price grabber, demonstrating the use of Alphavantage API" \
      org.label-schema.vcs-url=https://github.com/Charles546/spicker

RUN apk add --no-cache tzdata
WORKDIR /opt/spicker
COPY --from=build /go/bin/* ./

ENTRYPOINT ["./stock-price-grabber-server", "--port=8080", "--host=0.0.0.0"]
EXPOSE 8080
