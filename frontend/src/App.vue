<template>
  <div class="min-h-screen bg-ef-bg-dim text-ef-fg">
    <NavBar />
    <main class="max-w-3xl mx-auto px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import NavBar from '@/components/common/NavBar.vue'
</script>
