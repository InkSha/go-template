.PHONY: help build run swagger clean test deps docker-build docker-up docker-down

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## 生成Swagger文档
	swag init --parseInternal -g cmd/server/main.go -o docs

build: ## 编译项目
	go build -ldflags "-X server/pkg/version.Version=1.0.0 -X server/pkg/version.GitCommit=$(shell git rev-parse --short HEAD) -X server/pkg/version.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/server cmd/server/main.go

run: ## 运行项目
	go run cmd/server/main.go

dev: swagger run ## 生成文档并运行

clean: ## 清理编译文件
	rm -rf bin/
	rm -rf docs/

test: ## 运行测试
	go test -v ./...

fmt: ## 格式化代码
	go fmt ./...

lint: ## 代码检查
	golangci-lint run

docker-build: ## 构建 Docker 镜像
	docker build -t server:latest .

docker-up: ## 启动 Docker Compose
	docker-compose up -d

docker-down: ## 停止 Docker Compose
	docker-compose down

docker-logs: ## 查看 Docker 日志
	docker-compose logs -f app
