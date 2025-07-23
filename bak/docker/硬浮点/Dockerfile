############################################
# 硬浮点板子使用
# # ls /lib/ld*
# # /lib/ld-2.25.so     /lib/ld-linux.so.3
#  如果有 ld-linux-armhf.so.3 则可使用此镜像 （硬浮点）
#  如果是  /lib/ld-linux.so.3 刚不能使用此镜像
#############################################
FROM golang:1.24-bullseye

# 安装 arm-linux-gnueabihf 交叉编译工具链
RUN apt-get update && apt-get install -y \
    gcc-arm-linux-gnueabihf \
    g++-arm-linux-gnueabihf \
    libc6-dev-armhf-cross \
    && rm -rf /var/lib/apt/lists/*

# 设置交叉编译环境变量
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabihf-gcc
ENV CXX=arm-linux-gnueabihf-g++

WORKDIR /workspace