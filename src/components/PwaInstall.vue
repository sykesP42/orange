<template>
  <Transition name="pwa-install">
    <div v-if="showPrompt" class="pwa-install-prompt">
      <div class="pwa-install-card">
        <div class="pwa-install-icon">🍊</div>
        <div class="pwa-install-info">
          <div class="pwa-install-title">安装 Orange</div>
          <div class="pwa-install-desc">添加到主屏幕，获得更好的使用体验</div>
        </div>
        <div class="pwa-install-actions">
          <button class="pwa-install-dismiss" @click="dismiss">稍后</button>
          <button class="pwa-install-confirm" @click="install">安装</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const showPrompt = ref(false)
let deferredPrompt = null

const handleBeforeInstall = (e) => {
  e.preventDefault()
  deferredPrompt = e
  const dismissed = localStorage.getItem('pwa-install-dismissed')
  if (!dismissed || Date.now() - parseInt(dismissed) > 7 * 24 * 60 * 60 * 1000) {
    setTimeout(() => {
      showPrompt.value = true
    }, 3000)
  }
}

const install = async () => {
  if (!deferredPrompt) return
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  if (outcome === 'accepted') {
    localStorage.removeItem('pwa-install-dismissed')
  }
  deferredPrompt = null
  showPrompt.value = false
}

const dismiss = () => {
  showPrompt.value = false
  localStorage.setItem('pwa-install-dismissed', Date.now().toString())
}

onMounted(() => {
  window.addEventListener('beforeinstallprompt', handleBeforeInstall)
})

onUnmounted(() => {
  window.removeEventListener('beforeinstallprompt', handleBeforeInstall)
})
</script>

<style scoped>
.pwa-install-prompt {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  width: calc(100% - 32px);
  max-width: 420px;
}

.pwa-install-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.dark .pwa-install-card {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.pwa-install-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.pwa-install-info {
  flex: 1;
  min-width: 0;
}

.pwa-install-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.pwa-install-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.pwa-install-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.pwa-install-dismiss {
  padding: 6px 14px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.pwa-install-dismiss:hover {
  background: var(--bg-hover);
}

.pwa-install-confirm {
  padding: 6px 14px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.pwa-install-confirm:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.pwa-install-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.pwa-install-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.pwa-install-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}

.pwa-install-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}

@media (max-width: 768px) {
  .pwa-install-prompt {
    bottom: 76px;
  }
}
</style>
