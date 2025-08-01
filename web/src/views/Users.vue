<template>
  <a-card class="page-container" :bordered="false">
    <a-page-header title="用户管理" subtitle="管理系统用户">
      <template #extra>
        <a-button type="primary" @click="showCreateModal = true">
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

    <!-- 创建用户模态框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="添加用户"
      @ok="handleCreateUser"
      :confirm-loading="createLoading"
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
</style>