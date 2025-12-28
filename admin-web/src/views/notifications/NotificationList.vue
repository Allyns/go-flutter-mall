<template>
  <div class="notification-list">
    <el-card>
      <div class="header-actions">
        <h2>系统通知记录</h2>
      </div>
      <el-table :data="notifications" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="标题" width="200" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="接收用户" width="150">
          <template #default="scope">
            {{ scope.row.User ? scope.row.User.username : ('User ' + scope.row.user_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发送时间" width="180">
          <template #default="scope">
            {{ new Date(scope.row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="is_read" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_read ? 'success' : 'info'">
              {{ scope.row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const notifications = ref([])
const loading = ref(false)
const API_URL = 'http://localhost:8080/api'

const fetchNotifications = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await axios.get(`${API_URL}/notifications/admin/all`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    notifications.value = response.data
  } catch (error) {
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped>
.header-actions {
  margin-bottom: 20px;
}
</style>
