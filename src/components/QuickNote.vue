<template>
  <Teleport to="body">
    <Transition name="quick-note">
      <div v-if="visible" class="quick-note-overlay" @click.self="close">
        <div class="quick-note-panel">
          <div class="quick-note-header">
            <div class="header-left">
              <span class="header-icon">📝</span>
              <h3>快捷便签</h3>
            </div>
            <div class="header-actions">
              <button class="action-btn" @click="addNote" title="新建便签">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
              </button>
              <button class="action-btn" @click="clearAll" title="清空全部">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
              <button class="close-btn" @click="close">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/>
                  <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="quick-note-body">
            <div v-if="notes.length === 0" class="empty-notes">
              <div class="empty-icon">📝</div>
              <p>暂无便签</p>
              <p class="empty-hint">点击 + 创建第一条便签</p>
            </div>

            <div v-else class="notes-list">
              <div
                v-for="note in sortedNotes"
                :key="note.id"
                class="note-card"
                :class="{ editing: editingId === note.id, pinned: note.pinned }"
              >
                <div class="note-card-header">
                  <input
                    v-if="editingId === note.id"
                    v-model="note.title"
                    class="note-title-input"
                    placeholder="便签标题"
                  />
                  <div v-else class="note-title" @click="startEdit(note)">{{ note.title || '未命名便签' }}</div>
                  <div class="note-actions">
                    <button
                      class="note-action-btn"
                      :class="{ active: note.pinned }"
                      @click="togglePin(note)"
                      title="置顶"
                    >
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 2L12 22M12 2L8 6M12 2L16 6"/>
                      </svg>
                    </button>
                    <button class="note-action-btn delete" @click="deleteNote(note.id)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"/>
                        <line x1="6" y1="6" x2="18" y2="18"/>
                      </svg>
                    </button>
                  </div>
                </div>
                <textarea
                  v-if="editingId === note.id"
                  v-model="note.content"
                  class="note-content-input"
                  placeholder="输入内容..."
                  rows="4"
                  @blur="saveNote(note)"
                ></textarea>
                <div v-else class="note-content" @click="startEdit(note)">
                  {{ note.content || '点击编辑内容...' }}
                </div>
                <div class="note-meta">
                  <span class="note-time">{{ formatRelative(note.updatedAt) }}</span>
                  <span v-if="note.color" class="note-color-dot" :style="{ background: note.color }"></span>
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
import { ref, computed, watch } from 'vue'

const props = defineProps({
  visible: Boolean
})

const emit = defineEmits(['update:visible', 'close'])

const STORAGE_KEY = 'orange-quick-notes'

const notes = ref([])
const editingId = ref(null)

const loadNotes = () => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    notes.value = stored ? JSON.parse(stored) : []
  } catch {
    notes.value = []
  }
}

const saveToStorage = () => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(notes.value))
}

const sortedNotes = computed(() => {
  return [...notes.value].sort((a, b) => {
    if (a.pinned !== b.pinned) return b.pinned ? 1 : -1
    return new Date(b.updatedAt) - new Date(a.updatedAt)
  })
})

const addNote = () => {
  const note = {
    id: Date.now().toString(),
    title: '',
    content: '',
    pinned: false,
    color: '',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
  notes.value.unshift(note)
  saveToStorage()
  editingId.value = note.id
}

const startEdit = (note) => {
  editingId.value = note.id
}

const saveNote = (note) => {
  note.updatedAt = new Date().toISOString()
  editingId.value = null
  saveToStorage()
}

const deleteNote = (id) => {
  notes.value = notes.value.filter((n) => n.id !== id)
  if (editingId.value === id) editingId.value = null
  saveToStorage()
}

const togglePin = (note) => {
  note.pinned = !note.pinned
  saveToStorage()
}

const clearAll = () => {
  if (notes.value.length === 0) return
  notes.value = []
  editingId.value = null
  saveToStorage()
}

const close = () => {
  if (editingId.value) {
    const note = notes.value.find((n) => n.id === editingId.value)
    if (note) saveNote(note)
  }
  emit('update:visible', false)
  emit('close')
}

const formatRelative = (time) => {
  if (!time) return ''
  const diff = Date.now() - new Date(time).getTime()
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${Math.floor(diff / 86400000)}天前`
}

watch(
  () => props.visible,
  (val) => {
    if (val) loadNotes()
  }
)

loadNotes()
</script>

<style scoped>
.quick-note-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.quick-note-panel {
  width: calc(100% - 32px);
  max-width: 480px;
  max-height: 80vh;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.quick-note-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  font-size: 22px;
}

.header-left h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 4px;
}

.action-btn,
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

.action-btn:hover {
  background: var(--bg-hover);
  color: var(--primary);
}

.close-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.action-btn svg,
.close-btn svg {
  width: 18px;
  height: 18px;
}

.quick-note-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-notes {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.empty-notes p {
  margin: 4px 0;
  font-size: 14px;
}

.empty-hint {
  font-size: 12px !important;
  opacity: 0.7;
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.note-card {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 14px;
  transition: all 0.2s;
}

.note-card:hover {
  border-color: var(--primary);
}

.note-card.editing {
  border-color: var(--primary);
  box-shadow: 0 0 0 2px rgba(249, 115, 22, 0.1);
}

.note-card.pinned {
  border-left: 3px solid var(--primary);
}

.note-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.note-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  cursor: pointer;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-title-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  outline: none;
}

.note-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
}

.note-action-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.2s;
}

.note-action-btn svg {
  width: 14px;
  height: 14px;
}

.note-action-btn:hover {
  background: var(--bg-card);
  color: var(--text-primary);
}

.note-action-btn.active {
  color: var(--primary);
}

.note-action-btn.delete:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}

.note-content {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  cursor: pointer;
  white-space: pre-wrap;
  word-break: break-word;
}

.note-content-input {
  width: 100%;
  border: none;
  background: transparent;
  font-size: 13px;
  color: var(--text-primary);
  line-height: 1.5;
  resize: vertical;
  outline: none;
  font-family: inherit;
}

.note-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--border-color);
}

.note-time {
  font-size: 11px;
  color: var(--text-muted);
}

.note-color-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.quick-note-enter-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.quick-note-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.quick-note-enter-from,
.quick-note-leave-to {
  opacity: 0;
}

.quick-note-enter-from .quick-note-panel,
.quick-note-leave-to .quick-note-panel {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .quick-note-panel {
    width: 100%;
    max-width: 100%;
    max-height: 90vh;
    border-radius: 16px 16px 0 0;
  }
}
</style>
