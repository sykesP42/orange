<template>
  <div ref="containerRef" class="virtual-list" :style="{ height: height + 'px', overflow: 'auto' }" @scroll="onScroll">
    <div class="virtual-spacer" :style="{ height: totalHeight + 'px' }">
      <div
        v-for="(item, index) in visibleItems"
        :key="keyField ? item[keyField] : index"
        :data-index="startIndex + index"
        :style="{ position: 'absolute', top: (startIndex + index) * itemHeight + 'px', left: 0, width: '100%' }"
      >
        <slot :item="item" :index="startIndex + index" />
      </div>
    </div>
    <div v-if="items.length === 0 && $slots.empty" class="virtual-empty">
      <slot name="empty" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  itemHeight: { type: Number, required: true },
  buffer: { type: Number, default: 5 },
  keyField: { type: String, default: '' },
  height: { type: Number, default: 600 },
})

const containerRef = ref(null)
const scrollTop = ref(0)

const totalHeight = computed(() => props.items.length * props.itemHeight)
const startIndex = computed(() => {
  const rawStart = Math.floor(scrollTop.value / props.itemHeight)
  return Math.max(0, rawStart - props.buffer)
})
const endIndex = computed(() => {
  const visibleCount = Math.ceil(props.height / props.itemHeight)
  return Math.min(
    props.items.length,
    startIndex.value + visibleCount + props.buffer * 2
  )
})

const visibleItems = computed(() =>
  props.items.slice(startIndex.value, endIndex.value)
)

const onScroll = () => {
  if (containerRef.value) {
    scrollTop.value = containerRef.value.scrollTop
  }
}

watch(() => props.items, () => {
  if (containerRef.value) {
    scrollTop.value = containerRef.value.scrollTop
  }
}, { deep: false })

defineExpose({ scrollToTop: () => { if (containerRef.value) containerRef.value.scrollTop = 0 } })
</script>

<style scoped>
.virtual-list { position: relative; }
.virtual-spacer { position: relative; }
.virtual-empty { padding: 40px 16px; text-align: center; }
</style>
