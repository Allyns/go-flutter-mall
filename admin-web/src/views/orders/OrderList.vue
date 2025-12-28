<template>
  <div class="order-list">
    <el-card>
      <div class="header-actions">
        <h2>订单管理</h2>
      </div>
      <el-table :data="orders" style="width: 100%" v-loading="loading">
        <el-table-column prop="order_no" label="订单编号" min-width="220" />
        <el-table-column label="用户" min-width="120">
          <template #default="scope">
             {{ scope.row.user?.username || '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="金额" width="120">
          <template #default="scope">¥{{ scope.row.total_amount }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ new Date(scope.row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template #default="scope">
            <el-button-group>
              <el-button size="small" type="primary" @click="openStatusDialog(scope.row)">调整状态</el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 状态调整对话框 -->
    <el-dialog v-model="statusDialogVisible" title="调整订单状态" width="300px">
      <el-select v-model="currentStatus" placeholder="请选择状态">
        <el-option label="待付款" :value="0" />
        <el-option label="待发货" :value="1" />
        <el-option label="待收货" :value="2" />
        <el-option label="待评价" :value="3" />
        <el-option label="已完成" :value="4" />
        <el-option label="售后中" :value="5" />
        <el-option label="已取消" :value="-1" />
      </el-select>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="statusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmUpdateStatus">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const orders = ref([])
const loading = ref(false)
const API_URL = 'http://localhost:8080/api'

// Status Dialog
const statusDialogVisible = ref(false)
const currentOrder = ref(null)
const currentStatus = ref(0)

const openStatusDialog = (row) => {
  currentOrder.value = row
  currentStatus.value = row.status
  statusDialogVisible.value = true
}

const confirmUpdateStatus = async () => {
  if (!currentOrder.value) return
  await updateStatus(currentOrder.value, currentStatus.value)
  statusDialogVisible.value = false
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    // 使用新的管理员接口获取所有订单
    const response = await axios.get(`${API_URL}/orders/admin/all`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    orders.value = response.data
  } catch (error) {
    ElMessage.error('获取订单失败')
  } finally {
    loading.value = false
  }
}

const getStatusText = (status) => {
  const map = { 0: '待付款', 1: '待发货', 2: '待收货', 3: '待评价', 4: '已完成', 5: '售后中', '-1': '已取消' }
  return map[status] || '未知'
}

const getStatusType = (status) => {
  const map = { 0: 'warning', 1: 'primary', 2: 'success', 3: 'warning', 4: 'success', 5: 'danger', '-1': 'info' }
  return map[status] || 'info'
}

const updateStatus = async (row, status) => {
  try {
    const token = localStorage.getItem('admin_token')
    await axios.put(`${API_URL}/orders/${row.ID}/status`, { status }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success('状态更新成功')
    fetchOrders()
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该订单吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      try {
        const token = localStorage.getItem('admin_token')
        await axios.delete(`${API_URL}/orders/${row.ID}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        ElMessage.success('订单删除成功')
        fetchOrders()
      } catch (error) {
        ElMessage.error('删除失败')
      }
    })
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
