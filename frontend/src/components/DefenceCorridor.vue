<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { GameState, Group, Unit, UnitType } from '@/types/game'
import TowerSprite from '@/assets/tower.png'
import Icons from './Icons.vue'

const props = defineProps<{
  gameState: GameState | null
}>()

const BORDER_POSITION_PERCENT = 80

const getPosition = (ttl: number, maxTtl: number) => {
  if (!maxTtl || maxTtl <= 0) return '0%'
  const progress = (maxTtl - ttl) / maxTtl
  return `${progress * BORDER_POSITION_PERCENT}%`
}

const processedGroups = computed(() => {
  if (!props.gameState?.groupsInTransit) return []
  return Object.values(props.gameState.groupsInTransit).map((group) => {
    const unitCounts = group.units.reduce(
      (acc, unit) => {
        acc[unit.type] = (acc[unit.type] || 0) + 1
        return acc
      },
      {} as Record<string, number>,
    )
    return {
      id: group.id,
      position: getPosition(group.transitTtl, group.maxTransitTtl),
      unitTypes: Object.entries(unitCounts).map(([type, count]) => ({ type, count })),
    }
  })
})

const individualsInTransit = computed(() => {
  if (!props.gameState?.individualsInTransit) return []
  return Object.values(props.gameState.individualsInTransit)
})

const recentDestructions = ref<Set<string>>(new Set())
const recentBreaches = ref<Set<string>>(new Set())
const isShootingIndividual = ref(false)
const isShootingGroup = ref(false)

watch(
  () => props.gameState?.lastEvents,
  (newEvents) => {
    if (!newEvents || newEvents.length === 0) return

    newEvents.forEach((event) => {
      const [eventType, eventId] = event.split(':')
      switch (eventType) {
        case 'DESTROY_INDIVIDUAL':
          isShootingIndividual.value = true
          setTimeout(() => (isShootingIndividual.value = false), 200)
          recentDestructions.value.add(eventId)
          setTimeout(() => recentDestructions.value.delete(eventId), 500)
          break
        case 'DESTROY_GROUP':
          isShootingGroup.value = true
          setTimeout(() => (isShootingGroup.value = false), 200)
          recentDestructions.value.add(eventId)
          setTimeout(() => recentDestructions.value.delete(eventId), 500)
          break
        case 'BREACH_INDIVIDUAL':
        case 'BREACH_GROUP':
          recentBreaches.value.add(eventId)
          setTimeout(() => recentBreaches.value.delete(eventId), 500)
          break
      }
    })
  },
)

const isTargetedByIndividual = (id: string) => props.gameState?.processingIndividualUnit?.id === id
const isTargetedByGroup = (id: string) => props.gameState?.processingGroup?.id === id
const isDestroyed = (id: string) => recentDestructions.value.has(id)
const isBreached = (id: string) => recentBreaches.value.has(id)
</script>

