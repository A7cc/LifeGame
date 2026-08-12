import assert from 'node:assert/strict'
import { readdir } from 'node:fs/promises'
import path from 'node:path'
import test from 'node:test'
import { fileURLToPath } from 'node:url'

const imagesDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../public/images')

async function listFiles(directory, relativeDirectory = '') {
  const files = []
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const relativePath = path.join(relativeDirectory, entry.name)
    if (entry.isDirectory()) {
      files.push(...await listFiles(path.join(directory, entry.name), relativePath))
    } else {
      files.push(relativePath)
    }
  }
  return files
}

async function directoryNames(directory) {
  return (await readdir(directory, { withFileTypes: true }))
    .filter(entry => entry.isDirectory())
    .map(entry => entry.name)
    .sort()
}

test('images use four stable feature-domain directories without version suffixes', async () => {
  assert.deepEqual(await directoryNames(imagesDir), [
    'antiqueinfo',
    'carinfo',
    'datinginfo',
    'houseinfo',
  ])
  assert.deepEqual(await directoryNames(path.join(imagesDir, 'datinginfo')), [
    'dating-careers',
    'dating-moments',
    'dating-partner',
    'dating-scenes',
  ])
  assert.deepEqual(await directoryNames(path.join(imagesDir, 'carinfo')), ['car-moments', 'cars'])
  assert.deepEqual(await directoryNames(path.join(imagesDir, 'houseinfo')), ['house-moments', 'houses'])
  assert.deepEqual(await directoryNames(path.join(imagesDir, 'datinginfo', 'dating-partner')), ['female', 'male'])

  const files = await listFiles(imagesDir)
  assert.equal(files.length, 910)
  assert.deepEqual(files.filter(file => !file.endsWith('.webp')), [])
  assert.deepEqual(files.filter(file => /(?:^|[/\\])[^/\\]*v\d+(?:[/\\]|$)/i.test(file)), [])
})
