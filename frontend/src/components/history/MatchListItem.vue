<template>
  <div class="flex items-stretch gap-2">
    <RouterLink
      :to="`/history/${match.id}`"
      class="flex-1 block bg-ef-bg0 hover:bg-ef-bg2 rounded-xl p-4 transition-colors"
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

    <button
      @click="confirmDelete($event)"
      class="delete-btn bg-ef-bg0 text-ef-grey1 rounded-xl px-3 transition-colors shrink-0"
      title="Delete match"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="3 6 5 6 21 6" />
        <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
        <path d="M10 11v6M14 11v6" />
        <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useMatchStore } from '@/stores/matchStore'
import type { MatchSummary } from '@/types'

const props = defineProps<{ match: MatchSummary }>()
const matchStore = useMatchStore()

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

async function confirmDelete(e: MouseEvent) {
  const btn = e.currentTarget as HTMLButtonElement
  const confirmed = confirm('Delete this match? This cannot be undone.')
  btn.blur()
  if (!confirmed) return
  await matchStore.deleteMatch(props.match.id)
}
</script>

<style scoped>
@media (hover: hover) {
  .delete-btn:hover {
    background-color: #514045;
    color: #e67e80;
  }
}
</style>
