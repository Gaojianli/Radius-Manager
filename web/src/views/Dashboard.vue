<template>
  <div class="page-container">
    <a-page-header title="仪表板" subtitle="系统概览" :show-back="false" />
    
    <a-row :gutter="16" class="mb-24">
      <a-col :xs="12" :sm="12" :md="6" :lg="6" v-if="authStore.isAdmin">
        <a-card class="dashboard-card">
          <a-statistic
            title="总用户数"
            :value="stats.totalUsers"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <icon-user-group />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      
      <a-col :xs="12" :sm="12" :md="6" :lg="6" v-if="authStore.isAdmin">
        <a-card class="dashboard-card">
          <a-statistic
            title="活跃用户"
            :value="stats.activeUsers"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <icon-user />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      
      <a-col :xs="12" :sm="12" :md="6" :lg="6" v-if="authStore.isAdmin">
        <a-card class="dashboard-card">
          <a-statistic
            title="被封禁用户"
            :value="stats.bannedUsers"
            :value-style="{ color: '#ff4d4f' }"
          >
            <template #prefix>
              <icon-stop />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      
      <a-col :xs="12" :sm="12" :md="6" :lg="6">
        <a-card class="dashboard-card">
          <a-statistic
            title="授权次数"
            :value="stats.authCount"
            suffix="次"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <icon-safe />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :xs="24" :sm="24" :md="12" :lg="12">
        <a-card title="系统信息" class="mb-16 dashboard-card">
          <a-descriptions :column="1">
            <a-descriptions-item label="系统版本">v1.0.0</a-descriptions-item>
            <a-descriptions-item label="数据库">MySQL</a-descriptions-item>
            <a-descriptions-item label="认证方式">SHA256+Salt</a-descriptions-item>
            <a-descriptions-item label="RADIUS集成">已启用</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :sm="24" :md="12" :lg="12">
        <a-card title="欢迎使用 RADIUS 管理系统" class="dashboard-card">
          <p>您已成功登录 RADIUS 管理系统。本系统提供以下功能：</p>
          <ul>
            <li v-if="authStore.isAdmin">👥 <strong>用户管理</strong>：创建、编辑、删除用户</li>
            <li v-if="authStore.isAdmin">🔒 <strong>密码管理</strong>：重置用户密码</li>
            <li v-if="authStore.isAdmin">🚫 <strong>用户封禁</strong>：临时或永久封禁用户</li>
            <li>🔐 <strong>密码安全</strong>：SHA256+Salt 加密存储</li>
            <li>🌐 <strong>RADIUS集成</strong>：与FreeRADIUS无缝对接</li>
            <li>📊 <strong>实时监控</strong>：用户状态和连接监控</li>
          </ul>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { 
  IconUserGroup, 
  IconUser, 
  IconStop, 
  IconSafe 
} from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'
import { getStats, getAdminStats } from '@/services/api'
import { Message } from '@arco-design/web-vue'

const authStore = useAuthStore()

const stats = reactive({
  totalUsers: 0,
  activeUsers: 0,
  bannedUsers: 0,
  authCount: 0
})

const loadStats = async () => {
  if (authStore.isAdmin) {
    try {
      const response = await getAdminStats()
      stats.totalUsers = response.data.total_users
      stats.activeUsers = response.data.active_users
      stats.bannedUsers = response.data.banned_users
      stats.authCount = response.data.auth_count || 0
    } catch (error) {
      console.error('Failed to load admin stats:', error)
      Message.error('加载统计数据失败')
    }
  } else {
    // 非管理员用户只能查看自己的授权次数
    try {
      const response = await getStats()
      stats.authCount = response.data.auth_count || 0
    } catch (error) {
      console.error('Failed to load auth count:', error)
    }
  }
}

onMounted(loadStats)
</script>

<style scoped>
.dashboard-card {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.dashboard-card ul {
  margin: 16px 0;
  padding-left: 24px;
}

.dashboard-card li {
  margin-bottom: 8px;
  line-height: 1.6;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .mb-24 {
    margin-bottom: 16px !important;
  }
  
  .dashboard-card {
    margin-bottom: 16px;
  }
  
  .dashboard-card ul {
    margin: 12px 0;
  }
  
  .dashboard-card li {
    margin-bottom: 6px;
    font-size: 14px;
  }
}
</style>