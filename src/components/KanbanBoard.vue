<template>
  <div class="kanban-board">
    <div
      v-for="col in columns"
      :key="col.status"
      class="kanban-column"
      :class="col.cssClass"
      @dragover.prevent="onDragOver($event, col.status)"
      @dragleave="onDragLeave($event)"
      @drop="onDrop($event, col.status)"
    >
      <div class="kanban-column-header">
        <span class="kanban-status-dot" :style="{ background: col.color }"></span>
        <span class="kanban-status-name">{{ col.status }}</span>
        <span class="kanban-count">{{ getColumnBloggers(col.status).length }}</span>
      </div>
      <div class="kanban-column-body">
        <div
          v-for="blogger in getColumnBloggers(col.status)"
          :key="blogger.id"
          class="kanban-card"
          draggable="true"
          @dragstart="onDragStart($event, blogger)"
          @dragend="onDragEnd($event)"
          @click="$emit('goToDetail', blogger.id)"
        >
          <div class="kanban-card-header">
            <div class="kanban-avatar" v-if="blogger.avatar">
              <img :src="blogger.avatar" :alt="blogger.nickname" />
            </div>
            <div class="kanban-avatar kanban-avatar-placeholder" v-else>
              {{ (blogger.nickname || '?')[0] }}
            </div>
            <div class="kanban-card-info">
              <span class="kanban-nickname">{{ blogger.nickname }}</span>
              <span class="kanban-category" v-if="blogger.category">{{ blogger.category }}</span>
            </div>
            <div class="kanban-mobile-move" @click.stop="showMoveMenu($event, blogger)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="6 9 12 15 18 9"/>
              </svg>
            </div>
          </div>
          <div class="kanban-card-meta" v-if="blogger.platform || blogger.product">
            <span class="kanban-platform" v-if="blogger.platform">{{ blogger.platform }}</span>
            <span class="kanban-product" v-if="blogger.product">{{ blogger.product }}</span>
          </div>
          <div class="kanban-card-tags" v-if="blogger.tags && blogger.tags.length">
            <span
              v-for="tag in blogger.tags.slice(0, 3)"
              :key="tag"
              class="kanban-tag"
            >{{ tag }}</span>
            <span v-if="blogger.tags.length > 3" class="kanban-tag-more">+{{ blogger.tags.length - 3 }}</span>
          </div>
          <div class="kanban-card-footer">
            <span class="kanban-real-name">{{ blogger.real_name }}</span>
          </div>
        </div>
        <div v-if="getColumnBloggers(col.status).length === 0" class="kanban-empty">
          <span>拖拽博主到此处</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  bloggers: { type: Array, default: () => [] }
})

const emit = defineEmits(['goToDetail', 'updateStatus'])

const columns = [
  { status: '初次联系', color: 'var(--info)', cssClass: 'col-first' },
  { status: '洽谈中', color: 'var(--primary)', cssClass: 'col-negotiating' },
  { status: '已合作', color: 'var(--success)', cssClass: 'col-cooperated' },
  { status: '已拒绝', color: 'var(--danger)', cssClass: 'col-rejected' },
  { status: '暂停合作', color: '#6b7280', cssClass: 'col-paused' }
]

const getColumnBloggers = (status) => {
  return props.bloggers.filter(b => (b.status || '初次联系') === status)
}

let dragBlogger = null

const onDragStart = (e, blogger) => {
  dragBlogger = blogger
  e.dataTransfer.effectAllowed = 'move'
  e.target.classList.add('kanban-card-dragging')
}

const onDragEnd = (e) => {
  e.target.classList.remove('kanban-card-dragging')
  dragBlogger = null
  document.querySelectorAll('.kanban-column').forEach(col => {
    col.classList.remove('kanban-column-dragover')
  })
}

const onDragOver = (e, status) => {
  e.preventDefault()
  e.dataTransfer.dropEffect = 'move'
  e.currentTarget.classList.add('kanban-column-dragover')
}

const onDragLeave = (e) => {
  e.currentTarget.classList.remove('kanban-column-dragover')
}

const onDrop = (e, newStatus) => {
  e.preventDefault()
  e.currentTarget.classList.remove('kanban-column-dragover')
  if (!dragBlogger) return
  const oldStatus = dragBlogger.status || '初次联系'
  if (oldStatus !== newStatus) {
    emit('updateStatus', dragBlogger.id, newStatus, oldStatus)
  }
}

