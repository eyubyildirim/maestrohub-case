<script setup lang="ts">
import type { UnitType } from '@/types/game'

defineProps<{
  type: UnitType | 'group' | 'individual'
}>()

const icons: Record<string, string> = {
  soldier: '👤',
  tank: '🚗',
  helicopter: '🚁',
  group: '📦',
  individual: '🎯',
}
</script>

<template>
  <span>{{ icons[type] || '❓' }}</span>
</template>
