<template>
  <RouterLink
    :to="`/history/${match.id}`"
    class="block bg-ef-bg0 hover:bg-ef-bg2 rounded-xl p-4 transition-colors"
  >
    <div class="flex items-start justify-between gap-2">
      <div class="space-y-1">
        <div class="flex items-center gap-2 flex-wrap">
          <span class="font-semibold text-ef-fg">{{ match.player_count }} players</span>
          <span
            class="text-xs px-2 py-0.5 rounded-full font-medium"
            :class="match.is_standard ? 'bg-ef-bg-green text-ef-green' : 'bg-ef-bg-yellow text-ef-yellow'"
          >{{ match.is_standard ? 'Standard' : 'Custom' }}</span>
        </div>
        <div class="text-sm text-ef-grey1">
          <span v-if="match.status === 'finished'">
            Winner: <span class="text-ef-fg font-medium">{{ match.winners.map(w => w.name).join(', ') }}</span>
          </span>
          <span v-else-if="match.status === 'draw'" class="text-ef-yellow">Draw</span>
          <span v-else class="text-ef-blue">In Progress</span>
        </div>
      </div>
      <div class="text-right shrink-0">
        <p class="text-xs text-ef-grey0">{{ formatDate(match.started_at) }}</p>
        <p class="text-xs text-ef-grey0 mt-0.5">Win at {{ match.winning_score }} pts</p>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import type { MatchSummary } from '@/types'

defineProps<{ match: MatchSummary }>()

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>
