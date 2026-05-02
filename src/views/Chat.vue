<template>
  <div class="chat-page">
    <div class="chat-container">
      <div class="chat-sidebar">
        <div class="sidebar-header">
          <div class="header-left">
            <button class="back-btn" @click="goBack">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="19" y1="12" x2="5" y2="12"/>
                <polyline points="12 19 5 12 12 5"/>
              </svg>
            </button>
            <h2>消息中心</h2>
          </div>
          <button class="new-chat-btn" @click="showNewChatDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </button>
        </div>
        <div class="conversations-list">
          <div v-for="conv in conversations" :key="conv.user_id"
               class="conversation-item"
               :class="{ active: selectedUserId === conv.user_id, unread: conv.unread_count > 0 }"
               @click="selectConversation(conv)">
            <div class="conv-avatar">
              <img v-if="conv.avatar" :src="conv.avatar" :alt="conv.real_name" v-avatar />
              <span v-else>{{ conv.real_name?.charAt(0) || 'U' }}</span>
            </div>
            <div class="conv-info">
              <div class="conv-name">{{ conv.real_name }}</div>
              <div class="conv-last-msg">{{ conv.last_message || '暂无消息' }}</div>
            </div>
            <div class="conv-meta">
              <div class="conv-time">{{ formatRelativeTime(conv.last_message_time) }}</div>
              <div v-if="conv.unread_count > 0" class="conv-badge">{{ conv.unread_count }}</div>
            </div>
          </div>
          <div v-if="conversations.length === 0" class="empty-tip">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <p>暂无会话</p>
            <p class="hint">点击右上角开始新对话</p>
          </div>
        </div>
      </div>

      <div class="chat-main">
        <div v-if="!selectedUserId" class="chat-empty">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <h3>选择一个对话开始聊天</h3>
          <p>从左侧列表选择或创建新对话</p>
        </div>

        <div v-else class="chat-content">
          <div class="chat-header">
            <div class="header-info">
              <div class="header-avatar">
                <img v-if="selectedUser?.avatar" :src="selectedUser?.avatar" :alt="selectedUser?.real_name" v-avatar />
                <span v-else>{{ selectedUser?.real_name?.charAt(0) || 'U' }}</span>
              </div>
              <div class="header-text">
                <div class="header-name">{{ selectedUser?.real_name }}</div>
                <div class="header-username">@{{ selectedUser?.username }}</div>
              </div>
            </div>
          </div>

          <div class="messages-area" ref="messagesArea">
            <div v-for="(msg, index) in messages" :key="msg.id"
                 class="message-item"
                 :class="{ 'self': msg.from_user_id === currentUserId }">
              <div class="message-avatar">
                <img v-if="!isSelf(msg) && selectedUser?.avatar" :src="selectedUser?.avatar" :alt="selectedUser?.real_name" v-avatar />
                <span v-else-if="!isSelf(msg)">{{ selectedUser?.real_name?.charAt(0) || 'U' }}</span>
              </div>
              <div class="message-bubble">
                <div class="message-content">{{ msg.content }}</div>
                <div class="message-time">{{ formatTime(msg.create_time) }}</div>
              </div>
            </div>
          </div>

          <div class="chat-input-area">
            <div class="input-wrapper">
              <textarea v-model="newMessage"
                        placeholder="输入消息..."
                        @keydown.enter.prevent="sendMessage"
                        rows="1"
                        ref="messageInput"></textarea>
            </div>
            <button class="send-btn" @click="sendMessage" :disabled="!newMessage.trim()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="22" y1="2" x2="11" y2="13"/>
                <polygon points="22 2 15 22 11 13 2 9 22 2"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showNewChat" class="modal-overlay" @click="showNewChat = false">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>新对话</h3>
          <button class="close-btn" @click="showNewChat = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input type="text" v-model="searchUserQuery" placeholder="搜索用户..." />
          </div>
          <div class="users-list">
            <div v-for="user in filteredUsers" :key="user.id"
                 class="user-item"
                 @click="startNewChat(user)">
              <div class="user-avatar">
                <img v-if="user.avatar" :src="user.avatar" :alt="user.real_name" v-avatar />
                <span v-else>{{ user.real_name?.charAt(0) || 'U' }}</span>
              </div>
              <div class="user-info">
                <div class="user-name">{{ user.real_name }}</div>
                <div class="user-username">@{{ user.username }}</div>
              </div>
            </div>
            <div v-if="filteredUsers.length === 0" class="empty-tip">
              <p>未找到用户</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { getConversationsAPI, getMessagesAPI, sendMessageAPI, getPublicUsersAPI } from '../api/index'

const { success: notifySuccess, error: notifyError } = useNotification()

const router = useRouter()
const userStore = useUserStore()

const conversations = ref([])
const messages = ref([])
const selectedUserId = ref(null)
const selectedUser = ref(null)
const newMessage = ref('')
const showNewChat = ref(false)
const searchUserQuery = ref('')
const allUsers = ref([])
const messagesArea = ref(null)
const messageInput = ref(null)

const currentUserId = computed(() => userStore.id || parseInt(localStorage.getItem('userId') || '0'))

const filteredUsers = computed(() => {
  if (!searchUserQuery.value) return allUsers.value
  const query = searchUserQuery.value.toLowerCase()
  return allUsers.value.filter(u => 
    u.real_name?.toLowerCase().includes(query) || 
    u.username?.toLowerCase().includes(query)
  )
})

const formatRelativeTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'
  return date.toLocaleDateString()
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const isSelf = (msg) => msg.from_user_id === currentUserId.value

