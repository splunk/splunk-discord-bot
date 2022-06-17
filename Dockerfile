FROM golang:1.18.3-alpine3.16 as builder
RUN apk add --no-cache make
RUN mkdir -p /build
WORKDIR /build
COPY . .
RUN make

FROM alpine:3.16
COPY --from=builder /build/bin/bot /bin/bot

ENV CONFIG_FILE /etc/config.json

CMD ["/bin/bot"]
