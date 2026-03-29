<template>
  <div class="user-settings">
    <div class="page-header">
      <h1>个人设置</h1>
      <p>管理您的个人信息和账户安全</p>
    </div>

    <div class="settings-container">
      <!-- 头像设置 -->
      <div class="settings-card">
        <div class="card-header">
          <h3>头像设置</h3>
        </div>
        <div class="card-body">
          <div class="avatar-section">
            <div class="avatar-preview">
              <img v-if="profile.avatar" :src="profile.avatar" alt="头像" v-avatar />
              <div v-else class="avatar-placeholder">
                {{ profile.real_name?.charAt(0) || profile.username?.charAt(0) || '?' }}
              </div>
            </div>
            <div class="avatar-actions">
              <input
                type="file"
                ref="fileInput"
                accept="image/*"
                style="display: none"
                @change="handleFileChange"
              />
              <button class="btn btn-primary" @click="fileInput.click()">
                上传头像
              </button>
              <p class="hint">支持 JPG、PNG 格式，文件大小不超过 10MB</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 基本信息 -->
      <div class="settings-card">
        <div class="card-header">
          <h3>基本信息</h3>
        </div>
        <div class="card-body">
          <form @submit.prevent="updateProfile">
            <div class="form-group">
              <label>用户名</label>
              <input type="text" v-model="profile.username" disabled />
              <span class="hint">用户名不可修改</span>
            </div>
            <div class="form-group">
              <label>真实姓名</label>
              <input type="text" v-model="form.real_name" placeholder="请输入真实姓名" />
            </div>
            <div class="form-group">
              <label>邮箱</label>
              <input type="email" v-model="form.email" placeholder="请输入邮箱" />
            </div>
            <div class="form-group">
              <label>手机号</label>
              <input type="tel" v-model="form.phone" placeholder="请输入手机号" />
            </div>
            <div class="form-group">
              <label>个人简介</label>
              <textarea v-model="form.bio" rows="3" placeholder="请输入个人简介"></textarea>
            </div>
            <div class="form-group" v-if="profile.approved_by">
              <label>注册批准人</label>
              <input type="text" :value="profile.approved_by" disabled />
              <span class="hint">批准时间：{{ formatDateTime(profile.approved_time) }}</span>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="saving">
                {{ saving ? '保存中...' : '保存修改' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 修改密码 -->
      <div class="settings-card">
        <div class="card-header">
          <h3>修改密码</h3>
        </div>
        <div class="card-body">
          <form @submit.prevent="updatePassword">
            <div class="form-group">
              <label>当前密码</label>
              <input type="password" v-model="passwordForm.oldPassword" placeholder="请输入当前密码" />
            </div>
            <div class="form-group">
              <label>新密码</label>
              <input type="password" v-model="passwordForm.newPassword" placeholder="请输入新密码" />
            </div>
            <div class="form-group">
              <label>确认新密码</label>
              <input type="password" v-model="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="updatingPassword">
                {{ updatingPassword ? '修改中...' : '修改密码' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <AvatarCropper
      v-model:visible="cropDialogVisible"
      :image-file="selectedImageFile"
      :upload-api="uploadAvatarAPI"
      @success="handleAvatarUploadSuccess"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useNotification } from '../stores/notification'
import { getUserProfileAPI, updateUserProfileAPI, updatePasswordAPI, uploadAvatarAPI } from '../api'
import AvatarCropper from '../components/AvatarCropper.vue'

const { success, error: notifyError, warning } = useNotification()

const profile = ref({})
const form = ref({
  real_name: '',
  email: '',
  phone: '',
  bio: ''
})
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const saving = ref(false)
const updatingPassword = ref(false)
const fileInput = ref(null)
const cropDialogVisible = ref(false)
const selectedImageFile = ref(null)

const loadProfile = async () => {
  try {
    const res = await getUserProfileAPI()
    if (res.code === 200) {
      profile.value = res.data
      form.value = {
        real_name: res.data.real_name || '',
        email: res.data.email || '',
        phone: res.data.phone || '',
        bio: res.data.bio || ''
      }
    }
  } catch (error) {
    notifyError('加载用户信息失败')
  }
}

const updateProfile = async () => {
  saving.value = true
  try {
    const res = await updateUserProfileAPI(form.value)
    if (res.code === 200) {
      success('保存成功')
    } else {
      notifyError(res.message || '保存失败')
    }
  } catch (error) {
    notifyError('保存失败')
  } finally {
    saving.value = false
  }
}

const updatePassword = async () => {
  if (!passwordForm.value.oldPassword || !passwordForm.value.newPassword) {
    warning('请填写完整密码信息')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    warning('两次输入的新密码不一致')
    return
  }

  updatingPassword.value = true
  try {
    const res = await updatePasswordAPI({
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    })
    if (res.code === 200) {
      success('密码修改成功')
      passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
    } else {
      notifyError(res.message || '密码修改失败')
    }
  } catch (error) {
    notifyError('密码修改失败')
  } finally {
    updatingPassword.value = false
  }
}

const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    warning('文件大小不能超过10MB')
    return
  }

  selectedImageFile.value = file
  cropDialogVisible.value = true

  event.target.value = ''
}

const handleAvatarUploadSuccess = async (url) => {
  try {
    await updateUserProfileAPI({ avatar: url })
    success('头像上传成功')
    loadProfile()
  } catch (error) {
    notifyError('更新头像失败')
  }
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.user-settings {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.page-header p {
  color: #666;
  font-size: 14px;
}

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.card-body {
  padding: 20px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 36px;
  font-weight: 600;
  color: #fff;
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 100%);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.hint {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #ff6b35;
}

.form-group input:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
