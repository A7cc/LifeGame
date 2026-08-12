import { spawnSync } from 'node:child_process'

const minimumMajor = 20
const currentMajor = Number(process.versions.node.split('.')[0])
const playwrightArgs = [
  './node_modules/@playwright/test/cli.js',
  'test',
  ...process.argv.slice(2),
]

let result
if (currentMajor >= minimumMajor) {
  result = spawnSync(process.execPath, playwrightArgs, { stdio: 'inherit' })
} else {
  console.log(`[e2e] Node ${minimumMajor}+ is required; using an isolated Node 22 runtime.`)
  result = spawnSync(
    'npx',
    ['--yes', '--package=node@22', 'node', ...playwrightArgs],
    { stdio: 'inherit' },
  )
}

if (result.error) {
  console.error(result.error.message)
  process.exit(1)
}
process.exit(result.status ?? 1)
