const LOG_LEVELS = { error: 0, warn: 1, info: 2, debug: 3 }
const currentLevel = LOG_LEVELS[import.meta.env.VITE_LOG_LEVEL || 'warn']

function formatMessage(context, message) {
  const timestamp = new Date().toISOString().slice(11, 19)
  return `[${timestamp}][${context}] ${message}`
}

export const logger = {
  error(context, message, error) {
    if (currentLevel >= LOG_LEVELS.error) {
      console.error(formatMessage(context, message), error || '')
    }
  },
  warn(context, message) {
    if (currentLevel >= LOG_LEVELS.warn) {
      console.warn(formatMessage(context, message))
    }
  },
  info(context, message) {
    if (currentLevel >= LOG_LEVELS.info) {
      console.info(formatMessage(context, message))
    }
  },
  debug(context, message) {
    if (currentLevel >= LOG_LEVELS.debug) {
      console.debug(formatMessage(context, message))
    }
  }
}
