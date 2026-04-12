import { config } from '@vue/test-utils'

config.global.stubs = {
  'router-link': {
    template: '<a><slot /></a>',
  },
  'router-view': {
    template: '<div><slot /></div>',
  },
}

config.global.directives = {
  avatar: {
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
    },
  },
}

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
})

class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}
window.ResizeObserver = ResizeObserver

window.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  observe() {}
  unobserve() {}
  disconnect() {}
}

HTMLElement.prototype.scrollIntoView = () => {}

if (typeof globalThis.document !== 'undefined') {
  document.documentElement.classList.add = () => {}
  document.documentElement.classList.remove = () => {}
  document.documentElement.classList.contains = () => false
  document.documentElement.classList.toggle = () => false
}
