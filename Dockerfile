# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# 复制依赖文件并下载
COPY go.mod ./
# 如果有 go.sum 也要复制
# COPY go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# -ldflags="-s -w" 用于减小二进制体积（去除符号表和调试信息）
RUN go build -ldflags="-s -w" -o anyproxy main.go

# 运行阶段
FROM alpine:latest

# 安装基础证书（HTTPS 请求需要）和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/anyproxy .

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./anyproxy"]
