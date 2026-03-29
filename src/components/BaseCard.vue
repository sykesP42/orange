<template>
  <div class="base-card" :class="['base-card--' + type, { 'is-hoverable': hoverable, 'is-clickable': clickable }]" 
       @click="handleClick"
       @mouseenter="handleMouseEnter"
       @mouseleave="handleMouseLeave">
    <div class="base-card-header" v-if="$slots.header">
      <slot name="header"></slot>
    </div>
    <div class="base-card-body">
      <slot></slot>
    </div>
    <div class="base-card-footer" v-if="$slots.footer">
      <slot name="footer"></slot>
    </div>
    <div class="base-card-overlay" v-if="hoverable && isHovered">
      <slot name="overlay"></slot>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'primary', 'success', 'warning', 'danger', 'info'].includes(value)
  },
  hoverable: {
    type: Boolean,
    default: false
  },
  clickable: {
    type: Boolean,
    default: false
  },
  shadow: {
    type: String,
    default: 'always',
    validator: (value) => ['always', 'hover', 'never'].includes(value)
  }
})

const emit = defineEmits(['click'])

const isHovered = ref(false)

const handleMouseEnter = () => {
  if (props.hoverable) {
    isHovered.value = true
  }
}

const handleMouseLeave = () => {
  isHovered.value = false
}

const handleClick = (event) => {
  if (props.clickable) {
    emit('click', event)
  }
}
</script>

<style scoped>
.base-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 20px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.base-card-body {
  position: relative;
  z-index: 1;
}

.base-card-header {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-light);
}

.base-card-footer {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-light);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

/* 类型样式 */
.base-card--primary {
  border-left: 4px solid var(--primary);
}

.base-card--success {
  border-left: 4px solid var(--success);
}

.base-card--warning {
  border-left: 4px solid var(--warning);
}

.base-card--danger {
  border-left: 4px solid var(--danger);
}

.base-card--info {
  border-left: 4px solid #3b82f6;
}

/* 悬停效果 */
.base-card.is-hoverable {
  cursor: pointer;
}

.base-card.is-hoverable:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

/* 点击效果 */
.base-card.is-clickable {
  cursor: pointer;
}

.base-card.is-clickable:active {
  transform: scale(0.98);
}

/* 阴影控制 */
.base-card[shadow="never"] {
  box-shadow: none;
}

.base-card[shadow="hover"]:not(.is-hoverable) {
  box-shadow: none;
}

.base-card[shadow="hover"].is-hoverable:hover {
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

/* 覆盖层 */
.base-card-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
  opacity: 0;
  transition: opacity 0.3s;
}

.base-card.is-hovered .base-card-overlay {
  opacity: 1;
}

/* 移动端触摸反馈 */
@media (hover: none) {
  .base-card.is-hoverable:hover {
    transform: none;
  }
  
  .base-card.is-hoverable:active {
    transform: scale(0.98);
  }
}

/* 深色模式 */
html.dark .base-card:hover {
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
}
</style>
