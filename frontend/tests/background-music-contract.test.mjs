import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const readText = path => readFileSync(new URL(path, import.meta.url), 'utf8')
const read = path => readFileSync(new URL(path, import.meta.url))

test('background music has a persistent user control and follows app lifecycle', () => {
  const music = readText('../src/composables/useBackgroundMusic.js')
  const app = readText('../src/App.vue')
  const topBar = readText('../src/components/GameTopBar.vue')
  const audio = read('../public/audio/lifegame-theme.wav')

  assert.match(music, /lifegame-background-music/)
  assert.match(music, /\/audio\/lifegame-theme\.wav/)
  assert.match(music, /AudioContext \|\| window\.webkitAudioContext/)
  assert.match(music, /window\.fetch\(requestURL, \{ cache: 'no-store' \}\)/)
  assert.match(music, /context\.decodeAudioData\(bytes\)/)
  assert.match(music, /context\.createBufferSource\(\)/)
  assert.match(music, /source\.loop = true/)
  assert.match(music, /source\.loopStart = 0/)
  assert.match(music, /source\.loopEnd = buffer\.duration/)
  assert.match(music, /document\.addEventListener\('pointerdown', unlockAudio, true\)/)
  assert.match(music, /document\.addEventListener\('keydown', unlockAudio, true\)/)
  assert.match(music, /context\.close\(\)/)

  assert.equal(audio.subarray(0, 4).toString(), 'RIFF')
  assert.equal(audio.subarray(8, 12).toString(), 'WAVE')
  assert.ok(audio.length > 1_000_000, 'background music asset is unexpectedly small')
  assert.equal(audio.readInt16LE(44), audio.readInt16LE(audio.length - 4), 'left loop boundary is discontinuous')
  assert.equal(audio.readInt16LE(46), audio.readInt16LE(audio.length - 2), 'right loop boundary is discontinuous')

  assert.match(app, /activateBackgroundMusic\(\)/)
  assert.match(app, /deactivateBackgroundMusic\(\)/)
  assert.match(topBar, /data-testid="background-music-toggle"/)
  assert.match(topBar, /关闭背景音乐/)
  assert.match(topBar, /开启背景音乐/)
  assert.match(topBar, /播放背景音乐/)
  assert.match(topBar, /useBackgroundMusic\(\)/)
})
