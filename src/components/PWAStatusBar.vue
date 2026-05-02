<template>
  <div class="pwa-status-bar" :class="{ 'is-offline': !isOnline }">
    <Transition name="slide-down" mode="out-in">
      <div v-if="!isOnline && showOfflineBar" class="offline-banner" key="offline">
        <div class="offline-content">
          <span class="offline-icon">⚡</span>
          <span class="offline-text">{{ t('pwa.offline.banner') }}</span>
        </div>
      </div>

      <div v-else-if="showUpdateBanner" class="update-banner" key="update">
        <div class="update-content">
          <span class="update-icon">🆕</span>
          <span class="update-text">{{ t('pwa.update.text') }}</span>
          <button class="update-btn" @click="$emit('update')">{{ t('pwa.update.button') }}</button>
          <button class="update-dismiss" @click="dismissUpdate">×</button>
        </div>
      </div>

      <div v-else-if="isInstallable && !hideInstall" class="install-banner" key="install">
        <div class="install-content">
          <span class="install-icon">📲</span>
          <span class="install-text">{{ t('pwa.install.text') }}</span>
          <button class="install-btn" @click="$emit('install')">{{ t('pwa.install.button') }}</button>
          <button class="install-dismiss" @click="$emit('dismiss-install')">×</button>
        </div>
      </div>

      <div v-else-if="!isOnline && !showOfflineBar" key="dot" class="offline-dot-wrapper">
        <span class="offline-dot" :title="t('pwa.offline.tooltip')"></span>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  isOnline: { type: Boolean, default: true },
  showUpdateBanner: { type: Boolean, default: false },
  isInstallable: { type: Boolean, default: false },
  hideInstall: { type: Boolean, default: false },
  showOfflineBar: { type: Boolean, default: true },
})

defineEmits(['update', 'install', 'dismiss-update', 'dismiss-install'])
</script>

<style scoped>
.pwa-status-bar {
  position: relative;
  z-index: 900;
}

.offline-banner {
  background: linear-gradient(90deg, #dc2626, #b91c1c);
  color: var(--color-on-brand);
  padding: 0.4rem 1rem;
  text-align: center;
  font-size: 0.8rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
}
.offline-content { display: flex; align-items: center; gap: 0.5rem; }
.offline-icon { font-size: 0.9rem; }
.offline-text { letter-spacing: 0.02em; }

.update-banner {
  background: linear-gradient(90deg, #f97316, #ea580c);
  color: var(--color-on-brand);
  padding: 0.35rem 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}
.update-content { display: flex; align-items: center; gap: 0.6rem; }
.update-icon { font-size: 0.85rem; }
.update-text { font-size: 0.82rem; font-weight: 500; }
.update-btn {
  background: rgba(255,255,255,0.22);
  border: 1px solid rgba(255,255,255,0.35);
  color: var(--color-on-brand);
  padding: 2px 12px;
  border-radius: 6px;
  font-size: 0.78rem;
  cursor: pointer;
  transition: background 0.2s;
  font-weight: 600;
}
.update-btn:hover { background: rgba(255,255,255,0.32); }
.update-dismiss {
  background: none;
  border: none;
  color: rgba(255,255,255,0.7);
  font-size: 1.1rem;
  cursor: pointer;
  padding: 0 2px;
  line-height: 1;
}
.update-dismiss:hover { color: var(--color-on-brand); }

.install-banner {
  background: var(--bg-success-gradient);
  color: var(--color-on-brand);
  padding: 0.35rem 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}
.install-content { display: flex; align-items: center; gap: 0.6rem; }
.install-icon { font-size: 0.85rem; }
.install-text { font-size: 0.82rem; font-weight: 500; }
.install-btn {
  background: rgba(255,255,255,0.22);
  border: 1px solid rgba(255,255,255,0.35);
  color: var(--color-on-brand);
  padding: 2px 14px;
  border-radius: 6px;
  font-size: 0.78rem;
  cursor: pointer;
  transition: background 0.2s;
  font-weight: 600;
}
.install-btn:hover { background: rgba(255,255,255,0.32); }
.install-dismiss {
  background: none;
  border: none;
  color: rgba(255,255,255,0.7);
  font-size: 1.1rem;
  cursor: pointer;
  padding: 0 2px;
}
.install-dismiss:hover { color: var(--color-on-brand); }

.offline-dot-wrapper {
  position: absolute;
  top: -2px;
  right: 60px;
  z-index: 10;
}
.offline-dot {
  display: block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--danger);
  animation: pulse-dot 2s ease-in-out infinite;
  box-shadow: 0 0 6px rgba(239,68,68,0.5);
}
@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.85); }
}

.slide-down-enter-active,
.slide-down-leave-active { transition: all 0.3s ease; }
.slide-down-enter-from { transform: translateY(-100%); opacity: 0; }
.slide-down-leave-to { transform: translateY(-100%); opacity: 0; }

@media (max-width: 480px) {
  .offline-text, .update-text, .install-text { font-size: 0.75rem; }
  .update-btn, .install-btn { padding: 2px 10px; font-size: 0.74rem; }
  .offline-dot-wrapper { right: 48px; }
}
</style>
