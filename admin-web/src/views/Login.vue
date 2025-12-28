<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <h2 class="title">Go Mall 后台管理系统</h2>
        <p class="subtitle">专业的电商管理平台</p>
      </div>
      <el-card class="login-card">
        <el-form :model="form" @submit.prevent="handleLogin" size="large">
          <el-form-item>
            <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.password" type="password" placeholder="密码" :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" class="login-btn" @click="handleLogin">
              登 录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: ''
})

const loading = ref(false)

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    await authStore.login(form.value.username, form.value.password)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    ElMessage.error('登录失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #2d3a4b;
  background-image: linear-gradient(135deg, #001529 0%, #001529 100%);
}

.login-content {
  width: 400px;
  padding: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.title {
  color: #fff;
  font-size: 26px;
  font-weight: bold;
  margin: 0;
}

.subtitle {
  color: #eee;
  margin-top: 10px;
  font-size: 14px;
}

.login-card {
  border-radius: 8px;
}

.login-btn {
  width: 100%;
  font-size: 16px;
  padding: 12px 0;
}
</style>
