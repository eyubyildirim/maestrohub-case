<script setup lang="ts">
const props = defineProps<{
  selectionSize: number
}>()

const emit = defineEmits(['create-squad', 'deploy-individuals', 'clear-selection'])
</script>

<template>
  <div class="bg-slate-800 p-4 rounded-lg shadow-lg">
    <h3 class="text-lg font-bold text-slate-300 mb-2">Actions</h3>
    <p class="text-sm text-slate-400 mb-4">Selected Units: {{ selectionSize }}</p>
    <div class="space-y-2">
      <button
        @click="emit('create-squad')"
        :disabled="selectionSize === 0"
        class="w-full px-4 py-3 font-bold rounded-md transition-colors text-white bg-purple-600 hover:bg-purple-500 disabled:bg-slate-600 disabled:text-slate-400 disabled:cursor-not-allowed"
      >
        Create Squad Blueprint
      </button>
      <button
        @click="emit('deploy-individuals')"
        :disabled="selectionSize === 0"
        class="w-full px-4 py-3 font-bold rounded-md transition-colors text-white bg-green-600 hover:bg-green-500 disabled:bg-slate-600 disabled:text-slate-400 disabled:cursor-not-allowed"
      >
        Deploy Selected Individually
      </button>
      <button
        @click="emit('clear-selection')"
        :disabled="selectionSize === 0"
        class="w-full px-4 py-2 font-semibold rounded-md transition-colors text-slate-300 bg-slate-700 hover:bg-slate-600 disabled:opacity-50"
      >
        Clear Selection
      </button>
    </div>
  </div>
</template>
