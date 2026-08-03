import { defineStore } from 'pinia'
import { ref } from 'vue'
import { playerApi } from '@/api/client'
import type { Player } from '@/types'

export const usePlayerStore = defineStore('player', () => {
  const players = ref<Player[]>([])

  async function fetchPlayers() {
    players.value = await playerApi.list()
  }

  async function createPlayer(name: string) {
    const player = await playerApi.create(name)
    players.value.push(player)
    return player
  }

  async function deletePlayer(id: number) {
    await playerApi.delete(id)
    players.value = players.value.filter(p => p.id !== id)
  }

  return { players, fetchPlayers, createPlayer, deletePlayer }
})
