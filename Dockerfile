FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖包
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

# 复制项目的源代码到容器的工作目录
COPY . .

# 构建应用
RUN go build -o main .

# run it...
FROM alpine:latest

# 设置工作目录
WORKDIR /root

# 从构建阶段复制编译好的应用到当前镜像
COPY --from=builder /app/main .
RUN chmod +x ./main

# Copy static resources
COPY --from=builder /app/assets /root/assets

# Create bucket directory
RUN mkdir -p /root/bucket/tablecache

# 暴露端口号
EXPOSE 8080

CMD ["./main"]