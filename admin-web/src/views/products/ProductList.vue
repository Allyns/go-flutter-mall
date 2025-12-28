<template>
  <div class="product-list">
    <el-card>
      <div class="header-actions">
        <h2>商品管理</h2>
        <el-button type="primary" @click="openDialog()">添加商品</el-button>
      </div>

      <el-table :data="products" style="width: 100%" v-loading="loading">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column label="图片" width="100">
          <template #default="scope">
            <el-image 
              style="width: 50px; height: 50px" 
              :src="scope.row.cover_image" 
              fit="cover"
              :preview-src-list="[scope.row.cover_image]"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="scope">¥{{ scope.row.price }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="120" />
        <el-table-column label="操作" width="220">
          <template #default="scope">
            <el-button-group>
              <el-button size="small" @click="openDialog(scope.row)">编辑</el-button>
              <el-button size="small" type="primary" @click="handleViewReviews(scope.row)">评价</el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑商品' : '添加商品'">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="form.price" :precision="2" :step="0.1" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="封面图片链接">
          <el-input v-model="form.cover_image" />
        </el-form-item>
        <el-form-item label="分类 ID">
          <el-input-number v-model="form.category_id" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Reviews Dialog -->
    <el-dialog v-model="reviewsDialogVisible" title="商品评价" width="60%" append-to-body>
      <el-table :data="currentReviews" style="width: 100%" v-loading="reviewsLoading">
        <el-table-column prop="user.username" label="用户" width="120">
          <template #default="scope">
            {{ scope.row.user?.username || '匿名用户' }}
          </template>
        </el-table-column>
        <el-table-column prop="rating" label="评分" width="150">
           <template #default="scope">
             <el-rate v-model="scope.row.rating" disabled show-score text-color="#ff9900" />
           </template>
        </el-table-column>
        <el-table-column prop="content" label="评价内容" />
        <el-table-column prop="CreatedAt" label="时间" width="180">
            <template #default="scope">
                {{ new Date(scope.row.CreatedAt).toLocaleString() }}
            </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const API_URL = 'http://localhost:8080/api'
const products = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)

// Reviews state
const reviewsDialogVisible = ref(false)
const currentReviews = ref([])
const reviewsLoading = ref(false)

const form = ref({
  name: '',
  description: '',
  price: 0,
  stock: 0,
  cover_image: '',
  category_id: 1,
  status: 1
})

const fetchProducts = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${API_URL}/products`)
    products.value = response.data
  } catch (error) {
    ElMessage.error('获取商品失败')
  } finally {
    loading.value = false
  }
}

const openDialog = (product = null) => {
  if (product) {
    editingId.value = product.ID
    form.value = { ...product }
  } else {
    editingId.value = null
    form.value = {
      name: '',
      description: '',
      price: 0,
      stock: 0,
      cover_image: '',
      category_id: 1,
      status: 1
    }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const config = { headers: { Authorization: `Bearer ${token}` } }
    
    if (editingId.value) {
      // 更新
      await axios.put(`${API_URL}/products/${editingId.value}`, form.value, config)
      ElMessage.success('更新成功')
    } else {
      // 创建
      await axios.post(`${API_URL}/products`, form.value, config)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    fetchProducts()
  } catch (error) {
    ElMessage.error('操作失败: ' + (error.response?.data?.error || '未知错误'))
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该商品吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      try {
        const token = localStorage.getItem('admin_token')
        const config = { headers: { Authorization: `Bearer ${token}` } }
        await axios.delete(`${API_URL}/products/${row.ID}`, config)
        ElMessage.success('删除成功')
        fetchProducts()
      } catch (error) {
        ElMessage.error('删除失败')
      }
    })
}

const handleViewReviews = async (row) => {
  reviewsDialogVisible.value = true
  reviewsLoading.value = true
  currentReviews.value = []
  try {
    const response = await axios.get(`${API_URL}/products/${row.ID}/reviews`)
    currentReviews.value = response.data
  } catch (error) {
    ElMessage.error('获取评价失败')
  } finally {
    reviewsLoading.value = false
  }
}

onMounted(() => {
  fetchProducts()
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
