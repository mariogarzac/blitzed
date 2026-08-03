<template>
  <div class="h-full flex flex-col bg-ef-bg-dim text-ef-fg" style="padding-bottom: env(safe-area-inset-bottom); padding-left: env(safe-area-inset-left); padding-right: env(safe-area-inset-right);">
    <NavBar />
    <main class="flex-1 overflow-y-auto max-w-3xl w-full mx-auto px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import NavBar from '@/components/common/NavBar.vue'
</script>
