<template>
  <div class="import-panel">
    <div class="import-steps">
      <div
        v-for="(step, index) in steps"
        :key="index"
        class="step"
        :class="{ active: currentStep === index, completed: currentStep > index }"
      >
        <div class="step-number">
          <span v-if="currentStep <= index">{{ index + 1 }}</span>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <polyline points="20,6 9,17 4,12"/>
          </svg>
        </div>
        <span class="step-name">{{ step }}</span>
      </div>
    </div>

    <div class="step-content">
      <div v-if="currentStep === 0" class="step-upload">
        <div
          class="upload-zone"
          :class="{ dragging: isDragging }"
          @dragover.prevent="isDragging = true"
          @dragleave="isDragging = false"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <input
            type="file"
            ref="fileInput"
            accept=".csv,.json"
            style="display: none"
            @change="handleFileSelect"
          />
          <div class="upload-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17,8 12,3 7,8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
          </div>
          <h3>点击或拖拽上传文件</h3>
          <p>支持 CSV 或 JSON 格式</p>
        </div>

        <div class="template-section">
          <h4>下载导入模板</h4>
          <div class="template-buttons">
            <button class="template-btn" @click="downloadTemplate('csv')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
              CSV 模板
            </button>
            <button class="template-btn" @click="downloadTemplate('json')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
              JSON 模板
            </button>
          </div>
        </div>
      </div>

      <div v-if="currentStep === 1" class="step-preview">
        <div class="preview-stats">
          <div class="stat-card success">
            <span class="stat-number">{{ previewData.validRows || 0 }}</span>
            <span class="stat-label">有效数据</span>
          </div>
          <div class="stat-card error">
            <span class="stat-number">{{ previewData.invalidRows || 0 }}</span>
            <span class="stat-label">无效数据</span>
          </div>
          <div class="stat-card total">
            <span class="stat-number">{{ previewData.totalRows || 0 }}</span>
            <span class="stat-label">总行数</span>
          </div>
        </div>

        <div v-if="previewData.errors && previewData.errors.length > 0" class="error-list">
          <h4>数据问题</h4>
          <div class="error-items">
            <div v-for="(err, idx) in previewData.errors.slice(0, 10)" :key="idx" class="error-item">
              <span class="error-row">行 {{ err.row }}</span>
              <span class="error-field">{{ err.field }}</span>
              <span class="error-msg">{{ err.error }}</span>
            </div>
            <div v-if="previewData.errors.length > 10" class="more-errors">
              还有 {{ previewData.errors.length - 10 }} 个错误...
            </div>
          </div>
        </div>

        <div class="preview-table">
          <h4>数据预览（前10条）</h4>
          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>昵称</th>
                  <th>分类</th>
                  <th>平台</th>
                  <th>平台账号</th>
                  <th>联系方式</th>
                  <th>微信</th>
                  <th>产品</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, idx) in previewData.previewData" :key="idx">
                  <td>{{ row.nickname || '-' }}</td>
                  <td>{{ row.category || '-' }}</td>
                  <td>{{ row.platform || '-' }}</td>
                  <td>{{ row.platform_account || '-' }}</td>
                  <td>{{ row.contact || '-' }}</td>
                  <td>{{ row.wechat || '-' }}</td>
                  <td>{{ row.product || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="preview-actions">
          <button class="btn-secondary" @click="resetImport">重新选择文件</button>
          <button class="btn-primary" @click="startImport" :disabled="previewData.validRows === 0">
            确认导入 {{ previewData.validRows || 0 }} 条数据
          </button>
        </div>
      </div>

      <div v-if="currentStep === 2" class="step-result">
        <div v-if="importing" class="importing">
          <div class="spinner"></div>
          <p>正在导入数据...</p>
        </div>
        <div v-else class="result">
          <div class="result-icon success" v-if="importResult.failCount === 0">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polyline points="16,8 10,14 8,12"/>
            </svg>
          </div>
          <div class="result-icon partial" v-else>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>

          <h3>导入完成</h3>
          <div class="result-stats">
            <div class="result-stat">
              <span class="num success">{{ importResult.successCount }}</span>
              <span class="label">成功</span>
            </div>
            <div class="result-stat">
              <span class="num error">{{ importResult.failCount }}</span>
              <span class="label">失败</span>
            </div>
          </div>

          <div v-if="importResult.errors && importResult.errors.length > 0" class="fail-list">
            <h4>失败记录</h4>
            <div class="fail-items">
              <div v-for="(err, idx) in importResult.errors.slice(0, 10)" :key="idx" class="fail-item">
                <span>行 {{ err.row }}: {{ err.error }}</span>
              </div>
            </div>
            <button class="btn-sm" @click="exportErrors">导出错误列表</button>
          </div>

          <div class="result-actions">
            <button class="btn-secondary" @click="resetImport">继续导入</button>
            <button class="btn-primary" @click="finishImport">完成</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { previewImportAPI, doImportAPI, downloadTemplateAPI, exportFailedRowsAPI } from '../api'

const emit = defineEmits(['complete'])

const steps = ['上传文件', '预览确认', '导入完成']
const currentStep = ref(0)
const isDragging = ref(false)
const fileInput = ref(null)
const selectedFile = ref(null)
const selectedFilename = ref('')

const previewData = ref({})
const importResult = ref({})
const importing = ref(false)

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (e) => {
  const file = e.target.files[0]
  if (file) {
    processFile(file)
  }
}

const handleDrop = (e) => {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (file) {
    processFile(file)
  }
}

const processFile = async (file) => {
  selectedFile.value = file
  selectedFilename.value = file.name

  const formData = new FormData()
  formData.append('file', file)
  formData.append('filename', file.name)

  try {
    const res = await previewImportAPI(formData)
    if (res.code === 200) {
      previewData.value = res.data
      currentStep.value = 1
    }
  } catch (error) {
    console.error('预览失败:', error)
  }
}

const downloadTemplate = async (format) => {
  try {
    const blob = await downloadTemplateAPI(format)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `blogger_import_template.${format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('下载模板失败:', error)
  }
}

const startImport = async () => {
  importing.value = true
  try {
    const res = await doImportAPI({
      filename: selectedFilename.value,
      data: previewData.value.previewData
    })
    if (res.code === 200) {
      importResult.value = res.data
      emit('complete', res.data)
    }
  } catch (error) {
    console.error('导入失败:', error)
  } finally {
    importing.value = false
    currentStep.value = 2
  }
}

const resetImport = () => {
  currentStep.value = 0
  selectedFile.value = null
  selectedFilename.value = ''
  previewData.value = {}
  importResult.value = {}
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const finishImport = () => {
  resetImport()
}

const exportErrors = async () => {
  if (importResult.value.errors && importResult.value.errors.length > 0) {
    try {
      const blob = await exportFailedRowsAPI(importResult.value.errors)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'import_errors.csv'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('导出错误失败:', error)
    }
  }
}
</script>

<style scoped>
.import-panel {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 24px;
}

.import-steps {
  display: flex;
  justify-content: center;
  gap: 48px;
  margin-bottom: 32px;
}

.step {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-muted);
}

.step.active {
  color: var(--primary);
}

.step.completed {
  color: #22c55e;
}

.step-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--bg-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
}

.step.active .step-number {
  background: var(--primary);
  color: white;
}

.step.completed .step-number {
  background: #22c55e;
  color: white;
}

.step-number svg {
  width: 18px;
  height: 18px;
}

.step-name {
  font-size: 14px;
  font-weight: 500;
}

.step-content {
  min-height: 300px;
}

.upload-zone {
  border: 2px dashed var(--border-color);
  border-radius: 16px;
  padding: 60px 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 32px;
}

.upload-zone:hover,
.upload-zone.dragging {
  border-color: var(--primary);
  background: rgba(249, 115, 22, 0.05);
}

.upload-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-icon svg {
  width: 32px;
  height: 32px;
  color: white;
}

.upload-zone h3 {
  margin: 0 0 8px;
  font-size: 18px;
  color: var(--text-primary);
}

.upload-zone p {
  margin: 0;
  font-size: 14px;
  color: var(--text-muted);
}

.template-section {
  text-align: center;
}

.template-section h4 {
  margin: 0 0 16px;
  font-size: 14px;
  color: var(--text-secondary);
}

.template-buttons {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.template-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.template-btn:hover {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.template-btn svg {
  width: 18px;
  height: 18px;
}

.preview-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  flex: 1;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
}

.stat-card.success {
  background: rgba(34, 197, 94, 0.1);
}

.stat-card.error {
  background: rgba(239, 68, 68, 0.1);
}

.stat-card.total {
  background: rgba(59, 130, 246, 0.1);
}

.stat-number {
  display: block;
  font-size: 28px;
  font-weight: 700;
}

.stat-card.success .stat-number {
  color: #22c55e;
}

.stat-card.error .stat-number {
  color: #ef4444;
}

.stat-card.total .stat-number {
  color: #60a5fa;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
}

.error-list {
  background: rgba(239, 68, 68, 0.05);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 24px;
}

.error-list h4 {
  margin: 0 0 12px;
  font-size: 14px;
  color: #ef4444;
}

.error-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.error-item {
  display: flex;
  gap: 12px;
  font-size: 13px;
  padding: 8px 12px;
  background: var(--bg-card);
  border-radius: 6px;
}

.error-row {
  color: #ef4444;
  font-weight: 600;
  min-width: 60px;
}

.error-field {
  color: var(--text-secondary);
  min-width: 80px;
}

.error-msg {
  color: var(--text-muted);
}

.more-errors {
  text-align: center;
  color: var(--text-muted);
  font-size: 13px;
  padding: 8px;
}

.preview-table {
  margin-bottom: 24px;
}

.preview-table h4 {
  margin: 0 0 12px;
  font-size: 14px;
  color: var(--text-secondary);
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

th, td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

th {
  background: var(--bg-dark);
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

td {
  color: var(--text-primary);
}

tr:last-child td {
  border-bottom: none;
}

.preview-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-secondary,
.btn-primary {
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-secondary {
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.btn-secondary:hover {
  background: var(--bg-card-hover);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border: none;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.importing {
  text-align: center;
  padding: 60px;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.importing p {
  font-size: 14px;
  color: var(--text-muted);
}

.result {
  text-align: center;
  padding: 40px;
}

.result-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24px;
}

.result-icon.success {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.result-icon.partial {
  background: rgba(249, 115, 22, 0.1);
  color: var(--primary);
}

.result-icon svg {
  width: 40px;
  height: 40px;
}

.result h3 {
  margin: 0 0 24px;
  font-size: 20px;
  color: var(--text-primary);
}

.result-stats {
  display: flex;
  gap: 32px;
  justify-content: center;
  margin-bottom: 24px;
}

.result-stat {
  text-align: center;
}

.result-stat .num {
  display: block;
  font-size: 32px;
  font-weight: 700;
}

.result-stat .num.success {
  color: #22c55e;
}

.result-stat .num.error {
  color: #ef4444;
}

.result-stat .label {
  font-size: 13px;
  color: var(--text-muted);
}

.fail-list {
  background: var(--bg-dark);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 24px;
  text-align: left;
}

.fail-list h4 {
  margin: 0 0 12px;
  font-size: 14px;
  color: var(--text-secondary);
}

.fail-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.fail-item {
  font-size: 13px;
  color: var(--text-secondary);
  padding: 8px 12px;
  background: var(--bg-card);
  border-radius: 6px;
}

.btn-sm {
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
}

.btn-sm:hover {
  background: var(--bg-card-hover);
}

.result-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
