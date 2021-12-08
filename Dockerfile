FROM golang:alpine AS builder
WORKDIR /
RUN go env -w GOPRIVATE=github.com/shadowsocks

RUN  go build -o xray -ldflags "-s -w" ./main

FROM alpine
WORKDIR /
COPY --from=builder /xray /usr/bin/
RUN apk update  && apk add gettext

COPY update_config.sh /docker-entrypoint.d/update-config.sh
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /docker-entrypoint.d/update-config.sh
RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/bin/sh","/entrypoint.sh"]
CMD ["/usr/bin/xray", "-config" ,"/etc/xray/config.json"]
