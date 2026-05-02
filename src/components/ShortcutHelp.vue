<template>
  <Teleport to="body">
    <Transition name="shortcut-help">
      <div v-if="visible" class="shortcut-help-overlay" @click.self="close">
        <div class="shortcut-help-panel">
          <div class="shortcut-help-header">
            <h3>⌨️ 快捷键</h3>
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
defineProps({
  visible: Boolean
})

const emit = defineEmits(['update:visible', 'close'])

const shortcutGroups = [
  {
    label: '导航',
    items: [
      { desc: '命令面板', keys: ['Ctrl', 'K'] },
      { desc: '首页', keys: ['G', 'H'] },
      { desc: '录入博主', keys: ['G', 'A'] },
      { desc: '个人中心', keys: ['G', 'M'] },
      { desc: '我的博主', keys: ['G', 'B'] },
      { desc: '团队中心', keys: ['G', 'T'] },
      { desc: '公共论坛', keys: ['G', 'F'] },
      { desc: '私信', keys: ['G', 'C'] },
      { desc: '数据统计', keys: ['G', 'S'] },
      { desc: '管理后台', keys: ['G', 'D'] },
    ]
  },
  {
    label: '操作',
    items: [
      { desc: '快捷便签', keys: ['Ctrl', 'N'] },
      { desc: '切换主题', keys: ['Ctrl', 'D'] },
      { desc: '快捷键帮助', keys: ['?'] },
    ]
  },
  {
    label: '通用',
    items: [
      { desc: '搜索', keys: ['/'] },
      { desc: '关闭弹窗', keys: ['Esc'] },
    ]
  }
]

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
