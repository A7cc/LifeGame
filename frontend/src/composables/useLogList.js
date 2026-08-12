import { ref } from 'vue'

const formatTime = () => {
  const now = new Date()
  return `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
}

export const useLogList = (actionMap = {}) => {
  const logs = ref([])
  let nextId = 0

  const addLog = (detail, type = 'info', extra = {}) => {
    logs.value.unshift({
      id: nextId++,
      time: formatTime(),
      type,
      action: actionMap[type] || type,
      detail,
      ...extra,
    })
  }

  const clearLogs = () => {
    logs.value = []
  }

  return {
    addLog,
    clearLogs,
    logs,
  }
}
