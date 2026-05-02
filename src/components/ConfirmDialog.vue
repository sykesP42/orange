<template>
  <Teleport to="body">
    <Transition name="confirm-fade">
      <div v-if="isVisible" class="confirm-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="handleCancel">
        <div class="confirm-dialog" style="background:var(--bg-card);" :class="{ 'dark-mode': isDark }">
          <div class="confirm-icon" :class="dialogType">
            <svg v-if="dialogType === 'danger'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <svg v-else-if="dialogType === 'warning'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
              <line x1="12" y1="9" x2="12" y2="13"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
          </div>
          <h3 class="confirm-title">{{ title }}</h3>
          <p class="confirm-message">{{ message }}</p>
          <div class="confirm-actions">
            <button class="btn btn-secondary" @click="handleCancel">{{ cancelText }}</button>
            <button class="btn" :class="confirmButtonClass" @click="handleConfirm">{{ confirmText }}</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useConfirm } from '../utils/confirm'
import { useThemeStore } from '../stores/theme'

const { isVisible, title, message, dialogType, confirmText, cancelText, handleConfirm, handleCancel } = useConfirm()
const themeStore = useThemeStore()
const isDark = computed(() => themeStore.isDark)

const confirmButtonClass = computed(() => {
  return {
    'btn-danger': dialogType.value === 'danger',
    'btn-warning': dialogType.value === 'warning',
    'btn-primary': dialogType.value === 'info'
  }
})
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100001;
  backdrop-filter: blur(2px);
}

.confirm-dialog {
  background: var(--bg-card);
  border-radius: 16px;
  padding: 32px;
  width: 380px;
  max-width: 90vw;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  animation: confirmPop 0.2s ease;
}

@keyframes confirmPop {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.confirm-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
}

.confirm-icon svg {
  width: 28px;
  height: 28px;
}

.confirm-icon.info {
  background: linear-gradient(135deg, #4facfe, #00f2fe);
  color: white;
}

.confirm-icon.warning {
  background: linear-gradient(135deg, #fa709a, #fee140);
  color: white;
}

.confirm-icon.danger {
  background: linear-gradient(135deg, #ff416c, #ff4b2b);
  color: white;
}

.confirm-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 12px;
}

.confirm-message {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 28px;
  line-height: 1.6;
}

.confirm-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.confirm-actions .btn {
  flex: 1;
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--bg-tertiary);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-danger {
  background: linear-gradient(135deg, #ff4b2b, #ff416c);
  color: white;
}

.btn-danger:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255, 75, 75, 0.3);
}

.btn-warning {
  background: linear-gradient(135deg, #fa709a, #fee140);
  color: white;
}

.btn-warning:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(250, 112, 154, 0.3);
}

.confirm-fade-enter-active,
.confirm-fade-leave-active {
  transition: opacity 0.2s ease;
}

.confirm-fade-enter-from,
.confirm-fade-leave-to {
  opacity: 0;
}

html.dark .confirm-overlay {
  background: rgba(0, 0, 0, 0.8);
}

.confirm-fade-enter-active .confirm-dialog {
  animation: confirmPop 0.2s ease;
}

.confirm-fade-leave-active .confirm-dialog {
  animation: confirmPop 0.15s ease reverse;
}

.confirm-dialog.dark-mode {
  border: 1px solid var(--border-color);
}
</style>
