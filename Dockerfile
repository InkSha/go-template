# 构建阶段
FROM golang:1.25-alpine AS builder

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git make

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X server/pkg/version.Version=${VERSION} -X server/pkg/version.GitCommit=${GIT_COMMIT} -X server/pkg/version.BuildTime=${BUILD_TIME}" \
    -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache ca-certificates tzdata

# 复制编译好的二进制文件
COPY --from=builder /build/server .
COPY --from=builder /build/config.yaml .
COPY --from=builder /build/web ./web

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行
CMD ["./server"]
