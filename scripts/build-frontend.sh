#!/bin/bash

# 构建前端项目的脚本

set -e

echo "🚀 开始构建前端项目..."

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js"
    exit 1
fi

# 检查 npm 是否安装
if ! command -v pnpm &> /dev/null; then
    echo "❌ pnpm 未安装，请先安装 pnpm"
    exit 1
fi

# 进入前端目录
cd web

# 检查是否存在 node_modules，如果不存在则安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    pnpm install
else
    echo "✅ 前端依赖已存在，跳过安装"
fi

# 构建前端项目
echo "🔨 构建前端项目..."
pnpm run build

# 返回根目录
cd ..

# 检查构建结果
if [ -f "static/dist/index.html" ]; then
    echo "✅ 前端构建成功！静态文件已生成到 static/dist 目录"
    echo "📁 构建产物："
    ls -la static/dist/
else
    echo "❌ 前端构建失败！"
    exit 1
fi

echo "🎉 前端构建完成！"