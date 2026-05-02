<template>
  <div class="lang-switcher" :class="{ 'dark-mode': isDark }">
    <button class="lang-btn" @click="toggleDropdown" :title="$t('app.actions')">
      <span class="lang-icon">🌐</span>
      <span class="lang-label">{{ currentLabel }}</span>
      <svg class="lang-arrow" :class="{ open: showDropdown }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="6 9 12 15 18 9"/>
      </svg>
    </button>

    <Transition name="lang-dropdown">
      <div v-if="showDropdown" class="lang-menu">
        <button
          v-for="locale in locales"
          :key="locale.code"
          class="lang-option"
          :class="{ active: locale.code === currentLocale }"
          @click="switchLocale(locale.code)"
        >
          <span class="option-flag">{{ locale.flag }}</span>
          <span class="option-name">{{ locale.name }}</span>
          <svg v-if="locale.code === currentLocale" class="check-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({ isDark: { type: Boolean, default: false } })

const { locale: currentLocale } = useI18n()
const showDropdown = ref(false)

const locales = [
  { code: 'zh', name: '中文', flag: '🇨🇳' },
  { code: 'en', name: 'English', flag: '🇺🇸' },
]

const currentLabel = computed(() => {
  const found = locales.find(l => l.code === currentLocale.value)
  return found ? found.name : currentLocale.value.toUpperCase()
})

const switchLocale = (code) => {
  currentLocale.value = code
  localStorage.setItem('orange_locale', code)
  showDropdown.value = false
}

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const handleClickOutside = (e) => {
  if (!e.target.closest('.lang-switcher')) {
    showDropdown.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<style scoped>
.lang-switcher {
  position: relative;
  display: inline-flex;
  align-items: center;
}
.lang-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  background: transparent;
  border: 1px solid var(--border-color, rgba(128,128,128,0.15));
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.82rem;
  color: var(--text-primary, #374151);
  transition: all 0.2s ease;
  white-space: nowrap;
}
.lang-btn:hover {
  background: rgba(0,0,0,0.04);
  border-color: rgba(128,128,128,0.25);
}
.dark-mode .lang-btn:hover { background: rgba(255,255,255,0.06); }
.lang-icon { font-size: 1rem; }
.lang-label { font-weight: 500; letter-spacing: 0.02em; }
.lang-arrow {
  width: 14px; height: 14px;
  transition: transform 0.2s ease;
  color: var(--text-muted);
}
.lang-arrow.open { transform: rotate(180deg); }

.lang-menu {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  background: var(--bg-card);
  border-radius: 10px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.12), 0 2px 8px rgba(0,0,0,0.06), 0 0 0 1px rgba(0,0,0,0.04);
  overflow: hidden;
  z-index: 1000;
  min-width: 140px;
}
.dark-mode .lang-menu {
  background: var(--bg-card);
  border: 1px solid rgba(255,255,255,0.08);
}
.lang-option {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 9px 14px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--text-secondary);
  transition: background 0.15s;
  text-align: left;
}
.dark-mode .lang-option { color: var(--border-color); }
.lang-option:hover { background: var(--bg-hover); }
.dark-mode .lang-option:hover { background: rgba(255,255,255,0.06); }
.lang-option.active { background: var(--primary-50); color: var(--primary-700); font-weight: 600; }
.dark-mode .lang-option.active { background: rgba(249,115,22,0.12); color: var(--primary); }
.option-flag { font-size: 1.1rem; flex-shrink: 0; }
.option-name { flex: 1; }
.check-icon {
  width: 16px; height: 16px;
  color: var(--primary);
  flex-shrink: 0;
}

.lang-dropdown-enter-active,
.lang-dropdown-leave-active { transition: all 0.2s ease; }
.lang-dropdown-enter-from,
.lang-dropdown-leave-to { opacity: 0; transform: translateY(-4px); }

@media (max-width: 768px) {
  .lang-label { display: none; }
  .lang-btn { padding: 5px 7px; }
}
</style>
