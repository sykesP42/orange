<template>
  <div class="add-page">
    <div class="page-header">
      <h1>录入博主</h1>
      <p>谁先录入，归属谁，不可重复</p>
    </div>

    <form @submit.prevent="handleSubmit" class="form-card">
      <!-- 头像上传 -->
      <div class="form-group">
        <label>博主头像</label>
        <div class="avatar-upload">
          <div class="avatar-preview" @click="triggerFileInput">
            <img v-if="form.avatar" :src="form.avatar" alt="博主头像" v-avatar />
            <img v-else-if="userStore.avatar" :src="userStore.avatar" alt="对接人头像" class="fallback-avatar" v-avatar />
            <div v-else class="avatar-placeholder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              <span>点击上传头像</span>
            </div>
          </div>
          <input
            type="file"
            ref="fileInput"
            accept="image/*"
            style="display: none"
            @change="handleFileChange"
          />
          <button v-if="form.avatar" type="button" class="remove-avatar" @click="removeAvatar">
            移除头像
          </button>
        </div>
        <span class="hint avatar-hint">未添加头像则使用对接人头像</span>
      </div>

      <div class="form-group">
        <label>博主昵称 <span class="required">*</span></label>
        <div class="input-wrapper">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <input
            v-model="form.nickname"
            type="text"
            placeholder="输入博主昵称"
            required
          />
        </div>
        <span class="hint">昵称唯一，录入后不可更改</span>
      </div>

      <div class="form-group">
        <label>博主分类 <span class="required">*</span></label>
        <div class="input-wrapper">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
          <select v-model="form.category" required>
            <option value="">选择分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.name">
              {{ cat.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label>对接产品 <span class="required">*</span></label>
        <div class="input-wrapper">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
          </svg>
          <select v-model="form.product" required>
            <option value="">选择产品</option>
            <option v-for="prod in products" :key="prod.id" :value="prod.name">
              {{ prod.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label>活跃平台 <span class="required">*</span></label>
        <div class="input-wrapper">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="2" y1="12" x2="22" y2="12"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          </svg>
          <select v-model="form.platform" required>
            <option value="">选择平台</option>
            <option v-for="plat in platforms" :key="plat.id" :value="plat.name">
              {{ plat.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label>平台账号 <span class="required">*</span></label>
        <div class="input-wrapper">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <input
            v-model="form.platform_account"
            type="text"
            placeholder="输入平台账号/ID"
            required
          />
        </div>
        <span class="hint">同一平台账号不可重复录入</span>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label>联系方式</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
            </svg>
            <input
              v-model="form.contact"
              type="text"
              placeholder="输入联系方式"
            />
          </div>
        </div>

        <div class="form-group">
          <label>微信号</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
            </svg>
            <input
              v-model="form.wechat"
              type="text"
              placeholder="输入微信号"
            />
          </div>
        </div>
        <div class="form-group">
          <label>自定义联系方式</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
            </svg>
            <input
              v-model="form.custom_contact"
              type="text"
              placeholder="输入其他联系方式"
            />
          </div>
        </div>
      </div>

      <div class="form-group">
        <label>简要介绍</label>
        <div class="input-wrapper textarea">
          <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14,2 14,8 20,8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <polyline points="10,9 9,9 8,9"/>
          </svg>
          <textarea
            v-model="form.description"
            placeholder="博主简介、合作风格、粉丝量等"
            rows="3"
          ></textarea>
        </div>
      </div>

      <div class="form-group">
        <label>标签 Tags</label>
        <div class="tags-input-wrapper">
          <div class="tags-list">
            <span v-for="(tag, index) in form.tags" :key="index" class="tag-item">
              {{ tag }}
              <button type="button" class="tag-remove" @click="removeTag(index)">×</button>
            </span>
          </div>
          <div class="tag-input-row">
            <input
              v-model="tagInput"
              type="text"
              placeholder="输入标签后按回车添加"
              @keydown.enter.prevent="addTag"
            />
            <button type="button" class="add-tag-btn" @click="addTag" :disabled="!tagInput.trim()">添加</button>
            <button type="button" class="recommend-tags-btn" @click="recommendTags" :disabled="!canRecommend">
              <span v-if="recommending">分析中...</span>
              <span v-else>✨ 智能推荐</span>
            </button>
          </div>
          <div v-if="recommendedTags.length > 0" class="recommended-tags">
            <span class="recommended-label">推荐标签：</span>
            <span
              v-for="tag in recommendedTags"
              :key="tag.tag_id"
              class="recommended-tag"
              :class="{ selected: form.tags.includes(tag.tag_name) }"
              @click="toggleRecommendTag(tag.tag_name)"
            >
              {{ tag.tag_name }}
            </span>
          </div>
        </div>
        <span class="hint">添加相关标签，便于搜索和归类</span>
      </div>

      <div class="form-actions">
        <button type="button" class="btn-secondary" @click="router.push('/')">
          取消
        </button>
        <button type="submit" class="btn-primary" :disabled="loading">
          <span v-if="!loading">提交录入</span>
          <span v-else class="loading">
            <svg class="spinner" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" fill="none" stroke-dasharray="31.4 31.4"/>
            </svg>
            提交中...
          </span>
        </button>
      </div>
    </form>

    <AvatarCropper
      v-model:visible="cropDialogVisible"
      :image-file="selectedImageFile"
      @success="handleAvatarUploadSuccess"
    />

    <div class="info-card">
      <div class="info-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="16" x2="12" y2="12"/>
          <line x1="12" y1="8" x2="12.01" y2="8"/>
        </svg>
      </div>
      <div class="info-content">
        <h3>录入规则</h3>
        <ul>
          <li>博主昵称全局唯一，不可重复</li>
          <li>第一个提交成功的用户成为永久归属人</li>
          <li>普通用户不可编辑、删除他人录入的数据</li>
          <li>管理员可修正数据但不可更改归属人</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { categoryListAPI, productListAPI, bloggerAddAPI, platformListAPI, getRecommendTagsAPI } from '../api'
import AvatarCropper from '../components/AvatarCropper.vue'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()

const loading = ref(false)
const categories = ref([])
const products = ref([])
const platforms = ref([])
const tagInput = ref('')
const recommendedTags = ref([])
const recommending = ref(false)

const form = reactive({
  nickname: '',
  category: '',
  product: '',
  platform: '',
  platform_account: '',
  contact: '',
  wechat: '',
  custom_contact: '',
  description: '',
  tags: [],
  avatar: ''
})

const fileInput = ref(null)
const cropDialogVisible = ref(false)
const selectedImageFile = ref(null)

const canRecommend = computed(() => {
  return form.category || form.platform || form.description
})

const recommendTags = async () => {
  recommending.value = true
  try {
    const params = {
      category: form.category,
      platform: form.platform,
      product: form.product,
      description: form.description
    }
    const res = await getRecommendTagsAPI(params)
    if (res.code === 200 && res.data && res.data.recommendations) {
      recommendedTags.value = res.data.recommendations
      if (recommendedTags.value.length === 0) {
        notification.info('暂未找到相关推荐标签，请手动添加')
      }
    }
  } catch (error) {
    console.error('推荐标签失败', error)
  } finally {
    recommending.value = false
  }
}

const toggleRecommendTag = (tagName) => {
  if (form.tags.includes(tagName)) {
    form.tags = form.tags.filter(t => t !== tagName)
  } else {
    form.tags.push(tagName)
  }
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    notification.warning('文件大小不能超过10MB')
    return
  }

  selectedImageFile.value = file
  cropDialogVisible.value = true

  event.target.value = ''
}

const handleAvatarUploadSuccess = (url) => {
  form.avatar = url
}

const removeAvatar = () => {
  form.avatar = ''
}

const addTag = () => {
  const tag = tagInput.value.trim()
  if (tag && !form.tags.includes(tag)) {
    form.tags.push(tag)
    tagInput.value = ''
  }
}

const removeTag = (index) => {
  form.tags.splice(index, 1)
}

const loadCategories = async () => {
  try {
    const res = await categoryListAPI()
    if (res.code === 200) {
      categories.value = res.data || []
    }
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

const loadProducts = async () => {
  try {
    const res = await productListAPI()
    if (res.code === 200) {
      products.value = res.data || []
    }
  } catch (error) {
    console.error('加载产品失败', error)
  }
}

const loadPlatforms = async () => {
  try {
    const res = await platformListAPI()
    if (res.code === 200) {
      platforms.value = res.data || []
    }
  } catch (error) {
    console.error('加载平台失败', error)
  }
}

const handleSubmit = async () => {
  if (!form.nickname || !form.category || !form.product || !form.platform || !form.platform_account) {
    notification.warning('请填写必填项')
    return
  }

  loading.value = true
  try {
    const res = await bloggerAddAPI(form)
    if (res.code === 200) {
      const finalNickname = res.data?.nickname || form.nickname
      if (finalNickname !== form.nickname) {
        notification.success(`录入成功！由于昵称重复，系统已自动编号为：${finalNickname}`)
      } else {
        notification.success('录入成功！该博主已锁定为您归属')
      }
      router.push('/')
    } else if (res.code === 4003) {
      notification.error(res.message || '半个月内添加博主已达上限')
    } else {
      notification.error(res.message || '录入失败')
    }
  } catch (error) {
    const msg = error.response?.data?.message || '录入失败，请重试'
    notification.error(msg)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadProducts()
  loadPlatforms()
})
</script>

<style scoped>
.add-page {
  max-width: 640px;
  margin: 0 auto;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-header p {
  font-size: 14px;
  color: var(--text-muted);
}

.form-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-row .form-group {
  margin-bottom: 0;
}

.form-group:last-of-type {
  margin-bottom: 32px;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
}

.avatar-preview {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--bg-dark);
  border: 2px dashed var(--border-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.avatar-preview:hover {
  border-color: var(--primary);
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
}

.avatar-placeholder svg {
  width: 32px;
  height: 32px;
}

.avatar-placeholder span {
  font-size: 12px;
}

.remove-avatar {
  padding: 6px 12px;
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.remove-avatar:hover {
  background: #ef4444;
  color: white;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.required {
  color: #ef4444;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  width: 18px;
  height: 18px;
  color: var(--text-muted);
  pointer-events: none;
}

.input-wrapper input,
.input-wrapper select,
.input-wrapper.textarea textarea {
  width: 100%;
  height: 50px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0 16px 0 48px;
  font-size: 15px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.input-wrapper.textarea textarea {
  height: auto;
  padding-top: 14px;
  resize: vertical;
  min-height: 100px;
}

.input-wrapper select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 16px center;
}

.input-wrapper input::placeholder,
.input-wrapper textarea::placeholder {
  color: var(--text-muted);
}

.input-wrapper input:focus,
.input-wrapper select:focus,
.input-wrapper textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.15);
}

.hint {
  display: block;
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

.avatar-hint {
  font-style: italic;
  font-size: 11px;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-secondary,
.btn-primary {
  padding: 14px 28px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
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
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  border: none;
  color: white;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(249, 115, 22, 0.4);
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 18px;
  height: 18px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.loading {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-card {
  display: flex;
  gap: 16px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 20px;
}

.info-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  background: rgba(59, 130, 246, 0.2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #60a5fa;
}

.info-icon svg {
  width: 20px;
  height: 20px;
}

.info-content h3 {
  font-size: 14px;
  font-weight: 600;
  color: #60a5fa;
  margin-bottom: 10px;
}

.info-content ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.info-content li {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  padding-left: 16px;
  position: relative;
}

.info-content li::before {
  content: '•';
  position: absolute;
  left: 0;
  color: #60a5fa;
}

.tags-input-wrapper {
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 12px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 20px;
  font-size: 13px;
  color: white;
  font-weight: 500;
}

.tag-remove {
  background: none;
  border: none;
  color: white;
  font-size: 16px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  opacity: 0.8;
}

.tag-remove:hover {
  opacity: 1;
}

.tag-input-row {
  display: flex;
  gap: 8px;
}

.tag-input-row input {
  flex: 1;
  height: 40px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 0 14px;
  font-size: 14px;
  color: var(--text-primary);
}

.tag-input-row input:focus {
  outline: none;
  border-color: var(--primary);
}

.add-tag-btn {
  padding: 0 16px;
  background: var(--primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-tag-btn:hover:not(:disabled) {
  background: var(--primary-dark);
}

.add-tag-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .form-card {
    padding: 24px 20px;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-secondary,
  .btn-primary {
    width: 100%;
  }
}

.recommend-tags-btn {
  padding: 0 16px;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.recommend-tags-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #d97706 0%, #b45309 100%);
}

.recommend-tags-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.recommended-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
  padding: 12px;
  background: rgba(245, 158, 11, 0.1);
  border-radius: 8px;
  border: 1px dashed rgba(245, 158, 11, 0.3);
}

.recommended-label {
  width: 100%;
  font-size: 12px;
  color: #92400e;
  margin-bottom: 4px;
}

.recommended-tag {
  display: inline-block;
  padding: 6px 14px;
  background: white;
  border: 1px solid #f59e0b;
  border-radius: 20px;
  font-size: 13px;
  color: #d97706;
  cursor: pointer;
  transition: all 0.2s;
}

.recommended-tag:hover {
  background: #fef3c7;
}

.recommended-tag.selected {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  border-color: #f59e0b;
}
</style>