############################################
# 软浮点板子使用
# # ls /lib/ld*
# # /lib/ld-2.25.so     /lib/ld-linux.so.3
#  如果有 ld-linux-armhf.so.3 则可使用此镜像 （硬浮点）
#  如果是  /lib/ld-linux.so.3 刚不能使用此镜像
#############################################
#FROM golang:1.24-bullseye
FROM cgo-arm7
# 安装软浮点工具链
RUN apt-get update && apt-get install -y \
        gcc-arm-linux-gnueabi \
        libc6-dev-armel-cross \
    && rm -rf /var/lib/apt/lists/*

ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabi-gcc
ENV CXX=arm-linux-gnueabi-g++
ENV CGO_LDFLAGS="-static"

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /workspace
