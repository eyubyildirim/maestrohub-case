<!-- src/components/UnitPool.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import type { GameState, Unit, UnitType } from '@/types/game'

const props = defineProps<{
  gameState: GameState
  selectedUnits: Set<string>
}>()

const emit = defineEmits(['toggle-selection'])

const icons: Record<UnitType, string> = {
  soldier: '👤',
  tank: '🚗',
  helicopter: '🚁',
}

const unitsByType = computed(() => {
  const grouped: Record<UnitType, Unit[]> = {
    soldier: [],
    tank: [],
    helicopter: [],
  }
  if (props.gameState.availableUnits) {
    for (const unit of Object.values(props.gameState.availableUnits)) {
      if (grouped[unit.type]) {
        grouped[unit.type].push(unit)
      }
    }
  }
  return grouped
})

const isSelected = (unitId: string) => props.selectedUnits.has(unitId)
</script>

<template>
  <div class="bg-slate-800 p-4 rounded-lg shadow-lg">
    <h3 class="text-lg font-bold text-slate-300 mb-4">Available Unit Pool</h3>
    <div class="space-y-4">
      <div v-for="(units, type) in unitsByType" :key="type">
        <h4 class="font-bold text-slate-400 mb-2 flex items-center gap-2">
          <span class="text-2xl">{{ icons[type] }}</span>
          <span>{{ type.charAt(0).toUpperCase() + type.slice(1) }}s ({{ units.length }})</span>
        </h4>
        <div v-if="units.length > 0" class="flex flex-wrap gap-2">
          <div
            v-for="unit in units"
            :key="unit.id"
            @click="emit('toggle-selection', unit.id)"
            class="p-1.5 rounded-md cursor-pointer transition-colors"
            :class="isSelected(unit.id) ? 'bg-sky-500' : 'bg-slate-700 hover:bg-slate-600'"
          >
            <!-- Just show a small icon, the main info is in the header -->
            <span class="text-xl">{{ icons[unit.type] }}</span>
          </div>
        </div>
        <p v-else class="text-sm text-slate-500 italic">None available.</p>
      </div>
    </div>
  </div>
</template>
