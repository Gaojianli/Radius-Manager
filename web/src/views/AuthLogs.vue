<template>
  <a-card class="page-container" :bordered="false">
    <a-page-header title="认证日志" subtitle="查看系统认证日志记录" :show-back="false" />
    

    <!-- 日志表格 -->
    <a-table
      :data="logs"
      :columns="columns"
      :pagination="pagination"
      :loading="loading"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
      row-key="id"
    >
      <template #username-filter="{ filterValue, setFilterValue, handleFilterConfirm, handleFilterReset }">
        <a-card>
          <a-space direction="vertical">
            <a-input
              :model-value="filterValue?.[0] || ''"
              @input="(value) => setFilterValue(value ? [value] : [])"
              placeholder="输入用户名"
              allow-clear
            />
            <div class="custom-filter-footer">
              <a-button @click="handleFilterReset" size="mini">
                重置
              </a-button>
              <a-button type="primary" @click="handleFilterConfirm" size="mini">
                确认
              </a-button>
            </div>
          </a-space>
        </a-card>
      </template>
      <template #auth-type="{ record }">
        <a-tag :color="record.auth_type === 'authorize' ? 'blue' : 'green'">
          {{ record.auth_type === 'authorize' ? '授权' : '认证' }}
        </a-tag>
      </template>
      <template #success="{ record }">
        <a-tag :color="record.success ? 'green' : 'red'">
          {{ record.success ? '成功' : '失败' }}
        </a-tag>
      </template>
      <template #user-agent="{ record }">
        <div class="user-agent-cell" :title="record.user_agent">
          {{ record.user_agent }}
        </div>
      </template>
      <template #created-at="{ record }">
        {{ formatDate(record.created_at) }}
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { getAuthLogs, type AuthLog } from '@/services/api'
import { Message } from '@arco-design/web-vue'
import { onMounted, reactive, ref } from 'vue'


// 表格数据
const logs = ref<AuthLog[]>([])
const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100']
})

// 表格列配置
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80
  },
  {
    title: '用户名',
    dataIndex: 'username',
    width: 120,
    filterable: {
      filter: (value: string[], record: AuthLog) => {
        if (!value || !value.length) return true
        return record.username.toLowerCase().includes(value[0].toLowerCase())
      },
      slotName: 'username-filter'
    }
  },
  {
    title: '认证类型',
    dataIndex: 'auth_type',
    width: 120,
    slotName: 'auth-type'
  },
  {
    title: '结果',
    dataIndex: 'success',
    width: 100,
    slotName: 'success'
  },
  {
    title: 'IP地址',
    dataIndex: 'ip_address',
    width: 150
  },
  {
    title: 'User-Agent',
    dataIndex: 'user_agent',
    slotName: 'user-agent'
  },
  {
    title: '时间',
    dataIndex: 'created_at',
    width: 180,
    slotName: 'created-at'
  }
]

// 加载日志数据
const loadLogs = async () => {
  loading.value = true
  try {
    const response = await getAuthLogs(
      pagination.current, 
      pagination.pageSize
    )
    logs.value = response.data.logs
    pagination.total = response.data.pagination.total
  } catch (error) {
    console.error('Failed to load auth logs:', error)
    Message.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

// 分页处理
const handlePageChange = (page: number) => {
  pagination.current = page
  loadLogs()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadLogs()
}


// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  loadLogs()
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

.user-agent-cell {
  max-width: 300px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}


.custom-filter-footer {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .page-container {
    margin: 16px auto;
    padding: 16px;
  }
  
  /* 在手机端隐藏一些不重要的列 */
  :deep(.arco-table-th):nth-child(1),
  :deep(.arco-table-td):nth-child(1) {
    display: none; /* 隐藏ID列 */
  }
  
  :deep(.arco-table-th):nth-child(5),
  :deep(.arco-table-td):nth-child(5) {
    display: none; /* 隐藏IP地址列 */
  }
  
  :deep(.arco-table-th):nth-child(6),
  :deep(.arco-table-td):nth-child(6) {
    display: none; /* 隐藏User-Agent列 */
  }
  
  /* 调整表格列宽 */
  :deep(.arco-table-th):nth-child(2),
  :deep(.arco-table-td):nth-child(2) {
    min-width: 100px; /* 用户名列 */
  }
  
  :deep(.arco-table-th):nth-child(7),
  :deep(.arco-table-td):nth-child(7) {
    min-width: 140px; /* 时间列 */
  }
  
  /* 表格横向滚动 */
  :deep(.arco-table-container) {
    overflow-x: auto;
  }
  
  .user-agent-cell {
    max-width: 150px;
  }
}
</style>