<template>
  <div class="space-y-5">
    <div class="flex items-center gap-3">
      <RouterLink to="/history" class="text-ef-grey1 hover:text-ef-fg transition-colors">← Back</RouterLink>
      <h2 class="text-xl font-bold text-ef-fg">Match Detail</h2>
    </div>

    <div v-if="loading" class="text-center py-12 text-ef-grey1">Loading…</div>

    <template v-else-if="match">
      <!-- Meta -->
      <div class="bg-ef-bg0 rounded-xl p-4 flex flex-wrap gap-4 text-sm">
        <div>
          <p class="text-ef-grey1 text-xs">Date</p>
          <p class="text-ef-fg font-medium">{{ formatDate(match.started_at) }}</p>
        </div>
        <div>
          <p class="text-ef-grey1 text-xs">Winning Score</p>
          <p class="text-ef-fg font-medium">{{ match.winning_score }}</p>
        </div>
        <div>
          <p class="text-ef-grey1 text-xs">Type</p>
          <p class="font-medium" :class="match.is_standard ? 'text-ef-aqua' : 'text-ef-yellow'">
            {{ match.is_standard ? 'Standard' : 'Custom' }}
          </p>
        </div>
        <div>
          <p class="text-ef-grey1 text-xs">Result</p>
          <p class="text-ef-fg font-medium">
            <span v-if="match.status === 'finished'">{{ match.winners.map(w => w.name).join(', ') }} won</span>
            <span v-else-if="match.status === 'draw'">Draw</span>
            <span v-else>In Progress</span>
          </p>
        </div>
      </div>

      <!-- Chart -->
      <ScoreChart v-if="match.rounds.length > 0" :match="match" />
      <p v-else class="text-ef-grey0 text-sm text-center py-4">No rounds played yet.</p>

      <!-- Round table -->
      <div v-if="match.rounds.length > 0" class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-ef-grey1 text-xs border-b border-ef-bg1">
              <th class="text-left py-2 pr-4">Round</th>
              <th v-for="player in match.players" :key="player.id" class="text-center py-2 px-2">
                {{ player.name }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="round in match.rounds"
              :key="round.id"
              class="border-b border-stone-900 hover:bg-ef-bg0"
            >
              <td class="py-2 pr-4 text-ef-grey1">R{{ round.round_num }}</td>
              <td
                v-for="player in match.players"
                :key="player.id"
                class="text-center py-2 px-2"
              >
                <span class="text-ef-grey2">{{ getScore(round, player.id).delta >= 0 ? '+' : '' }}{{ getScore(round, player.id).delta }}</span>
                <span class="text-ef-grey0 text-xs ml-1">({{ getScore(round, player.id).cumulative }})</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useMatchStore } from '@/stores/matchStore'
import { storeToRefs } from 'pinia'
import ScoreChart from '@/components/history/ScoreChart.vue'
import type { Round } from '@/types'

const route = useRoute()
const matchStore = useMatchStore()
const { currentMatchDetail: match } = storeToRefs(matchStore)
const loading = ref(true)

onMounted(async () => {
  await matchStore.fetchMatchDetail(Number(route.params.id))
  loading.value = false
})

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
}

function getScore(round: Round, playerId: number) {
  return round.scores.find(s => s.player_id === playerId) ?? { delta: 0, cumulative: 0 }
}
</script>
