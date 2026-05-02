<template>
  <span v-if="!text || !keyword" class="highlight-text">{{ text }}</span>
  <span v-else class="highlight-text" v-html="highlightedText"></span>
</template>

<script setup>
import { computed } from 'vue'
import { highlightMatch } from '../utils/pinyin-search'

const props = defineProps({
  text: {
    type: String,
    default: ''
  },
  keyword: {
    type: String,
    default: ''
  },
  caseSensitive: {
    type: Boolean,
    default: false
  },
  className: {
    type: String,
    default: ''
  }
})

const highlightedText = computed(() => {
  if (!props.text || !props.keyword) {
    return props.text
  }
  
  return highlightMatch(props.text, props.keyword)
})
</script>

<style scoped>
.highlight-text {
  :deep(mark) {
    background: linear-gradient(120deg, #f97316 0%, #fb923c 100%);
    color: var(--color-on-brand);
    padding: 1px 4px;
    border-radius: 3px;
    font-weight: 600;
    box-shadow: 0 1px 2px rgba(249, 115, 22, 0.3);
  }
}

/* 暗色主题适配 */
:global([data-theme='dark']) .highlight-text {
  :deep(mark) {
    background: linear-gradient(120deg, #f97316 0%, #fb923c 100%);
    color: #1a1a1a;
  }
}
</style>
