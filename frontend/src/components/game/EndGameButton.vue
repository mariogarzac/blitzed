<template>
  <div class="flex gap-2">
    <button
      @click="confirmEnd = true"
      class="flex-1 py-2 rounded-lg border border-ef-bg4 text-ef-grey2 hover:border-ef-green hover:text-ef-green text-sm font-medium transition-colors"
    >End Game Early</button>
  </div>

  <div v-if="confirmEnd" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4">
    <div class="bg-ef-bg2 rounded-2xl p-6 max-w-sm w-full space-y-4">
      <h3 class="text-lg font-bold text-ef-fg">End game early?</h3>
      <p class="text-sm text-ef-grey1">The player with the most points will be declared the winner.</p>
      <div class="flex gap-3">
        <button
          @click="confirmEnd = false"
          class="flex-1 py-2 rounded-lg border border-ef-bg4 text-ef-grey2 hover:bg-ef-bg3 transition-colors"
        >Cancel</button>
        <button
          @click="end"
          :disabled="loading"
          class="flex-1 py-2 rounded-lg bg-ef-green hover:bg-ef-aqua text-ef-bg0 font-bold transition-colors disabled:opacity-40"
        >{{ loading ? '…' : 'End Game' }}</button>
      </div>
      <p v-if="error" class="text-sm text-ef-green">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMatchStore } from '@/stores/matchStore'

const props = defineProps<{ matchId: number }>()

const matchStore = useMatchStore()
const confirmEnd = ref(false)
const loading = ref(false)
const error = ref('')

async function end() {
  error.value = ''
  loading.value = true
  try {
    await matchStore.endMatch(props.matchId)
    confirmEnd.value = false
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Failed to end game'
  } finally {
    loading.value = false
  }
}
</script>
