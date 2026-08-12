import { readonly, ref } from 'vue'

const MUSIC_PREFERENCE_KEY = 'lifegame-background-music'
const MUSIC_SOURCE = '/audio/lifegame-theme.wav'
const MUSIC_VOLUME = 0.36

const enabled = ref(true)
const playing = ref(false)
const supported = ref(true)

let preferenceLoaded = false
let active = false
let audioContext = null
let masterGain = null
let musicBuffer = null
let bufferPromise = null
let sourceNode = null

const loadPreference = () => {
  if (preferenceLoaded || typeof window === 'undefined') return
  preferenceLoaded = true
  try {
    enabled.value = window.localStorage.getItem(MUSIC_PREFERENCE_KEY) !== 'off'
  } catch {
    enabled.value = true
  }
}

const savePreference = () => {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(MUSIC_PREFERENCE_KEY, enabled.value ? 'on' : 'off')
  } catch {
    // WebView 禁用本地存储时仍允许本次运行控制音乐。
  }
}

const removeUnlockListeners = () => {
  if (typeof document === 'undefined') return
  document.removeEventListener('pointerdown', unlockAudio, true)
  document.removeEventListener('keydown', unlockAudio, true)
}

const addUnlockListeners = () => {
  if (typeof document === 'undefined') return
  document.addEventListener('pointerdown', unlockAudio, true)
  document.addEventListener('keydown', unlockAudio, true)
}

const ensureAudioContext = () => {
  if (audioContext && audioContext.state !== 'closed') return audioContext
  if (typeof window === 'undefined') return null
  const AudioContextClass = window.AudioContext || window.webkitAudioContext
  if (!AudioContextClass) {
    supported.value = false
    return null
  }
  audioContext = new AudioContextClass()
  masterGain = audioContext.createGain()
  masterGain.gain.setValueAtTime(MUSIC_VOLUME, audioContext.currentTime)
  masterGain.connect(audioContext.destination)
  supported.value = true
  return audioContext
}

const loadMusicBuffer = (context, reload = false) => {
  if (reload && bufferPromise) return bufferPromise
  if (reload) {
    musicBuffer = null
    bufferPromise = null
  }
  if (musicBuffer) return Promise.resolve(musicBuffer)
  if (bufferPromise) return bufferPromise

  const requestURL = `${MUSIC_SOURCE}?reload=${reload ? Date.now() : 'startup'}`
  bufferPromise = window.fetch(requestURL, { cache: 'no-store' })
    .then(response => {
      if (!response.ok) throw new Error(`背景音乐读取失败：${response.status}`)
      return response.arrayBuffer()
    })
    .then(bytes => context.decodeAudioData(bytes))
    .then(buffer => {
      if (context !== audioContext) throw new Error('音频上下文已经更新')
      musicBuffer = buffer
      supported.value = true
      return buffer
    })
    .catch(error => {
      supported.value = false
      console.error(error)
      throw error
    })
    .finally(() => {
      bufferPromise = null
    })
  return bufferPromise
}

const stopSource = () => {
  const currentSource = sourceNode
  sourceNode = null
  if (currentSource) {
    currentSource.onended = null
    try {
      currentSource.stop()
    } catch {
      // 已结束的 BufferSource 无需再次停止。
    }
    currentSource.disconnect()
  }
  playing.value = false
}

export const playBackgroundMusic = async (reload = false) => {
  loadPreference()
  if (!active || !enabled.value) return false
  const context = ensureAudioContext()
  if (!context) return false

  // resume 必须在点击/按键事件栈中立即调用，才能通过 WebKit 媒体策略。
  const resumePromise = context.resume()
  const pendingBuffer = loadMusicBuffer(context, reload)
  try {
    await resumePromise
    const buffer = await pendingBuffer
    if (!active || !enabled.value || context !== audioContext || context.state !== 'running') {
      addUnlockListeners()
      return false
    }

    stopSource()
    const source = context.createBufferSource()
    source.buffer = buffer
    source.loop = true
    source.loopStart = 0
    source.loopEnd = buffer.duration
    source.connect(masterGain)
    source.onended = () => {
      if (sourceNode === source) {
        sourceNode = null
        playing.value = false
      }
    }
    sourceNode = source
    source.start(0)
    playing.value = true
    removeUnlockListeners()
    return true
  } catch {
    playing.value = false
    addUnlockListeners()
    return false
  }
}

function unlockAudio() {
  if (!active || !enabled.value || playing.value) return
  void playBackgroundMusic(false)
}

export const activateBackgroundMusic = () => {
  active = true
  loadPreference()
  if (!enabled.value) return
  ensureAudioContext()
  addUnlockListeners()
  void playBackgroundMusic(false)
}

export const deactivateBackgroundMusic = () => {
  active = false
  removeUnlockListeners()
  stopSource()
  musicBuffer = null
  bufferPromise = null
  const context = audioContext
  audioContext = null
  masterGain = null
  if (context && context.state !== 'closed') void context.close()
}

export const toggleBackgroundMusic = () => {
  loadPreference()

  // 未播放时作为明确的播放/重新读取入口，方便替换 ~/.lifegame/audio 后试听。
  if (enabled.value && !playing.value) {
    void playBackgroundMusic(true)
    return
  }

  enabled.value = !enabled.value
  savePreference()
  if (enabled.value) {
    addUnlockListeners()
    void playBackgroundMusic(true)
  } else {
    stopSource()
    removeUnlockListeners()
  }
}

export const useBackgroundMusic = () => ({
  enabled: readonly(enabled),
  playing: readonly(playing),
  supported: readonly(supported),
  play: playBackgroundMusic,
  toggle: toggleBackgroundMusic,
})
