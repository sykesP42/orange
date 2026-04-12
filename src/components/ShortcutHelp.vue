<template>
  <Teleport to="body">
    <Transition name="shortcut-help">
      <div v-if="visible" class="shortcut-help-overlay" @click.self="close">
        <div class="shortcut-help-panel">
          <div class="shortcut-help-header">
            <h3>{{ t('shortcuts.title') }}</h3>
            <button class="close-btn" @click="close">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>

          <div class="shortcut-help-body">
            <div v-for="group in shortcutGroups" :key="group.label" class="shortcut-group">
              <div class="shortcut-group-label">{{ group.label }}</div>
              <div class="shortcut-list">
                <div v-for="item in group.items" :key="item.key" class="shortcut-item">
                  <div class="shortcut-desc">{{ item.desc }}</div>
                  <div class="shortcut-keys">
                    <kbd v-for="key in item.keys" :key="key">{{ key }}</kbd>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  visible: Boolean
})

const emit = defineEmits(['update:visible', 'close'])

const shortcutGroups = computed(() => [
  {
    label: t('shortcuts.groups.navigation'),
    items: [
      { desc: t('shortcuts.items.commandPalette'), keys: ['Ctrl+Alt', 'K'] },
      { desc: t('shortcuts.items.home'), keys: ['G', 'H'] },
      { desc: t('shortcuts.items.addBlogger'), keys: ['Ctrl+Alt', 'N'] },
      { desc: t('shortcuts.items.profile'), keys: ['G', 'M'] },
      { desc: t('shortcuts.items.myBloggers'), keys: ['G', 'B'] },
      { desc: t('shortcuts.items.teamCenter'), keys: ['G', 'T'] },
      { desc: t('shortcuts.items.forums'), keys: ['G', 'F'] },
      { desc: t('shortcuts.items.chat'), keys: ['G', 'C'] },
      { desc: t('shortcuts.items.statistics'), keys: ['G', 'S'] },
      { desc: t('shortcuts.items.admin'), keys: ['G', 'D'] },
    ]
  },
  {
    label: t('shortcuts.groups.actions'),
    items: [
      { desc: t('shortcuts.items.quickNote'), keys: ['Ctrl+Alt', 'Q'] },
      { desc: t('shortcuts.items.toggleTheme'), keys: ['Ctrl+Alt', 'D'] },
      { desc: t('shortcuts.items.help'), keys: ['?'] },
    ]
  },
  {
    label: t('shortcuts.groups.general'),
    items: [
      { desc: t('shortcuts.items.search'), keys: ['/'] },
      { desc: t('shortcuts.items.closeDialog'), keys: ['Esc'] },
    ]
  }
])

const close = () => {
  emit('update:visible', false)
  emit('close')
}
</script>

<style scoped>
.shortcut-help-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.shortcut-help-panel {
  width: calc(100% - 32px);
  max-width: 520px;
  max-height: 80vh;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.shortcut-help-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.shortcut-help-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.close-btn svg {
  width: 18px;
  height: 18px;
}

.shortcut-help-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 24px;
}

.shortcut-group {
  margin-bottom: 20px;
}

.shortcut-group:last-child {
  margin-bottom: 0;
}

.shortcut-group-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}

.shortcut-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.shortcut-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background 0.15s;
}

.shortcut-item:hover {
  background: var(--bg-hover);
}

.shortcut-desc {
  font-size: 14px;
  color: var(--text-primary);
}

.shortcut-keys {
  display: flex;
  gap: 4px;
}

.shortcut-keys kbd {
  padding: 3px 8px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 12px;
  font-family: inherit;
  color: var(--text-secondary);
  min-width: 24px;
  text-align: center;
}

.shortcut-help-enter-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.shortcut-help-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.shortcut-help-enter-from,
.shortcut-help-leave-to {
  opacity: 0;
}

.shortcut-help-enter-from .shortcut-help-panel,
.shortcut-help-leave-to .shortcut-help-panel {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .shortcut-help-panel {
    width: 100%;
    max-width: 100%;
    max-height: 90vh;
    border-radius: 16px 16px 0 0;
  }
}
</style>
