<template>
  <div class="login-page" ref="pageRef" @mousemove="handleMouseMove">
    <canvas ref="canvasRef" class="geo-canvas"></canvas>

    <div class="login-container">
      <div class="login-card">
        <div class="card-header">
          <div class="logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
            </div>
            <span class="logo-text">Orange</span>
          </div>
          <h1 class="title">博主管理系统</h1>
          <p class="subtitle">v{{ version }}</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form" v-if="!showRegister && !showForgot">
          <div class="form-item">
            <label>用户名</label>
            <input v-model="form.username" type="text" placeholder="请输入用户名" autocomplete="username" />
          </div>
          <div class="form-item">
            <label>密码</label>
            <div class="password-input">
              <input v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="请输入密码" autocomplete="current-password" />
              <button type="button" class="toggle-btn" @click="showPassword = !showPassword">
                <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
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
          <div class="form-item"><label>用户名</label><input v-model="registerForm.username" type="text" placeholder="请输入用户名" required /></div>
          <div class="form-item"><label>真实姓名</label><input v-model="registerForm.real_name" type="text" placeholder="请输入真实姓名" required /></div>
          <div class="form-item"><label>密码</label><input v-model="registerForm.password" type="password" placeholder="请输入密码（至少6位）" required /></div>
          <div class="form-item"><label>确认密码</label><input v-model="registerForm.confirmPassword" type="password" placeholder="请再次输入密码" required /></div>
          <div class="form-tip"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg><span>注册申请提交后，需管理员审核通过方可使用</span></div>
          <button type="submit" class="login-btn" :disabled="loading"><span v-if="loading" class="btn-loading"></span>{{ loading ? '提交中...' : '提交注册' }}</button>
          <button type="button" class="back-btn" @click="showRegister = false">返回登录</button>
        </form>

        <form @submit.prevent="handleForgot" class="login-form" v-if="showForgot">
          <div class="form-item"><label>用户名</label><input v-model="forgotForm.username" type="text" placeholder="请输入用户名" required /></div>
          <div class="form-item"><label>真实姓名</label><input v-model="forgotForm.real_name" type="text" placeholder="请输入注册时的真实姓名" required /></div>
          <div class="form-tip"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg><span>申请提交后，管理员将重置密码并通知您</span></div>
          <button type="submit" class="login-btn" :disabled="loading"><span v-if="loading" class="btn-loading"></span>{{ loading ? '提交中...' : '提交申请' }}</button>
          <button type="button" class="back-btn" @click="showForgot = false">返回登录</button>
        </form>

        <div class="card-footer">
          <div class="footer-shortcuts">
            <kbd>Ctrl+K</kbd> 命令面板
            <span class="login-hint" @click="showLoginHint">未登录</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { loginAPI, registerAPI, forgotPasswordAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()
const version = '1.0.0'

const loading = ref(false)
const showPassword = ref(false)
const showRegister = ref(false)
const showForgot = ref(false)

const form = reactive({ username: '', password: '' })
const registerForm = reactive({ username: '', real_name: '', password: '', confirmPassword: '' })
const forgotForm = reactive({ username: '', real_name: '' })

const pageRef = ref(null)
const canvasRef = ref(null)
let ctx = null
let particles = []
let ambientOrbs = []
let mouseX = 0.5, mouseY = 0.5
let targetMouseX = 0.5, targetMouseY = 0.5
let rafId = null
let timeColor = { bg1: '#0f0c29', bg2: 'var(--purple)', bg3: '#24243e', accent: 'var(--primary)', accentAlt: 'var(--purple)', gridAlpha: 0.03 }

function getTimeTheme() {
  const h = new Date().getHours() + new Date().getMinutes() / 60
  if (h >= 5 && h < 7) return {
    bg1: '#1a1428', bg2: '#2d1b3d', bg3: '#1e1528',
    accent: 'var(--warning)', accentAlt: 'var(--primary)',
    gridAlpha: 0.025, orbAlpha: 0.04, particleAlpha: [0.03, 0.08],
    name: 'dawn'
  }
  if (h >= 7 && h < 10) return {
    bg1: '#1c1410', bg2: '#3d2518', bg3: '#261a10',
    accent: 'var(--primary)', accentAlt: 'var(--primary)',
    gridAlpha: 0.02, orbAlpha: 0.035, particleAlpha: [0.02, 0.07],
    name: 'morning'
  }
  if (h >= 10 && h < 15) return {
    bg1: '#0f1419', bg2: '#1a2026', bg3: '#111720',
    accent: 'var(--primary)', accentAlt: 'var(--primary-dark)',
    gridAlpha: 0.018, orbAlpha: 0.03, particleAlpha: [0.018, 0.055],
    name: 'midday'
  }
  if (h >= 15 && h < 17.5) return {
    bg1: '#17120d', bg2: '#332218', bg3: '#1f1810',
    accent: 'var(--warning)', accentAlt: 'var(--danger)',
    gridAlpha: 0.022, orbAlpha: 0.04, particleAlpha: [0.022, 0.065],
    name: 'afternoon'
  }
  if (h >= 17.5 && h < 20) return {
    bg1: '#1a0f14', bg2: '#3d1828', bg3: '#26141c',
    accent: 'var(--primary)', accentAlt: 'var(--purple)',
    gridAlpha: 0.028, orbAlpha: 0.05, particleAlpha: [0.025, 0.07],
    name: 'sunset'
  }
  if (h >= 20 && h < 23) return {
    bg1: '#0f0c1a', bg2: '#1e1838', bg3: '#151220',
    accent: 'var(--purple)', accentAlt: '#f472b6',
    gridAlpha: 0.025, orbAlpha: 0.04, particleAlpha: [0.02, 0.06],
    name: 'evening'
  }
  return {
    bg1: '#090a0f', bg2: '#0f1020', bg3: '#080812',
    accent: 'var(--purple)', accentAlt: 'var(--purple)',
    gridAlpha: 0.02, orbAlpha: 0.028, particleAlpha: [0.015, 0.045],
    name: 'night'
  }
}

function lerp(a, b, t) { return a + (b - a) * t }

class Particle {
  constructor(w, h, theme) { this.init(w, h, theme) }

  init(w, h, theme) {
    this.x = Math.random() * w
    this.y = Math.random() * h
    const types = ['circle', 'ring', 'dot', 'arc']
    this.type = types[Math.floor(Math.random() * types.length)]
    this.size = Math.random() * 18 + 6
    this.baseSize = this.size
    this.rotation = Math.random() * Math.PI * 2
    this.rotSpeed = (Math.random() - 0.5) * 0.004
    this.vx = (Math.random() - 0.5) * 0.12
    this.vy = (Math.random() - 0.5) * 0.12
    this.opacity = Math.random() * (theme.particleAlpha[1] - theme.particleAlpha[0]) + theme.particleAlpha[0]
    this.baseOpacity = this.opacity
    this.hueShift = Math.random() * 20 - 10
    this.mouseInfluence = Math.random() * 40 + 20
    this.phase = Math.random() * Math.PI * 2
    this.phaseSpeed = Math.random() * 0.004 + 0.001
    this.arcStart = Math.random() * Math.PI * 2
    this.arcLen = Math.random() * Math.PI * 1.2 + 0.4
  }

  update(w, h, mx, my, theme, time) {
    this.x += this.vx
    this.y += this.vy
    this.rotation += this.rotSpeed
    this.phase += this.phaseSpeed

    const dx = mx * w - this.x, dy = my * h - this.y
    const dist = Math.sqrt(dx * dx + dy * dy) || 1
    const force = Math.min(this.mouseInfluence / dist, 1.5)
    this.x += (dx / dist) * force * 0.008
    this.y += (dy / dist) * force * 0.008

    if (this.x < -50) this.x = w + 50
    if (this.x > w + 50) this.x = -50
    if (this.y < -50) this.y = h + 50
    if (this.y > h + 50) this.y = -50

    const breathe = Math.sin(this.phase) * 0.06 + 1
    this.size = this.baseSize * breathe
    this.opacity = this.baseOpacity * (0.85 + Math.sin(this.phase * 0.8) * 0.15)
  }

  draw(ctx, theme) {
    ctx.save()
    ctx.translate(this.x, this.y)
    ctx.rotate(this.rotation)
    ctx.globalAlpha = this.opacity

    switch (this.type) {
      case 'circle':
        ctx.beginPath()
        ctx.arc(0, 0, this.size / 2, 0, Math.PI * 2)
        ctx.strokeStyle = theme.accent
        ctx.lineWidth = 0.6
        ctx.stroke()
        break
      case 'ring':
        ctx.beginPath()
        ctx.arc(0, 0, this.size / 2, 0, Math.PI * 2)
        ctx.strokeStyle = theme.accentAlt
        ctx.lineWidth = 0.5
        ctx.globalAlpha *= 0.3
        ctx.stroke()
        break
      case 'dot':
        ctx.beginPath()
        ctx.arc(0, 0, this.size / 8, 0, Math.PI * 2)
        ctx.fillStyle = theme.accent
        ctx.fill()
        break
      case 'arc':
        ctx.beginPath()
        ctx.arc(0, 0, this.size / 2, this.arcStart, this.arcStart + this.arcLen)
        ctx.strokeStyle = theme.accent
        ctx.lineWidth = 0.7
        ctx.stroke()
        break
    }
    ctx.restore()
  }
}

class AmbientOrb {
  constructor(w, h, theme) { this.init(w, h, theme) }

  init(w, h, theme) {
    this.x = Math.random() * w
    this.y = Math.random() * h
    this.radius = Math.random() * Math.min(w, h) * 0.25 + 100
    this.targetX = Math.random() * w
    this.targetY = Math.random() * h
    this.speed = Math.random() * 0.0003 + 0.0001
    this.progress = Math.random()
    this.useAccent = Math.random() > 0.4
    this.alpha = Math.random() * theme.orbAlpha + theme.orbAlpha * 0.3
  }

  update(w, h, theme, time) {
    this.progress += this.speed
    if (this.progress > 1) {
      this.progress = 0
      this.targetX = Math.random() * w
      this.targetY = Math.random() * h
    }
    const ease = this.progress < 0.5
      ? 2 * this.progress * this.progress
      : 1 - Math.pow(-2 * this.progress + 2, 2) / 2
    this.x = lerp(this.x, this.targetX, 0.002)
    this.y = lerp(this.y, this.targetY, 0.002)

    const pulse = Math.sin(time * 0.0004 + this.progress * Math.PI * 2) * 0.15 + 1
    this.currentRadius = this.radius * pulse
    this.currentAlpha = this.alpha * (0.7 + Math.sin(time * 0.0003) * 0.3)
  }

  draw(ctx, theme) {
    const grad = ctx.createRadialGradient(this.x, this.y, 0, this.x, this.y, this.currentRadius)
    const color = this.useAccent ? theme.accent : theme.accentAlt
    
    // 处理 CSS 变量和十六进制颜色
    const getRgbaColor = (c, alpha) => {
      if (c.startsWith('var(')) {
        // 对于 CSS 变量，使用默认的橙色
        return `rgba(249,115,22,${alpha})`
      }
      if (c.startsWith('#')) {
        return this.hexToRgba(c, alpha)
      }
      return `rgba(249,115,22,${alpha})`
    }
    
    grad.addColorStop(0, getRgbaColor(color, this.currentAlpha))
    grad.addColorStop(0.5, getRgbaColor(color, this.currentAlpha * 0.3))
    grad.addColorStop(1, getRgbaColor(color, 0))
    ctx.fillStyle = grad
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height)
  }

  hexToRgba(hex, alpha) {
    const r = parseInt(hex.slice(1, 3), 16), g = parseInt(hex.slice(3, 5), 16), b = parseInt(hex.slice(5, 7), 16)
    return `rgba(${r},${g},${b},${alpha})`
  }
}

function initCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return
  ctx = canvas.getContext('2d')
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)
}

function resizeCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return
  const dpr = window.devicePixelRatio || 1
  canvas.width = window.innerWidth * dpr
  canvas.height = window.innerHeight * dpr
  canvas.style.width = window.innerWidth + 'px'
  canvas.style.height = window.innerHeight + 'px'
  if (ctx) ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  createScene(window.innerWidth, window.innerHeight)
}

function createScene(w, h) {
  const theme = getTimeTheme()
  particles = []
  for (let i = 0; i < 14; i++) particles.push(new Particle(w, h, theme))

  ambientOrbs = []
  for (let i = 0; i < 3; i++) ambientOrbs.push(new AmbientOrb(w, h, theme))
}

function handleMouseMove(e) {
  const rect = pageRef.value?.getBoundingClientRect()
  if (rect) {
    targetMouseX = (e.clientX - rect.left) / rect.width
    targetMouseY = (e.clientY - rect.top) / rect.height
  }
}

function animate(time) {
  if (!ctx || !canvasRef.value) return
  const w = window.innerWidth, h = window.innerHeight
  const theme = getTimeTheme()

  mouseX += (targetMouseX - mouseX) * 0.05
  mouseY += (targetMouseY - mouseY) * 0.05

  ctx.clearRect(0, 0, w, h)

  const bgGrad = ctx.createLinearGradient(0, 0, w, h)
  bgGrad.addColorStop(0, theme.bg1)
  bgGrad.addColorStop(0.5, theme.bg2)
  bgGrad.addColorStop(1, theme.bg3)
  ctx.fillStyle = bgGrad
  ctx.fillRect(0, 0, w, h)

  for (let i = 0; i < ambientOrbs.length; i++) {
    ambientOrbs[i].update(w, h, theme, time)
    ambientOrbs[i].draw(ctx, theme)
  }

  const mouseGlow = ctx.createRadialGradient(mouseX * w, mouseY * h, 0, mouseX * w, mouseY * h, 280)
  const mAlpha = 0.035 + Math.sin(time * 0.0005) * 0.012
  mouseGlow.addColorStop(0, `rgba(249,115,22,${mAlpha})`)
  mouseGlow.addColorStop(0.4, `rgba(249,115,22,${mAlpha * 0.3})`)
  mouseGlow.addColorStop(1, 'rgba(249,115,22,0)')
  ctx.fillStyle = mouseGlow
  ctx.fillRect(0, 0, w, h)

  const dotSpacing = 64
  const dotAlpha = theme.gridAlpha + Math.sin(time * 0.0003) * 0.005
  ctx.fillStyle = `rgba(255,255,255,${dotAlpha})`
  for (let dx = dotSpacing / 2; dx < w; dx += dotSpacing) {
    for (let dy = dotSpacing / 2; dy < h; dy += dotSpacing) {
      ctx.beginPath()
      ctx.arc(dx, dy, 0.6, 0, Math.PI * 2)
      ctx.fill()
    }
  }

  for (let i = 0; i < particles.length; i++) {
    particles[i].update(w, h, mouseX, mouseY, theme, time)
    particles[i].draw(ctx, theme)
  }

  ctx.save()
  ctx.globalAlpha = 0.04
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x
      const dy = particles[i].y - particles[j].y
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist < 120) {
        const la = (1 - dist / 120) * 0.035
        // 处理 CSS 变量和十六进制颜色
        const accentColor = theme.accent.startsWith('var(') || theme.accent.startsWith('#')
          ? `rgba(249,115,22,${la})` : `rgba(255,255,255,${la})`
        ctx.strokeStyle = accentColor
        ctx.lineWidth = 0.4
        ctx.beginPath()
        ctx.moveTo(particles[i].x, particles[i].y)
        ctx.lineTo(particles[j].x, particles[j].y)
        ctx.stroke()
      }
    }
  }
  ctx.restore()

  rafId = requestAnimationFrame(animate)
}

