<template>
  <div class="page-container">
    <a-page-header title="个人资料" subtitle="管理您的账户信息" :show-back="false" />
    
    <a-row :gutter="24">
      <a-col :xs="24" :sm="24" :md="8" :lg="8">
        <a-card title="基本信息" class="profile-card">
          <a-descriptions :column="1" v-if="authStore.user">
            <a-descriptions-item label="用户名">
              {{ authStore.user.username }}
            </a-descriptions-item>
            <a-descriptions-item label="邮箱">
              {{ authStore.user.email }}
            </a-descriptions-item>
            <a-descriptions-item label="角色">
              <a-tag :color="authStore.user.is_admin ? 'red' : 'blue'">
                {{ authStore.user.is_admin ? '管理员' : '普通用户' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag :color="authStore.user.banned ? 'gray' : 'green'">
                {{ authStore.user.banned ? '禁用' : '正常' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">
              {{ formatDate(authStore.user.created_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="更新时间">
              {{ formatDate(authStore.user.updated_at) }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :sm="24" :md="16" :lg="16">
        <a-card title="修改密码" class="profile-card">
          <a-alert type="info" show-icon class="security-tips">
            <template #icon><icon-info-circle /></template>
            <div>
              <h4>密码安全建议：</h4>
              <ul>
                <li>密码长度至少6个字符</li>
                <li>包含字母、数字和特殊字符</li>
                <li>避免使用常见密码</li>
                <li>定期更换密码</li>
                <li>不要在多个系统中使用相同密码</li>
              </ul>
            </div>
          </a-alert>
          
          <a-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            layout="vertical"
            style="max-width: 400px"
          >
            <a-form-item field="old_password" label="当前密码">
              <a-input-password
                v-model="passwordForm.old_password"
                placeholder="请输入当前密码"
              />
            </a-form-item>
            
            <a-form-item field="new_password" label="新密码">
              <a-input-password
                v-model="passwordForm.new_password"
                placeholder="请输入新密码"
              />
            </a-form-item>
            
            <a-form-item field="confirm_password" label="确认新密码">
              <a-input-password
                v-model="passwordForm.confirm_password"
                placeholder="请再次输入新密码"
              />
            </a-form-item>
            
            <a-form-item>
              <a-button
                type="primary"
                @click="handleChangePassword"
                :loading="loading"
              >
                修改密码
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { changePassword, type ChangePasswordRequest } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Message } from '@arco-design/web-vue'
import { IconInfoCircle } from '@arco-design/web-vue/es/icon'
import { reactive, ref } from 'vue'

const authStore = useAuthStore()

const loading = ref(false)
const passwordFormRef = ref()

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const passwordRules = {
  old_password: [
    { required: true, message: '请输入当前密码' }
  ],
  new_password: [
    { required: true, message: '请输入新密码' },
    { min: 6, message: '密码至少6个字符' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码' },
    {
      validator: (value: string, cb: (error?: string) => void) => {
        if (value !== passwordForm.new_password) {
          cb('两次输入的密码不一致')
        } else {
          cb()
        }
      }
    }
  ]
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
  } catch (error) {
    return // 验证失败
  }
  
  loading.value = true
  
  try {
    const changePasswordData: ChangePasswordRequest = {
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    }
    
    await changePassword(changePasswordData)
    Message.success('密码修改成功')
    
    // 重置表单
    Object.assign(passwordForm, {
      old_password: '',
      new_password: '',
      confirm_password: ''
    })
    
  } catch (error) {
    console.error('Failed to change password:', error)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>

.profile-card {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.security-tips {
  margin-bottom: 24px;
}

.security-tips ul {
  margin: 8px 0 0 0;
  padding-left: 24px;
}

.security-tips li {
  margin-bottom: 4px;
}

.security-tips h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .page-container {
    margin: 16px auto;
    padding: 16px;
    min-height: calc(100vh - 80px);
  }
  
  .profile-card {
    margin-bottom: 16px;
  }
  
  .security-tips ul {
    padding-left: 20px;
  }
  
  .security-tips li {
    font-size: 14px;
  }
}
</style>