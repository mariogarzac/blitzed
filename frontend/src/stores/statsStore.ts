import { defineStore } from 'pinia'
import { ref } from 'vue'
import { statsApi } from '@/api/client'
import type { HeadToHeadResult } from '@/types'

export const useStatsStore = defineStore('stats', () => {
  const headToHead = ref<HeadToHeadResult | null>(null)

  async function fetchHeadToHead(playerAId: number, playerBId: number) {
    headToHead.value = await statsApi.headToHead(playerAId, playerBId)
  }

  function clear() {
    headToHead.value = null
  }

  return { headToHead, fetchHeadToHead, clear }
})
