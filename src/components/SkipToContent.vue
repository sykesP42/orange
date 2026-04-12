<template>
  <a href="#main-content" class="skip-link" @click="handleClick">
    {{ $t('nav.home') || '跳到主内容' }}
    <span class="skip-hint">→</span>
  </a>
</template>

<script setup>
const handleClick = (e) => {
  const target = document.getElementById('main-content')
  if (target) {
    e.preventDefault()
    target.tabIndex = -1
    target.focus()
    target.removeAttribute('tabindex')
  }
}
</script>

<style scoped>
.skip-link {
  position: fixed;
  top: -100px;
  left: 16px;
  z-index: 99999;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 18px;
  background: #f97316;
  color: white;
  font-weight: 600;
  font-size: 0.9rem;
  border-radius: 0 0 8px 8px;
  text-decoration: none;
  transition: top 0.2s ease, box-shadow 0.2s ease;
  box-shadow: 0 4px 12px rgba(249,115,22,0.3);
}
.skip-link:focus,
.skip-link:focus-visible {
  top: 0;
  outline: none;
  box-shadow: 0 4px 20px rgba(249,115,22,0.5);
}
.skip-hint { font-size: 1.1rem; }
@media (max-width: 480px) {
  .skip-link { font-size: 0.82rem; padding: 8px 14px; }
}
</style>
