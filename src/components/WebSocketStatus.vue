<template>
  <div v-if="showIndicator" class="ws-indicator" :class="statusClass" @click="toggleDetails">
    <span class="ws-dot"></span>
    <span class="ws-text">{{ statusText }}</span>

    <div v-if="showDetails" class="ws-details">
      <div class="details-header">
        <span>WebSocket 状态</span>
        <button class="close-btn" @click.stop="showDetails = false">✕</button>
      </div>
      <div class="details-body">
        <div class="status-row">
          <span class="label">连接状态:</span>
          <span class="value" :class="{ connected: isConnected }">
            {{ isConnected ? '已连接' : '未连接' }}
          </span>
        </div>
        <div class="status-row">
          <span class="label">最后消息:</span>
          <span class="value">{{ lastMessageTime || '无' }}</span>
        </div>
        <div class="status-row">
          <span class="label">重连次数:</span>
          <span class="value">{{ reconnectCount }}</span>
        </div>
      </div>
      <div class="details-footer">
        <button class="reconnect-btn" @click="handleReconnect" :disabled="isConnected">
          {{ isConnected ? '已连接' : '重新连接' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useWebSocket } from '../composables/useWebSocket'

const { isConnected, lastMessage, connect, disconnect } = useWebSocket()

const showIndicator = ref(true)
const showDetails = ref(false)
const reconnectCount = ref(0)
const lastMessageTime = ref(null)

const statusClass = computed(() => {
  return isConnected.value ? 'connected' : 'disconnected'
})

const statusText = computed(() => {
  return isConnected.value ? '在线' : '离线'
})

const toggleDetails = () => {
  showDetails.value = !showDetails.value
}

const handleReconnect = () => {
  if (!isConnected.value) {
    reconnectCount.value++
    connect()
  }
}

let unsubscribe = null

onMounted(() => {
  connect()

  unsubscribe = require('../composables/useWebSocket').useWebSocket().onMessage((message) => {
    if (message) {
      lastMessageTime.value = new Date().toLocaleTimeString()
    }
  })
})

onUnmounted(() => {
  if (unsubscribe) {
    unsubscribe()
  }
})
</script>

<style scoped>
.ws-indicator {
  position: fixed;
  bottom: 20px;
  left: 20px;
  z-index: 9998;
  background: white;
  border-radius: 20px;
  padding: 8px 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  transition: all 0.2s;
}

.ws-indicator:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.ws-indicator.connected {
  border: 2px solid #10b981;
}

.ws-indicator.disconnected {
  border: 2px solid #ef4444;
}

.ws-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.ws-indicator.connected .ws-dot {
  background: #10b981;
}

.ws-indicator.disconnected .ws-dot {
  background: #ef4444;
  animation: none;
}

.ws-text {
  font-weight: 500;
}

.ws-indicator.connected .ws-text {
  color: #10b981;
}

.ws-indicator.disconnected .ws-text {
  color: #ef4444;
}

.ws-details {
  position: absolute;
  bottom: 100%;
  left: 0;
  margin-bottom: 12px;
  background: white;
  border-radius: 12px;
  width: 280px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.details-header {
  padding: 12px 16px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.details-header .close-btn {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  cursor: pointer;
  font-size: 12px;
}

.details-body {
  padding: 16px;
}

.status-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.status-row .label {
  color: #6b7280;
}

.status-row .value {
  color: #374151;
  font-weight: 500;
}

.status-row .value.connected {
  color: #10b981;
}

.details-footer {
  padding: 12px 16px;
  border-top: 1px solid #e5e7eb;
}

.reconnect-btn {
  width: 100%;
  padding: 10px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.reconnect-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.reconnect-btn:disabled {
  background: #d1d5db;
  cursor: not-allowed;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