const showMoveMenu = (e, blogger) => {
  const currentStatus = blogger.status || '初次联系'
  const otherStatuses = columns.filter(c => c.status !== currentStatus)
  const menu = document.createElement('div')
  menu.className = 'kanban-move-menu'
  otherStatuses.forEach(col => {
    const btn = document.createElement('button')
    btn.className = 'kanban-move-btn'
    const dot = document.createElement('span')
    dot.style.cssText = 'display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:8px'
    dot.style.background = col.color
    btn.appendChild(dot)
    btn.appendChild(document.createTextNode(col.status))
    btn.addEventListener('click', () => {
      emit('updateStatus', blogger.id, col.status, currentStatus)
      menu.remove()
    })
    menu.appendChild(btn)
  })
  const existing = document.querySelector('.kanban-move-menu')
  if (existing) existing.remove()
  document.body.appendChild(menu)
  const rect = e.currentTarget.getBoundingClientRect()
  menu.style.top = `${rect.bottom + 4}px`
  menu.style.left = `${Math.min(rect.left, window.innerWidth - 160)}px`
  setTimeout(() => {
    const close = (ev) => {
      if (!menu.contains(ev.target)) {
        menu.remove()
        document.removeEventListener('click', close)
      }
    }
    document.addEventListener('click', close)
  }, 0)
}
</script>

<style scoped>
.kanban-board {
  display: flex;
  gap: 12px;
  overflow-x: auto;
  padding-bottom: 16px;
  min-height: 400px;
}

.kanban-column {
  flex: 1;
  min-width: 220px;
  max-width: 300px;
  background: var(--bg-primary, #f5f7fa);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  transition: background 0.2s, box-shadow 0.2s;
  border: 2px solid transparent;
}

.kanban-column-dragover {
  background: var(--bg-hover, #eff6ff);
  border-color: var(--primary, #f97316);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.15);
}

.kanban-column-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  border-bottom: 1px solid var(--border-color, #e5e7eb);
  font-weight: 600;
  font-size: 14px;
}

.kanban-status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.kanban-status-name {
  flex: 1;
  color: var(--text-primary, #1f2937);
}

.kanban-count {
  background: var(--bg-card, #fff);
  color: var(--text-muted, #6b7280);
  font-size: 12px;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: 10px;
  border: 1px solid var(--border-color, #e5e7eb);
}

.kanban-column-body {
  flex: 1;
  padding: 8px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 200px;
}

.kanban-card {
  background: var(--bg-card, #fff);
  border-radius: 10px;
  padding: 12px;
  cursor: grab;
  transition: transform 0.15s, box-shadow 0.15s;
  border: 1px solid var(--border-color, #e5e7eb);
}

.kanban-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.kanban-card:active {
  cursor: grabbing;
}

.kanban-card-dragging {
  opacity: 0.4;
  transform: rotate(2deg);
}

.kanban-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.kanban-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.kanban-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.kanban-avatar-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary, #f97316);
  color: var(--color-on-brand);
  font-weight: 600;
  font-size: 14px;
}

.kanban-card-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.kanban-nickname {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-primary, #1f2937);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kanban-category {
  font-size: 12px;
  color: var(--text-muted, #6b7280);
}

.kanban-card-meta {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}

.kanban-platform,
.kanban-product {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-primary, #f5f7fa);
  color: var(--text-secondary, #4b5563);
}

.kanban-card-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-bottom: 6px;
}

.kanban-tag {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--primary-50, #fff7ed);
  color: var(--primary, #f97316);
}

.kanban-tag-more {
  font-size: 11px;
  color: var(--text-muted, #6b7280);
}

.kanban-card-footer {
  display: flex;
  justify-content: flex-end;
  font-size: 11px;
  color: var(--text-muted, #6b7280);
}

.kanban-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 80px;
  color: var(--text-muted, #9ca3af);
  font-size: 13px;
  border: 2px dashed var(--border-color, #e5e7eb);
  border-radius: 8px;
  margin: 4px;
}

.kanban-mobile-move {
  display: none;
  width: 28px;
  height: 28px;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: var(--bg-primary, #f5f7fa);
  color: var(--text-muted, #6b7280);
  cursor: pointer;
  margin-left: auto;
  flex-shrink: 0;
}

.kanban-mobile-move svg {
  width: 16px;
  height: 16px;
}

@media (max-width: 768px) {
  .kanban-mobile-move {
    display: flex;
  }
  .kanban-board {
    flex-direction: column;
    gap: 12px;
    overflow-x: visible;
    min-height: auto;
  }
  .kanban-column {
    min-width: 100%;
    max-width: 100%;
  }
  .kanban-column-body {
    flex-direction: row;
    flex-wrap: wrap;
    min-height: auto;
    gap: 8px;
  }
  .kanban-card {
    flex: 1 1 calc(50% - 8px);
    min-width: 140px;
  }
  .kanban-empty {
    width: 100%;
    min-height: 48px;
    border-width: 1px;
  }
}

@media (max-width: 480px) {
  .kanban-card {
    flex: 1 1 100%;
  }
  .kanban-card-header {
    gap: 8px;
  }
  .kanban-avatar {
    width: 32px;
    height: 32px;
  }
  .kanban-nickname {
    font-size: 13px;
  }
  .kanban-platform,
  .kanban-product {
    font-size: 12px;
    padding: 3px 8px;
  }
  .kanban-tag {
    font-size: 12px;
    padding: 2px 8px;
  }
  .kanban-card-footer {
    font-size: 12px;
  }
  .kanban-status-name {
    font-size: 13px;
  }
}
</style>
