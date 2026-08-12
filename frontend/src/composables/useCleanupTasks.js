import { onBeforeUnmount } from 'vue'

export const useCleanupTasks = () => {
  const timers = new Set()
  const cleanupTasks = new Set()

  const setManagedTimeout = (handler, timeout) => {
    const id = window.setTimeout(() => {
      timers.delete(id)
      handler()
    }, timeout)
    timers.add(id)
    return id
  }

  const setManagedInterval = (handler, timeout) => {
    const id = window.setInterval(handler, timeout)
    timers.add(id)
    return id
  }

  const clearManagedTimer = (id) => {
    if (!id) return
    window.clearTimeout(id)
    window.clearInterval(id)
    timers.delete(id)
  }

  const addCleanupTask = (task) => {
    cleanupTasks.add(task)
    return () => cleanupTasks.delete(task)
  }

  const cleanup = () => {
    timers.forEach(clearManagedTimer)
    timers.clear()
    cleanupTasks.forEach(task => task())
    cleanupTasks.clear()
  }

  onBeforeUnmount(cleanup)

  return {
    addCleanupTask,
    cleanup,
    clearManagedTimer,
    setManagedInterval,
    setManagedTimeout,
  }
}
