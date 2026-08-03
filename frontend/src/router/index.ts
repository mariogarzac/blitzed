import { createRouter, createWebHistory } from 'vue-router'
import GameView from '@/views/GameView.vue'
import MatchHistoryView from '@/views/MatchHistoryView.vue'
import MatchDetailView from '@/views/MatchDetailView.vue'
import StatsView from '@/views/StatsView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'Game', component: GameView },
    { path: '/history', name: 'MatchHistory', component: MatchHistoryView },
    { path: '/history/:id', name: 'MatchDetail', component: MatchDetailView },
    { path: '/stats', name: 'Stats', component: StatsView },
  ],
})