const showLoginHint = () => notification.info('请先登录后使用命令面板')

function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') { e.preventDefault(); e.stopPropagation(); showLoginHint() }
  if (e.key === '?' && !e.target.matches('input, textarea')) { e.preventDefault(); e.stopPropagation(); notification.info('登录后可使用快捷键面板') }
}

onMounted(() => { initCanvas(); animate(0); window.addEventListener('keydown', handleKeydown) })
onUnmounted(() => { cancelAnimationFrame(rafId); window.removeEventListener('resize', resizeCanvas); window.removeEventListener('keydown', handleKeydown) })

const handleLogin = async () => {
  if (!form.username || !form.password) { notification.warning('请输入用户名和密码'); return }
  loading.value = true
  try {
    const res = await loginAPI(form)
    if (res.code === 0 || res.code === 200 || res.token) { userStore.setUser(res.data || res); notification.success('登录成功'); router.push('/') }
    else notification.error(res.message || '用户名或密码错误')
  } catch (error) { notification.error(error?.response?.data?.message || error?.message || '用户名或密码错误') }
  finally { loading.value = false }
}

const handleRegister = async () => {
  if (!registerForm.username || !registerForm.real_name || !registerForm.password) { notification.warning('请填写所有必填项'); return }
  if (registerForm.password !== registerForm.confirmPassword) { notification.warning('两次输入的密码不一致'); return }
  if (registerForm.password.length < 6) { notification.warning('密码长度至少6位'); return }
  loading.value = true
  try {
    const res = await registerAPI({ username: registerForm.username, real_name: registerForm.real_name, password: registerForm.password })
    if (res.code === 0 || res.code === 200) { notification.success('注册申请已提交，请等待管理员审核'); showRegister.value = false }
    else notification.error(res.message || '注册失败')
  } catch (error) { notification.error(error.response?.data?.message || '注册失败') }
  finally { loading.value = false }
}

