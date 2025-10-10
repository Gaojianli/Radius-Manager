/**
 * Storage工具类测试
 * 验证localStorage前缀功能是否正常工作
 */

import Storage, { STORAGE_KEYS } from '../storage'

// Mock localStorage for testing
const localStorageMock = (() => {
  let store: Record<string, string> = {}

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value },
    removeItem: (key: string) => { delete store[key] },
    clear: () => { store = {} },
    key: (index: number) => Object.keys(store)[index] || null,
    get length() { return Object.keys(store).length }
  }
})()

// 替换全局localStorage
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('Storage 工具类测试', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  test('应该设置带前缀的localStorage项', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'test-token')

    expect(localStorage.getItem('radius_mgnt.auth_token')).toBe('test-token')
    expect(localStorage.getItem('auth_token')).toBeNull()
  })

  test('应该获取带前缀的localStorage项', () => {
    localStorage.setItem('radius_mgnt.auth_token', 'test-token')

    const token = Storage.getItem(STORAGE_KEYS.AUTH_TOKEN)
    expect(token).toBe('test-token')
  })

  test('应该返回默认值当项不存在时', () => {
    const token = Storage.getItem('non-existent-key', 'default-value')
    expect(token).toBe('default-value')
  })

  test('应该设置和获取JSON对象', () => {
    const user = { id: 1, username: 'test', email: 'test@example.com' }

    Storage.setJSON(STORAGE_KEYS.AUTH_USER, user)
    const retrievedUser = Storage.getJSON(STORAGE_KEYS.AUTH_USER)

    expect(retrievedUser).toEqual(user)
    expect(localStorage.getItem('radius_mgnt.auth_user')).toBe(JSON.stringify(user))
  })

  test('应该移除带前缀的localStorage项', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'test-token')
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(true)

    Storage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(false)
  })

  test('应该检查项是否存在', () => {
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(false)

    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'test-token')
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(true)
  })

  test('应该获取所有项目相关的key', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'token')
    Storage.setItem(STORAGE_KEYS.AUTH_USER, 'user')
    localStorage.setItem('other-app.data', 'other-data') // 其他应用的数据

    const projectKeys = Storage.getProjectKeys()
    expect(projectKeys).toContain('auth_token')
    expect(projectKeys).toContain('auth_user')
    expect(projectKeys).not.toContain('other-app.data')
  })

  test('应该清除所有项目相关的localStorage项', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'token')
    Storage.setItem(STORAGE_KEYS.AUTH_USER, 'user')
    localStorage.setItem('other-app.data', 'other-data') // 其他应用的数据

    Storage.clear()

    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(false)
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_USER)).toBe(false)
    expect(localStorage.getItem('other-app.data')).toBe('other-data') // 其他应用数据应该保留
  })

  test('应该正确统计存储使用情况', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'token')
    Storage.setItem(STORAGE_KEYS.AUTH_USER, 'user')
    localStorage.setItem('other-app.data', 'other-data')

    const stats = Storage.getStorageStats()
    expect(stats.totalItems).toBe(3)
    expect(stats.projectItems).toBe(2)
    expect(stats.projectKeys).toContain('auth_token')
    expect(stats.projectKeys).toContain('auth_user')
  })

  test('应该处理JSON解析错误', () => {
    localStorage.setItem('radius_mgnt.invalid_json', 'invalid-json-string')

    const result = Storage.getJSON('invalid_json', { default: 'value' })
    expect(result).toEqual({ default: 'value' })
  })

  test('应该处理undefined字符串', () => {
    localStorage.setItem('radius_mgnt.undefined_value', 'undefined')

    const result = Storage.getJSON('undefined_value', { default: 'value' })
    expect(result).toEqual({ default: 'value' })
  })

  test('应该处理null字符串', () => {
    localStorage.setItem('radius_mgnt.null_value', 'null')

    const result = Storage.getJSON('null_value', { default: 'value' })
    expect(result).toEqual({ default: 'value' })
  })

  test('应该处理空字符串', () => {
    localStorage.setItem('radius_mgnt.empty_value', '')

    const result = Storage.getJSON('empty_value', { default: 'value' })
    expect(result).toEqual({ default: 'value' })
  })

  test('设置null值时应该移除该项', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'test-token')
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(true)

    Storage.setJSON(STORAGE_KEYS.AUTH_USER, null)
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_USER)).toBe(false)
  })

  test('设置undefined值时应该移除该项', () => {
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, 'test-token')
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_TOKEN)).toBe(true)

    Storage.setJSON(STORAGE_KEYS.AUTH_USER, undefined)
    expect(Storage.hasItem(STORAGE_KEYS.AUTH_USER)).toBe(false)
  })
})