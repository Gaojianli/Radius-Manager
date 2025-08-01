<template>
  <a-layout class="layout-container">
    <a-layout-sider
      class="layout-sider"
      breakpoint="lg"
      :collapsed="collapsed"
      :collapsible="true"
      @collapse="onCollapse"
    >
      <div class="logo">
        <h3 v-if="!collapsed">RADIUS 管理</h3>
        <h3 v-else>R</h3>
      </div>
      
      <a-menu
        :default-selected-keys="[route.name as string]"
        :style="{width: '100%'}"
        @menu-item-click="onMenuClick"
      >
        <a-menu-item key="Dashboard">
          <template #icon><icon-dashboard /></template>
          仪表板
        </a-menu-item>
        
        <a-menu-item key="Users" v-if="authStore.isAdmin">
          <template #icon><icon-user-group /></template>
          用户管理
        </a-menu-item>
        
        <a-menu-item key="AuthLogs" v-if="authStore.isAdmin">
          <template #icon><icon-history /></template>
          认证日志
        </a-menu-item>
        
        <a-menu-item key="Profile">
          <template #icon><icon-user /></template>
          个人资料
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    
    <a-layout>
      <a-layout-header class="layout-header">
        <div class="header-left">
          <a-button
            type="text"
            size="small"
            @click="collapsed = !collapsed"
          >
            <template #icon>
              <icon-menu-unfold v-if="collapsed" />
              <icon-menu-fold v-else />
            </template>
          </a-button>
        </div>
        
        <div class="header-right">
          <a-dropdown @select="onUserMenuSelect">
            <a-button type="text">
              <template #icon><icon-user /></template>
              {{ authStore.user?.username }}
              <template #suffix><icon-down /></template>
            </a-button>
            <template #content>
              <a-doption value="profile">
                <template #icon><icon-user /></template>
                个人资料
              </a-doption>
              <a-doption value="logout">
                <template #icon><icon-export /></template>
                退出登录
              </a-doption>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>
      
      <a-layout-content class="layout-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, Message } from '@arco-design/web-vue'
import {
  IconDashboard,
  IconUserGroup,
  IconUser,
  IconMenuFold,
  IconMenuUnfold,
  IconDown,
  IconExport,
  IconHistory
} from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const collapsed = ref(false)

const onCollapse = (val: boolean) => {
  collapsed.value = val
}

const onMenuClick = (key: string) => {
  router.push({ name: key })
}

const onUserMenuSelect = (value: string) => {
  if (value === 'profile') {
    router.push({ name: 'Profile' })
  } else if (value === 'logout') {
    Modal.confirm({
      title: '确认退出',
      content: '您确定要退出登录吗？',
      onOk: () => {
        authStore.logout()
        Message.success('已退出登录')
        router.push('/login')
      }
    })
  }
}

onMounted(async () => {
  try {
    await authStore.loadUserInfo()
  } catch (error) {
    router.push('/login')
  }
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-sider {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  z-index: 999;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  margin-bottom: 8px;
}

.logo h3 {
  color: #fff;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.layout-header {
  position: fixed;
  top: 0;
  left: 200px;
  right: 0;
  height: 64px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 998;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  transition: left 0.2s ease;
}

.layout-header.collapsed {
  left: 48px;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
}

.layout-content {
  margin-left: 200px;
  margin-top: 64px;
  padding: 24px;
  transition: margin-left 0.2s ease;
  min-height: calc(100vh - 64px);
  background: #f5f5f5;
}

.layout-content.collapsed {
  margin-left: 48px;
}

:deep(.arco-layout-sider-collapsed .layout-content) {
  margin-left: 48px;
}

:deep(.arco-layout-sider-collapsed + .arco-layout .layout-header) {
  left: 48px;
}
</style>