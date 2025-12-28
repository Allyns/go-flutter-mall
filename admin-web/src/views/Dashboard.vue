<template>
  <div class="dashboard-welcome">
    <h2>欢迎使用 Go Mall 后台管理系统</h2>
    <p>请从左侧菜单选择模块开始管理您的店铺。</p>
    
    <el-row :gutter="20" class="mt-4">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>用户总数</template>
          <div class="stat-value">{{ stats.total_users }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>商品总数</template>
          <div class="stat-value">{{ stats.total_products }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>订单总数</template>
          <div class="stat-value">{{ stats.total_orders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>销售总额</template>
          <div class="stat-value">¥ {{ stats.total_sales.toFixed(2) }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const stats = ref({
  total_users: 0,
  total_products: 0,
  total_orders: 0,
  total_sales: 0.0
})

const fetchStats = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await axios.get('http://localhost:8080/api/auth/admin/stats', {
      headers: { Authorization: `Bearer ${token}` }
    })
    stats.value = response.data
  } catch (error) {
    ElMessage.error('获取统计数据失败')
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.dashboard-welcome {
  /* padding: 20px; handled by el-main */
}
.mt-4 {
  margin-top: 20px;
}
.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}
</style>