const handleForgot = async () => {
  if (!forgotForm.username || !forgotForm.real_name) { notification.warning('请填写所有必填项'); return }
  loading.value = true
  try {
    const res = await forgotPasswordAPI(forgotForm)
    if (res.code === 0 || res.code === 200) { notification.success('申请已提交，管理员将重置密码并通知您'); showForgot.value = false }
    else notification.error(res.message || '申请失败')
  } catch (error) { notification.error(error.response?.data?.message || '申请失败') }
  finally { loading.value = false }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: #0c0a09;
  isolation: auto;
}

.geo-canvas { position: absolute; inset: 0; width: 100%; height: 100%; display: block; z-index: 0; pointer-events: none; }

.login-container {
  position: relative; z-index: 1;
  width: 100%; max-width: 420px; padding: 20px;
}

.login-card {
  background: rgba(255, 255, 255, 0.035);
  backdrop-filter: blur(56px);
  -webkit-backdrop-filter: blur(56px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 24px;
  padding: 40px 32px;
  box-shadow:
    0 32px 64px rgba(0, 0, 0, 0.35),
    0 0 0 1px rgba(255, 255, 255, 0.03) inset,
    0 1px 0 rgba(255, 255, 255, 0.05) inset;
}

.card-header { text-align: center; margin-bottom: 32px; }

.logo { display: flex; align-items: center; justify-content: center; gap: 10px; margin-bottom: 16px; }

.logo-icon {
  width: 44px; height: 44px;
  background: var(--bg-primary-gradient-alt);
  border-radius: 13px;
  display: flex; align-items: center; justify-content: center;
  color: var(--color-on-brand);
  box-shadow: 0 8px 24px rgba(249, 115, 22, 0.25);
}

.logo-icon svg { width: 24px; height: 24px; }

.logo-text {
  font-size: 26px; font-weight: 700;
  background: linear-gradient(135deg, #fb923c, #f97316, #ea580c);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
  letter-spacing: -0.3px;
}

.title { font-size: 17px; font-weight: 600; color: rgba(255,255,255,0.85); margin-bottom: 4px; letter-spacing: -0.2px; }
.subtitle { font-size: 13px; color: rgba(255,255,255,0.32); letter-spacing: 0.2px; }

.login-form { display: flex; flex-direction: column; gap: 20px; }
.form-item { display: flex; flex-direction: column; gap: 6px; }
.form-item label { font-size: 13px; font-weight: 500; color: rgba(255,255,255,0.48); letter-spacing: 0.2px; }
.form-item input {
  height: 46px; padding: 0 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.07);
  border-radius: 12px; font-size: 15px; color: var(--color-on-brand);
  transition: all 0.25s cubic-bezier(0.4,0,0.2,1);
}
.form-item input::placeholder { color: rgba(255,255,255,0.2); }
.form-item input:focus { outline: none; border-color: rgba(249,115,22,0.4); background: rgba(255,255,255,0.06); box-shadow: 0 0 0 3px rgba(249,115,22,0.06); }

.password-input { position: relative; }
.password-input input { width: 100%; padding-right: 48px; }

.toggle-btn {
  position: absolute; right: 12px; top: 50%; transform: translateY(-50%);
  background: none; border: none; color: rgba(255,255,255,0.28); cursor: pointer; padding: 4px; display: flex; transition: color 0.2s;
}
.toggle-btn:hover { color: rgba(255,255,255,0.55); }
.toggle-btn svg { width: 18px; height: 18px; }

.form-actions { display: flex; justify-content: space-between; font-size: 13px; }
.form-actions a { color: rgba(249,115,22,0.65); text-decoration: none; transition: color 0.2s; }
.form-actions a:hover { color: var(--primary); }

.login-btn {
  height: 48px;
  background: var(--bg-primary-gradient-alt);
  border: none; border-radius: 12px; color: var(--color-on-brand);
  font-size: 16px; font-weight: 600; cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4,0,0.2,1);
  display: flex; align-items: center; justify-content: center; gap: 8px;
  position: relative; overflow: hidden;
}
.login-btn::before {
  content: ''; position: absolute; top: 0; left: -100%; width: 100%; height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.16), transparent);
  transition: left 0.5s;
}
.login-btn:hover:not(:disabled)::before { left: 100%; }
.login-btn:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 10px 30px rgba(249,115,22,0.3); }
.login-btn:active:not(:disabled) { transform: translateY(0); }
.login-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.btn-loading {
  width: 18px; height: 18px;
  border: 2px solid rgba(255,255,255,0.3); border-top-color: var(--color-on-brand); border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.back-btn {
  height: 46px; background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.07); border-radius: 12px;
  color: rgba(255,255,255,0.48); font-size: 15px; cursor: pointer; transition: all 0.2s;
}
.back-btn:hover { background: rgba(255,255,255,0.07); color: rgba(255,255,255,0.75); }

