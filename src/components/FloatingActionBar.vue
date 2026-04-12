<template>
  <div class="fab-container" :class="{ open: isOpen, 'dark-mode': isDark }">
    <Transition name="fab-overlay">
      <div v-if="isOpen" class="fab-backdrop" @click="close"></div>
    </Transition>

    <div class="fab-menu" :class="{ expanded: isOpen }">
      <TransitionGroup name="fab-item">
        <button
          v-for="(item, index) in items"
          v-show="isOpen"
          :key="item.key"
          class="fab-action"
          :style="{ '--index': index, '--total': items.length }"
          @click="handleAction(item)"
        >
          <span class="fab-action-icon" :class="item.iconClass || ''">
            <component :is="item.icon" v-if="typeof item.icon === 'object'" />
            <svg v-else-if="item.svgIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="item.svgIcon"></svg>
            <span v-else>{{ item.label.charAt(0) }}</span>
          </span>
          <span class="fab-action-label">{{ item.label }}</span>
          <kbd v-if="item.shortcut" class="fab-shortcut">{{ item.shortcut }}</kbd>
        </button>
      </TransitionGroup>
    </div>

    <button
      class="fab-main"
      :class="{ active: isOpen, pulse: showPulse }"
      @click="toggle"
      :title="isOpen ? t('fab.closeTooltip') : t('fab.tooltip')"
      :aria-label="t('fab.ariaLabel')"
    >
      <Transition name="fab-icon" mode="out-in">
        <svg v-if="!isOpen" key="plus" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        <svg v-else key="close" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </Transition>
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  isDark: { type: Boolean, default: false },
})

const emit = defineEmits(['action'])

const router = useRouter()
const isOpen = ref(false)
const showPulse = ref(true)

const items = computed(() => [
  {
    key: 'add',
    label: t('fab.items.add'),
    route: '/add',
    shortcut: 'N',
    svgIcon: '<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/>',
    color: '#f97316',
  },
  {
    key: 'search',
    label: t('fab.items.search'),
    action: 'command-palette',
    shortcut: 'K',
    svgIcon: '<circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>',
    color: '#3b82f6',
  },
  {
    key: 'calendar',
    label: t('fab.items.calendar'),
    route: '/calendar',
    svgIcon: '<rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>',
    color: '#8b5cf6',
  },
  {
    key: 'workflow',
    label: t('fab.items.workflow'),
    route: '/workflow',
    svgIcon: '<path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>',
    color: '#06b6d4',
  },
  {
    key: 'kanban',
    label: t('fab.items.kanban'),
    action: 'open-kanban',
    svgIcon: '<rect x="3" y="3" width="7" height="9" rx="1"/><rect x="14" y="3" width="7" height="5" rx="1"/><rect x="14" y="12" width="7" height="9" rx="1"/><rect x="3" y="16" width="7" height="5" rx="1"/>',
    color: '#10b981',
  },
])

const toggle = () => {
  isOpen.value = !isOpen.value
  if (!isOpen.value && showPulse.value) {
    setTimeout(() => { showPulse.value = false }, 300)
  }
}

const close = () => {
  isOpen.value = false
}

const handleAction = (item) => {
  close()

  if (item.route) {
    router.push(item.route)
  } else if (item.action === 'command-palette') {
    window.dispatchEvent(new CustomEvent('open-command-palette'))
  } else if (item.action === 'open-kanban') {
    window.dispatchEvent(new CustomEvent('open-kanban'))
  }

  emit('action', item.key)
}
</script>

<style scoped>
.fab-container {
  position: fixed;
  z-index: 800;
  bottom: 24px;
  right: 24px;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.fab-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(2px);
  z-index: -1;
  border-radius: 0;
}

.fab-menu {
  display: flex;
  flex-direction: column-reverse;
  align-items: flex-end;
  gap: 10px;
  pointer-events: none;
}

.fab-menu.expanded .fab-action {
  pointer-events: auto;
}

