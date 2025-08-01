# 前端构建阶段
FROM node:alpine AS frontend-builder

# 安装 pnpm
RUN corepack enable pnpm

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package.json web/pnpm-lock.yaml ./
# 安装依赖，使用缓存优化
RUN pnpm install --frozen-lockfile

# 复制前端源代码
COPY web/ ./
# 构建前端
RUN pnpm run build

# 后端构建阶段
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# 复制Go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .
# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./static/dist

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux go build -o radius_mgnt .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# 复制可执行文件和配置
COPY --from=backend-builder /app/radius_mgnt .
COPY --from=backend-builder /app/.env.example .env.example

EXPOSE 8080

CMD ["./radius_mgnt"]