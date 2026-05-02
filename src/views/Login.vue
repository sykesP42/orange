<template>
  <div class="login-page">
    <img v-if="bgImageUrl" :src="bgImageUrl" class="bg-wallpaper" alt="" @load="onBgLoaded" @error="onBgError" />
    <div class="login-bg">
      <div class="bg-overlay"></div>
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
      <div class="bg-grid"></div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <div class="card-header">
          <div class="logo">
            <span class="logo-icon">🍊</span>
            <span class="logo-text">Orange</span>
          </div>
          <h1 class="title">博主管理系统</h1>
          <p class="subtitle">小而美的博主归属锁定管理工具</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form" v-if="!showRegister && !showForgot">
          <div class="form-item">
            <label>用户名</label>
            <input
              v-model="form.username"
              type="text"
              placeholder="请输入用户名"
              autocomplete="username"
            />
          </div>

          <div class="form-item">
            <label>密码</label>
            <div class="password-input">
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                autocomplete="current-password"
              />
              <button type="button" class="toggle-btn" @click="showPassword = !showPassword">
                <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="form-actions">
            <a href="#" @click.prevent="showForgot = true">忘记密码？</a>
            <a href="#" @click.prevent="showRegister = true">注册账号</a>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
            <span v-if="loading" class="btn-loading"></span>
            {{ loading ? '登录中...' : '登 录' }}
          </button>
        </form>

        <form @submit.prevent="handleRegister" class="login-form" v-if="showRegister">
          <div class="form-item">
            <label>用户名</label>
            <input
              v-model="registerForm.username"
              type="text"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div class="form-item">
            <label>真实姓名</label>
            <input
              v-model="registerForm.real_name"
              type="text"
              placeholder="请输入真实姓名"
              required
            />
          </div>

          <div class="form-item">
            <label>密码</label>
            <input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              required
            />
          </div>

          <div class="form-item">
            <label>确认密码</label>
            <input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              required
            />
          </div>

          <div class="form-item">
            <label>邀请码（选填）</label>
            <input
              v-model="registerForm.invite_code"
              type="text"
              placeholder="输入邀请码，如 INV1A2B3C"
            />
          </div>

          <div class="form-tip">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="16" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            <span>填写邀请码可直接注册通过，无需管理员审核</span>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
            <span v-if="loading" class="btn-loading"></span>
            {{ loading ? '提交中...' : '提交注册' }}
          </button>

          <button type="button" class="back-btn" @click="showRegister = false">
            返回登录
          </button>
        </form>

        <form @submit.prevent="handleForgot" class="login-form" v-if="showForgot">
          <div class="form-item">
            <label>用户名</label>
            <input
              v-model="forgotForm.username"
              type="text"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div class="form-item">
            <label>真实姓名</label>
            <input
              v-model="forgotForm.real_name"
              type="text"
              placeholder="请输入注册时的真实姓名"
              required
            />
          </div>

          <div class="form-tip">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="16" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            <span>申请提交后，管理员将重置密码并通知您</span>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
            <span v-if="loading" class="btn-loading"></span>
            {{ loading ? '提交中...' : '提交申请' }}
          </button>

          <button type="button" class="back-btn" @click="showForgot = false">
            返回登录
          </button>
        </form>

        <div class="card-footer">
          <div class="footer-shortcuts">
            <kbd>Ctrl+K</kbd> 命令面板
            <kbd>?</kbd> 快捷键
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { loginAPI, registerAPI, forgotPasswordAPI } from '../api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const notification = useNotification()

const loading = ref(false)
const showPassword = ref(false)
const showRegister = ref(false)
const showForgot = ref(false)
const bgImageUrl = ref('')
const bgLoaded = ref(false)

const fetchBingWallpaper = async () => {
  const cachedUrl = localStorage.getItem('login_bg_url')
  const cachedTime = localStorage.getItem('login_bg_time')
  const cacheKey = localStorage.getItem('login_bg_key')
  if (cachedUrl && cachedTime && Date.now() - parseInt(cachedTime) < 86400000 && cacheKey) {
    bgImageUrl.value = cachedUrl
    return
  }
  try {
    const ts = Date.now()
    const url = `/api/bing-wallpaper?t=${ts}`
    const res = await fetch(url, { method: 'GET' })
    if (res.ok) {
      const blob = await res.blob()
      if (blob.size > 1000) {
        const objectUrl = URL.createObjectURL(blob)
        bgImageUrl.value = objectUrl
        return
      }
    }
  } catch { }
}

const onBgLoaded = (e) => { e.target.classList.add('loaded'); bgLoaded.value = true }
const onBgError = () => { bgImageUrl.value = '' }

const form = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  real_name: '',
  password: '',
  confirmPassword: '',
  invite_code: ''
})

const forgotForm = reactive({
  username: '',
  real_name: ''
})

