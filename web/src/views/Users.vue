<template>
  <a-card class="page-container" :bordered="false">
    <a-page-header title="用户管理" subtitle="管理系统用户" :show-back="false">
      <template #extra>
        <a-button type="primary" @click="showCreateModal = true" class="desktop-add-btn">
          <template #icon><icon-plus /></template>
          添加用户
        </a-button>
      </template>
    </a-page-header>

    <a-table
      :data="users"
      :pagination="pagination"
      :loading="loading"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
      row-key="id"
    >
      <template #columns>
        <a-table-column title="ID" data-index="id" />
        <a-table-column title="用户名" data-index="username" />
        <a-table-column title="邮箱" data-index="email" />
        <a-table-column title="角色" data-index="is_admin">
          <template #cell="{ record }">
            <a-tag :color="record.is_admin ? 'red' : 'blue'">
              {{ record.is_admin ? '管理员' : '普通用户' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="封禁状态" data-index="banned">
          <template #cell="{ record }">
            <a-tag :color="record.banned ? 'red' : 'green'">
              {{ record.banned ? '已封禁' : '正常' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="created_at">
          <template #cell="{ record }">
            {{ formatDate(record.created_at) }}
          </template>
        </a-table-column>
        <a-table-column title="操作" align="center" :width="200">
          <template #cell="{ record }">
            <a-space>
              <a-tooltip content="修改密码">
                <a-button 
                  type="text" 
                  size="small"
                  shape="circle"
                  class="action-btn edit-btn"
                  @click="router.push({ name: 'Profile' })"
                >
                  <template #icon><icon-edit /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip :content="record.banned ? '解封用户' : '封禁用户'">
                <a-button 
                  type="text" 
                  size="small"
                  shape="circle"
                  :class="['action-btn', record.banned ? 'unban-btn' : 'ban-btn']"
                  @click="toggleUserBan(record)"
                >
                  <template #icon><icon-stop /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip :content="record.id === authStore.user?.id ? '不能删除自己' : '删除用户'">
                <a-button 
                  type="text" 
                  size="small"
                  shape="circle"
                  class="action-btn delete-btn"
                  @click="showDeleteModal(record)"
                  :disabled="record.id === authStore.user?.id"
                >
                  <template #icon><icon-delete /></template>
                </a-button>
              </a-tooltip>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 移动端悬浮添加按钮 -->
    <a-button 
      type="primary" 
      shape="circle" 
      size="large"
      @click="showCreateModal = true"
      class="mobile-fab"
    >
      <template #icon><icon-plus /></template>
    </a-button>

    <!-- 创建用户模态框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="添加用户"
      @ok="handleCreateUser"
      :confirm-loading="createLoading"
      class="create-user-modal"
      :width="400"
    >
      <a-form ref="createFormRef" :model="createForm" :rules="createRules">
        <a-form-item field="username" label="用户名">
          <a-input v-model="createForm.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item field="email" label="邮箱">
          <a-input v-model="createForm.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="createForm.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item field="is_admin" label="角色">
          <a-switch v-model="createForm.is_admin">
            <template #checked>管理员</template>
            <template #unchecked>普通用户</template>
          </a-switch>
        </a-form-item>
      </a-form>
    </a-modal>

  </a-card>
</template>

<script setup lang="ts">
import {
  deleteUser as apiDeleteUser,
  toggleUserBan as apiToggleUserBan,
  createUser,
  getUsers,
  type CreateUserRequest,
  type User
} from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconEdit, IconStop, IconDelete } from '@arco-design/web-vue/es/icon'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

// 表格数据
const users = ref<User[]>([])
const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100']
})

// 创建用户
const showCreateModal = ref(false)
const createLoading = ref(false)
const createFormRef = ref()
const createForm = reactive<CreateUserRequest & { is_admin: boolean }>({
  username: '',
  email: '',
  password: '',
  is_admin: false
})

const createRules = {
  username: [
    { required: true, message: '请输入用户名' },
    { min: 3, max: 50, message: '用户名长度在3-50个字符' }
  ],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6个字符' }
  ]
}


// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await getUsers(pagination.current, pagination.pageSize)
    users.value = response.data.users
    pagination.total = response.data.pagination.total
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

// 分页处理
const handlePageChange = (page: number) => {
  pagination.current = page
  loadUsers()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadUsers()
}

// 创建用户
const handleCreateUser = async () => {
  if (!createFormRef.value) return
  
  try {
    await createFormRef.value.validate()
  } catch (error) {
    return // 验证失败
  }
  
  createLoading.value = true
  try {
    await createUser(createForm)
    Message.success('用户创建成功')
    showCreateModal.value = false
    resetCreateForm()
    loadUsers()
  } catch (error) {
    console.error('Failed to create user:', error)
  } finally {
    createLoading.value = false
  }
}

const resetCreateForm = () => {
  Object.assign(createForm, {
    username: '',
    email: '',
    password: '',
    is_admin: false
  })
}


// 切换封禁状态
const toggleUserBan = async (user: User) => {
  Modal.confirm({
    title: `确认${user.banned ? '解封' : '封禁'}用户`,
    content: `您确定要${user.banned ? '解封' : '封禁'}用户 "${user.username}" 吗？`,
    onOk: async () => {
      try {
        const updatedUser = await apiToggleUserBan(user.id)
        Message.success(`用户已${updatedUser.banned ? '封禁' : '解封'}`)
        loadUsers()
      } catch (error) {
        console.error('Failed to toggle user ban:', error)
      }
    }
  })
}

