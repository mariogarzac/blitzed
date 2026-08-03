<template>
  <div class="space-y-5">
    <div class="grid grid-cols-2 gap-4">
      <div class="space-y-1">
        <label class="block text-xs font-medium text-ef-grey1">Player A</label>
        <select
          v-model="playerAId"
          class="w-full bg-ef-bg2 border border-ef-bg4 rounded-lg px-3 py-2 text-ef-fg text-sm focus:outline-none focus:border-ef-green"
        >
          <option :value="null" disabled>Select…</option>
          <option v-for="p in players" :key="p.id" :value="p.id" :disabled="p.id === playerBId">
            {{ p.name }}
          </option>
        </select>
      </div>
      <div class="space-y-1">
        <label class="block text-xs font-medium text-ef-grey1">Player B</label>
        <select
          v-model="playerBId"
          class="w-full bg-ef-bg2 border border-ef-bg4 rounded-lg px-3 py-2 text-ef-fg text-sm focus:outline-none focus:border-ef-green"
        >
          <option :value="null" disabled>Select…</option>
          <option v-for="p in players" :key="p.id" :value="p.id" :disabled="p.id === playerAId">
            {{ p.name }}
          </option>
        </select>
      </div>
    </div>

    <div v-if="result" class="bg-ef-bg0 rounded-xl p-6 text-center space-y-4">
      <p class="text-ef-grey1 text-sm">{{ result.matches_played }} standard game{{ result.matches_played !== 1 ? 's' : '' }}</p>
      <div class="flex items-center justify-center gap-4">
        <div class="text-center">
          <p class="text-3xl font-bold text-ef-fg">{{ result.player_a_wins }}</p>
          <p class="text-sm text-ef-grey1 mt-1">{{ result.player_a.name }}</p>
        </div>
        <div class="text-center px-4">
          <p class="text-3xl font-bold text-ef-grey0">{{ result.draws }}</p>
          <p class="text-xs text-ef-grey0 mt-1">Draws</p>
        </div>
        <div class="text-center">
          <p class="text-3xl font-bold text-ef-fg">{{ result.player_b_wins }}</p>
          <p class="text-sm text-ef-grey1 mt-1">{{ result.player_b.name }}</p>
        </div>
      </div>
    </div>

    <p v-if="error" class="text-sm text-ef-green">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useStatsStore } from '@/stores/statsStore'
import { storeToRefs } from 'pinia'
import type { Player } from '@/types'

defineProps<{ players: Player[] }>()

const statsStore = useStatsStore()
const { headToHead: result } = storeToRefs(statsStore)

const playerAId = ref<number | null>(null)
const playerBId = ref<number | null>(null)
const error = ref('')

watch([playerAId, playerBId], async ([a, b]) => {
  if (!a || !b) return
  error.value = ''
  try {
    await statsStore.fetchHeadToHead(a, b)
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Failed to load stats'
  }
})
</script>
