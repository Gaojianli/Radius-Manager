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
                  @click="showChangePasswordModal(record)"
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

    <!-- 修改密码模态框 -->
    <a-modal
      v-model:visible="showPasswordModal"
      :title="`修改用户密码 - ${selectedUser?.username}`"
      @ok="handleChangeUserPassword"
      :confirm-loading="passwordLoading"
      class="change-password-modal"
      :width="500"
    >
      <div class="password-change-content">
        <div class="user-info-section">
          <div class="section-header">
            <icon-user class="section-icon" />
            <span class="section-title">当前操作用户</span>
          </div>
          <a-descriptions :column="1" size="small" class="user-descriptions">
            <a-descriptions-item label="用户名">
              {{ selectedUser?.username }}
            </a-descriptions-item>
            <a-descriptions-item label="邮箱">
              {{ selectedUser?.email }}
            </a-descriptions-item>
            <a-descriptions-item label="角色">
              <a-tag :color="selectedUser?.is_admin ? 'red' : 'blue'" size="small">
                {{ selectedUser?.is_admin ? '管理员' : '普通用户' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag :color="selectedUser?.banned ? 'gray' : 'green'" size="small">
                {{ selectedUser?.banned ? '已封禁' : '正常' }}
              </a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </div>

        <a-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" layout="vertical">
          <a-form-item field="new_password" label="新密码">
            <a-input-password
              v-model="passwordForm.new_password"
              placeholder="请输入新密码（至少6位字符）"
              size="large"
            />
          </a-form-item>
          <a-form-item field="confirm_password" label="确认新密码">
            <a-input-password
              v-model="passwordForm.confirm_password"
              placeholder="请再次输入新密码"
              size="large"
            />
          </a-form-item>
        </a-form>

        <a-alert type="warning" show-icon class="security-warning">
          <template #icon><icon-exclamation-circle /></template>
          <div>
            <div class="alert-title">安全提醒</div>
            <ul class="security-tips">
              <li>新密码将立即生效，用户需要使用新密码重新登录</li>
              <li>建议使用包含字母、数字和特殊字符的强密码</li>
              <li>请确认用户身份后再进行密码修改操作</li>
            </ul>
          </div>
        </a-alert>
      </div>
    </a-modal>

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
        <a-form-item field="is_admin" label="用户角色">
          <a-radio-group v-model="createForm.is_admin" direction="vertical">
            <a-radio :value="false">
              <template #radio="{ checked }">
                <a-space 
                  align="start" 
                  class="custom-radio-card" 
                  :class="{ 'custom-radio-card-checked': checked }"
                >
                  <div class="custom-radio-button" :class="{ 'custom-radio-button-checked': checked }">
                    <div class="custom-radio-inner" v-if="checked"></div>
                  </div>
                  <div class="custom-radio-content">
                    <div class="custom-radio-header">
                      <icon-user class="custom-radio-icon" :class="{ 'custom-radio-icon-checked': checked }" />
                      <div class="custom-radio-title">普通用户</div>
                    </div>
                    <div class="custom-radio-description">只能查看个人信息和修改密码，无管理权限</div>
                  </div>
                </a-space>
              </template>
            </a-radio>
            <a-radio :value="true">
              <template #radio="{ checked }">
                <a-space 
                  align="start" 
                  class="custom-radio-card" 
                  :class="{ 'custom-radio-card-checked': checked }"
                >
                  <div class="custom-radio-button" :class="{ 'custom-radio-button-checked': checked }">
                    <div class="custom-radio-inner" v-if="checked"></div>
                  </div>
                  <div class="custom-radio-content">
                    <div class="custom-radio-header">
                      <icon-user-group class="custom-radio-icon" :class="{ 'custom-radio-icon-checked': checked }" />
                      <div class="custom-radio-title">管理员</div>
                    </div>
                    <div class="custom-radio-description">拥有完整系统权限，可管理用户、查看日志等</div>
                  </div>
                </a-space>
              </template>
            </a-radio>
          </a-radio-group>
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
  adminChangePassword,
  type CreateUserRequest,
  type AdminChangePasswordRequest,
  type User
} from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconEdit, IconStop, IconDelete, IconUser, IconUserGroup, IconExclamationCircle } from '@arco-design/web-vue/es/icon'
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

