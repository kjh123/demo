FROM golang:1.21-alpine as builder

WORKDIR /usr/src/hello

COPY . .

RUN apk update && \
    go build -o server.out ./server/... && \
    go build -o client.out ./client/...

FROM scratch as server

COPY --from=builder /usr/src/hello/server.out ./hello-server

CMD ["hello-server"]

FROM scratch as client

COPY --from=builder /usr/src/hello/client.out ./hello-client

CMD ["hello-client"]