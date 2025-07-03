# 多阶段构建
FROM golang:1.24-alpine AS builder

# 设置Go环境变量
ENV GO111MODULE=on

# 安装必要的包
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY ../go.mod ../go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY .. .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 最终镜像
FROM alpine:latest


WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 9527

# 运行应用
CMD ["./main"] 