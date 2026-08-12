import assert from 'node:assert/strict'
import { readFile, readdir, stat } from 'node:fs/promises'
import path from 'node:path'
import test from 'node:test'
import { fileURLToPath } from 'node:url'

const frontendDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const projectDir = path.resolve(frontendDir, '..')
const antiquesDir = path.join(frontendDir, 'public', 'images', 'antiqueinfo')

test('all 50 antiques use unique English-named optimized WebP assets', async () => {
  const seedSource = await readFile(path.join(projectDir, 'internal', 'db', 'seeds.go'), 'utf8')
  const antiqueSection = seedSource.slice(
    seedSource.indexOf('func seedAntiques'),
    seedSource.indexOf('func seedStocks'),
  )
  const imagePaths = [...antiqueSection.matchAll(/"(\/images\/antiqueinfo\/[0-9]{2}-[a-z0-9-]+\.webp)"/g)]
    .map((match) => match[1])

  assert.equal(imagePaths.length, 50)
  assert.equal(new Set(imagePaths).size, 50)
  imagePaths.forEach((imagePath, index) => {
    assert.match(path.basename(imagePath), new RegExp(`^${String(index + 1).padStart(2, '0')}-`))
  })
  assert.doesNotMatch(antiqueSection, /\/images\/antiqueinfo\/[^"\s]+\.(?:png|svg)/)

  const filenames = (await readdir(antiquesDir)).sort()
  assert.equal(filenames.length, 51)
  assert.deepEqual(filenames.filter((filename) => !filename.endsWith('.webp')), [])
  assert.ok(filenames.includes('default.webp'))

  for (const imagePath of [...imagePaths, '/images/antiqueinfo/default.webp']) {
    const localPath = path.join(frontendDir, 'public', imagePath)
    const [buffer, metadata] = await Promise.all([readFile(localPath), stat(localPath)])
    assert.equal(buffer.subarray(0, 4).toString('ascii'), 'RIFF', imagePath)
    assert.equal(buffer.subarray(8, 12).toString('ascii'), 'WEBP', imagePath)
    assert.equal(buffer.subarray(12, 16).toString('ascii'), 'VP8 ', imagePath)
    assert.equal(buffer.readUInt16LE(26) & 0x3fff, 1024, imagePath)
    assert.equal(buffer.readUInt16LE(28) & 0x3fff, 1024, imagePath)
    assert.ok(metadata.size > 40_000, `${imagePath} is unexpectedly small`)
  }
})
