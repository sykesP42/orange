<template>
  <img
    v-if="isVisible || !lazy"
    :src="src"
    :alt="alt"
    :style="imageStyle"
    @load="handleLoad"
    @error="handleError"
  />
  <div v-else class="lazy-placeholder" :style="placeholderStyle">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
      <circle cx="8.5" cy="8.5" r="1.5"/>
      <polyline points="21 15 16 10 5 21"/>
    </svg>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  src: { type: String, required: true },
  alt: { type: String, default: '' },
  lazy: { type: Boolean, default: true },
  placeholderColor: { type: String, default: '#f0f0f0' }
})

const emit = defineEmits(['load', 'error'])

const isVisible = ref(!props.lazy)
const isLoaded = ref(false)
const hasError = ref(false)
const observer = ref(null)

const imageStyle = computed(() => ({
  opacity: isLoaded.value ? 1 : 0,
  transition: 'opacity 0.3s ease'
}))

const placeholderStyle = computed(() => ({
  backgroundColor: props.placeholderColor
}))

const handleLoad = (e) => {
  isLoaded.value = true
  emit('load', e)
}

const handleError = (e) => {
  hasError.value = true
  emit('error', e)
}

onMounted(() => {
  if (props.lazy && 'IntersectionObserver' in window) {
    observer.value = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            isVisible.value = true
            if (observer.value) {
              observer.value.disconnect()
            }
          }
        })
      },
      { rootMargin: '50px' }
    )

    const element = document.querySelector(`[data-lazy-id="${props.src}"]`)
    if (element) {
      observer.value.observe(element)
    }
  }
})

onUnmounted(() => {
  if (observer.value) {
    observer.value.disconnect()
  }
})
</script>

<style scoped>
.lazy-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100px;
}

.lazy-placeholder svg {
  width: 32px;
  height: 32px;
  color: #ccc;
}
</style>