.fab-action {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0;
  background: white;
  border: none;
  border-radius: 14px;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1), 0 1px 4px rgba(0, 0, 0, 0.06);
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
  transform: translateY(calc((var(--total) - var(--index) - 1) * 60px)) scale(0.3);
  opacity: 0;
  animation: fab-in 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  animation-delay: calc(var(--index) * 45ms);
}

@keyframes fab-in {
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

.fab-action:hover {
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.15), 0 2px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px) scale(1.03) !important;
}

.fab-action:active {
  transform: scale(0.96) !important;
}

.fab-action-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  min-width: 42px;
  border-radius: 12px;
  font-size: 18px;
  color: white;
  background: var(--action-color, #f97316);
  transition: background 0.2s;
}

.fab-action-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--text-secondary);
  padding-right: 14px;
  white-space: nowrap;
  letter-spacing: 0.01em;
}

.fab-shortcut {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 18px;
  padding: 0 5px;
  margin-right: 12px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 0.68rem;
  color: var(--text-tertiary);
  font-weight: 600;
  font-family: inherit;
}

.fab-main {
  position: relative;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(135deg, #f97316, #ea580c);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 6px 20px rgba(249, 115, 22, 0.35),
    0 2px 6px rgba(249, 115, 22, 0.2),
    inset 0 -2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 2;
}

.fab-main:hover {
  transform: scale(1.08);
  box-shadow:
    0 8px 30px rgba(249, 115, 22, 0.45),
    0 3px 10px rgba(249, 115, 22, 0.25);
}

.fab-main:active {
  transform: scale(0.94);
}

.fab-main.active {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
  box-shadow:
    0 4px 16px rgba(220, 38, 38, 0.35),
    0 2px 6px rgba(220, 38, 38, 0.2);
  transform: rotate(90deg) scale(1);
}

.fab-main.pulse::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  background: rgba(249, 115, 22, 0.3);
  animation: pulse-ring 2s ease-out infinite;
}

@keyframes pulse-ring {
  0% { transform: scale(0.85); opacity: 0.7; }
  100% { transform: scale(1.6); opacity: 0; }
}

.fab-main svg {
  width: 26px;
  height: 26px;
}

.dark-mode .fab-action {
  background: #1f2937;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}
.dark-mode .fab-action-label { color: #e5e7eb; }
.dark-mode .fab-shortcut {
  background: #374151;
  border-color: var(--text-secondary);
  color: #9ca3af;
}

.fab-overlay-enter-active,
.fab-overlay-leave-active { transition: opacity 0.25s ease; }
.fab-overlay-enter-from,
.fab-overlay-leave-to { opacity: 0; }

.fab-icon-enter-active,
.fab-icon-leave-active { transition: all 0.2s ease; }
.fab-icon-enter-from { transform: rotate(-90deg) scale(0); opacity: 0; }
.fab-icon-leave-to { transform: rotate(90deg) scale(0); opacity: 0; }

.fab-item-enter-active { transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1); }
.fab-item-leave-active { transition: all 0.2s ease; }
.fab-item-enter-from {
  transform: translateY(20px) scale(0.5);
  opacity: 0;
}
.fab-item-leave-to {
  transform: translateY(-10px) scale(0.8);
  opacity: 0;
}

@media (max-width: 768px) {
  .fab-container { bottom: 84px; right: 16px; }
  .fab-main { width: 52px; height: 52px; }
  .fab-main svg { width: 24px; height: 24px; }
  .fab-action-icon { width: 38px; height: 38px; min-width: 38px; font-size: 16px; }
  .fab-action-label { font-size: 0.82rem; padding-right: 10px; }
  .fab-shortcut { display: none; }
  .fab-action { gap: 8px; }
}

@media (max-width: 480px) {
  .fab-container { bottom: 78px; right: 12px; }
  .fab-main { width: 48px; height: 48px; }
  .fab-main svg { width: 22px; height: 22px; }
  .fab-action-icon { width: 36px; height: 36px; min-width: 36px; border-radius: 10px; }
  .fab-action-label { font-size: 0.78rem; }
  .fab-action { border-radius: 12px; }
}
</style>
