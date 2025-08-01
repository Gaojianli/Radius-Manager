<template>
  <a-layout class="layout-container">
    <a-layout-sider
      class="layout-sider"
      breakpoint="md"
      :collapsed="collapsed"
      :collapsible="true"
      :collapsed-width="0"
      @collapse="onCollapse"
    >
      <div class="logo">
        <h3 v-if="!collapsed">RADIUS 管理</h3>
        <h3 v-else">R</h3>
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
    
    <!-- 移动端遮罩层 -->
    <transition name="overlay-fade">
      <div 
        v-if="isMobile && !collapsed" 
        class="mobile-overlay"
        @click="collapsed = true"
      ></div>
    </transition>
    
    <a-layout>
      <a-layout-header :class="['layout-header', { collapsed }]">
        <div class="header-left">
          <a-button
            type="text"
            size="small"
            @click="collapsed = !collapsed"
            class="mobile-menu-btn"
            v-if="isMobile"
          >
            <template #icon>
              <icon-menu />
            </template>
          </a-button>
          <h2 class="header-title">Radius Manager</h2>
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
      
      <a-layout-content :class="['layout-content', { collapsed }]">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
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
  IconHistory,
  IconMenu
} from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const collapsed = ref(false)
const isMobile = ref(false)

// 检测是否为移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
  // 移动端默认折叠侧边栏
  if (isMobile.value) {
    collapsed.value = true
  }
}

// 监听窗口大小变化
const handleResize = () => {
  checkMobile()
}

const onCollapse = (val: boolean) => {
  collapsed.value = val
}

const onMenuClick = (key: string) => {
  router.push({ name: key })
  // 手机端点击菜单项后自动收起侧边栏
  if (isMobile.value) {
    collapsed.value = true
  }
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
  
  // 初始化移动端检测
  checkMobile()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
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
  overflow-y: auto;
}

.layout-content.collapsed {
  margin-left: 48px;
}

.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: none;
}

/* 遮罩层过渡动画 */
.overlay-fade-enter-active,
.overlay-fade-leave-active {
  transition: opacity 0.2s ease;
}

.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

.mobile-menu-btn {
  margin-right: 8px;
}

.header-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1d2129;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .layout-sider {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    height: 100vh !important;
    z-index: 1002 !important;
  }

  .layout-header {
    left: 0 !important;
    padding: 0 16px !important;
  }

  .layout-header.collapsed {
    left: 0 !important;
  }

  .layout-content {
    margin-left: 0 !important;
    padding: 16px !important;
    min-height: calc(100vh - 64px) !important;
  }

  .layout-content.collapsed {
    margin-left: 0 !important;
  }

  .header-title {
    font-size: 16px !important;
  }
  
  /* 移动端侧边栏折叠时完全隐藏 */
  .layout-container :deep(.arco-layout-sider-collapsed) {
    width: 0 !important;
    min-width: 0 !important;
  }
  
  /* 移动端显示遮罩层 */
  .mobile-overlay {
    display: block !important;
  }
  
  /* 隐藏移动端的侧边栏折叠按钮 */
  .layout-sider :deep(.arco-layout-sider-trigger) {
    display: none !important;
  }
}
</style>