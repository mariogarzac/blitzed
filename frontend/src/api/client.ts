import axios from 'axios'
import type { Player, Match, MatchSummary, HeadToHeadResult, RoundScoreInput } from '@/types'

const api = axios.create({ baseURL: '/api/v1' })

export const playerApi = {
  list: () => api.get<Player[]>('/players').then(r => r.data),
  create: (name: string) => api.post<Player>('/players', { name }).then(r => r.data),
  delete: (id: number) => api.delete(`/players/${id}`),
}

export const matchApi = {
  list: (status?: string) =>
    api.get<MatchSummary[]>('/matches', { params: status ? { status } : {} }).then(r => r.data),
  create: (playerIds: number[], winningScore: number) =>
    api.post<Match>('/matches', { player_ids: playerIds, winning_score: winningScore }).then(r => r.data),
  get: (id: number) => api.get<Match>(`/matches/${id}`).then(r => r.data),
  submitRound: (id: number, scores: RoundScoreInput[]) =>
    api.post<Match>(`/matches/${id}/rounds`, { scores }).then(r => r.data),
  undoRound: (id: number) => api.delete<Match>(`/matches/${id}/rounds/last`).then(r => r.data),
  end: (id: number) => api.post<Match>(`/matches/${id}/end`).then(r => r.data),
  delete: (id: number) => api.delete(`/matches/${id}`),
}

export const statsApi = {
  headToHead: (playerA: number, playerB: number) =>
    api
      .get<HeadToHeadResult>('/stats/head-to-head', { params: { player_a: playerA, player_b: playerB } })
      .then(r => r.data),
}
