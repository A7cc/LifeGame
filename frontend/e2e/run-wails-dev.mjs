import { spawn } from 'node:child_process'
import { mkdtempSync, rmSync } from 'node:fs'
import { homedir, tmpdir } from 'node:os'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const repoRoot = dirname(dirname(dirname(fileURLToPath(import.meta.url))))
const realHome = homedir()
const isolatedHome = mkdtempSync(join(tmpdir(), 'lifegame-e2e-'))
const goPath = process.env.GOPATH || join(realHome, 'go')

const env = {
  ...process.env,
  HOME: isolatedHome,
  GOPATH: goPath,
  GOMODCACHE: process.env.GOMODCACHE || join(goPath, 'pkg', 'mod'),
  GOCACHE: process.env.GOCACHE || join(realHome, '.cache', 'go-build'),
  GOMAXPROCS: process.env.GOMAXPROCS || '1',
  GOFLAGS: `${process.env.GOFLAGS || ''} -p=1`.trim(),
  npm_config_cache: process.env.npm_config_cache || join(realHome, '.npm'),
}

const args = [
  'run',
  'github.com/wailsapp/wails/v2/cmd/wails@v2.10.2',
  'dev',
  '-tags',
  'webkit2_41',
  '-devserver',
  'localhost:34115',
  '-compiler',
  join(repoRoot, 'scripts', 'go-wails-e2e'),
  '-m',
  '-nosyncgomod',
  '-skipbindings',
  '-s',
  '-noreload',
  '-nogorebuild',
  '-nocolour',
]

console.log(`[e2e] Starting Wails with isolated game data: ${isolatedHome}`)
const child = spawn('go', args, {
  cwd: repoRoot,
  env,
  stdio: 'inherit',
  detached: process.platform !== 'win32',
})

let stopping = false

const removeIsolatedHome = () => {
  rmSync(isolatedHome, { recursive: true, force: true })
}

const stop = (signal = 'SIGINT') => {
  if (stopping) return
  stopping = true
  try {
    if (process.platform === 'win32') {
      child.kill(signal)
    } else {
      process.kill(-child.pid, signal)
    }
  } catch (error) {
    if (error?.code !== 'ESRCH') console.error(error)
  }
}

process.once('SIGINT', () => stop('SIGINT'))
process.once('SIGTERM', () => stop('SIGTERM'))

child.once('error', error => {
  console.error(error)
  removeIsolatedHome()
  process.exit(1)
})

child.once('exit', (code, signal) => {
  removeIsolatedHome()
  if (signal) process.kill(process.pid, signal)
  process.exit(code ?? 1)
})
