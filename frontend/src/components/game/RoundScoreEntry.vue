<template>
  <div class="bg-ef-bg0 rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="font-semibold text-ef-fg">Round {{ nextRoundNum }}</h3>
      <button
        v-if="match.rounds.length > 0"
        @click="undo"
        :disabled="loading"
        class="text-xs text-ef-grey1 hover:text-ef-green transition-colors disabled:opacity-40"
      >Undo Round {{ match.rounds.length }}</button>
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div v-for="player in match.players" :key="player.id" class="space-y-1">
        <label class="block text-xs text-ef-grey1">{{ player.name }}</label>
        <input
          v-model.number="scores[player.id]"
          type="number"
          placeholder="0"
          class="w-full bg-ef-bg2 border border-ef-bg4 rounded-lg px-3 py-2 text-ef-fg text-center focus:outline-none focus:border-ef-green"
        />
      </div>
    </div>

    <p v-if="match.hard_mode" class="text-xs text-ef-yellow font-medium">Hard mode: scores can go below 0</p>
    <p v-else class="text-xs text-ef-aqua font-medium">Scores can't go below 0</p>

    <p v-if="error" class="text-sm text-ef-green">{{ error }}</p>

    <button
      @click="submit"
      :disabled="loading"
      class="w-full py-2.5 rounded-lg font-bold text-ef-bg0 bg-ef-green hover:bg-ef-aqua disabled:opacity-40 transition-colors"
    >{{ loading ? 'Saving…' : 'Submit Round' }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMatchStore } from '@/stores/matchStore'
import type { Match } from '@/types'

const props = defineProps<{ match: Match }>()

const matchStore = useMatchStore()
const loading = ref(false)
const error = ref('')

const scores = ref<Record<number, number>>(
  Object.fromEntries(props.match.players.map(p => [p.id, 0]))
)

const nextRoundNum = computed(() => props.match.rounds.length + 1)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const inputs = props.match.players.map(p => ({
      player_id: p.id,
      delta: scores.value[p.id] ?? 0,
    }))
    await matchStore.submitRound(props.match.id, inputs)
    // Reset scores
    props.match.players.forEach(p => { scores.value[p.id] = 0 })
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Failed to submit round'
  } finally {
    loading.value = false
  }
}

async function undo() {
  error.value = ''
  loading.value = true
  try {
    await matchStore.undoLastRound(props.match.id)
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Failed to undo round'
  } finally {
    loading.value = false
  }
}
</script>
