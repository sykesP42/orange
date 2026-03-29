/**
 * 移动端手势工具函数
 */

export function initSwipeGesture(element, onSwipeLeft, onSwipeRight) {
  let startX = 0
  let startY = 0
  let isSwiping = false
  const threshold = 50 // 滑动阈值（像素）

  const handleTouchStart = (e) => {
    startX = e.touches[0].clientX
    startY = e.touches[0].clientY
    isSwiping = false
  }

  const handleTouchMove = (e) => {
    if (!startX || !startY) return

    const diffX = e.touches[0].clientX - startX
    const diffY = e.touches[0].clientY - startY

    // 判断是否为水平滑动
    if (Math.abs(diffX) > Math.abs(diffY) && Math.abs(diffX) > 10) {
      isSwiping = true
      e.preventDefault()
    }
  }

  const handleTouchEnd = (e) => {
    if (!isSwiping) return

    const endX = e.changedTouches[0].clientX
    const diffX = endX - startX

    if (Math.abs(diffX) > threshold) {
      if (diffX > 0 && onSwipeRight) {
        onSwipeRight()
      } else if (diffX < 0 && onSwipeLeft) {
        onSwipeLeft()
      }
    }

    startX = 0
    startY = 0
    isSwiping = false
  }

  if (element) {
    element.addEventListener('touchstart', handleTouchStart, { passive: true })
    element.addEventListener('touchmove', handleTouchMove, { passive: false })
    element.addEventListener('touchend', handleTouchEnd, { passive: true })

    // 清理函数
    return () => {
      element.removeEventListener('touchstart', handleTouchStart)
      element.removeEventListener('touchmove', handleTouchMove)
      element.removeEventListener('touchend', handleTouchEnd)
    }
  }
}

/**
 * 长按手势
 */
export function initLongPress(element, onLongPress, duration = 800) {
  let timer = null
  let isPressed = false

  const handleTouchStart = () => {
    isPressed = true
    timer = setTimeout(() => {
      if (isPressed) {
        onLongPress()
        // 震动反馈
        if (navigator.vibrate) {
          navigator.vibrate(50)
        }
      }
    }, duration)
  }

  const handleTouchEnd = () => {
    isPressed = false
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  if (element) {
    element.addEventListener('touchstart', handleTouchStart)
    element.addEventListener('touchend', handleTouchEnd)
    element.addEventListener('touchcancel', handleTouchEnd)

    // 清理函数
    return () => {
      element.removeEventListener('touchstart', handleTouchStart)
      element.removeEventListener('touchend', handleTouchEnd)
      element.removeEventListener('touchcancel', handleTouchEnd)
    }
  }
}

/**
 * 双击手势
 */
export function initDoubleTap(element, onDoubleTap, maxDelay = 300) {
  let lastTap = 0

  const handleTouchEnd = (e) => {
    const currentTime = new Date().getTime()
    const tapLength = currentTime - lastTap

    if (tapLength < maxDelay && tapLength > 0) {
      onDoubleTap(e)
      e.preventDefault()
    }

    lastTap = currentTime
  }

  if (element) {
    element.addEventListener('touchend', handleTouchEnd, { passive: false })

    // 清理函数
    return () => {
      element.removeEventListener('touchend', handleTouchEnd)
    }
  }
}

/**
 * 捏合缩放手势
 */
export function initPinchZoom(element, onPinch) {
  let initialDistance = 0
  let isPinching = false

  const getDistance = (touches) => {
    const dx = touches[0].clientX - touches[1].clientX
    const dy = touches[0].clientY - touches[1].clientY
    return Math.sqrt(dx * dx + dy * dy)
  }

  const handleTouchStart = (e) => {
    if (e.touches.length === 2) {
      initialDistance = getDistance(e.touches)
      isPinching = true
    }
  }

  const handleTouchMove = (e) => {
    if (!isPinching || e.touches.length !== 2) return

    const currentDistance = getDistance(e.touches)
    const scale = currentDistance / initialDistance

    onPinch(scale)
    e.preventDefault()
  }

  const handleTouchEnd = () => {
    isPinching = false
    initialDistance = 0
  }

  if (element) {
    element.addEventListener('touchstart', handleTouchStart, { passive: true })
    element.addEventListener('touchmove', handleTouchMove, { passive: false })
    element.addEventListener('touchend', handleTouchEnd, { passive: true })

    // 清理函数
    return () => {
      element.removeEventListener('touchstart', handleTouchStart)
      element.removeEventListener('touchmove', handleTouchMove)
      element.removeEventListener('touchend', handleTouchEnd)
    }
  }
}
