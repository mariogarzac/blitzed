<template>
  <div class="space-y-4">
    <h2 class="text-2xl font-bold text-ef-fg">Match History</h2>

    <div v-if="loading" class="text-center py-12 text-ef-grey1">Loading…</div>

    <div v-else-if="matches.length === 0" class="text-center py-12 text-ef-grey0">
      No matches yet. Go play some Dutch Blitz!
    </div>

    <div v-else class="space-y-2">
      <MatchListItem v-for="match in matches" :key="match.id" :match="match" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMatchStore } from '@/stores/matchStore'
import { storeToRefs } from 'pinia'
import MatchListItem from '@/components/history/MatchListItem.vue'

const matchStore = useMatchStore()
const { matches } = storeToRefs(matchStore)
const loading = ref(true)

onMounted(async () => {
  await matchStore.fetchMatches()
  loading.value = false
})
</script>
