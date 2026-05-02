<template>
  <div class="backup-manager" :class="{ 'dark-mode': isDark }">
    <div class="backup-header">
      <h3>📦 数据备份与恢�?/h3>
      <div class="backup-stats">
        <span class="stat-item">�?{{ stats.totalBackups }} 个备�?/span>
        <span class="stat-item">最�? {{ stats.lastBackup }}</span>
        <span class="stat-item">上限: {{ stats.maxBackups }}</span>
      </div>
    </div>

    <div class="backup-actions-bar">
      <button class="action-btn primary" @click="handleCreateBackup" :disabled="backup.isBackingUp">
        {{ backup.isBackingUp ? '备份�?..' : '�?创建备份' }}
      </button>
      <label class="action-btn" title="从文件导入备�?>
        📥 导入
        <input type="file" accept=".json" @change="handleImportFile" hidden />
      </label>
      <button class="action-btn danger" @click="handleClearAll" :disabled="backup.backupList.length === 0">
        🗑�?清空
      </button>
    </div>

    <transition-group name="backup-list" tag="div" class="backup-list">
      <div v-for="item in backup.backupList" :key="item.id" class="backup-item">
        <div class="backup-info">
          <div class="backup-label">{{ item.label || '手动备份' }}</div>
          <div class="backup-meta">
            <span>{{ backup.timeAgo(item.timestamp) }}</span>
            <span>{{ backup.formatSize(item.dataSize) }}</span>
            <span>{{ item.recordCount }} 条记�?/span>
            <span class="backup-version">v{{ item.version }}</span>
          </div>
        </div>
        <div class="backup-item-actions">
          <button class="item-btn restore" @click="handleRestore(item)" :disabled="backup.isRestoring" title="恢复此备�?>↩️</button>
          <button class="item-btn download" @click="handleDownload(item)" title="下载备份文件">⬇️</button>
          <button class="item-btn delete" @click="handleDelete(item)" title="删除备份">�?/button>
        </div>
      </div>
    </transition-group>

    <div v-if="backup.backupList.length === 0" class="backup-empty">
      <span class="empty-icon">💾</span>
      <p>暂无备份记录</p>
      <p class="empty-hint">点击上方按钮创建第一个备�?/p>
    </div>
  </div>
</template>

<script setup>
import { useDataBackup } from '../composables/useDataBackup'
import { useToastStore } from '../stores/toast'
import { computed } from 'vue'

const props = defineProps({
  isDark: Boolean,
  dataSource: {
    type: Function,
    default: null
  }
})

const toastStore = useToastStore()
const backup = useDataBackup()
const stats = computed(() => backup.stats)

async function handleCreateBackup() {
  if (!props.dataSource) {
    toastStore.warning('未提供数据源，使用默认空备份')
  }
  const data = props.dataSource ? await props.dataSource() : { bloggers: [], timestamp: new Date().toISOString(), source: 'manual' }
  try {
    const result = await backup.createBackup(data)
    toastStore.success('备份创建成功')
  } catch (e) {
    toastStore.error('备份失败', e.message)
  }
}

async function handleDownload(item) {
  const stored = localStorage.getItem(`orange_backup_${item.id}`)
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      await backup.downloadBackup(parsed.data, `orange_backup_${item.label || 'backup'}_${item.timestamp.slice(0, 10)}.json`)
      toastStore.success('备份已下�?)
    } catch {
      toastStore.error('下载失败')
    }
  }
}

async function handleRestore(item) {
  try {
    const data = await backup.restoreFromLocal(item.id)
    window.dispatchEvent(new CustomEvent('data-restore', { detail: data }))
    toastStore.success('数据已恢复，请刷新页�?)
  } catch (e) {
    toastStore.error('恢复失败', e.message)
  }
}

async function handleImportFile(e) {
  const file = e.target.files?.[0]
  if (!file) return
  try {
    const data = await backup.importFromFile(file)
    await backup.createBackup(data, `导入-${file.name}`)
    window.dispatchEvent(new CustomEvent('data-import', { detail: data }))
    toastStore.success('数据导入成功并已创建备份')
  } catch (err) {
    toastStore.error('导入失败', err.message)
  }
  e.target.value = ''
}

function handleDelete(item) {
  backup.deleteBackup(item.id)
  toastStore.success('备份已删�?)
}

function handleClearAll() {
  backup.clearAllBackups()
  toastStore.success('所有备份已清空')
}
</script>

<style scoped>
.backup-manager {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
}

.backup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}

.backup-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.backup-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-muted);
}

.stat-item {
  padding: 2px 8px;
  background: var(--bg-secondary);
  border-radius: 10px;
}

.backup-actions-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.action-btn {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.action-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.action-btn.primary {
  background: var(--primary);
  color: var(--color-on-brand);
  border-color: var(--primary);
}

.action-btn.primary:hover {
  opacity: 0.9;
}

.action-btn.danger:hover {
  border-color: #ef4444;
  color: #ef4444;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.backup-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 400px;
  overflow-y: auto;
}

.backup-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: all 0.2s;
}

.backup-item:hover {
  border-color: var(--primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.backup-info {
  flex: 1;
  min-width: 0;
}

.backup-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.backup-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-muted);
  flex-wrap: wrap;
}

.backup-version {
  padding: 0 4px;
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
  border-radius: 4px;
  font-size: 10px;
}

.backup-item-actions {
  display: flex;
  gap: 4px;
  margin-left: 12px;
}

.item-btn {
  width: 30px;
  height: 30px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.item-btn:hover {
  transform: scale(1.1);
}

.item-btn.restore:hover {
  border-color: #22c55e;
  background: rgba(34, 197, 94, 0.08);
}

.item-btn.download:hover {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.08);
}

.item-btn.delete:hover {
  border-color: #ef4444;
  background: rgba(239, 68, 68, 0.08);
}

.item-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.backup-empty {
  padding: 48px 16px;
  text-align: center;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 12px;
}

.empty-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.backup-list-enter-active {
  transition: all 0.3s ease;
}
.backup-list-leave-active {
  transition: all 0.2s ease;
}
.backup-list-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}
.backup-list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

@media (max-width: 640px) {
  .backup-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .backup-actions-bar {
    flex-direction: column;
  }
  .action-btn {
    width: 100%;
    justify-content: center;
  }
  .backup-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .backup-item-actions {
    margin-left: 0;
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
