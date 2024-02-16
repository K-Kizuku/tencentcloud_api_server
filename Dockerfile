ARG GO_VERSION=1.22
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as go-builder

WORKDIR /go/src/tencentcloud_api_server

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o api main.go

FROM alpine:${ALPINE_VERSION}

WORKDIR /usr/src/tencentcloud_api_server

COPY --from=go-builder /go/src/tencentcloud_api_server/api api
RUN chmod +x "/usr/src/tencentcloud_api_server/api"

ENTRYPOINT ["/usr/src/tencentcloud_api_server/api"]