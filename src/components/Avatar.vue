<template>
  <div class="avatar-wrapper" :class="[sizeClass, { 'avatar-square': square }]">
    <img
      v-if="src && !hasError"
      :src="src"
      :alt="alt"
      class="avatar-img"
      :class="{ 'is-loaded': isLoaded }"
      @load="onLoad"
      @error="onError"
    />
    <div v-if="!src || hasError" class="avatar-fallback">
      <span v-if="fallbackText" class="avatar-initials">{{ fallbackText }}</span>
      <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
        <circle cx="12" cy="7" r="4"/>
      </svg>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  size: { type: String, default: 'md' },
  square: { type: Boolean, default: false },
  name: { type: String, default: '' }
})

const isLoaded = ref(false)
const hasError = ref(false)

const sizeClass = computed(() => `avatar-${props.size}`)

const fallbackText = computed(() => {
  if (!props.name) return ''
  const parts = props.name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  return parts[0].substring(0, 2).toUpperCase()
})

watch(() => props.src, () => {
  isLoaded.value = false
  hasError.value = false
})

function onLoad() {
  isLoaded.value = true
}

function onError() {
  hasError.value = true
}
</script>

<style scoped>
.avatar-wrapper {
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  flex-shrink: 0;
}

.avatar-wrapper:not(.avatar-square) {
  border-radius: 50%;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 0.35s ease;
}

.avatar-img.is-loaded {
  opacity: 1;
}

.avatar-fallback {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-initials {
  font-weight: 600;
  letter-spacing: 0.02em;
  color: var(--text-secondary);
}

.avatar-xs { width: 28px; height: 28px; }
.avatar-sm { width: 36px; height: 36px; }
.avatar-md { width: 44px; height: 44px; }
.avatar-lg { width: 56px; height: 56px; }
.avatar-xl { width: 72px; height: 72px; }
.avatar-2xl { width: 96px; height: 96px; }

.avatar-xs .avatar-initials { font-size: 10px; }
.avatar-sm .avatar-initials { font-size: 12px; }
.avatar-md .avatar-initials { font-size: 14px; }
.avatar-lg .avatar-initials { font-size: 16px; }
.avatar-xl .avatar-initials { font-size: 20px; }
.avatar-2xl .avatar-initials { font-size: 26px; }
</style>
