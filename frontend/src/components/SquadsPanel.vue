<script setup lang="ts">
import type { GameState, Squad, UnitType } from '@/types/game'
import Icons from './Icons.vue'

const props = defineProps<{
  gameState: GameState
}>()

const emit = defineEmits(['deploy-squad'])
</script>

<template>
  <div class="bg-slate-800 p-4 rounded-lg shadow-lg">
    <h2 class="text-xl font-bold text-slate-200 mb-4">Saved Squads</h2>
    <div v-if="Object.keys(props.gameState.squads).length > 0" class="space-y-4">
      <div
        v-for="squad in props.gameState.squads"
        :key="squad.id"
        class="bg-slate-700/50 p-3 rounded-lg flex justify-between items-center"
      >
        <div class="flex gap-3 items-center">
          <span
            v-for="(count, type) in squad.composition"
            :key="type"
            class="flex items-center gap-1 font-mono text-slate-300"
          >
            <Icons :type="type as UnitType" />
            <span>x{{ count }}</span>
          </span>
        </div>
        <button
          @click="emit('deploy-squad', squad.id)"
          class="px-4 py-2 font-bold rounded-md transition-colors text-white bg-blue-600 hover:bg-blue-500"
        >
          Deploy
        </button>
      </div>
    </div>
    <p v-else class="text-slate-400 text-center">
      No squad blueprints saved. Select units and create one.
    </p>
  </div>
</template>
