# 阶段1：交叉编译
FROM golang:1.24 AS builder

# 安装交叉编译工具链
RUN apt-get update && apt-get install -y gcc-arm-linux-gnueabihf

# 设置工作目录
WORKDIR /app
COPY . .

# 设置交叉编译环境变量
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabihf-gcc

RUN go mod tidy && go build -o app ./cmd

# 阶段2：运行时镜像
FROM arm32v7/alpine
WORKDIR /app
COPY --from=builder /app/app /app/app
CMD ["./app"]