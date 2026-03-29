<template>
  <Teleport to="body">
    <TransitionGroup name="notification" tag="div" class="notification-container">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="notification-item"
        :class="notification.type"
      >
        <div class="notification-icon">
          <span v-if="notification.type === 'success'">✓</span>
          <span v-else-if="notification.type === 'error'">✕</span>
          <span v-else-if="notification.type === 'warning'">⚠</span>
          <span v-else>ℹ</span>
        </div>
        <div class="notification-content">
          <div class="notification-message">{{ notification.message }}</div>
        </div>
        <button class="notification-close" @click="removeNotification(notification.id)">×</button>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<script setup>
import { useNotification } from '../stores/notification'

const { notifications, removeNotification } = useNotification()
</script>

<style scoped>
.notification-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}

.notification-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  min-width: 280px;
  max-width: 360px;
  pointer-events: auto;
  border-left: 4px solid;
}

.dark .notification-item {
  background: rgba(30, 41, 59, 0.98);
  color: #f1f5f9;
}

.notification-item.success {
  border-color: #22c55e;
}

.notification-item.error {
  border-color: #ef4444;
}

.notification-item.warning {
  border-color: #f97316;
}

.notification-item.info {
  border-color: #3b82f6;
}

.notification-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  flex-shrink: 0;
}

.notification-item.success .notification-icon {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.notification-item.error .notification-icon {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.notification-item.warning .notification-icon {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.notification-item.info .notification-icon {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-message {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  line-height: 1.4;
}

.dark .notification-message {
  color: #f1f5f9;
}

.notification-close {
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  color: #94a3b8;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.notification-close:hover {
  background: rgba(0, 0, 0, 0.08);
  color: #64748b;
}

.dark .notification-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.notification-enter-active {
  animation: slideIn 0.3s ease;
}

.notification-leave-active {
  animation: slideOut 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}
</style>
