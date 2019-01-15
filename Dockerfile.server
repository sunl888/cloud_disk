FROM golang:1.11.1-alpine3.7 as builder

RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		bash \
		musl-dev \
        openssl \
        go

COPY . /go/src/github.com/wq1019/cloud_disk

RUN go build -v -o /app/server /go/src/github.com/wq1019/cloud_disk/cmd/server/main.go && \
    go build -v -o /app/cli /go/src/github.com/wq1019/cloud_disk/cmd/cli/main.go


FROM alpine:3.7

RUN apk update && apk --no-cache add mailcap ca-certificates tzdata

ENV TZ=Asia/Shanghai
#设置时区
#RUN /bin/cp /usr/share/zoneinfo/$TZ /etc/localtime \
#  && echo '$TZ' >/etc/timezone


COPY --from=builder /app/server /app/server
COPY --from=builder /app/cli /app/cli
COPY --from=builder /go/src/github.com/wq1019/cloud_disk/config/config.yml /app/config/config.yml

WORKDIR /app

RUN chmod +x /app/server /app/cli

CMD ["./server"]