<template>
  <div class="space-y-3">
    <div
      v-for="(entry, i) in sorted"
      :key="entry.player.id"
      class="rounded-xl p-4 transition-all"
      :class="i === 0 ? 'bg-ef-bg2 ring-1 ring-ef-green' : 'bg-ef-bg0'"
    >
      <div class="flex justify-between items-center mb-2">
        <div class="flex items-center gap-2">
          <span class="text-ef-grey1 text-sm w-5">{{ i + 1 }}.</span>
          <span class="font-semibold text-ef-fg">{{ entry.player.name }}</span>
          <span v-if="i === 0" class="text-xs text-ef-green font-medium">Leading</span>
        </div>
        <span class="text-xl font-bold" :class="entry.cumulative < 0 ? 'text-ef-grey1' : 'text-ef-fg'">
          {{ entry.cumulative }}
        </span>
      </div>
      <div class="w-full bg-ef-bg3 rounded-full h-2 overflow-hidden">
        <div
          class="h-2 rounded-full transition-all duration-500"
          :class="i === 0 ? 'bg-ef-green' : 'bg-stone-500'"
          :style="{ width: barWidth(entry.cumulative) }"
        />
      </div>
      <div class="flex justify-between text-xs text-ef-grey0 mt-1">
        <span>0</span>
        <span>{{ winningScore }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Match, Player } from '@/types'

const props = defineProps<{ match: Match }>()

const winningScore = computed(() => props.match.winning_score)

const sorted = computed(() => {
  const cumulatives: Record<number, number> = {}
  for (const round of props.match.rounds) {
    for (const score of round.scores) {
      cumulatives[score.player_id] = score.cumulative
    }
  }
  return props.match.players
    .map(player => ({ player, cumulative: cumulatives[player.id] ?? 0 }))
    .sort((a, b) => b.cumulative - a.cumulative)
})

function barWidth(cumulative: number) {
  const pct = Math.max(0, Math.min(100, (cumulative / winningScore.value) * 100))
  return `${pct}%`
}
</script>
