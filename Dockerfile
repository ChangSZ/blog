FROM golang:latest AS builder

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /work
COPY . /work
RUN go build -o blog main.go

# 多级构建
FROM alpine:latest

WORKDIR /work
ENV TZ=Asia/Shanghai
RUN set -eux; \
    # 替换源
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    # 解决时区问题
    apk add --no-cache tzdata; \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime; \
    # 二进制文件无法通过./执行 ; \
    apk add --no-cache libc6-compat; \
    # 在go程序中无法访问https链接; \
    apk add --no-cache ca-certificates; \
    # 主要用来进行mysqldump操作
    apk add --no-cache mysql-client; \
    # go程序在alpine下不解析hosts文件; \
    echo "hosts: files dns" > /etc/nsswitch.conf

COPY --from=builder /work/blog /work/blog
COPY --from=builder /work/static /work/static
COPY --from=builder /work/template /work/template
COPY --from=builder /work/env.prod.yaml /work/env.prod.yaml

EXPOSE 8081
ENTRYPOINT ["./blog"]