############################################
# 软浮点板子使用
# # ls /lib/ld*
# # /lib/ld-2.25.so     /lib/ld-linux.so.3
#  如果有 ld-linux-armhf.so.3 则可使用此镜像 （硬浮点）
#  如果是  /lib/ld-linux.so.3 刚不能使用此镜像
#############################################
FROM golang:1.24-bullseye

# 安装 musl armv7 工具链（静态链接更干净）
RUN apt-get update && apt-get install -y \
    musl-tools musl-dev musl:armhf \
    gcc-arm-linux-gnueabihf \
    && rm -rf /var/lib/apt/lists/*

# 设置交叉编译环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabihf-gcc
ENV CXX=arm-linux-gnueabihf-g++
ENV CGO_LDFLAGS="-static"

WORKDIR /workspace