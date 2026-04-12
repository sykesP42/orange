export const Platform = {
  isMobile: () => {
    return typeof (window as any).Capacitor !== 'undefined' || 
           window.innerWidth < 768
  },
  
  isDesktop: () => {
    return '__TAURI__' in window || window.innerWidth >= 1024
  },
  
  isTauri: () => {
    return '__TAURI__' in window
  },
  
  isCapacitor: () => {
    return typeof (window as any).Capacitor !== 'undefined'
  },
  
  isWeb: () => {
    return !Platform.isCapacitor() && !Platform.isTauri()
  },

  getOS: (): string => {
    const userAgent = navigator.userAgent.toLowerCase()
    
    if (userAgent.includes('android')) return 'Android'
    if (userAgent.includes('iphone') || userAgent.includes('ipad')) return 'iOS'
    if (userAgent.includes('mac')) return 'macOS'
    if (userAgent.includes('win')) return 'Windows'
    if (userAgent.includes('linux')) return 'Linux'
    
    return 'Unknown'
  },

  getType: (): 'mobile' | 'desktop' | 'web' => {
    if (Platform.isTauri()) return 'desktop'
    if (Platform.isCapacitor()) return 'mobile'
    return 'web'
  }
}

export const usePlatform = () => {
  const platform = ref(Platform.getType())
  const os = ref(Platform.getOS())
  const isMobile = computed(() => Platform.isMobile())
  const isDesktop = computed(() => Platform.isDesktop())
  const isNative = computed(() => Platform.isCapacitor() || Platform.isTauri())

  onMounted(() => {
    // 监听窗口大小变化
    const handleResize = () => {
      platform.value = Platform.getType()
      os.value = Platform.getOS()
    }
    
    window.addEventListener('resize', handleResize)
    onUnmounted(() => {
      window.removeEventListener('resize', handleResize)
    })
  })

  return {
    platform,
    os,
    isMobile,
    isDesktop,
    isNative
  }
}
