<script setup lang="ts">
import { computed } from 'vue'
import type { GameState, Unit } from '@/types/game'
import Icons from './Icons.vue'

const props = defineProps<{
  gameState: GameState
  selectedUnits: Set<string>
}>()

const emit = defineEmits(['toggle-selection'])

const units = computed<Unit[]>(() => {
  if (!props.gameState?.availableUnits) return []
  const allUnits = Object.values(props.gameState.availableUnits)
  allUnits.sort((a, b) => a.creationTime - b.creationTime)
  return allUnits
})

const getTtlColor = (ttl: number): string => {
  if (ttl <= 5) return 'text-red-500 animate-pulse'
  if (ttl <= 10) return 'text-yellow-400'
  return 'text-slate-400'
}

const isSelected = (unitId: string) => props.selectedUnits.has(unitId)
</script>

<template>
  <div class="bg-slate-800 p-4 rounded-lg shadow-lg">
    <h2 class="text-xl font-bold text-slate-200 mb-4">
      Battlefield Status ({{ units.length }} units)
    </h2>
    <div
      class="bg-slate-900/50 rounded p-4 min-h-[200px] flex flex-wrap items-start justify-start gap-4"
    >
      <template v-if="units.length > 0">
        <div
          v-for="unit in units"
          :key="unit.id"
          @click="emit('toggle-selection', unit.id)"
          class="relative transition-all duration-150 cursor-pointer p-2 rounded-md flex flex-col items-center"
          :class="{
            'bg-sky-500/50 scale-110 ring-2 ring-sky-400': isSelected(unit.id),
            'hover:bg-sky-500/20': !isSelected(unit.id),
          }"
        >
          <span class="text-4xl"><Icons :type="unit.type" /></span>
          <span class="text-xs font-mono font-bold" :class="getTtlColor(unit.ttl)">
            {{ unit.ttl.toFixed(1) }}s
          </span>
        </div>
      </template>
      <p v-else class="text-slate-400 text-center w-full self-center">Battlefield Clear.</p>
    </div>
  </div>
</template>
