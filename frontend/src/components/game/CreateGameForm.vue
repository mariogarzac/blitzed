<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-ef-fg mb-1">New Game</h2>
      <p class="text-sm text-ef-grey1">Select players and set the winning score.</p>
    </div>

    <PlayerPicker :players="players" v-model="selectedPlayerIds" />

    <div class="space-y-1">
      <label class="block text-sm font-medium text-ef-grey2">Winning Score</label>
      <div class="flex items-center gap-3">
        <input
          v-model.number="winningScore"
          type="number"
          min="10"
          class="w-28 bg-ef-bg2 border border-ef-bg4 rounded-lg px-3 py-2 text-ef-fg focus:outline-none focus:border-ef-green"
        />
        <span v-if="winningScore !== 75" class="text-xs text-ef-yellow font-medium"> Allows score to be negative</span>
        <span v-else class="text-xs text-ef-aqua font-medium">Standard game</span>
      </div>
    </div>

    <div class="space-y-1">
      <label id="hard-mode-label" class="block text-sm font-medium text-ef-grey2">Difficulty Setting</label>
      <div class="flex items-center gap-3">
        <button
          type="button"
          role="switch"
          :aria-checked="hardMode"
          aria-labelledby="hard-mode-label"
          @click="toggleHardMode"
          :class="hardMode ? 'bg-ef-bg-yellow border-ef-yellow' : 'bg-ef-bg2 border-ef-bg4'"
          class="relative inline-flex shrink-0 w-11 h-6 border rounded-full transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ef-green focus-visible:ring-offset-2 focus-visible:ring-offset-ef-bg0"
        >
          <span
            :class="hardMode ? 'translate-x-[22px] bg-ef-yellow' : 'translate-x-0.5 bg-ef-grey1'"
            class="absolute top-1/2 -translate-y-1/2 left-0 w-[18px] h-[18px] rounded-full shadow transition-transform"
          />
        </button>
        <span
          @click="toggleHardMode"
          :class="hardMode ? 'text-ef-yellow' : 'text-ef-grey1'"
          class="text-sm font-medium cursor-pointer select-none w-20"
        >{{ hardMode ? 'Hard' : 'Normal' }}</span>
        <span v-if="hardMode" class="text-xs text-ef-yellow font-medium">Allows score to be negative</span>
        <span v-else class="text-xs text-ef-aqua font-medium">No negative score</span>
      </div>
    </div>

    <div class="space-y-2">
      <button
        @click="submit"
        :disabled="selectedPlayerIds.length < 2 || loading"
        class="w-full py-3 rounded-xl font-bold text-ef-bg0 bg-ef-green hover:bg-ef-aqua disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ loading ? 'Starting…' : 'Start Game' }}
      </button>
      <p v-if="error" class="text-sm text-ef-green">{{ error }}</p>
    </div>

    <!-- Player management -->
    <div class="border-t border-ef-bg1 pt-4 space-y-3">
      <p class="text-sm font-medium text-ef-grey1">Manage Players</p>
      <div class="flex gap-2">
        <input
          v-model="newPlayerName"
          @keyup.enter="addPlayer"
          type="text"
          placeholder="Player name"
          class="flex-1 bg-ef-bg2 border border-ef-bg4 rounded-lg px-3 py-2 text-ef-fg text-sm focus:outline-none focus:border-ef-green"
        />
        <button
          @click="addPlayer"
          :disabled="!newPlayerName.trim()"
          class="px-4 py-2 rounded-lg bg-ef-bg3 hover:bg-ef-bg3 text-ef-fg text-sm font-medium disabled:opacity-40 transition-colors"
        >Add</button>
      </div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="player in players"
          :key="player.id"
          class="flex items-center gap-1 px-2 py-1 bg-ef-bg2 rounded-full text-xs text-ef-grey2"
        >
          {{ player.name }}
          <button @click="removePlayer(player.id)" class="text-ef-grey0 hover:text-ef-green ml-1">✕</button>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { usePlayerStore } from '@/stores/playerStore'
import { useMatchStore } from '@/stores/matchStore'
import { storeToRefs } from 'pinia'
import PlayerPicker from './PlayerPicker.vue'

const playerStore = usePlayerStore()
const matchStore = useMatchStore()
const { players } = storeToRefs(playerStore)

const selectedPlayerIds = ref<number[]>([])
const winningScore = ref(75)
const hardMode = ref(false)
const newPlayerName = ref('')
const loading = ref(false)
const error = ref('')

async function addPlayer() {
  const name = newPlayerName.value.trim()
  if (!name) return
  await playerStore.createPlayer(name)
  newPlayerName.value = ''
}

async function removePlayer(id: number) {
  try {
    await playerStore.deletePlayer(id)
    selectedPlayerIds.value = selectedPlayerIds.value.filter(pid => pid !== id)
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Cannot delete player'
  }
}

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await matchStore.createMatch(selectedPlayerIds.value, winningScore.value, hardMode.value)
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Failed to create game'
  } finally {
    loading.value = false
  }
}

async function toggleHardMode() {
  hardMode.value = !hardMode.value
}
</script>
