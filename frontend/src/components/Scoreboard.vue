<script setup lang="ts">
import type { GameState } from '@/types/game'

defineProps<{
  gameState: GameState
}>()
</script>

<template>
  <div class="bg-slate-800 p-6 rounded-lg shadow-lg mt-4">
    <h2 class="text-xl font-bold text-slate-200 mb-4">Commander Stats</h2>
    <div class="grid grid-cols-3 gap-4 text-center">
      <div>
        <p class="text-slate-400 text-sm">Score</p>
        <p class="text-3xl font-mono font-bold text-cyan-400">{{ gameState.score }}</p>
      </div>
      <div>
        <p class="text-slate-400 text-sm">Escaped</p>
        <p class="text-3xl font-mono font-bold text-yellow-500">{{ gameState.escapedCount }}</p>
      </div>
      <div>
        <p class="text-slate-400 text-sm">Breached</p>
        <p class="text-3xl font-mono font-bold text-red-500">{{ gameState.breachCount }}</p>
      </div>
    </div>
  </div>
</template>
