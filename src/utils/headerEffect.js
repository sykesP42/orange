/**
 * Header 效果工具函数
 * 用于触发页面 Header 的光效动画
 */

export function triggerHeaderEffect(headerSelector = '.page-header') {
  const header = document.querySelector(headerSelector)
  if (header) {
    header.classList.add('header-glow')
    header.classList.add('header-shine')
    setTimeout(() => {
      header.classList.remove('header-shine')
    }, 600)
  }
}

export function removeHeaderGlow(headerSelector = '.page-header') {
  const header = document.querySelector(headerSelector)
  if (header) {
    header.classList.remove('header-glow')
  }
}