const handleLogin = async () => {
  if (!form.username || !form.password) {
    notification.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const res = await loginAPI(form)
    if (res.code === 0 || res.code === 200 || res.token) {
      userStore.setUser(res.data || res)
      notification.success('登录成功')
      router.push('/')
    } else {
      notification.error(res.message || '登录失败')
    }
  } catch (error) {
    notification.error(error.response?.data?.message || '登录失败，请检查网络')
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (!registerForm.username || !registerForm.real_name || !registerForm.password) {
    notification.warning('请填写所有必填项')
    return
  }

  if (registerForm.password !== registerForm.confirmPassword) {
    notification.warning('两次输入的密码不一致')
    return
  }

  if (registerForm.password.length < 6) {
    notification.warning('密码长度至少6位')
    return
  }

  loading.value = true
  try {
    const res = await registerAPI({
      username: registerForm.username,
      real_name: registerForm.real_name,
      password: registerForm.password,
      invite_code: registerForm.invite_code
    })
    if (res.code === 0 || res.code === 200) {
      const msg = res.message || (res.data?.status === 'active' ? '注册成功，欢迎使用！' : '注册申请已提交，请等待管理员审核')
      notification.success(msg)
      showRegister.value = false
    } else {
      notification.error(res.message || '注册失败')
    }
  } catch (error) {
    notification.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}

const handleForgot = async () => {
  if (!forgotForm.username || !forgotForm.real_name) {
    notification.warning('请填写所有必填项')
    return
  }

  loading.value = true
  try {
    const res = await forgotPasswordAPI(forgotForm)
    if (res.code === 0 || res.code === 200) {
      notification.success('申请已提交，管理员将重置密码并通知您')
      showForgot.value = false
    } else {
      notification.error(res.message || '申请失败')
    }
  } catch (error) {
    notification.error(error.response?.data?.message || '申请失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const inviteCode = route.query.invite
  if (inviteCode) {
    showRegister.value = true
    registerForm.invite_code = inviteCode
  }
  fetchBingWallpaper()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: #0f172a;
}

.bg-wallpaper {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0;
  transition: opacity 1s ease;
}

.bg-wallpaper.loaded {
  opacity: 1;
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.bg-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(15, 23, 42, 0.75) 0%,
    rgba(15, 23, 42, 0.55) 40%,
    rgba(15, 23, 42, 0.65) 70%,
    rgba(15, 23, 42, 0.8) 100%
  );
  z-index: 1;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: orbFloat 20s ease-in-out infinite;
}

.bg-orb-1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, #f97316, #ea580c);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.bg-orb-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, #ec4899, #be185d);
  bottom: -150px;
  left: -100px;
  animation-delay: -7s;
}

.bg-orb-3 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, #8b5cf6, #6d28d9);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -30px) scale(1.05); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(20px, 10px) scale(1.02); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

.login-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
  padding: 20px;
}

.login-card {
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 24px;
  padding: 40px 32px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.3);
}

.card-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 40px;
  filter: drop-shadow(0 4px 12px rgba(249, 115, 22, 0.4));
}

.logo-text {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #f97316, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 4px;
}

.subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-item label {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.form-item input {
  height: 46px;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  font-size: 15px;
  color: white;
  transition: all 0.25s;
}

.form-item input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.form-item input:focus {
  outline: none;
  border-color: rgba(249, 115, 22, 0.5);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.password-input {
  position: relative;
}

.password-input input {
  width: 100%;
  padding-right: 48px;
}

.toggle-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  padding: 4px;
  display: flex;
  transition: color 0.2s;
}

.toggle-btn:hover {
  color: rgba(255, 255, 255, 0.7);
}

.toggle-btn svg {
  width: 18px;
  height: 18px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.form-actions a {
  color: rgba(249, 115, 22, 0.8);
  text-decoration: none;
  transition: color 0.2s;
}

.form-actions a:hover {
  color: #f97316;
}

.login-btn {
  height: 48px;
  background: linear-gradient(135deg, #f97316, #ea580c);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  position: relative;
  overflow: hidden;
}

.login-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.5s;
}

.login-btn:hover:not(:disabled)::before {
  left: 100%;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(249, 115, 22, 0.4);
}

.login-btn:active:not(:disabled) {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-loading {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.back-btn {
  height: 46px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: 10px;
  font-size: 13px;
  color: rgba(245, 158, 11, 0.9);
}

.form-tip svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.card-footer {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.footer-shortcuts {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
}

.footer-shortcuts kbd {
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-size: 10px;
  font-family: inherit;
  color: rgba(255, 255, 255, 0.4);
}

@media (max-width: 480px) {
  .login-card {
    padding: 32px 24px;
    border-radius: 20px;
  }

  .logo-icon {
    font-size: 32px;
  }

  .logo-text {
    font-size: 24px;
  }
}
</style>
