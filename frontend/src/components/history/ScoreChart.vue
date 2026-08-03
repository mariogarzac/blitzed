<template>
  <div class="bg-ef-bg0 rounded-xl p-4">
    <canvas ref="chartCanvas" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import {
  Chart,
  LineController,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Legend,
  Tooltip,
} from 'chart.js'
import type { Match } from '@/types'

Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Legend, Tooltip)

const COLORS = [
  '#a7c080', // green
  '#7fbbb3', // blue
  '#dbbc7f', // yellow
  '#e67e80', // red
  '#83c092', // aqua
  '#d699b6', // purple
  '#e69875', // orange
  '#859289', // grey1
]

const props = defineProps<{ match: Match }>()

const chartCanvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

function buildChartData() {
  const rounds = props.match.rounds
  const labels = ['0', ...rounds.map(r => `R${r.round_num}`)]

  const datasets = props.match.players.map((player, i) => {
    const data: number[] = [0]
    for (const round of rounds) {
      const score = round.scores.find(s => s.player_id === player.id)
      data.push(score ? score.cumulative : (data[data.length - 1] ?? 0))
    }
    return {
      label: player.name,
      data,
      borderColor: COLORS[i % COLORS.length],
      backgroundColor: COLORS[i % COLORS.length] + '33',
      tension: 0.3,
      pointRadius: 4,
    }
  })

  // Dashed winning score line
  datasets.push({
    label: `Win (${props.match.winning_score})`,
    data: Array(labels.length).fill(props.match.winning_score),
    borderColor: '#4f5b62',
    backgroundColor: 'transparent',
    borderDash: [6, 4],
    tension: 0,
    pointRadius: 0,
  } as any)

  return { labels, datasets }
}

function initChart() {
  if (!chartCanvas.value) return
  const { labels, datasets } = buildChartData()
  chart = new Chart(chartCanvas.value, {
    type: 'line',
    data: { labels, datasets },
    options: {
      responsive: true,
      plugins: {
        legend: { position: 'bottom', labels: { color: '#9da9a0', boxWidth: 12 } },
        tooltip: { mode: 'index', intersect: false },
      },
      scales: {
        x: { ticks: { color: '#7a8478' }, grid: { color: '#323d43' } },
        y: { ticks: { color: '#7a8478' }, grid: { color: '#323d43' } },
      },
    },
  })
}

onMounted(initChart)

watch(() => props.match.rounds.length, () => {
  if (!chart) return
  const { labels, datasets } = buildChartData()
  chart.data.labels = labels
  chart.data.datasets = datasets
  chart.update()
})

onUnmounted(() => { chart?.destroy() })
</script>
