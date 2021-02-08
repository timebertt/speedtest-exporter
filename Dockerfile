FROM golang:1.15.8 AS builder

WORKDIR /go/src/github.com/timebertt/speedtest-exporter

COPY . .
RUN CGO_ENABLED=0 go install .

FROM alpine:3.10.5 AS speedtest-exporter

LABEL org.opencontainers.image.source=https://github.com/timebertt/speedtest-exporter

RUN apk add --no-cache speedtest-cli

COPY --from=builder /go/bin/speedtest-exporter /speedtest-exporter
EXPOSE 8080

ENTRYPOINT [ "/speedtest-exporter" ]
