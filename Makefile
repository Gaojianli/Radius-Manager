.PHONY: help build build-frontend build-backend clean dev run install test

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	@echo "📦 安装前端依赖..."
	@cd web && pnpm install
	@echo "📦 安装Go依赖..."
	@go mod tidy
	@echo "✅ 依赖安装完成"

build-frontend: ## 构建前端
	@echo "🚀 构建前端项目..."
	@./scripts/build-frontend.sh

build-backend: ## 构建后端
	@echo "🔨 构建后端项目..."
	@go build -o radius_mgnt .
	@echo "✅ 后端构建完成"

build: build-frontend build-backend ## 构建整个项目

clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	@rm -rf static/dist
	@rm -f radius_mgnt
	@rm -rf web/node_modules
	@echo "✅ 清理完成"

dev-frontend: ## 启动前端开发服务器
	@echo "🚀 启动前端开发服务器..."
	@cd web && npm run dev

dev-backend: ## 启动后端开发服务器
	@echo "🚀 启动后端开发服务器..."
	@go run .

run: build ## 构建并运行项目
	@echo "🚀 启动 RADIUS 管理系统..."
	@./radius_mgnt

docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	@docker build -t radius-mgnt:latest .

docker-run: ## 运行 Docker 容器
	@echo "🐳 启动 Docker 容器..."
	@docker-compose up -d

test: ## 运行测试
	@echo "🧪 运行测试..."
	@go test -v ./...

format: ## 格式化代码
	@echo "📝 格式化 Go 代码..."
	@go fmt ./...
	@echo "📝 格式化前端代码..."
	@cd web && npm run format || echo "前端格式化命令不存在，跳过"

lint: ## 检查代码规范
	@echo "🔍 检查 Go 代码规范..."
	@go vet ./...
	@echo "🔍 检查前端代码规范..."
	@cd web && npm run lint || echo "前端检查命令不存在，跳过"