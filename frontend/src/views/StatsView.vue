<template>
  <div class="space-y-5">
    <div>
      <h2 class="text-2xl font-bold text-ef-fg">Head-to-Head</h2>
      <p class="text-sm text-ef-grey1">Standard games only</p>
    </div>

    <div v-if="loading" class="text-center py-12 text-ef-grey1">Loading players…</div>
    <HeadToHead v-else :players="players" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePlayerStore } from '@/stores/playerStore'
import { storeToRefs } from 'pinia'
import HeadToHead from '@/components/stats/HeadToHead.vue'

const playerStore = usePlayerStore()
const { players } = storeToRefs(playerStore)
const loading = ref(true)

onMounted(async () => {
  await playerStore.fetchPlayers()
  loading.value = false
})
</script>
