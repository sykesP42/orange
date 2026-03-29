<template>
  <div class="avatar-cropper-dialog" v-if="visible" @click.self="handleCloseClick">
    <div class="avatar-cropper-content">
      <div class="avatar-cropper-header">
        <h3>裁剪头像</h3>
        <button class="close-btn" @click="handleCloseClick" title="关闭 (Esc)">×</button>
      </div>
      <div class="avatar-cropper-body">
        <div class="crop-container" ref="cropContainer">
          <canvas ref="cropCanvas" class="crop-canvas"></canvas>
          <div class="crop-hint" v-if="!isDragging">拖动图片调整位置</div>
        </div>

        <div class="preview-section">
          <div class="preview-label">预览效果</div>
          <div class="preview-box">
            <canvas ref="previewCanvas" class="preview-canvas"></canvas>
          </div>
        </div>

        <div class="crop-controls">
          <div class="control-group">
            <div class="control-header">
              <label>缩放</label>
              <span class="control-value">{{ Math.round(scale * 100) }}%</span>
            </div>
            <input
              type="range"
              v-model.number="scale"
              :min="0.1"
              :max="3"
              :step="0.05"
              @input="handleScaleChange"
            />
          </div>
          <div class="control-group">
            <div class="control-header">
              <label>旋转</label>
              <span class="control-value">{{ rotation }}°</span>
            </div>
            <input
              type="range"
              v-model.number="rotation"
              :min="0"
              :max="360"
              :step="1"
              @input="handleRotationChange"
            />
          </div>
        </div>

        <div class="quick-actions">
          <button class="quick-btn" @click="resetCrop" title="重置">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
              <path d="M3 3v5h5"/>
            </svg>
            重置
          </button>
          <button class="quick-btn" @click="rotateLeft" title="向左旋转">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M2.5 2v6h6M2.66 15.57a10 10 0 1 0 .57-8.38"/>
            </svg>
            左旋
          </button>
          <button class="quick-btn" @click="rotateRight" title="向右旋转">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38"/>
            </svg>
            右旋
          </button>
        </div>
      </div>
      <div class="avatar-cropper-actions">
        <button class="btn-secondary" @click="handleCloseClick">取消</button>
        <button class="btn-primary" @click="applyCrop" :disabled="uploading">
          <span v-if="uploading" class="loading-spinner"></span>
          {{ uploading ? '上传中...' : '确认裁剪' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { uploadBloggerAvatarAPI } from '../api'
import { ElMessage } from 'element-plus'
import { useConfirm } from '../utils/confirm'

const { confirm } = useConfirm()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  imageFile: {
    type: Object,
    default: null
  },
  uploadAPI: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'success'])

const cropCanvas = ref(null)
const previewCanvas = ref(null)
const cropContainer = ref(null)
const originalImage = ref(null)
const scale = ref(1)
const rotation = ref(0)
const cropOffset = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const isDirty = ref(false)
const uploading = ref(false)
const dragStart = ref({ x: 0, y: 0 })
const cropSize = 200

watch(() => props.visible, (newVal) => {
  if (newVal && props.imageFile) {
    isDirty.value = false
    loadImage()
  }
})