// 修改密码
const showPasswordModal = ref(false)
const passwordLoading = ref(false)
const passwordFormRef = ref()
const selectedUser = ref<User | null>(null)
const passwordForm = reactive({
  new_password: '',
  confirm_password: ''
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

const passwordRules = {
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

// 显示修改密码模态框
const showChangePasswordModal = (user: User) => {
  selectedUser.value = user
  showPasswordModal.value = true
}

// 处理修改用户密码
const handleChangeUserPassword = async () => {
  if (!passwordFormRef.value || !selectedUser.value) return

  try {
    await passwordFormRef.value.validate()
  } catch (error) {
    return // 验证失败
  }

  passwordLoading.value = true
  try {
    const changePasswordData: AdminChangePasswordRequest = {
      new_password: passwordForm.new_password
    }

    await adminChangePassword(selectedUser.value.id, changePasswordData)
    Message.success(`用户 ${selectedUser.value.username} 的密码修改成功`)
    showPasswordModal.value = false
    resetPasswordForm()
  } catch (error) {
    console.error('Failed to change user password:', error)
  } finally {
    passwordLoading.value = false
  }
}

const resetPasswordForm = () => {
  Object.assign(passwordForm, {
    new_password: '',
    confirm_password: ''
  })
  selectedUser.value = null
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

/* 自定义单选框样式 */
.custom-radio-card {
  width: 100%;
  padding: 16px;
  border: 1px solid var(--color-border-2);
  border-radius: 6px;
  transition: all 0.2s ease;
  cursor: pointer;
  background-color: var(--color-fill-1);
}

.custom-radio-card:hover {
  border-color: var(--color-primary-light-4);
  background-color: var(--color-primary-light-6);
}

.custom-radio-card-checked {
  border-color: #165dff !important;
  background-color: var(--color-primary-light-5) !important;
}

.custom-radio-button {
  width: 16px;
  height: 16px;
  border: 2px solid var(--color-border-3);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
  margin-top: 2px;
}

.custom-radio-button-checked {
  border-color: #165dff;
  background-color: #165dff;
}

.custom-radio-inner {
  width: 8px;
  height: 8px;
  background-color: white;
  border-radius: 50%;
}

.custom-radio-content {
  flex: 1;
}

.custom-radio-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.custom-radio-icon {
  font-size: 16px;
  color: var(--color-text-3);
  transition: color 0.2s ease;
}

.custom-radio-icon-checked {
  color: #165dff;
}

.custom-radio-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  line-height: 20px;
}

.custom-radio-description {
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 16px;
  margin-left: 24px;
}

.custom-radio-card-checked .custom-radio-title {
  color: #165dff;
  font-weight: 600;
}

.custom-radio-card-checked .custom-radio-description {
  color: var(--color-text-2);
}

/* 覆盖原始单选框样式 */
:deep(.arco-radio) {
  display: block !important;
  margin-bottom: 12px;
}

:deep(.arco-radio:last-child) {
  margin-bottom: 0;
}

:deep(.arco-radio .arco-radio-button) {
  display: none !important;
}

:deep(.arco-radio-group-direction-vertical .arco-radio) {
  margin-right: 0 !important;
}

/* 密码修改模态框样式 */
.password-change-content {
  padding: 0;
}

.user-info-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.section-icon {
  color: var(--color-primary);
  font-size: 16px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
}

.user-descriptions :deep(.arco-descriptions-item-label) {
  color: var(--color-text-2);
  font-weight: 500;
  width: 80px;
  padding: 8px 0;
}

.user-descriptions :deep(.arco-descriptions-item-value) {
  color: var(--color-text-1);
  padding: 8px 0;
}

.security-warning {
  margin-top: 24px;
  border-radius: 8px;
}

.security-tips {
  margin: 8px 0 0 0;
  padding-left: 16px;
}

.security-tips li {
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 4px;
  color: var(--color-text-2);
}

/* 模态框响应式优化 */
.change-password-modal :deep(.arco-modal) {
  max-width: 90vw;
}

.change-password-modal :deep(.arco-modal-body) {
  padding: 24px 32px;
}

.change-password-modal :deep(.arco-form-item) {
  margin-bottom: 20px;
}

.change-password-modal :deep(.arco-input-password),
.change-password-modal :deep(.arco-input) {
  border-radius: 6px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .change-password-modal :deep(.arco-modal) {
    width: calc(100vw - 32px) !important;
    margin: 16px !important;
  }

  .change-password-modal :deep(.arco-modal-body) {
    padding: 20px 24px !important;
  }

  .user-details {
    gap: 8px;
  }

  .user-field {
    font-size: 12px;
  }

  .security-tips li {
    font-size: 12px;
  }
}
</style>