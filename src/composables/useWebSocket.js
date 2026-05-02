import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '../stores/user'

let ws = null
let reconnectTimer = null
let heartbeatTimer = null

const isConnected = ref(false)
const lastMessage = ref(null)
const messageListeners = []

const WS_URL = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const token = localStorage.getItem('token')
  return `${protocol}//${host}/api/ws?token=${token}`
}

export const useWebSocket = () => {
  const connect = () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      return
    }

    try {
      ws = new WebSocket(WS_URL())

      ws.onopen = () => {
        isConnected.value = true

        sendMessage({ type: 'subscribe' })

        if (heartbeatTimer) clearInterval(heartbeatTimer)
        heartbeatTimer = setInterval(() => {
          sendMessage({ type: 'ping' })
        }, 30000)
      }

      ws.onmessage = (event) => {
        try {
          const messages = event.data.split('\n').filter(Boolean)
          messages.forEach(data => {
            const message = JSON.parse(data)
            lastMessage.value = message
            messageListeners.forEach(listener => listener(message))
          })
        } catch (error) {
          console.error('WebSocket message parse error:', error)
        }
      }

      ws.onclose = () => {
        isConnected.value = false
        if (heartbeatTimer) clearInterval(heartbeatTimer)

        if (reconnectTimer) clearTimeout(reconnectTimer)
        reconnectTimer = setTimeout(() => {
          connect()
        }, 3000)
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('WebSocket connection error:', error)
      if (reconnectTimer) clearTimeout(reconnectTimer)
      reconnectTimer = setTimeout(() => {
        connect()
      }, 5000)
    }
  }

  const disconnect = () => {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    if (heartbeatTimer) clearInterval(heartbeatTimer)
    if (ws) {
      ws.close()
      ws = null
    }
    isConnected.value = false
  }

  const sendMessage = (message) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message))
      return true
    }
    return false
  }

  const onMessage = (callback) => {
    messageListeners.push(callback)
    return () => {
      const index = messageListeners.indexOf(callback)
      if (index > -1) {
        messageListeners.splice(index, 1)
      }
    }
  }

  const subscribeToNotifications = (callback) => {
    return onMessage((message) => {
      if (message.type === 'notification') {
        callback(message.content)
      }
    })
  }

  const subscribeToMessages = (callback) => {
    return onMessage((message) => {
      if (message.type === 'new_message') {
        callback(message.content)
      }
    })
  }

  return {
    isConnected: computed(() => isConnected.value),
    lastMessage: computed(() => lastMessage.value),
    connect,
    disconnect,
    sendMessage,
    onMessage,
    subscribeToNotifications,
    subscribeToMessages
  }
}

export default useWebSocket