const handleKeydown = (e) => {
  if (!props.visible) return
  if (e.key === 'Escape') {
    handleCloseClick()
  } else if (e.key === 'Enter' && !uploading.value) {
    applyCrop()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

const loadImage = () => {
  const reader = new FileReader()
  reader.onload = (e) => {
    const img = new Image()
    img.onload = () => {
      originalImage.value = img
      nextTick(() => {
        initCropCanvas()
      })
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(props.imageFile)
}

const initCropCanvas = () => {
  const canvas = cropCanvas.value
  if (!canvas) return

  canvas.width = cropSize
  canvas.height = cropSize

  const container = cropContainer.value
  if (container) {
    container.style.width = cropSize + 'px'
    container.style.height = cropSize + 'px'
  }

  scale.value = 1
  rotation.value = 0
  cropOffset.value = { x: 0, y: 0 }

  drawCropImage()
  updatePreview()

  canvas.addEventListener('mousedown', startDrag)
  canvas.addEventListener('mousemove', drag)
  canvas.addEventListener('mouseup', endDrag)
  canvas.addEventListener('mouseleave', endDrag)

  canvas.addEventListener('touchstart', startDragTouch, { passive: false })
  canvas.addEventListener('touchmove', dragTouch, { passive: false })
  canvas.addEventListener('touchend', endDrag)
}

const drawCropImage = () => {
  const canvas = cropCanvas.value
  if (!canvas || !originalImage.value) return

  const ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, cropSize, cropSize)

  ctx.fillStyle = '#f3f4f6'
  ctx.fillRect(0, 0, cropSize, cropSize)

  ctx.save()
  ctx.translate(cropSize / 2, cropSize / 2)
  ctx.rotate(rotation.value * Math.PI / 180)
  ctx.scale(scale.value, scale.value)

  const img = originalImage.value
  const drawWidth = img.width
  const drawHeight = img.height

  ctx.drawImage(
    img,
    -drawWidth / 2 + cropOffset.value.x,
    -drawHeight / 2 + cropOffset.value.y,
    drawWidth,
    drawHeight
  )

  ctx.restore()

  ctx.strokeStyle = '#3b82f6'
  ctx.lineWidth = 2
  ctx.setLineDash([])
  ctx.strokeRect(0, 0, cropSize, cropSize)

  ctx.strokeStyle = 'rgba(59, 130, 246, 0.3)'
  ctx.lineWidth = 1
  ctx.setLineDash([5, 5])
  ctx.beginPath()
  ctx.moveTo(cropSize / 3, 0)
  ctx.lineTo(cropSize / 3, cropSize)
  ctx.moveTo(cropSize * 2 / 3, 0)
  ctx.lineTo(cropSize * 2 / 3, cropSize)
  ctx.moveTo(0, cropSize / 3)
  ctx.lineTo(cropSize, cropSize / 3)
  ctx.moveTo(0, cropSize * 2 / 3)
  ctx.lineTo(cropSize, cropSize * 2 / 3)
  ctx.stroke()
}

const updatePreview = () => {
  const canvas = cropCanvas.value
  const preview = previewCanvas.value
  if (!canvas || !preview) return

  preview.width = 80
  preview.height = 80
  const ctx = preview.getContext('2d')
  ctx.drawImage(canvas, 0, 0, 80, 80)
}

const handleScaleChange = () => {
  isDirty.value = true
  drawCropImage()
  updatePreview()
}

const handleRotationChange = () => {
  isDirty.value = true
  drawCropImage()
  updatePreview()
}

const startDrag = (e) => {
  isDragging.value = true
  dragStart.value = {
    x: e.clientX - cropOffset.value.x,
    y: e.clientY - cropOffset.value.y
  }
}

const startDragTouch = (e) => {
  if (e.touches.length !== 1) return
  e.preventDefault()
  isDragging.value = true
  dragStart.value = {
    x: e.touches[0].clientX - cropOffset.value.x,
    y: e.touches[0].clientY - cropOffset.value.y
  }
}

const drag = (e) => {
  if (!isDragging.value) return
  isDirty.value = true
  cropOffset.value = {
    x: e.clientX - dragStart.value.x,
    y: e.clientY - dragStart.value.y
  }
  drawCropImage()
}

const dragTouch = (e) => {
  if (!isDragging.value || e.touches.length !== 1) return
  e.preventDefault()
  isDirty.value = true
  cropOffset.value = {
    x: e.touches[0].clientX - dragStart.value.x,
    y: e.touches[0].clientY - dragStart.value.y
  }
  drawCropImage()
}

const endDrag = () => {
  isDragging.value = false
  updatePreview()
}

const resetCrop = () => {
  scale.value = 1
  rotation.value = 0
  cropOffset.value = { x: 0, y: 0 }
  isDirty.value = true
  drawCropImage()
  updatePreview()
}

const rotateLeft = () => {
  rotation.value = (rotation.value - 90 + 360) % 360
  isDirty.value = true
  drawCropImage()
  updatePreview()
}

const rotateRight = () => {
  rotation.value = (rotation.value + 90) % 360
  isDirty.value = true
  drawCropImage()
  updatePreview()
}

const handleCloseClick = async () => {
  if (isDirty.value) {
    if (await confirm('确定要放弃当前的裁剪结果吗？')) {
      close()
    }
  } else {
    close()
  }
}

const applyCrop = async () => {
  const canvas = cropCanvas.value
  if (!canvas) return

  uploading.value = true

  const finalCanvas = document.createElement('canvas')
  finalCanvas.width = 200
  finalCanvas.height = 200
  const finalCtx = finalCanvas.getContext('2d')
  finalCtx.drawImage(canvas, 0, 0, 200, 200)

  finalCanvas.toBlob(async (blob) => {
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const uploadFunction = props.uploadAPI || uploadBloggerAvatarAPI
        const res = await uploadFunction({ image: e.target.result })
        if (res.code === 200) {
          emit('success', res.data.url)
          ElMessage.success('头像上传成功')
          close()
        } else {
          ElMessage.error(res.message || '上传失败')
        }
      } catch (error) {
        ElMessage.error('上传失败')
      } finally {
        uploading.value = false
      }
    }
    reader.readAsDataURL(blob)
  }, 'image/jpeg', 0.9)
}

const close = () => {
  emit('update:visible', false)
  originalImage.value = null
  scale.value = 1
  rotation.value = 0
  cropOffset.value = { x: 0, y: 0 }
  isDirty.value = false
}
</script>

<style scoped>
.avatar-cropper-dialog {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.avatar-cropper-content {
  background: var(--bg-card);
  border-radius: 20px;
  width: 90%;
  max-width: 420px;
  max-height: 95vh;
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.avatar-cropper-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.avatar-cropper-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.avatar-cropper-header .close-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: var(--bg-page);
  border-radius: 10px;
  font-size: 26px;
  line-height: 1;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.avatar-cropper-header .close-btn:hover {
  background: #fee2e2;
  color: #ef4444;
}

.avatar-cropper-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.crop-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-page);
  border-radius: 16px;
  padding: 16px;
}

.crop-canvas {
  border-radius: 12px;
  cursor: grab;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.crop-canvas:active {
  cursor: grabbing;
}

.crop-hint {
  position: absolute;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.9);
  padding: 4px 10px;
  border-radius: 10px;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s;
}

.crop-container:hover .crop-hint {
  opacity: 1;
}

.preview-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.preview-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.preview-box {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 3px solid white;
}

.preview-canvas {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.crop-controls {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.control-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.control-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.control-value {
  font-size: 13px;
  font-weight: 600;
  color: var(--primary);
  min-width: 50px;
  text-align: right;
}

.control-group input[type="range"] {
  width: 100%;
  height: 8px;
  border-radius: 4px;
  background: linear-gradient(to right, var(--primary-light), var(--primary));
  outline: none;
  -webkit-appearance: none;
  cursor: pointer;
}

.control-group input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  transition: all 0.2s;
}

.control-group input[type="range"]::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.quick-actions {
  display: flex;
  gap: 12px;
  width: 100%;
}

.quick-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 12px;
  background: var(--bg-page);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: var(--border-color);
  color: var(--text-primary);
}

.quick-btn svg {
  width: 16px;
  height: 16px;
}

.avatar-cropper-actions {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-page);
}

.btn-secondary,
.btn-primary {
  flex: 1;
  padding: 14px 24px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-secondary {
  background: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.btn-secondary:hover {
  background: var(--bg-card);
  color: var(--text-primary);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border: none;
  color: white;
  box-shadow: 0 4px 14px rgba(59, 130, 246, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>