<script setup lang="ts">
import type { GameState } from '@/types/game'
import Icons from './Icons.vue'

defineProps<{
  gameState: GameState
}>()

const getHealthPercent = (hp: number, maxHp: number) => {
  if (!maxHp || maxHp <= 0) return 0
  return (hp / maxHp) * 100
}
</script>

<template>
  <div class="bg-slate-800 p-4 rounded-lg shadow-lg">
    <h2 class="text-xl font-bold text-slate-200 mb-4">Weapon Systems</h2>
    <div class="space-y-4">
      <!-- Individual Weapon Status -->
      <div>
        <h3 class="font-bold text-slate-300 flex items-center gap-2">
          <Icons type="individual" /> Individual Weapon
        </h3>
        <div v-if="gameState.processingIndividualUnit" class="mt-2">
          <div class="flex justify-between items-center text-sm text-slate-400">
            <span>Target: <Icons :type="gameState.processingIndividualUnit.type" /></span>
            <span class="font-mono">
              HP: {{ gameState.processingIndividualUnit.hp }} /
              {{ gameState.processingIndividualUnit.maxHp }}
            </span>
          </div>
          <div class="w-full bg-slate-700 rounded-full h-2.5 mt-1">
            <div
              class="bg-green-500 h-2.5 rounded-full transition-all duration-200"
              :style="{
                width:
                  getHealthPercent(
                    gameState.processingIndividualUnit.hp,
                    gameState.processingIndividualUnit.maxHp,
                  ) + '%',
              }"
            ></div>
          </div>
        </div>
        <p v-else class="text-slate-500 italic mt-2">Status: IDLE</p>
      </div>
      <!-- Group Weapon Status -->
      <div>
        <h3 class="font-bold text-slate-300 flex items-center gap-2">
          <Icons type="group" /> Group Weapon
        </h3>
        <div v-if="gameState.processingGroup" class="mt-2 space-y-2">
          <div v-for="unit in gameState.processingGroup.units" :key="unit.id">
            <div class="flex justify-between items-center text-sm text-slate-400">
              <span>Target: <Icons :type="unit.type" /></span>
              <span class="font-mono">HP: {{ unit.hp }} / {{ unit.maxHp }}</span>
            </div>
            <div class="w-full bg-slate-700 rounded-full h-2.5 mt-1">
              <div
                class="bg-purple-500 h-2.5 rounded-full transition-all duration-200"
                :style="{ width: getHealthPercent(unit.hp, unit.maxHp) + '%' }"
              ></div>
            </div>
          </div>
        </div>
        <p v-else class="text-slate-500 italic mt-2">Status: IDLE</p>
      </div>
    </div>
  </div>
</template>
