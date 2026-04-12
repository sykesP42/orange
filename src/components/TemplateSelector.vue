<template>
  <div class="template-selector">
    <div class="selector-trigger">
      <button class="trigger-btn" @click="showPicker = !showPicker" :disabled="disabled">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
          <polyline points="10 9 9 9 8 9"/>
        </svg>
        <span v-if="selectedTemplate">{{ selectedTemplate.name }}</span>
        <span v-else>选择模板</span>
        <svg class="chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>
    </div>

    <div v-if="showPicker" class="picker-overlay" @click.self="showPicker = false">
      <div class="picker-dialog">
        <div class="picker-header">
          <h3>📝 选择跟进模板</h3>
          <button class="close-btn" @click="showPicker = false">✕</button>
        </div>
        <div class="picker-content">
          <div v-if="loading" class="loading">加载中...</div>
          <div v-else-if="templates.length === 0" class="empty">暂无模板</div>
          <div v-else class="template-list">
            <div
              v-for="template in templates"
              :key="template.id"
              class="template-item"
              :class="{ selected: selectedId === template.id, default: template.is_default }"
              @click="selectTemplate(template)"
            >
              <div class="template-header">
                <span class="template-name">{{ template.name }}</span>
                <span v-if="template.is_default" class="default-badge">默认</span>
                <span v-if="template.category" class="category-badge">{{ template.category }}</span>
              </div>
              <div class="template-preview">{{ getPreview(template.content) }}</div>
              <div class="template-footer">
                <span class="use-count">使用 {{ template.use_count }} 次</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTemplatesAPI, useTemplateAPI } from '../api'

const props = defineProps({
  modelValue: Object,
  disabled: Boolean
})

const emit = defineEmits(['update:modelValue'])

const showPicker = ref(false)
const templates = ref([])
const loading = ref(false)
const selectedId = ref(null)

const selectedTemplate = computed(() => props.modelValue)

const loadTemplates = async () => {
  loading.value = true
  try {
    const res = await getTemplatesAPI()
    if (res.code === 200) {
      templates.value = res.data || []
      if (templates.value.length > 0 && !selectedId.value) {
        const defaultTemplate = templates.value.find(t => t.is_default)
        if (defaultTemplate) {
          selectTemplate(defaultTemplate)
        }
      }
    }
  } catch (error) {
    console.error('加载模板失败', error)
  } finally {
    loading.value = false
  }
}

const selectTemplate = async (template) => {
  selectedId.value = template.id
  emit('update:modelValue', template)
  showPicker.value = false

  try {
    await useTemplateAPI(template.id)
    template.use_count++
  } catch (error) {
    console.error('记录模板使用失败', error)
  }
}

const getPreview = (content) => {
  if (!content) return ''
  const lines = content.split('\n').filter(l => l.trim())
  return lines.slice(0, 2).join('\n').substring(0, 100) + (content.length > 100 ? '...' : '')
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.template-selector {
  position: relative;
}

.selector-trigger {
  display: inline-block;
}

.trigger-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.trigger-btn:hover:not(:disabled) {
  background: var(--border-color);
  border-color: var(--primary);
  color: var(--primary);
}

.trigger-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.trigger-btn svg {
  width: 16px;
  height: 16px;
}

.trigger-btn .chevron {
  width: 14px;
  height: 14px;
  margin-left: 4px;
}

.picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.picker-dialog {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.3s ease;
}

.picker-header {
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-radius: 16px 16px 0 0;
}

.picker-header h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.picker-content {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #9ca3af;
  font-size: 14px;
}

.template-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.template-item {
  padding: 16px;
  background: var(--bg-card-hover);
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.template-item:hover {
  background: var(--bg-hover);
  border-color: #d1d5db;
  transform: translateY(-2px);
}

.template-item.selected {
  background: #eff6ff;
  border-color: var(--primary);
}

.template-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.template-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.template-item.selected .template-name {
  color: var(--primary);
}

.default-badge {
  font-size: 10px;
  padding: 2px 8px;
  background: #fef3c7;
  color: var(--warning);
  border-radius: 4px;
}

.category-badge {
  font-size: 10px;
  padding: 2px 8px;
  background: var(--border-color);
  color: var(--text-tertiary);
  border-radius: 4px;
  margin-left: auto;
}

.template-preview {
  font-size: 13px;
  color: var(--text-tertiary);
  line-height: 1.5;
  white-space: pre-wrap;
  overflow: hidden;
}

.template-footer {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #e5e7eb;
}

.use-count {
  font-size: 11px;
  color: #9ca3af;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
