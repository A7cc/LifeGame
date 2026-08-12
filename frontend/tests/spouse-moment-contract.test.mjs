import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const read = path => readFileSync(new URL(path, import.meta.url))
const readText = path => read(path).toString('utf8')

test('spouse date and bath shorts include production visual assets', () => {
  for (const asset of [
    '../public/images/datinginfo/dating-moments/date-sunset.webp',
    '../public/images/datinginfo/dating-moments/bath-spa.webp',
  ]) {
    const bytes = read(asset)
    assert.equal(bytes.subarray(0, 4).toString(), 'RIFF', `${asset} is not a WebP`)
    assert.equal(bytes.subarray(8, 12).toString(), 'WEBP', `${asset} is not a WebP`)
    assert.ok(bytes.length > 100_000, `${asset} is unexpectedly small`)
  }

  const dialog = readText('../src/components/Dialog/DialogSpouseMoment.vue')
  const scenes = readText('../src/utils/datingScenes.js')
  assert.match(scenes, /date-sunset\.webp/)
  assert.match(dialog, /bath-spa\.webp/)
  assert.match(dialog, /getDatingScene\(props\.location\)/)
  assert.match(dialog, /props\.spouse\?\.dimage/)
  assert.doesNotMatch(dialog, /moment-progress|progress-run/)
  assert.match(dialog, /data-testid="moment-interaction-sequence"/)
  assert.match(dialog, /data-testid="moment-relationship-status"/)
  assert.match(dialog, /relationshipStatus \|\| '陌生人'/)
  assert.match(dialog, /class="moment-relationship-dot"/)
  assert.doesNotMatch(dialog, /class="moment-ending"|class="spouse-name"/)
  const interactionBeatStyle = dialog.match(/\.interaction-beat\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(interactionBeatStyle, /border|background/)
  assert.match(interactionBeatStyle, /beat-show/)
  assert.match(dialog, /你轻轻抚摸/)
  assert.match(dialog, /交换了一个温柔的吻/)
  assert.match(dialog, /variant.*chat.*kiss.*intimacy.*argument/s)
  assert.match(dialog, /两位成年人再次确认/)
  assert.match(dialog, /isTenseMoment.*argument.*gift-disliked.*gift-rejected/s)
  assert.match(dialog, /moment-intimacy/)
  assert.match(dialog, /data-testid="moment-overlay-actions"/)
  assert.match(dialog, /data-testid="replay-spouse-moment"/)
  assert.match(dialog, /data-testid="close-spouse-moment"/)
  assert.ok(dialog.indexOf('data-testid="close-spouse-moment"') < dialog.indexOf('data-testid="replay-spouse-moment"'))
  assert.match(dialog, /\.moment-overlay-actions\s*\{[\s\S]*?flex-direction:\s*column;/)
  assert.match(dialog, /class="\{ 'outfit-action-anchor': action\.value === 'outfit' \}"/)
  assert.match(dialog, /\.outfit-picker\s*\{[\s\S]*?top:\s*127px;[\s\S]*?right:\s*48px;/)
  assert.match(dialog, /\.outfit-picker::after/)
  assert.doesNotMatch(dialog, /👗 指定本次造型|class="outfit-badge"/)
  const outfitPickerStyle = dialog.match(/\.outfit-picker\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(outfitPickerStyle, /border:/)
  assert.doesNotMatch(dialog, /class="moment-action-panel"|class="video-interaction-dock"|class="interaction-state"/)
  assert.match(dialog, /:show-close="false"/)
  assert.match(dialog, /aria-label="约会互动场景"/)
  assert.match(dialog, /:global\(\.spouse-moment-dialog\.el-dialog\)[\s\S]*?background:\s*transparent;[\s\S]*?box-shadow:\s*none;/)
  assert.match(dialog, /\.moment-stage\s*\{[\s\S]*?border:\s*4px solid #d5a43a;[\s\S]*?0 0 0 3px #f1cc72/)
  assert.match(dialog, /\.moment-stage::after\s*\{[\s\S]*?inset:\s*5px;[\s\S]*?border:\s*1px solid rgb\(255 235 169 \/ 88%\);[\s\S]*?inset 0 0 8px/)
  assert.doesNotMatch(dialog, /约 9 秒|momentActionHint/)
  for (const action of ['chat', 'outfit', 'caress', 'kiss', 'intimacy']) {
    assert.match(dialog, new RegExp(`date-interaction-\\$\\{action\\.value\\}`))
    assert.match(dialog, new RegExp(`value: '${action}'`))
  }
  assert.match(dialog, /emit\('interact', action, ''\)/)
  assert.match(dialog, /emit\('interact', 'outfit', outfitOption\.key\)/)
  assert.match(dialog, /重播短片/)
})

test('married cards expose bath interaction and every successful date opens the short', () => {
  const menu = readText('../src/components/Menu/MenuDating.vue')
  const dialog = readText('../src/components/Dialog/DialogSpouseMoment.vue')
  assert.match(menu, /bath-with-spouse-/)
  assert.match(menu, /BatheWithSpouse\(dating\.did\)/)
  assert.match(menu, /openSpouseMoment\('bath'/)
  assert.match(menu, /if \(result\.success && dating\)/)
  assert.match(menu, /result\.datinginfo\?\.dstatus \|\| dating\.daffinitylevel/)
  assert.match(menu, /result\.scene\?\.moment \|\| 'chat'/)
  assert.match(menu, /DoDatingInteraction\(momentSpouse\.value\.did, action, outfitKey\)/)
  assert.match(menu, /momentInteractionCompleted\.value = true/)
  assert.match(menu, /if \(momentInteractionCompleted\.value\) \{/)
  assert.match(menu, /momentVariant\.value = action/)
  assert.match(menu, /result\.outfitimage/)
  assert.match(menu, /result\.outfitvariant/)
  assert.match(menu, /datingOutfitImage/)
  assert.match(dialog, /portraitSource/)
  assert.match(dialog, /:disabled="interactionPending \|\| !action\.enabled"/)
  assert.doesNotMatch(dialog, /:disabled="[^"]*interactionCompleted/)
  assert.match(dialog, /interactionCompleted/)
  assert.match(menu, /openSpouseMoment\('date'/)
})
