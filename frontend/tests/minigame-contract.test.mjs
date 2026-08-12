import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const read = (path) => readFileSync(new URL(path, import.meta.url), 'utf8')

const entertainmentGames = {
  rps: 'GamblingGameRockPaperScissors.vue',
  guess: 'CasualGameGuessNumber.vue',
  dice: 'CasualGameDice.vue',
  slot: 'CasualGameSlotMachine.vue',
  gobang: 'BoardGameGobang.vue',
  jungle: 'BoardGameJungle.vue',
  go: 'BoardGameGo.vue',
  othello: 'BoardGameOthello.vue',
  landbattle: 'BoardGameLandBattleChess.vue',
  chess: 'BoardGameChess.vue',
  fps: 'CompetitiveGameFPS.vue',
  moba: 'CompetitiveGameMOBA.vue',
  racing: 'CompetitiveGameRacing.vue',
  fighting: 'CompetitiveGameFighting.vue',
  war: 'CompetitiveGameWar.vue',
  poker: 'GamblingGamePoker.vue',
  horseracing: 'GamblingGameHorseRacing.vue',
  roulette: 'GamblingGameRoulette.vue',
  baccarat: 'GamblingGameBaccarat.vue',
  blackjack: 'GamblingGameBlackjack.vue',
  lottery: 'GamblingGameLottery.vue',
}

test('all 21 entertainment mini-games have a start-to-settlement UI contract', () => {
  const menu = read('../src/components/Menu/MenuEntertainment.vue')
  assert.equal(Object.keys(entertainmentGames).length, 21)

  for (const [gameID, filename] of Object.entries(entertainmentGames)) {
    assert.match(menu, new RegExp(`\\b${gameID}:\\s*show`), `${gameID} is missing from gameDialogRefs`)
    assert.match(menu, new RegExp(filename.replace('.', '\\.')), `${filename} is not loaded by the entertainment page`)

    const component = read(`../src/components/GameList/${filename}`)
    assert.match(component, /useMiniGameBase/, `${filename} does not use the shared session lifecycle`)
    assert.match(component, /\bstart(?:Game|MiniGame)\(/, `${filename} has no start path`)
    assert.match(component, /\bend(?:Game|MiniGame)\(/, `${filename} has no settlement path`)
    assert.match(component, /emit\('complete'/, `${filename} does not report completion`)
  }
})

test('both work mini-games call backend start, settlement and cancellation APIs', () => {
  for (const filename of ['WorkGameBank.vue', 'WorkGameTaxiDriver.vue']) {
    const component = read(`../src/components/GameList/${filename}`)
    assert.match(component, /StartMiniGame\(/, `${filename} has no start API call`)
    assert.match(component, /EndMiniGame\(/, `${filename} has no settlement API call`)
    assert.match(component, /CancelMiniGame\(/, `${filename} has no cancellation API call`)
  }
})
