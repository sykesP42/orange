<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="card-header">
          <div class="logo">
            <span class="logo-icon">🍊</span>
            <span class="logo-text">Orange</span>
          </div>
          <h1 class="title">博主管理系统</h1>
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
                {{ showPassword ? '隐藏' : '显示' }}
              </button>
            </div>
          </div>

          <div class="form-actions">
            <a href="#" @click.prevent="showForgot = true">忘记密码？</a>
            <a href="#" @click.prevent="showRegister = true">注册账号</a>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
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
              placeholder="请输入密码"
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

          <div class="form-tip">
            <span>📋 注册申请提交后，需管理员审核通过方可使用</span>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
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
            <span>📋 申请提交后，管理员将重置密码并通知您</span>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
            {{ loading ? '提交中...' : '提交申请' }}
          </button>

          <button type="button" class="back-btn" @click="showForgot = false">
            返回登录
          </button>
        </form>
      </div>

      <div class="log-link">
        <router-link to="/log">📜 查看操作记录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { loginAPI, registerAPI, forgotPasswordAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()

const loading = ref(false)
const showPassword = ref(false)
const showRegister = ref(false)
const showForgot = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  real_name: '',
  password: '',
  confirmPassword: ''
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
      password: registerForm.password
    })
    if (res.code === 0 || res.code === 200) {
      notification.success('注册申请已提交，请等待管理员审核')
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
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  background: white;
  border-radius: 16px;
  padding: 40px 32px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.card-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 36px;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
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
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-item input {
  height: 44px;
  padding: 0 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 15px;
  color: #1e293b;
  transition: all 0.2s;
}

.form-item input:focus {
  outline: none;
  border-color: #f97316;
  box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
}

.form-item input::placeholder {
  color: #9ca3af;
}

.password-input {
  position: relative;
}

.password-input input {
  width: 100%;
  padding-right: 60px;
}

.toggle-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #6b7280;
  font-size: 14px;
  cursor: pointer;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.form-actions a {
  color: #f97316;
  text-decoration: none;
}

.form-actions a:hover {
  text-decoration: underline;
}

.login-btn {
  height: 48px;
  background: linear-gradient(135deg, #f97316, #ea580c);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.back-btn {
  height: 44px;
  background: transparent;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  color: #6b7280;
  font-size: 15px;
  cursor: pointer;
}

.back-btn:hover {
  background: #f9fafb;
}

.form-tip {
  padding: 12px;
  background: #fef3c7;
  border-radius: 8px;
  font-size: 13px;
  color: #92400e;
}

.divider {
  text-align: center;
  margin: 24px 0 16px;
  position: relative;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #e5e7eb;
}

.divider span {
  background: white;
  padding: 0 12px;
  color: #9ca3af;
  font-size: 13px;
  position: relative;
}

.accounts {
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: center;
}

.account {
  font-size: 13px;
  color: #6b7280;
}

.log-link {
  text-align: center;
  margin-top: 24px;
}

.log-link a {
  color: #6b7280;
  font-size: 14px;
  text-decoration: none;
}

.log-link a:hover {
  color: #f97316;
}
</style>