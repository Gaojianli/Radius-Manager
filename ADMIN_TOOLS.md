# 管理员工具使用说明

本文档介绍如何使用管理员命令行工具来管理系统中的admin用户。

## 功能概述

管理员工具提供以下功能：
- 重置admin用户密码为默认密码
- 创建默认admin用户（如果不存在）
- 初始化admin用户（创建或重置密码）

## 使用方法

### 方法一：使用 Makefile （推荐）

```bash
# 查看所有可用命令
make help

# 重置admin密码为默认密码(admin123)
make admin-reset-password

# 创建默认admin用户
make admin-create

# 初始化admin用户(如果不存在则创建，如果存在则重置密码)
make admin-init
```

### 方法二：直接运行 Go 命令

```bash
# 重置admin密码
go run cmd/admin_tool.go reset-password

# 创建admin用户
go run cmd/admin_tool.go create-admin

# 初始化admin用户
go run cmd/admin_tool.go init

# 查看帮助信息
go run cmd/admin_tool.go
```

## 默认凭据

- **用户名**: `admin`
- **密码**: `admin123`
- **邮箱**: `admin@example.com`

## 使用场景

### 1. 首次部署系统

在首次部署系统时，推荐使用 `admin-init` 命令来创建初始管理员用户：

```bash
make admin-init
```

### 2. 忘记admin密码

如果忘记了admin密码，可以使用 `admin-reset-password` 命令重置为默认密码：

```bash
make admin-reset-password
```

### 3. 重新创建admin用户

如果admin用户被意外删除，可以使用 `admin-create` 命令重新创建：

```bash
make admin-create
```

## 安全注意事项

⚠️ **重要安全提示**：

1. 在使用默认密码登录后，**请立即修改密码**
2. 建议使用强密码（包含大小写字母、数字和特殊字符）
3. 定期更换管理员密码
4. 在生产环境中，确保只有授权人员可以访问这些管理员工具

## 错误处理

### 常见错误及解决方案

1. **数据库连接失败**
   - 检查数据库配置是否正确
   - 确保数据库服务正在运行

2. **admin用户不存在**（在使用reset-password时）
   - 先运行 `make admin-create` 创建用户
   - 或者使用 `make admin-init` 自动处理

3. **权限问题**
   - 确保有足够的权限访问数据库
   - 检查配置文件权限

## 配置要求

在使用管理员工具之前，请确保：

1. 已正确配置数据库连接
2. 配置文件存在且可读
3. 系统具有必要的网络和文件权限

## 示例输出

### 成功重置密码
```
🔐 重置admin密码...
✅ Admin user 'admin' password has been reset to: admin123
⚠️  Please change the default password after login for security reasons.
```

### 成功创建用户
```
👤 创建admin用户...
✅ Admin user created successfully:
   Username: admin
   Password: admin123
   Email: admin@example.com
⚠️  Please change the default password after login for security reasons.
```

### 初始化用户
```
🚀 初始化admin用户...
Admin user 'admin' found, resetting password...
✅ Admin user 'admin' password has been reset to: admin123
⚠️  Please change the default password after login for security reasons.
```