<template>
  <div class="home-container">
    <el-container>
      <el-aside width="220px">
        <div class="sidebar-logo">
          <img src="https://element-plus.org/images/element-plus-logo.svg" alt="Logo" class="logo-img" />
          <span class="logo-text">Go Mall Admin</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical-demo"
          background-color="#001529"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          :router="true"
        >
          <el-menu-item index="/">
            <el-icon><DataLine /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/products">
            <el-icon><Goods /></el-icon>
            <span>商品管理</span>
          </el-menu-item>
          <el-menu-item index="/orders">
            <el-icon><List /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="/chat">
            <el-icon><ChatDotRound /></el-icon>
            <span>客服消息</span>
          </el-menu-item>
          <el-menu-item index="/notifications">
            <el-icon><Bell /></el-icon>
            <span>系统通知</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-left">
            <!-- Breadcrumb could go here -->
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="el-dropdown-link">
                欢迎, {{ authStore.admin?.username }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { DataLine, Goods, List, ChatDotRound, ArrowDown, Bell } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const activeMenu = computed(() => route.path)

const handleCommand = (command) => {
  if (command === 'logout') {
    authStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.home-container {
  height: 100vh;
}
.el-container {
  height: 100%;
}
.el-aside {
  background-color: #001529;
  border-right: none;
  box-shadow: 2px 0 6px rgba(0,21,41,.35);
  z-index: 10;
}
.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding-left: 20px;
  background-color: #002140;
  overflow: hidden;
}
.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}
.logo-text {
  color: #fff;
  font-weight: 600;
  font-size: 18px;
  white-space: nowrap;
}
.el-menu {
  border-right: none;
}
.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  z-index: 9;
}
.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #333;
}
.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
