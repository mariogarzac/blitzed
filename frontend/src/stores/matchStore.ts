import { defineStore } from 'pinia'
import { ref } from 'vue'
import { matchApi } from '@/api/client'
import type { Match, MatchSummary, RoundScoreInput } from '@/types'

export const useMatchStore = defineStore('match', () => {
  const activeMatch = ref<Match | null>(null)
  const matches = ref<MatchSummary[]>([])
  const currentMatchDetail = ref<Match | null>(null)

  async function fetchActiveMatch() {
    const list = await matchApi.list('in_progress')
    activeMatch.value = list.length > 0 ? await matchApi.get(list[0].id) : null
  }

  async function fetchMatches() {
    matches.value = await matchApi.list()
  }

  async function fetchMatchDetail(id: number) {
    currentMatchDetail.value = await matchApi.get(id)
  }

  async function createMatch(playerIds: number[], winningScore: number, hardMode: boolean) {
    const match = await matchApi.create(playerIds, winningScore, hardMode)
    activeMatch.value = match
    return match
  }

  async function submitRound(matchId: number, scores: RoundScoreInput[]) {
    const match = await matchApi.submitRound(matchId, scores)
    activeMatch.value = match
    return match
  }

  async function undoLastRound(matchId: number) {
    const match = await matchApi.undoRound(matchId)
    activeMatch.value = match
    return match
  }

  async function endMatch(matchId: number) {
    const match = await matchApi.end(matchId)
    activeMatch.value = null
    return match
  }

  async function deleteMatch(matchId: number) {
    await matchApi.delete(matchId)
    matches.value = matches.value.filter(m => m.id !== matchId)
    if (activeMatch.value?.id === matchId) activeMatch.value = null
  }

  function clearActiveMatch() {
    activeMatch.value = null
  }

  return {
    activeMatch,
    matches,
    currentMatchDetail,
    fetchActiveMatch,
    fetchMatches,
    fetchMatchDetail,
    createMatch,
    submitRound,
    undoLastRound,
    endMatch,
    deleteMatch,
    clearActiveMatch,
  }
})
