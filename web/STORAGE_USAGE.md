# localStorage 使用指南

本项目已升级localStorage管理方式，统一使用带前缀的Storage工具类，避免与同域名下其他项目的存储冲突。

## 存储前缀

所有localStorage键都会自动添加前缀：`radius_mgnt.`

例如：
- `auth_token` → `radius_mgnt.auth_token`
- `auth_user` → `radius_mgnt.auth_user`

## 使用方法

### 1. 导入Storage工具类

```typescript
import Storage, { STORAGE_KEYS } from '@/utils/storage'
```

### 2. 基本操作

#### 设置值
```typescript
// 设置字符串
Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'your-token')

// 设置JSON对象
const user = { id: 1, username: 'admin' }
Storage.setJSON(STORAGE_KEYS.AUTH_USER, user)
```

#### 获取值
```typescript
// 获取字符串（带默认值）
const token = Storage.getItem(STORAGE_KEYS.AUTH_TOKEN, '')

// 获取JSON对象
const user = Storage.getJSON<User>(STORAGE_KEYS.AUTH_USER)
```

#### 删除值
```typescript
// 删除单个项
Storage.removeItem(STORAGE_KEYS.AUTH_TOKEN)

// 清除所有项目相关的localStorage
Storage.clear()
```

#### 检查存在性
```typescript
if (Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)) {
  // token 存在
}
```

### 3. 预定义的存储键

项目中已预定义了常用的存储键常量：

```typescript
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  AUTH_USER: 'auth_user',
  USER_PREFERENCES: 'user_preferences',
  THEME: 'theme',
  LANGUAGE: 'language'
} as const
```

### 4. 高级功能

#### 获取项目相关的所有键
```typescript
const projectKeys = Storage.getProjectKeys()
console.log(projectKeys) // ['auth_token', 'auth_user', ...]
```

#### 存储使用情况统计
```typescript
const stats = Storage.getStorageStats()
console.log(stats)
// {
//   totalItems: 10,      // localStorage中的总项数
//   projectItems: 3,     // 本项目的项数
//   projectKeys: [...]   // 本项目的所有键
// }
```

## 迁移说明

### 旧代码 ❌
```typescript
// 直接使用localStorage（已废弃）
localStorage.setItem('auth_token', token)
const token = localStorage.getItem('auth_token')
localStorage.removeItem('auth_token')
```

### 新代码 ✅
```typescript
// 使用Storage工具类
Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, token)
const token = Storage.getItem(STORAGE_KEYS.AUTH_TOKEN)
Storage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
```

## 实际存储效果

在浏览器开发者工具的Application > Local Storage中，你会看到：

```
radius_mgnt.auth_token: "eyJhbGciOiJIUzI1NiIs..."
radius_mgnt.auth_user: '{"id":1,"username":"admin"}'
```

而不是：
```
auth_token: "eyJhbGciOiJIUzI1NiIs..."
auth_user: '{"id":1,"username":"admin"}'
```

## 优势

1. **避免冲突**：防止与同域名下其他项目的localStorage冲突
2. **统一管理**：所有存储操作通过统一接口
3. **类型安全**：使用TypeScript提供类型支持
4. **易于维护**：集中管理存储键名，避免硬编码
5. **错误处理**：内置错误处理机制
6. **调试友好**：提供统计和查看功能

## 错误处理和边界情况

Storage工具类已经内置了完善的错误处理：

### JSON解析保护
```typescript
// 安全处理各种异常情况
const user = Storage.getJSON<User>(STORAGE_KEYS.AUTH_USER, null)
// 自动处理以下情况：
// - localStorage中的值为 null
// - localStorage中的值为 "undefined" 字符串
// - localStorage中的值为 "null" 字符串
// - localStorage中的值为空字符串
// - localStorage中的值为无效JSON
```

### 设置null/undefined值
```typescript
// 设置null或undefined会自动清除该项
Storage.setJSON(STORAGE_KEYS.AUTH_USER, null)      // 等同于 removeItem
Storage.setJSON(STORAGE_KEYS.AUTH_USER, undefined) // 等同于 removeItem
```

## 注意事项

1. **只使用预定义的键**：尽量使用 `STORAGE_KEYS` 中定义的常量
2. **添加新键**：需要新的存储键时，请先在 `STORAGE_KEYS` 中定义
3. **不要直接使用localStorage**：始终通过Storage工具类操作
4. **JSON存储**：对象类型数据使用 `setJSON/getJSON` 方法
5. **默认值**：始终为getJSON方法提供合理的默认值

## 开发建议

在开发新功能时：

1. 如果需要新的存储键，先在 `STORAGE_KEYS` 中定义
2. 使用TypeScript类型来确保类型安全
3. 考虑是否需要默认值
4. 在组件卸载时考虑是否需要清理存储