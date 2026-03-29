import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './style.css'
import './styles/header.css'
import './styles/transitions.css'

const app = createApp(App)

// 注册全局头像错误处理指令
app.directive('avatar', {
  mounted(el) {
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

app.mount('#app')