.form-tip {
  display: flex; align-items: center; gap: 8px; padding: 12px 16px;
  background: rgba(245,158,11,0.06); border: 1px solid rgba(245,158,11,0.12); border-radius: 10px;
  font-size: 13px; color: rgba(245,158,11,0.75);
}
.form-tip svg { width: 16px; height: 16px; flex-shrink: 0; }

.card-footer { margin-top: 24px; padding-top: 16px; border-top: 1px solid rgba(255,255,255,0.04); }
.footer-shortcuts { display: flex; align-items: center; justify-content: center; gap: 6px; font-size: 11px; color: rgba(255,255,255,0.2); }
.footer-shortcuts kbd {
  padding: 2px 6px; background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.06);
  border-radius: 4px; font-size: 10px; font-family: inherit; color: rgba(255,255,255,0.25);
}
.login-hint { font-size: 10px; color: rgba(249,115,22,0.4); cursor: pointer; transition: color 0.2s; }
.login-hint:hover { color: rgba(249,115,22,0.7); }

@media (max-width: 480px) {
  .login-card { padding: 32px 24px; border-radius: 20px; }
  .logo-icon { width: 38px; height: 38px; border-radius: 11px; }
  .logo-icon svg { width: 20px; height: 20px; }
  .logo-text { font-size: 22px; }
}
</style>
