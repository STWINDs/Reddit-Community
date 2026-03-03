# 第一阶段：构建
FROM golang:1.25.7-alpine AS builder

# 设置工作目录
WORKDIR /build

# 设置环境变量，加速构建
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct\
    GOOS=linux


# 复制源代码
COPY . .

# 构建可执行文件
RUN go build -o bluebell .

# 第二阶段：运行
FROM alpine:latest

# 安装必要的基础包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /build/bluebell .
# 复制配置文件
COPY --from=builder /build/config.yaml .
# 复制.env文件（如果需要）
COPY --from=builder /build/.env .

# 暴露端口
EXPOSE 8080

# 启动程序
ENTRYPOINT ["./bluebell"]
