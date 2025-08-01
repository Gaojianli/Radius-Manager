# 🌐 RADIUS Manager

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker Pulls](https://img.shields.io/docker/pulls/gaojianli2333/radius-manager)](https://hub.docker.com/r/gaojianli2333/radius-manager)
[![Docker Image Size](https://img.shields.io/docker/image-size/gaojianli2333/radius-manager/latest)](https://hub.docker.com/r/gaojianli2333/radius-manager)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-green.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-blue.svg)](https://www.typescriptlang.org/)

现代化的 RADIUS 用户管理系统，基于 **CloudWeGo Hertz** 和 **Vue 3** 构建。为 FreeRADIUS 提供完整的用户认证、授权和管理功能，配备响应式的 Web 管理界面。

## ✨ 项目特色

- 🎨 **现代化 UI**：基于 Arco Design 的响应式界面，完美支持桌面端和移动端
- 🔐 **安全可靠**：SHA256+Salt 密码加密，JWT 认证，基于角色的权限控制
- ⚡ **高性能**：CloudWeGo Hertz 高性能框架，支持并发处理
- 🌍 **标准兼容**：完全兼容 FreeRADIUS 认证流程和 RADIUS 协议
- 📦 **容器化部署**：Docker Compose 一键部署，开箱即用
- 📱 **移动友好**：完整的移动端适配，支持触摸操作和响应式布局

## 🚀 功能特性

### 🖥️ **现代化 Web 管理界面**
- ✨ Vue 3 + TypeScript + Arco Design 构建的现代化界面
- 📱 完整的移动端适配，支持触摸操作和响应式布局
- 🎨 直观友好的用户管理界面，操作简单易懂
- ⚡ 实时数据展示和操作反馈，支持热更新

### 🔧 **强大的后端功能**
- 🔐 安全的用户认证和密码管理系统  
- 👥 全面的用户生命周期管理（创建、编辑、封禁、删除）
- 🛡️ 基于角色的精细权限控制系统
- 🔒 SHA256+Salt 密码加密，符合安全最佳实践
- 🎯 JWT 无状态认证，支持分布式部署
- 🔌 RESTful API 设计，易于集成和扩展
- 📊 详细的认证日志和统计分析
- 🌐 完全兼容 FreeRADIUS 认证流程
- 📦 Docker Compose 一键部署，支持容器化环境

## 💻 技术栈

### 🎨 前端技术栈
- **核心框架**: Vue 3.4+ + TypeScript 5.3+
- **UI 组件库**: Arco Design 2.56+ (官方推荐的企业级组件库)
- **状态管理**: Pinia (新一代 Vuex)
- **路由管理**: Vue Router 4.2+
- **HTTP 请求**: Axios 1.6+ (支持拦截器和自动重试)
- **开发工具**: Vite 5.0+ (闪电式热更新和构建)
- **代码质量**: TypeScript 类型检查 + Unplugin 自动导入

### ⚙️ 后端技术栈
- **高性能框架**: CloudWeGo Hertz 0.10+ (字节跳动开源的企业级框架)
- **数据存储**: MySQL 8.0+ (支持事务和并发控制)
- **ORM 框架**: GORM 1.30+ (功能强大的 Go ORM)
- **认证授权**: JWT (HMAC-SHA256 签名 + 自动过期刷新)
- **密码安全**: SHA256 + 32字节随机盐 加密
- **部署方式**: Docker & Docker Compose + Makefile 自动化
- **跨域支持**: CORS 中间件 + RESTful API 设计

## 🚀 快速开始

### ⚡ Docker 快速启动

30秒极速部署：

```bash
# 下载并启动完整服务
curl -fsSL https://raw.githubusercontent.com/Gaojianli/radius-manager/main/docker-compose.yml -o docker-compose.yml && docker-compose up -d
```

然后访问 http://localhost:8080 并使用 `admin` / `admin123` 登录

### 📦 使用 Docker Compose (推荐)

1. 创建 docker-compose.yml 文件
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: radius_mysql
    environment:
      MYSQL_ROOT_PASSWORD: radius123
      MYSQL_DATABASE: radius_mgnt
      MYSQL_USER: radius
      MYSQL_PASSWORD: radius123
    volumes:
      - mysql_data:/var/lib/mysql
    command: --default-authentication-plugin=mysql_native_password
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pradius123"]
      timeout: 20s
      retries: 10
      interval: 10s
      start_period: 40s
    networks:
      - radius_net

  radius_mgnt:
    image: gaojianli2333/radius-manager:latest
    container_name: radius_mgnt_app
    environment:
      # 数据库配置
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: radius
      DB_PASSWORD: radius123
      DB_NAME: radius_mgnt
      
      # 应用配置
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
      SERVER_PORT: :8080
      
      # 默认管理员用户配置（可根据需要自定义）
      DEFAULT_ADMIN_USER: admin
      DEFAULT_ADMIN_PASSWORD: admin123
      DEFAULT_ADMIN_EMAIL: admin@example.com
    ports:
      - "8080:8080"
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - radius_net
    restart: unless-stopped

volumes:
  mysql_data:

networks:
  radius_net:
    driver: bridge
```

2. 启动服务
```bash
docker-compose up -d
```

3. 访问 Web 界面
- 管理界面: http://localhost:8080
- API文档: http://localhost:8080/api/v1

默认管理员账户（可通过环境变量自定义）：
- 用户名: `admin` (或 `DEFAULT_ADMIN_USER`)
- 密码: `admin123` (或 `DEFAULT_ADMIN_PASSWORD`)

### 🔒 生产环境部署

生产环境使用时，请创建 `.env` 文件或修改环境变量：

```yaml
environment:
  # 安全：生产环境中必须修改这些值！
  JWT_SECRET: "your-production-jwt-secret-key-min-32-chars"
  DEFAULT_ADMIN_PASSWORD: "YourSecurePassword123!"
  
  # 可选：自定义管理员用户
  DEFAULT_ADMIN_USER: "admin"
  DEFAULT_ADMIN_EMAIL: "admin@yourcompany.com"
```

### 🛠️ 开发环境搭建

如需开发或从源码构建：

1. 克隆仓库
```bash
git clone https://github.com/Gaojianli/radius-manager.git
cd radius-manager
```

2. 安装依赖
```bash
make install
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，设置数据库连接信息
```

4. 启动数据库 (使用 Docker)
```bash
docker-compose up -d mysql
```

5. 构建并运行
```bash
# 构建前端
make build-frontend

# 运行完整应用
make run

# 或者分别启动前后端进行开发
make dev-frontend  # 启动前端开发服务器 (端口 3000)
make dev-backend   # 启动后端开发服务器 (端口 8080)
```

6. 构建自定义 Docker 镜像（可选）
```bash
# 构建自定义镜像
make docker-build

# 或手动构建
docker build -t your-registry/radius-manager .
```

## 🎯 Web 界面功能

### 📱 移动端优化
- **响应式设计**: 自适应桌面端、平板和手机屏幕
- **触摸优化**: 支持触摸手势、点击反馈和滑动操作
- **移动导航**: 侧边栏叠加设计，支持手势控制和遮罩层
- **移动表格**: 智能列隐藏、横向滚动和触摸优化
- **浮动操作**: FAB (浮动操作按钮) 设计，符合 Material Design 规范
- **模态框适配**: 移动端模态框自动适配屏幕尺寸

### 🔑 登录页面
- 🔒 安全的用户认证和会话管理
- ❌ 友好的错误提示和表单验证
- 💾 自动记住登录状态，支持自动登录
- 📱 完美的移动端适配和触摸优化

### 📊 仪表板
- 📈 实时系统统计和数据可视化
- 👥 用户活跃度、授权次数统计（支持按角色显示）
- ⚙️ 系统信息和版本展示
- 🎨 响应式卡片布局，自适应屏幕尺寸

### 👥 用户管理 (管理员)
- 📋 分页的用户列表，支持搜索和筛选
- ➕ 可视化的用户创建表单，实时验证
- 🔐 一键重置用户密码（管理员不可见密码）
- 🚫 灵活的用户封禁/解封功能
- 🗑️ 安全的用户删除（支持确认对话框）
- 🏷️ 用户状态实时标识和角色管理
- 🚀 移动端 FAB 浮动按钮，支持手动操作

### 👤 个人资料
- 🗺️ 详细的个人信息展示（用户名、邮箱、角色、状态）
- 🔐 安全的密码修改功能，包含密码强度检查
- 🛡️ 内置安全提示和密码安全建议
- 📅 创建时间和更新时间追踪

### 📋 认证日志 (管理员)
- 📊 详细的认证日志记录和分析
- 🔍 支持按用户名筛选和搜索
- 📅 时间范围查询和排序
- 📱 移动端适配的表格展示

## 📡 API 接口文档

### 认证相关

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

响应：
```json
{
  "code": 200,
  "expire": "2024-01-02T15:04:05Z",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 刷新 Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

### 用户管理

#### 获取个人信息
```http
GET /api/v1/user/profile
Authorization: Bearer <token>
```

#### 修改密码
```http
PUT /api/v1/user/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

### 管理员功能

#### 创建用户
```http
POST /api/v1/admin/users
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "is_admin": false
}
```

#### 获取所有用户列表
```http
GET /api/v1/admin/users?page=1&limit=20
Authorization: Bearer <admin-token>
```

#### 管理员修改用户密码
```http
PUT /api/v1/admin/users/{user_id}/password
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "new_password": "newpass123"
}
```

#### 启用/禁用用户
```http
PUT /api/v1/admin/users/{user_id}/status
Authorization: Bearer <admin-token>
```

#### 封禁/解封用户
```http
PUT /api/v1/admin/users/{user_id}/ban
Authorization: Bearer <admin-token>
```

响应：
```json
{
  "code": 200,
  "message": "User banned successfully",
  "data": {
    "id": 2,
    "username": "testuser",
    "email": "test@example.com",
    "is_admin": false,
    "status": true,
    "banned": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### 删除用户
```http
DELETE /api/v1/admin/users/{user_id}
Authorization: Bearer <admin-token>
```

### FreeRADIUS 接口

FreeRADIUS 集成采用标准的 Authorize → Authenticate 两阶段验证流程：

#### 授权接口（Authorization）
用于检查用户是否有权限进行认证，不验证密码。

```http
POST /api/v1/radius/authorize
Content-Type: application/json

{
  "username": "testuser"
}
```

响应：
```json
{
  "result": "accept",
  "username": "testuser", 
  "message": "User authorized to authenticate",
  "attrs": {
    "Auth-Type": "REST",
    "User-ID": "testuser",
    "User-Role": "user"
  }
}
```

#### 认证接口（Authentication）
用于验证用户凭据（用户名和密码）。

```http
POST /api/v1/radius/auth
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "result": "accept",
  "username": "testuser",
  "message": "Authentication successful",
  "attrs": {
    "User-ID": "testuser",
    "User-Role": "user",
    "Session-Type": "authenticated"
  }
}
```

#### 计费接口
```http
POST /api/v1/radius/accounting
Content-Type: application/json

{
  "username": "testuser",
  "acct_type": "start",
  "session_id": "session123",
  "session_time": "3600",
  "input_octets": "1024000",
  "output_octets": "2048000"
}
```

## 🔧 FreeRADIUS 集成原理

### 两阶段验证流程

FreeRADIUS 采用两阶段验证模式，确保安全性和性能：

1. **Authorize 阶段（授权）**
   - 检查用户是否存在
   - 验证账户状态（是否被禁用、封禁）
   - 设置认证类型
   - **不验证密码**，快速筛选无效用户

2. **Authenticate 阶段（认证）**
   - 验证用户凭据（密码）
   - 只对通过授权的用户执行
   - 返回最终认证结果

### 优势

- **性能优化**：无效用户在授权阶段被拦截，减少密码验证开销
- **安全增强**：被封禁用户无法进入认证阶段
- **符合标准**：遵循 RADIUS 协议设计理念
- **易于监控**：可以分别统计授权和认证的成功/失败率

## FreeRADIUS 配置

### 1. 安装 FreeRADIUS

#### Ubuntu/Debian:
```bash
sudo apt update
sudo apt install freeradius freeradius-utils
```

#### CentOS/RHEL:
```bash
sudo yum install freeradius freeradius-utils
```

### 2. 配置 REST 模块

编辑 `/etc/freeradius/3.0/mods-available/rest`:

```
rest {
    tls {
    }
    
    connect_uri = "http://localhost:8080"
    
    # 授权阶段：检查用户状态，设置认证类型
    authorize {
        uri = "${..connect_uri}/api/v1/radius/authorize"
        method = 'post'
        body = 'json'
        data = '{
            "username": "%{User-Name}"
        }'
        tls = ${..tls}
    }
    
    # 认证阶段：验证用户密码
    authenticate {
        uri = "${..connect_uri}/api/v1/radius/auth"
        method = 'post'
        body = 'json'
        data = '{
            "username": "%{User-Name}",
            "password": "%{User-Password}"
        }'
        tls = ${..tls}
    }
    
    accounting {
        uri = "${..connect_uri}/api/v1/radius/accounting"
        method = 'post'
        body = 'json'
        data = '{
            "username": "%{User-Name}",
            "acct_type": "%{Acct-Status-Type}",
            "session_id": "%{Acct-Session-Id}",
            "session_time": "%{Acct-Session-Time}",
            "input_octets": "%{Acct-Input-Octets}",
            "output_octets": "%{Acct-Output-Octets}"
        }'
        tls = ${..tls}
    }
    
    pool {
        start = ${thread[pool].start_servers}
        min = ${thread[pool].min_spare_servers}
        max = ${thread[pool].max_servers}
        spare = ${thread[pool].max_spare_servers}
        uses = 0
        retry_delay = 30
        lifetime = 0
        idle_timeout = 60
    }
}
```

### 3. 启用 REST 模块

```bash
sudo ln -s /etc/freeradius/3.0/mods-available/rest /etc/freeradius/3.0/mods-enabled/
```

### 4. 配置站点

编辑 `/etc/freeradius/3.0/sites-available/default`:

在 `authorize` 段落中添加：
```
authorize {
    # 其他配置...
    rest
    if (ok) {
        update control {
            Auth-Type := rest
        }
    }
}
```

在 `authenticate` 段落中添加：
```
authenticate {
    # 其他配置...
    Auth-Type rest{
		rest
		if(control:Auth-Type == "Accept"){
			ok
		}
	}
}
```

在 `accounting` 段落中添加：
```
accounting {
    # 其他配置...
    rest
}
```

### 5. 启动 FreeRADIUS

```bash
# 测试配置
sudo freeradius -X

# 启动服务
sudo systemctl enable freeradius
sudo systemctl start freeradius
```

### 6. 测试认证

使用 radtest 工具测试：
```bash
radtest testuser password123 localhost 0 testing123
```

## 📋 环境变量配置

| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| **数据库配置** | | |
| DB_HOST | localhost | MySQL 主机地址 |
| DB_PORT | 3306 | MySQL 端口 |
| DB_USER | root | MySQL 用户名 |
| DB_PASSWORD | - | MySQL 密码 |
| DB_NAME | radius_mgnt | 数据库名 |
| **应用配置** | | |
| JWT_SECRET | your-secret-key | JWT 签名密钥 |
| SERVER_PORT | :8080 | 服务器端口 |
| **默认管理员配置** | | |
| DEFAULT_ADMIN_USER | admin | 默认管理员用户名（仅在无管理员时创建） |
| DEFAULT_ADMIN_PASSWORD | admin123 | 默认管理员密码（仅在无管理员时创建） |
| DEFAULT_ADMIN_EMAIL | admin@example.com | 默认管理员邮箱（仅在无管理员时创建） |

### 🔐 安全注意事项

- **修改默认凭据**：生产环境中务必修改 `DEFAULT_ADMIN_PASSWORD` 和 `JWT_SECRET`
- **一次性创建**：默认管理员仅在数据库中无管理员用户时创建
- **环境变量覆盖**：使用环境变量自定义管理员凭据，无需修改代码

## 🔒 安全特性

- **密码加密**: 使用 SHA256 + 32字节随机盐
- **JWT 认证**: 使用 HMAC-SHA256 签名
- **权限控制**: 基于角色的访问控制
- **密码策略**: 最少6位字符
- **API 安全**: 所有敏感操作需要认证
- **前端路由守卫**: 未认证用户自动跳转登录页
- **自动登录过期**: Token过期自动退出

## 🏗️ 项目结构

```
radius_mgnt/
├── web/                    # Vue3 前端项目
│   ├── src/
│   │   ├── components/     # 可复用组件
│   │   ├── views/          # 页面组件
│   │   ├── layouts/        # 布局组件
│   │   ├── stores/         # Pinia状态管理
│   │   ├── services/       # API服务
│   │   └── router/         # 路由配置
│   ├── package.json
│   └── vite.config.ts
├── static/                 # 静态文件服务
│   ├── dist/               # 前端构建产物
│   └── static.go           # Go embed静态文件
├── config/                 # 配置文件
├── controllers/            # 控制器
├── dao/                    # 数据访问层
├── database/               # 数据库连接
├── middleware/             # 中间件
├── models/                 # 数据模型
├── routes/                 # 路由配置
├── scripts/                # 构建脚本
├── docker-compose.yml      # Docker编排
├── Dockerfile              # 容器构建
├── Makefile               # 构建命令
└── README.md              # 项目文档
```

## 🚀 开发指南

### 前端开发

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 启动开发服务器 (带热重载)
npm run dev

# 构建生产版本
npm run build
```

### 后端开发

```bash
# 运行后端服务
go run .

# 或使用 make 命令
make dev-backend
```

### 添加新的 API 接口

1. 在 `controllers/` 目录下创建控制器
2. 在 `routes/routes.go` 中添加路由
3. 如需数据库操作，在 `dao/` 中添加数据访问方法
4. 在 `models/` 中定义数据模型

### 添加新的前端页面

1. 在 `web/src/views/` 中创建页面组件
2. 在 `web/src/router/index.ts` 中添加路由
3. 在 `web/src/services/api.ts` 中添加API调用
4. 根据需要更新状态管理

## 🛠️ 可用命令

```bash
make help              # 显示所有可用命令
make install          # 安装前后端依赖
make build            # 构建整个项目
make build-frontend   # 只构建前端
make build-backend    # 只构建后端
make run              # 构建并运行项目
make dev-frontend     # 启动前端开发服务器
make dev-backend      # 启动后端开发服务器
make clean            # 清理构建文件
make docker-build     # 构建Docker镜像
make docker-run       # 使用Docker Compose启动
make test             # 运行测试
make format           # 格式化代码
make lint             # 代码规范检查
```

## 🐛 故障排除

### 前端无法加载

1. 确保前端已构建：`make build-frontend`
2. 检查 `static/dist/` 目录是否有文件
3. 确认服务器端口配置正确

### API 请求失败

1. 检查后端服务是否运行在正确端口
2. 验证 JWT Token 是否有效
3. 查看浏览器网络面板的错误信息

### 数据库连接问题

1. 确认 MySQL 服务运行状态
2. 检查 `.env` 文件中的数据库配置
3. 验证数据库用户权限

### FreeRADIUS 认证失败

1. 检查 REST 模块配置是否正确
2. 验证 API 端点可访问性
3. 查看 FreeRADIUS 调试输出：`sudo freeradius -X`

## 🤝 贡献指南

我们欢迎任何形式的贡献！无论是 bug 报告、功能请求、文档改进还是代码贡献。

### 📋 贡献流程

1. **Fork 项目**：点击右上角的 Fork 按钮
2. **克隆到本地**：`git clone https://github.com/Gaojianli/radius-manager.git`
3. **创建功能分支**：`git checkout -b feature/amazing-feature`
4. **本地开发**：
   ```bash
   # 安装依赖
   make install
   
   # 启动开发环境
   make dev-frontend    # 前端开发服务器
   make dev-backend     # 后端开发服务器
   ```
5. **代码规范**：
   ```bash
   make format         # 格式化代码
   make lint          # 代码规范检查
   make test          # 运行测试
   ```
6. **提交更改**：
   ```bash
   git add .
   git commit -m "✨ Add amazing feature: 功能描述"
   git push origin feature/amazing-feature
   ```
7. **创建 Pull Request**：回到 GitHub 创建 PR

### 🎯 贡献类型

- 🐛 **Bug 修复**：修复现有功能的问题
- ✨ **新功能**：添加新的功能特性
- 📚 **文档**：改进文档和说明
- 🎨 **UI/UX**：界面和用户体验改进
- ⚡ **性能**：性能优化和改进
- 🔧 **重构**：代码重构和架构优化

### 📝 代码规范

- **Go 代码**：遵循 `gofmt` 和 `golint` 标准
- **Vue 代码**：使用 TypeScript，遵循 Vue 3 组合式 API 风格
- **提交信息**：使用 [约定式提交](https://www.conventionalcommits.org/zh-hans/) 格式
- **文档**：所有公开 API 必须有对应的文档

## 🌟 致谢

感谢所有为本项目贡献的开发者！

## 📄 许可证

本项目采用 [MIT License](https://opensource.org/licenses/MIT) 开源许可证。

```
MIT License

Copyright (c) 2024 RADIUS Manager Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 📞 联系方式

- 📧 **Issues**: [GitHub Issues](https://github.com/Gaojianli/radius-manager/issues)
- 💬 **讨论**: [GitHub Discussions](https://github.com/Gaojianli/radius-manager/discussions)
- 🐛 **Bug 报告**: [Bug Report Template](https://github.com/Gaojianli/radius-manager/issues/new?template=bug_report.md)
- ✨ **功能请求**: [Feature Request Template](https://github.com/Gaojianli/radius-manager/issues/new?template=feature_request.md)

---

⭐ 如果这个项目对你有帮助，请给我们一个 **Star**！

**Happy Coding! 🎉**