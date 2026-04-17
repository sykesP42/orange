import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import './style.css'
import './styles/header.css'
import './styles/transitions.css'

const app = createApp(App)

// 注册全局头像错误处理指令
app.directive('avatar', {
  mounted(el) {
    el.setAttribute('loading', 'lazy')
    el.addEventListener('error', () => {
      el.style.display = 'none'
      const parent = el.parentElement
      if (parent) {
        const span = parent.querySelector('span')
        if (span) {
          span.style.display = 'flex'
        }
      }
    })
  }
})

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(i18n)

app.mount('#app')

// PWA 注册
if ('serviceWorker' in navigator) {
  window.addEventListener('load', async () => {
    try {
      const registration = await navigator.serviceWorker.register('/sw.js', { scope: '/' })

      registration.addEventListener('updatefound', () => {
        const newWorker = registration.installing
        newWorker.addEventListener('statechange', () => {
          if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
            window.dispatchEvent(new CustomEvent('sw-update-available'))
          }
        })
      })
    } catch (err) {
      // SW registration failed silently
    }
  })
}

// 全局键盘快捷键 (使用 Ctrl+Alt 避免与浏览器冲突)
window.addEventListener('keydown', (e) => {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.isContentEditable) return

  if (e.ctrlKey && e.altKey && e.key === 'k') {
    e.preventDefault()
    window.dispatchEvent(new CustomEvent('open-command-palette'))
  }

  if (e.ctrlKey && e.altKey && e.key === 'n') {
    e.preventDefault()
    router.push('/add')
  }

  if (e.ctrlKey && e.altKey && e.key === 'd') {
    e.preventDefault()
    document.documentElement.classList.toggle('dark')
    localStorage.setItem('theme', document.documentElement.classList.contains('dark') ? 'dark' : 'light')
  }

  if (e.ctrlKey && e.altKey && e.key === 'q') {
    e.preventDefault()
    window.dispatchEvent(new CustomEvent('open-quick-note'))
  }
})