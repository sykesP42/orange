<template>
  <div class="empty-state" :style="{ height: height + 'px' }">
    <div class="empty-icon">
      <slot name="icon">
        <el-icon :size="iconSize">
          <component :is="icon" />
        </el-icon>
      </slot>
    </div>
    <div class="empty-title">{{ title }}</div>
    <div class="empty-desc" v-if="description">{{ description }}</div>
    <div class="empty-action" v-if="$slots.action">
      <slot name="action"></slot>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  icon: {
    type: String,
    default: 'DocumentDelete'
  },
  title: {
    type: String,
    default: () => t('app.noData')
  },
  description: {
    type: String,
    default: ''
  },
  height: {
    type: Number,
    default: 300
  },
  iconSize: {
    type: Number,
    default: 64
  }
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
}

.empty-icon {
  margin-bottom: 24px;
  opacity: 0.5;
  animation: float 3s ease-in-out infinite;
}

.empty-title {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.empty-desc {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 24px;
  max-width: 400px;
  line-height: 1.6;
}

.empty-action {
  margin-top: 16px;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

/* 移动端优化 */
@media (max-width: 768px) {
  .empty-state {
    padding: 30px 16px;
  }
  
  .empty-title {
    font-size: 16px;
  }
  
  .empty-desc {
    font-size: 13px;
  }
}
</style>
