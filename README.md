# RADIUS 管理系统

基于 Hertz 框架和 GORM 构建的 RADIUS 用户管理后端系统，为 FreeRADIUS 提供用户认证服务。

## 功能特性

- ✅ 用户登录、密码管理
- ✅ 管理员创建、删除用户
- ✅ 管理员可以修改所有用户密码（无法查看密码）
- ✅ 基于 SHA256+Salt 的安全密码存储
- ✅ JWT 认证和权限管理
- ✅ RESTful API 接口
- ✅ FreeRADIUS 认证集成
- ✅ Docker 容器化部署

## 技术栈

- **框架**: CloudWeGo Hertz
- **数据库**: MySQL 8.0
- **ORM**: GORM
- **认证**: JWT
- **密码加密**: SHA256+Salt
- **容器化**: Docker & Docker Compose

## 快速开始

### 使用 Docker Compose (推荐)

1. 克隆项目并进入目录
```bash
git clone <repository-url>
cd radius_mgnt
```

2. 启动服务
```bash
docker-compose up -d
```

3. 检查服务状态
```bash
docker-compose ps
```

服务将在以下端口启动：
- API 服务: http://localhost:8080
- MySQL 数据库: localhost:3306

默认管理员账户：
- 用户名: `admin`
- 密码: `admin123`

### 手动部署

1. 安装依赖
```bash
go mod tidy
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，设置数据库连接信息
```

3. 启动 MySQL 数据库

4. 运行应用
```bash
go run main.go
# 或者使用启动脚本
./scripts/start.sh
```

## API 接口文档

### 认证相关

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
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

## FreeRADIUS 集成原理

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
    rest
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

## 环境变量配置

| 变量名 | 默认值 | 描述 |
|--------|--------|------|
| DB_HOST | localhost | MySQL 主机地址 |
| DB_PORT | 3306 | MySQL 端口 |
| DB_USER | root | MySQL 用户名 |
| DB_PASSWORD | - | MySQL 密码 |
| DB_NAME | radius_mgnt | 数据库名 |
| JWT_SECRET | your-secret-key | JWT 签名密钥 |
| SERVER_PORT | :8080 | 服务器端口 |

## 安全特性

- **密码加密**: 使用 SHA256 + 32字节随机盐
- **JWT 认证**: 使用 HMAC-SHA256 签名
- **权限控制**: 基于角色的访问控制
- **密码策略**: 最少6位字符
- **API 安全**: 所有敏感操作需要认证

## 项目结构

```
radius_mgnt/
├── config/          # 配置文件
├── controllers/     # 控制器
├── database/        # 数据库连接
├── middleware/      # 中间件
├── models/         # 数据模型
├── routes/         # 路由配置
├── scripts/        # 脚本文件
├── docker-compose.yml
├── Dockerfile
├── main.go
└── README.md
```

## 开发指南

### 添加新的 API 接口

1. 在 `controllers/` 目录下创建控制器
2. 在 `routes/routes.go` 中添加路由
3. 如需数据库操作，在 `models/` 中定义模型

### 数据库迁移

应用启动时会自动创建和更新数据库表结构。

### 日志

应用使用内置日志记录重要操作和错误。

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 MySQL 服务是否运行
   - 验证数据库连接参数
   - 确认数据库用户权限

2. **FreeRADIUS 认证失败**
   - 检查 REST 模块配置
   - 验证 API 端点可访问性
   - 查看 FreeRADIUS 调试输出

3. **JWT Token 无效**
   - 检查 JWT_SECRET 配置
   - 确认 token 未过期
   - 验证 token 格式

### 调试模式

以调试模式启动 FreeRADIUS：
```bash
sudo freeradius -X
```

查看应用日志：
```bash
docker-compose logs -f radius_mgnt
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。