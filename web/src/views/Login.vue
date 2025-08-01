<template>
  <div class="login-container full-height center-content">
    <a-card class="login-card" title="" :bordered="false">
      <div class="login-header">
        <h2>RADIUS 管理系统</h2>
        <p>请登录以继续</p>
      </div>
      
      <a-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="rules"
        @submit="handleSubmit"
      >
        <a-form-item field="username" hide-label>
          <a-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
          >
            <template #prefix>
              <IconUser />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item field="password" hide-label>
          <a-input-password
            v-model="loginForm.password"
            placeholder="密码"
            size="large"
          >
            <template #prefix>
              <IconLock />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item hide-label style="display: flex; justify-content: center;">
          <a-button
            type="primary"
            html-type="submit"
            long
            size="large"
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
      
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type { LoginRequest } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Message, type FormInstance } from '@arco-design/web-vue'
import { IconLock, IconUser } from '@arco-design/web-vue/es/icon'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const loginFormRef = ref<FormInstance>()

const loginForm = reactive<LoginRequest>({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名' },
    { min: 3, message: '用户名至少3个字符' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6个字符' }
  ]
}

const handleSubmit = async (data: { errors?: any; values: any }) => {
  
  const { errors, values } = data
  const errorFields = errors ? Object.keys(errors) : []
  if (errorFields.length > 0) {
    loginFormRef.value?.scrollToField(errorFields[0])
    return
  }
  
  const { username, password } = values
  if (loading.value) return
  loading.value = true
  
  try {
    await authStore.loginUser({ username, password })
    // 登录成功后加载用户信息
    await authStore.loadUserInfo()
    Message.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('Login failed:', error)
    loginFormRef.value?.setFields({
      password: {
        status: 'error',
        message: `密码错误`
      }
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.login-card {
  width: 400px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h2 {
  margin: 0 0 8px 0;
  color: #1d2129;
  font-size: 24px;
  font-weight: 600;
}

.login-header p {
  margin: 0;
  color: #86909c;
  font-size: 14px;
}


:deep(.arco-card-body) {
  padding: 32px;
}
</style>