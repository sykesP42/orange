<template>
  <div id="app">
    <OnlineStatus />
    <router-view />
    <MessageNotify />
    <ConfirmDialog />
    <PwaInstall />
    <Toast />
    <CommandPalette v-model:visible="showCommandPalette" @close="showCommandPalette = false" />
    <QuickNote v-model:visible="showQuickNote" @close="showQuickNote = false" />
    <ShortcutHelp v-model:visible="showShortcutHelp" @close="showShortcutHelp = false" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { useThemeStore } from './stores/theme'
import { useToastStore } from './stores/toast'

const MessageNotify = defineAsyncComponent(() => import('./components/MessageNotify.vue'))
const ConfirmDialog = defineAsyncComponent(() => import('./components/ConfirmDialog.vue'))
const PwaInstall = defineAsyncComponent(() => import('./components/PwaInstall.vue'))
const Toast = defineAsyncComponent(() => import('./components/Toast.vue'))
const CommandPalette = defineAsyncComponent(() => import('./components/CommandPalette.vue'))
const QuickNote = defineAsyncComponent(() => import('./components/QuickNote.vue'))
const ShortcutHelp = defineAsyncComponent(() => import('./components/ShortcutHelp.vue'))
const OnlineStatus = defineAsyncComponent(() => import('./components/OnlineStatus.vue'))
import { useKeyboardShortcut } from './composables/useKeyboardShortcut'

const toastStore = useToastStore()

const router = useRouter()
const themeStore = useThemeStore()

const showCommandPalette = ref(false)
const showQuickNote = ref(false)
const showShortcutHelp = ref(false)

const gSequence = ref([])
let gTimer = null

const handleGSequence = (key) => {
  gSequence.value.push(key)
  clearTimeout(gTimer)
  gTimer = setTimeout(() => {
    gSequence.value = []
  }, 1000)

  if (gSequence.value.length === 2) {
    const [, second] = gSequence.value
    const routeMap = {
      h: '/',
      a: '/add',
      m: '/my',
      b: '/my-bloggers',
      t: '/team',
      f: '/users',
      c: '/chat',
      s: '/statistics',
      d: '/admin'
    }
    if (routeMap[second]) {
      router.push(routeMap[second])
    }
    gSequence.value = []
  }
}

useKeyboardShortcut({
  'ctrl+k': () => {
    showCommandPalette.value = !showCommandPalette.value
  },
  'ctrl+n': () => {
    showQuickNote.value = !showQuickNote.value
  },
  'ctrl+d': () => {
    themeStore.toggleTheme()
  },
  '/': () => {
    showCommandPalette.value = true
  },
  '?': () => {
    showShortcutHelp.value = !showShortcutHelp.value
  },
  'g': () => {
    handleGSequence('g')
  },
  'escape': () => {
    showCommandPalette.value = false
    showQuickNote.value = false
    showShortcutHelp.value = false
  }
})

onMounted(() => {
  if ('launchQueue' in window) {
    window.launchQueue.setConsumer((launchParams) => {
      if (launchParams.targetURL) {
        const url = new URL(launchParams.targetURL)
        router.push(url.pathname)
      }
    })
  }

  window.addEventListener('open-command-palette', () => {
    showCommandPalette.value = true
  })
  window.addEventListener('open-quick-note', () => {
    showQuickNote.value = true
  })
  window.addEventListener('open-shortcut-help', () => {
    showShortcutHelp.value = true
  })
})
</script>

<style>
/* 全局样式保留在 style.css 和 transitions.css 中 */
</style>
