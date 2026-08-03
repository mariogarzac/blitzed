<template>
  <div class="space-y-2">
    <label class="block text-sm font-medium text-ef-grey2">Players (2–8)</label>
    <div class="grid grid-cols-2 gap-2">
      <button
        v-for="player in players"
        :key="player.id"
        type="button"
        @click="toggle(player.id)"
        class="flex items-center gap-2 px-3 py-2 rounded-lg border text-sm font-medium transition-colors"
        :class="modelValue.includes(player.id)
          ? 'bg-ef-green border-ef-green text-ef-bg0'
          : 'border-ef-bg4 text-ef-grey2 hover:border-ef-grey1'"
      >
        <span
          class="w-4 h-4 rounded border flex items-center justify-center text-xs"
          :class="modelValue.includes(player.id) ? 'bg-white border-white text-ef-green' : 'border-stone-500'"
        >{{ modelValue.includes(player.id) ? '✓' : '' }}</span>
        {{ player.name }}
      </button>
    </div>
    <p v-if="modelValue.length > 0" class="text-xs text-ef-grey1">{{ modelValue.length }} selected</p>
  </div>
</template>

<script setup lang="ts">
import type { Player } from '@/types'

const props = defineProps<{ players: Player[]; modelValue: number[] }>()
const emit = defineEmits<{ 'update:modelValue': [value: number[]] }>()

function toggle(id: number) {
  const curr = [...props.modelValue]
  const idx = curr.indexOf(id)
  if (idx === -1) {
    if (curr.length < 8) curr.push(id)
  } else {
    curr.splice(idx, 1)
  }
  emit('update:modelValue', curr)
}
</script>
