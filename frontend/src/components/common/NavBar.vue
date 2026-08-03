<template>
  <nav class="bg-ef-bg0 text-ef-fg shadow-lg">
    <div class="max-w-3xl mx-auto px-4 py-3 flex items-center justify-between">
      <span class="text-xl font-bold tracking-tight text-ef-green">🃏 Blitzed</span>
      <div class="flex gap-1">
        <RouterLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
          :class="$route.path === link.to || $route.path.startsWith(link.prefix ?? link.to + '/')
            ? 'bg-ef-green text-ef-bg0'
            : 'text-ef-grey2 hover:bg-ef-bg2'"
        >
          {{ link.label }}
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
const links = [
  { to: '/', label: 'Game' },
  { to: '/history', label: 'History', prefix: '/history' },
  { to: '/stats', label: 'Stats' },
]
</script>