<template>
  <div class="bg-slate-800 p-2 rounded-lg shadow-lg h-full">
    <h3 class="text-center font-bold text-slate-300 mb-2">Defense Corridor</h3>
    <div class="relative w-full h-[700px] bg-slate-900/50 rounded overflow-hidden">
      <!-- Deployed Individual Units (Right side of corridor) -->
      <div
        v-for="unit in individualsInTransit"
        :key="unit.id"
        class="absolute -translate-x-1/2 transition-all duration-200"
        :class="{
          'opacity-0 scale-150': isDestroyed(unit.id),
          'animate-breach': isBreached(unit.id),
        }"
        :style="{ top: getPosition(unit.transitTtl, unit.maxTransitTtl), left: '66%' }"
      >
        <div
          v-if="isTargetedByIndividual(unit.id)"
          class="absolute -inset-2 border-2 border-cyan-400 rounded-full animate-ping"
        ></div>
        <span class="text-4xl relative z-10"><Icons :type="unit.type" /></span>
      </div>

      <!-- Deployed Groups (Left side of corridor) -->
      <div
        v-for="group in processedGroups"
        :key="group.id"
        class="absolute -translate-x-1/2 transition-all duration-200"
        :class="{
          'opacity-0 scale-150': isDestroyed(group.id),
          'animate-breach': isBreached(group.id),
        }"
        :style="{ top: group.position, left: '33%' }"
      >
        <div
          v-if="isTargetedByGroup(group.id)"
          class="absolute -inset-2 border-2 border-purple-400 rounded-lg animate-ping"
        ></div>
        <div
          v-for="unitType in group.unitTypes"
          :key="unitType.type"
          class="relative z-10 text-center"
        >
          <span class="text-4xl"><Icons :type="unitType.type as UnitType" /></span>
          <span
            v-if="unitType.count > 1"
            class="absolute -top-1 -right-2 text-xs font-bold bg-red-600 text-white rounded-full px-1.5 py-0.5"
          >
            x{{ unitType.count }}
          </span>
        </div>
      </div>

      <!-- The Border Line -->
      <div
        class="absolute w-full h-1 bg-red-500/70 shadow-[0_0_10px_theme(colors.red.500)]"
        :style="{ top: `${BORDER_POSITION_PERCENT}%` }"
      ></div>

      <!-- The Tower and its Animations -->
      <div
        class="absolute bottom-0 left-1/2 -translate-x-1/2 w-36"
        :class="{ 'animate-shake': isShootingIndividual || isShootingGroup }"
      >
        <img :src="TowerSprite" alt="Defense Tower" class="w-full pointer-events-none" />

        <!-- LOCK-ON GLOW: Group Weapon (Left Prong) -->
        <div
          v-if="gameState?.processingGroup"
          class="absolute top-[18px] left-[15px] w-12 h-12 animate-lock-on-glow"
        ></div>

        <!-- LOCK-ON GLOW: Individual Weapon (Right Prong) -->
        <div
          v-if="gameState?.processingIndividualUnit"
          class="absolute top-[40px] right-[13px] w-10 h-10 animate-lock-on-glow"
        ></div>

        <!-- MUZZLE FLASH: Group Weapon (Left Prong) -->
        <div
          v-if="isShootingGroup"
          class="absolute top-[14px] left-[-15px] w-12 h-12 animate-muzzle-flash"
        ></div>

        <!-- MUZZLE FLASH: Individual Weapon (Right Prong) -->
        <div
          v-if="isShootingIndividual"
          class="absolute top-[43px] right-[-15px] w-10 h-10 animate-muzzle-flash"
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Keyframes for Breach Effect */
@keyframes breach-effect {
  50% {
    transform: scale(1.2) rotate(10deg);
    filter: drop-shadow(0 0 0.75rem theme(colors.red.500));
  }
  100% {
    transform: scale(0) rotate(-10deg);
    opacity: 0;
  }
}
.animate-breach {
  animation: breach-effect 0.5s ease-in-out forwards;
}

/* Keyframes for Tower Shake (Recoil) */
@keyframes shake-effect {
  0%,
  100% {
    transform: translate(0, 0);
  }
  25% {
    transform: translate(1px, -1px);
  }
  75% {
    transform: translate(-1px, 1px);
  }
}
.animate-shake {
  animation: shake-effect 0.1s linear;
}

/* Keyframes for Muzzle Flash */
@keyframes muzzle-flash-effect {
  from {
    opacity: 1;
    transform: scale(1.5);
    background-color: theme(colors.yellow.200);
    border-radius: 50%;
    filter: blur(5px);
  }
  to {
    opacity: 0;
    transform: scale(0);
  }
}
.animate-muzzle-flash {
  animation: muzzle-flash-effect 0.15s ease-out forwards;
}

/* Keyframes for Weapon Lock-on Glow */
@keyframes lock-on-glow-effect {
  0%,
  100% {
    box-shadow: 0 0 15px 5px theme(colors.cyan.500 / 0.3);
    opacity: 0.7;
  }
  50% {
    box-shadow: 0 0 25px 10px theme(colors.cyan.400 / 0.5);
    opacity: 1;
  }
}
.animate-lock-on-glow {
  background-color: theme(colors.cyan.500 / 0.1);
  border-radius: 50%;
  animation: lock-on-glow-effect 1.5s infinite;
}
</style>