const loadConversations = async () => {
  try {
    const res = await getConversationsAPI()
    if (res.code === 200) {
      conversations.value = res.data || []
    }
  } catch (error) {
    console.error('加载会话失败', error)
  }
}

const loadMessages = async (userId) => {
  try {
    const res = await getMessagesAPI(userId)
    if (res.code === 200) {
      messages.value = res.data || []
      await nextTick()
      scrollToBottom()
    }
  } catch (error) {
    console.error('加载消息失败', error)
  }
}

const selectConversation = async (conv) => {
  selectedUserId.value = conv.user_id
  selectedUser.value = conv
  await loadMessages(conv.user_id)
  await loadConversations()
}

const sendMessage = async () => {
  if (!newMessage.value.trim() || !selectedUserId.value) return
  
  try {
    const res = await sendMessageAPI({
      to_user_id: selectedUserId.value,
      content: newMessage.value.trim()
    })
    if (res.code === 200) {
      messages.value.push(res.data)
      newMessage.value = ''
      await nextTick()
      scrollToBottom()
      await loadConversations()
    } else {
      notifyError(res.message || '发送失败')
    }
  } catch (error) {
    notifyError('发送失败')
    console.error(error)
  }
}

const scrollToBottom = () => {
  if (messagesArea.value) {
    messagesArea.value.scrollTop = messagesArea.value.scrollHeight
  }
}

const showNewChatDialog = async () => {
  showNewChat.value = true
  searchUserQuery.value = ''
  try {
    const res = await getPublicUsersAPI()
    if (res.code === 200) {
      allUsers.value = (res.data || []).filter(u => u.id !== currentUserId.value)
    }
  } catch (error) {
    console.error('加载用户失败', error)
  }
}

const startNewChat = (user) => {
  showNewChat.value = false
  selectedUserId.value = user.id
  selectedUser.value = user
  messages.value = []
  
  const existingConv = conversations.value.find(c => c.user_id === user.id)
  if (existingConv) {
    loadMessages(user.id)
  }
}

const goBack = () => {
  router.push('/team')
}

onMounted(() => {
  loadConversations()
})
</script>

<style scoped>
.chat-page {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.chat-container {
  display: flex;
  width: 100%;
  max-width: 1400px;
  height: calc(100vh - 40px);
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.chat-sidebar {
  width: 360px;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  background: #f9fafb;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: #f3f4f6;
  color: #4b5563;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: #e5e7eb;
  color: #1f2937;
  transform: translateX(-2px);
}

.sidebar-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.new-chat-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, box-shadow 0.2s;
}

.new-chat-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.new-chat-btn svg {
  width: 20px;
  height: 20px;
}

.conversations-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.conversation-item {
  display: flex;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
  margin-bottom: 4px;
}

.conversation-item:hover {
  background: #e5e7eb;
}

.conversation-item.active {
  background: #667eea;
  color: white;
}

.conversation-item.unread:not(.active) {
  background: #dbeafe;
}

.conv-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
  font-weight: 600;
  flex-shrink: 0;
  overflow: hidden;
}

.conv-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conv-info {
  flex: 1;
  margin-left: 12px;
  min-width: 0;
}

.conv-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.conv-last-msg {
  font-size: 13px;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conv-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.conv-time {
  font-size: 12px;
  opacity: 0.6;
}

.conv-badge {
  background: #ef4444;
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  min-width: 20px;
  text-align: center;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.chat-empty svg {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
  opacity: 0.3;
}

.chat-empty h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #4b5563;
}

.chat-empty p {
  margin: 0;
  font-size: 14px;
}

.chat-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
  background: white;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
  overflow: hidden;
}

.header-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.header-username {
  font-size: 13px;
  color: #6b7280;
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #f9fafb;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
  align-items: flex-end;
}

.message-item.self {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
  overflow: hidden;
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-bubble {
  max-width: 60%;
  margin: 0 12px;
}

.message-item:not(.self) .message-bubble {
  margin-left: 12px;
  margin-right: 0;
}

.message-item.self .message-bubble {
  margin-right: 12px;
  margin-left: 0;
}

.message-content {
  padding: 12px 16px;
  border-radius: 16px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.message-item.self .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.message-time {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 4px;
  padding: 0 4px;
}

.message-item.self .message-time {
  text-align: right;
}

.chat-input-area {
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  background: white;
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-wrapper {
  flex: 1;
  background: #f3f4f6;
  border-radius: 24px;
  padding: 8px 16px;
}

.input-wrapper textarea {
  width: 100%;
  border: none;
  background: transparent;
  resize: none;
  font-size: 14px;
  line-height: 1.5;
  outline: none;
  max-height: 120px;
}

.send-btn {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, box-shadow 0.2s;
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-btn svg {
  width: 20px;
  height: 20px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #f3f4f6;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #e5e7eb;
}

.close-btn svg {
  width: 18px;
  height: 18px;
}

.modal-body {
  padding: 20px 24px;
  overflow-y: auto;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #f3f4f6;
  padding: 12px 16px;
  border-radius: 12px;
  margin-bottom: 16px;
}

.search-box svg {
  width: 20px;
  height: 20px;
  color: #9ca3af;
}

.search-box input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  outline: none;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.user-item:hover {
  background: #f3f4f6;
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
  overflow: hidden;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: 600;
  font-size: 14px;
  color: #1f2937;
}

.user-username {
  font-size: 13px;
  color: #6b7280;
}

.empty-tip {
  text-align: center;
  color: #9ca3af;
  padding: 40px 20px;
}

.empty-tip svg {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  opacity: 0.3;
}

.empty-tip p {
  margin: 4px 0;
  font-size: 14px;
}

.empty-tip .hint {
  font-size: 13px;
  opacity: 0.7;
}
</style>
