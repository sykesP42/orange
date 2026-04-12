<template>
  <Transition name="online-status">
    <div v-if="!isOnline" class="online-status-bar offline">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="1" y1="1" x2="23" y2="23"/>
        <path d="M16.72 11.06A10.94 10.94 0 0 1 19 12.55"/>
        <path d="M5 12.55a10.94 10.94 0 0 1 5.17-2.39"/>
        <path d="M10.71 5.05A16 16 0 0 1 22.56 9"/>
        <path d="M1.42 9a15.91 15.91 0 0 1 4.7-2.88"/>
        <path d="M8.53 16.11a6 6 0 0 1 6.95 0"/>
        <line x1="12" y1="20" x2="12.01" y2="20"/>
      </svg>
      <span>网络连接已断开，部分功能可能不可用</span>
    </div>
    <div v-else-if="wasOffline" class="online-status-bar online">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M5 12.55a11 11 0 0 1 14.08 0"/>
        <path d="M1.42 9a16 16 0 0 1 21.16 0"/>
        <path d="M8.53 16.11a6 6 0 0 1 6.95 0"/>
        <line x1="12" y1="20" x2="12.01" y2="20"/>
      </svg>
      <span>网络已恢复</span>
    </div>
  </Transition>
</template>

<script setup>
import { watch } from 'vue'
import { useOnlineStatus } from '../composables/useOnlineStatus'

const { isOnline, wasOffline } = useOnlineStatus()

watch(wasOffline, (val) => {
  if (val && isOnline.value) {
    setTimeout(() => {
      wasOffline.value = false
    }, 3000)
  }
})
</script>

<style scoped>
.online-status-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 10001;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.online-status-bar svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.online-status-bar.offline {
  background: rgba(239, 68, 68, 0.9);
  color: white;
}

.online-status-bar.online {
  background: rgba(34, 197, 94, 0.9);
  color: white;
}

.online-status-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.online-status-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.online-status-enter-from {
  transform: translateY(-100%);
  opacity: 0;
}

.online-status-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}
</style>