// 删除用户
const showDeleteModal = (user: User) => {
  Modal.confirm({
    title: '确认删除用户',
    content: `您确定要删除用户 "${user.username}" 吗？此操作不可撤销！`,
    onOk: async () => {
      try {
        await apiDeleteUser(user.id)
        Message.success('用户删除成功')
        loadUsers()
      } catch (error) {
        console.error('Failed to delete user:', error)
      }
    }
  })
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.page-container {
  margin: 24px auto;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  max-width: 1200px;
}

.action-btn {
  color: #8c8c8c !important;
  transition: all 0.2s ease;
}

.action-btn:hover:not(:disabled) {
  transform: scale(1.1);
}

.edit-btn:hover:not(:disabled) {
  color: #1890ff !important;
  background-color: rgba(24, 144, 255, 0.1) !important;
}

.ban-btn:hover:not(:disabled) {
  color: #fa8c16 !important;
  background-color: rgba(250, 140, 22, 0.1) !important;
}

.unban-btn:hover:not(:disabled) {
  color: #52c41a !important;
  background-color: rgba(82, 196, 26, 0.1) !important;
}

.delete-btn:hover:not(:disabled) {
  color: #ff4d4f !important;
  background-color: rgba(255, 77, 79, 0.1) !important;
}

.action-btn:disabled {
  color: #d9d9d9 !important;
}

/* 悬浮按钮样式 */
.mobile-fab {
  display: none;
  position: fixed !important;
  bottom: 80px;
  right: 24px;
  width: 56px !important;
  height: 56px !important;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15), 0 8px 24px rgba(0, 0, 0, 0.1);
  border: none !important;
  transition: all 0.3s ease;
}

.mobile-fab:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2), 0 12px 32px rgba(0, 0, 0, 0.15);
}

/* 桌面端添加按钮默认显示 */
.desktop-add-btn {
  display: inline-flex;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .page-container {
    margin: 16px auto;
    padding: 16px;
    overflow-x: hidden;
  }
  
  /* 隐藏桌面端添加按钮 */
  .desktop-add-btn {
    display: none !important;
  }
  
  /* 显示移动端悬浮按钮 */
  .mobile-fab {
    display: inline-flex !important;
  }
  
  /* 在手机端隐藏不重要的列 */
  :deep(.arco-table-th):nth-child(1),
  :deep(.arco-table-td):nth-child(1) {
    display: none; /* 隐藏ID列 */
  }
  
  :deep(.arco-table-th):nth-child(3),
  :deep(.arco-table-td):nth-child(3) {
    display: none; /* 隐藏邮箱列 */
  }
  
  :deep(.arco-table-th):nth-child(6),
  :deep(.arco-table-td):nth-child(6) {
    display: none; /* 隐藏创建时间列 */
  }
  
  /* 调整剩余列的宽度 */
  :deep(.arco-table-th):nth-child(2),
  :deep(.arco-table-td):nth-child(2) {
    width: 30%; /* 用户名列 - 弹性宽度 */
    min-width: 80px;
  }
  
  :deep(.arco-table-th):nth-child(4),
  :deep(.arco-table-td):nth-child(4) {
    width: 25%; /* 角色列 - 弹性宽度 */
    min-width: 60px;
  }
  
  :deep(.arco-table-th):nth-child(5),
  :deep(.arco-table-td):nth-child(5) {
    width: 25%; /* 封禁状态列 - 弹性宽度 */
    min-width: 60px;
  }
  
  /* 调整操作列 */
  :deep(.arco-table-th):last-child,
  :deep(.arco-table-td):last-child {
    width: 20% !important;
    min-width: 90px !important;
  }
  
  /* 表格整体设置 */
  :deep(.arco-table) {
    font-size: 14px;
    min-width: 100%;
    width: 100%;
  }
  
  :deep(.arco-table-container) {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    width: 100%;
  }
  
  :deep(.arco-table-tbody),
  :deep(.arco-table-thead) {
    width: 100%;
  }
  
  :deep(.arco-table-th) {
    padding: 8px 4px !important;
  }
  
  :deep(.arco-table-td) {
    padding: 12px 4px !important;
  }
  
  /* 调整标签样式 */
  :deep(.arco-tag) {
    font-size: 12px;
    padding: 2px 6px;
  }
  
  /* 调整按钮大小 */
  .action-btn {
    width: 26px !important;
    height: 26px !important;
    margin: 0 2px;
  }
  
  /* 用户名列内容换行 */
  :deep(.arco-table-td):nth-child(2) {
    word-break: break-all;
  }
  
  /* 移动端模态框适配 */
  .create-user-modal :deep(.arco-modal) {
    width: calc(100vw - 32px) !important;
    max-width: none !important;
    margin: 16px !important;
  }
  
  .create-user-modal :deep(.arco-modal-container) {
    padding: 16px !important;
  }
  
  .create-user-modal :deep(.arco-modal-header) {
    padding: 16px 20px 12px !important;
  }
  
  .create-user-modal :deep(.arco-modal-body) {
    padding: 12px 20px 20px !important;
    max-height: calc(100vh - 200px);
    overflow-y: auto;
  }
  
  .create-user-modal :deep(.arco-modal-footer) {
    padding: 12px 20px 16px !important;
  }
  
  .create-user-modal :deep(.arco-form-item) {
    margin-bottom: 16px !important;
  }
  
  .create-user-modal :deep(.arco-form-item-label) {
    font-size: 14px !important;
  }
  
  .create-user-modal :deep(.arco-input),
  .create-user-modal :deep(.arco-input-password) {
    font-size: 16px !important; /* 防止iOS缩放 */
  }
}
</style>