<template>
  <div>
    <!-- Winner / Draw banner -->
    <div
      v-if="finishedMatch"
      class="mb-6 bg-ef-bg0 rounded-2xl p-6 text-center space-y-3"
    >
      <p v-if="finishedMatch.status === 'draw'" class="text-2xl font-bold text-ef-yellow">It's a Draw!</p>
      <div v-else>
        <p class="text-xs text-ef-grey1 uppercase tracking-widest mb-1">Winner</p>
        <p class="text-3xl font-bold text-ef-green">{{ finishedMatch.winners.map(w => w.name).join(' & ') }}</p>
      </div>
      <p class="text-sm text-ef-grey0">
        {{ finishedMatch.is_standard ? 'Standard game — win counted!' : 'Custom game' }}
      </p>
      <button
        @click="startNew"
        class="mt-2 px-6 py-2.5 rounded-xl bg-ef-green hover:bg-ef-aqua text-ef-bg0 font-bold transition-colors"
      >New Game</button>
    </div>

    <!-- Active game -->
    <div v-else-if="activeMatch" class="space-y-5">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-bold text-ef-fg">Game in Progress</h2>
          <p class="text-sm text-ef-grey1">
            First to {{ activeMatch.winning_score }} pts
            <span v-if="!activeMatch.is_standard" class="text-ef-yellow"> · Custom</span>
          </p>
        </div>
        <EndGameButton :match-id="activeMatch.id" />
      </div>

      <LiveStandings :match="activeMatch" />
      <RoundScoreEntry :match="activeMatch" :key="activeMatch.rounds.length" />
    </div>

    <!-- No active game -->
    <div v-else-if="!loading">
      <CreateGameForm />
    </div>

    <div v-if="loading" class="text-center py-12 text-ef-grey1">Loading…</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMatchStore } from '@/stores/matchStore'
import { usePlayerStore } from '@/stores/playerStore'
import { storeToRefs } from 'pinia'
import CreateGameForm from '@/components/game/CreateGameForm.vue'
import LiveStandings from '@/components/game/LiveStandings.vue'
import RoundScoreEntry from '@/components/game/RoundScoreEntry.vue'
import EndGameButton from '@/components/game/EndGameButton.vue'

const matchStore = useMatchStore()
const playerStore = usePlayerStore()
const { activeMatch } = storeToRefs(matchStore)
const loading = ref(true)

const finishedMatch = computed(() => {
  const m = activeMatch.value
  if (m && (m.status === 'finished' || m.status === 'draw')) return m
  return null
})

onMounted(async () => {
  await Promise.all([matchStore.fetchActiveMatch(), playerStore.fetchPlayers()])
  loading.value = false
})

function startNew() {
  matchStore.clearActiveMatch()
}
</script>
