<template>
  <div class="chat-container">
    <div class="user-list">
      <div 
        v-for="user in users" 
        :key="user.id" 
        class="user-item" 
        :class="{ active: currentUserId === user.id }"
        @click="selectUser(user)"
      >
        <el-avatar :size="40" :src="user.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
        <div class="user-info">
          <div class="name">{{ user.username || '用户 ' + user.id }}</div>
          <div class="last-msg">点击开始聊天</div>
        </div>
      </div>
    </div>
    <div class="chat-main" v-if="currentUserId">
      <div class="chat-header">
        <span>正在与 {{ currentUserName }} 聊天</span>
        <el-button link type="primary" size="small" @click="viewUserNotifications">查看历史通知</el-button>
      </div>
      <div class="message-list" ref="msgListRef">
        <div 
          v-for="(msg, index) in currentMessages" 
          :key="index" 
          class="message-item"
          :class="{ 'my-message': msg.sender_type === 'admin' }"
        >
          <div class="content">{{ msg.content }}</div>
        </div>
      </div>
      <div class="input-area">
        <el-input v-model="inputText" placeholder="请输入消息..." @keyup.enter="sendMessage" />
        <el-button type="primary" @click="sendMessage">发送</el-button>
        <el-button type="warning" @click="showNotificationDialog">发送系统通知</el-button>
      </div>
    </div>
    <div class="empty-state" v-else>
      请选择一个用户开始聊天
    </div>

    <!-- 系统通知对话框 -->
    <el-dialog v-model="notificationDialogVisible" title="发送系统通知" width="500px">
      <el-form :model="notificationForm">
        <el-form-item label="标题">
          <el-input v-model="notificationForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input type="textarea" v-model="notificationForm.content" placeholder="请输入通知内容" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="notificationDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="sendNotification">发送</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 用户通知记录对话框 -->
    <el-dialog v-model="userHistoryDialogVisible" title="用户历史通知" width="600px">
      <el-table :data="userNotifications" style="width: 100%" max-height="400">
        <el-table-column prop="title" label="标题" width="150" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="160">
          <template #default="scope">
            {{ new Date(scope.row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="is_read" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.is_read ? 'success' : 'info'" size="small">
              {{ scope.row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import { useAuthStore } from '../../stores/auth'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const socket = ref(null)
const messages = ref([]) // 所有消息
const users = ref([]) // 活跃用户列表
const currentUserId = ref(null)
const inputText = ref('')
const msgListRef = ref(null)
const API_URL = 'http://localhost:8080/api'

// Notification Send
const notificationDialogVisible = ref(false)
const notificationForm = ref({ title: '', content: '' })

// User Notification History
const userHistoryDialogVisible = ref(false)
const userNotifications = ref([])

const showNotificationDialog = () => {
  notificationDialogVisible.value = true
}

const viewUserNotifications = async () => {
  if (!currentUserId.value) return
  try {
    const token = localStorage.getItem('admin_token')
    const response = await axios.get(`${API_URL}/notifications/admin/user/${currentUserId.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    userNotifications.value = response.data
    userHistoryDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取历史通知失败')
  }
}

const sendNotification = async () => {
  if (!notificationForm.value.title || !notificationForm.value.content) {
    ElMessage.warning('请填写完整信息')
    return
  }

  try {
    const token = localStorage.getItem('admin_token')
    await axios.post(`${API_URL}/chat/notification`, {
      user_id: currentUserId.value,
      title: notificationForm.value.title,
      content: notificationForm.value.content
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    ElMessage.success('系统通知发送成功')
    notificationDialogVisible.value = false
    notificationForm.value = { title: '', content: '' }
  } catch (error) {
    ElMessage.error('发送失败')
  }
}

// 获取最近联系人
const fetchChatUsers = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await axios.get(`${API_URL}/chat/users`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    // 转换后端数据格式 (增加 defensive check)
    const userList = response.data || []
    users.value = userList.map(u => ({
      id: u.ID,
      username: u.username,
      avatar: u.avatar
    }))
  } catch (error) {
    console.error('Failed to fetch chat users:', error)
  }
}

// 获取聊天记录
const fetchMessages = async (userId) => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await axios.get(`${API_URL}/chat/messages/${userId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    // 过滤掉已存在的历史消息 (简单去重)
    const historyMsgs = response.data || []
    const newMsgs = historyMsgs.map(m => ({
      sender_id: m.sender_id,
      sender_type: m.sender_type,
      receiver_id: m.receiver_id,
      content: m.content,
      type: m.type || 'text',
      id: m.ID
    }))
    
    // 将新消息合并到 messages 中
    newMsgs.forEach(newMsg => {
      if (!messages.value.find(m => m.id === newMsg.id)) {
        messages.value.push(newMsg)
      }
    })
    
    scrollToBottom()
  } catch (error) {
    console.error('Failed to fetch messages:', error)
  }
}

const currentUserName = computed(() => {
  const u = users.value.find(u => u.id === currentUserId.value)
  return u ? u.username : '未知用户'
})

const currentMessages = computed(() => {
  return messages.value.filter(m => 
    (m.sender_id === currentUserId.value && m.sender_type === 'user') ||
    (m.receiver_id === currentUserId.value && m.sender_type === 'admin')
  )
})

const connectWebSocket = () => {
  const adminId = authStore.admin.ID
  // 连接 WS
  socket.value = new WebSocket(`ws://localhost:8080/api/ws?user_id=${adminId}&type=admin`)
  
  socket.value.onopen = () => {
    console.log('WS 已连接')
  }
  
  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    
    // 后端返回的格式可能是直接的消息对象，或者是 { type, payload }
    // 根据 Hub.go 的实现，广播出来的是 ChatMessage 结构体，但也可能被包装
    // 假设后端 Hub.go 广播的是 raw ChatMessage JSON
    
    // 如果是标准 WSMessage 格式
    let msg = data
    if (data.type === 'message' && data.payload) {
      msg = data.payload
    }

    // 确保是消息对象
    if (msg.content) {
      messages.value.push(msg)
      
      // 如果是新用户的消息，添加到用户列表
      if (msg.sender_type === 'user') {
        const existingUser = users.value.find(u => u.id === msg.sender_id)
        if (!existingUser) {
          users.value.push({ 
            id: msg.sender_id, 
            username: '用户 ' + msg.sender_id,
            avatar: ''
          })
          
          // 如果当前没有选中用户，自动选中这个新用户
          if (!currentUserId.value) {
             selectUser({ id: msg.sender_id, username: '用户 ' + msg.sender_id })
          }
        } else if (currentUserId.value === msg.sender_id) {
          // 如果已经在聊天中，不需要做额外操作，scrollToBottom 已经有了
        }
      }
      
      scrollToBottom()
    }
  }
}

const selectUser = (user) => {
  currentUserId.value = user.id
  fetchMessages(user.id) // 加载历史记录
  scrollToBottom()
}

const sendMessage = () => {
  if (!inputText.value.trim() || !currentUserId.value) return
  
  const msg = {
    sender_id: authStore.admin.ID,
    sender_type: 'admin',
    receiver_id: currentUserId.value, // 发给选中的用户
    content: inputText.value,
    type: 'text'
  }
  
  socket.value.send(JSON.stringify({
    type: 'message',
    payload: msg
  }))
  
  // 本地立即显示 (后端也会广播回来，去重逻辑暂略，简单起见依赖广播或本地push)
  // 这里依赖后端广播回来自己
  inputText.value = ''
}

const scrollToBottom = () => {
  nextTick(() => {
    if (msgListRef.value) {
      msgListRef.value.scrollTop = msgListRef.value.scrollHeight
    }
  })
}

onMounted(() => {
  connectWebSocket()
  fetchChatUsers()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  height: calc(100vh - 100px);
  border: 1px solid #e6e6e6;
  background: #fff;
  border-radius: 4px;
  overflow: hidden;
}
.user-list {
  width: 250px;
  border-right: 1px solid #e6e6e6;
  overflow-y: auto;
}
.user-item {
  padding: 15px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-bottom: 1px solid #f5f5f5;
}
.user-item:hover, .user-item.active {
  background-color: #f0f9eb;
}
.user-info {
  margin-left: 10px;
}
.name {
  font-weight: bold;
}
.last-msg {
  font-size: 12px;
  color: #999;
}
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.chat-header {
  padding: 15px;
  border-bottom: 1px solid #e6e6e6;
  font-weight: bold;
  background-color: #fafafa;
}
.message-list {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #f5f7fa;
}
.message-item {
  margin-bottom: 10px;
  display: flex;
}
.message-item.my-message {
  justify-content: flex-end;
}
.content {
  background: #fff;
  padding: 10px 15px;
  border-radius: 4px;
  max-width: 70%;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}
.my-message .content {
  background: #95ec69;
}
.input-area {
  padding: 15px;
  border-top: 1px solid #e6e6e6;
  display: flex;
  gap: 10px;
  background-color: #fff;
}
.empty-state {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #999;
}
</style>
