<script setup lang="ts">
import { ref } from 'vue'
import { useGameWebSocket } from './composables/useGameWebSocket'
import type { CreateSquadPayload, DeployIndividualsPayload, DeploySquadPayload } from './types/game'

import Battlefield from './components/Battlefield.vue'
import SquadsPanel from './components/SquadsPanel.vue'
import Scoreboard from './components/Scoreboard.vue'
import ControlPanel from './components/ControlPanel.vue'
import WeaponStatus from './components/WeaponStatus.vue'
import DefenseCorridor from './components/DefenceCorridor.vue'

const { gameState, isConnected, sendCommand } = useGameWebSocket()

const selectedUnitIds = ref<Set<string>>(new Set())

const handleToggleSelection = (unitId: string) => {
  const newSet = new Set(selectedUnitIds.value)
  if (newSet.has(unitId)) {
    newSet.delete(unitId)
  } else {
    newSet.add(unitId)
  }
  selectedUnitIds.value = newSet
}

const handleClearSelection = () => {
  selectedUnitIds.value = new Set()
}

const handleCreateSquad = () => {
  if (selectedUnitIds.value.size === 0) return
  const payload: CreateSquadPayload = {
    unit_ids: Array.from(selectedUnitIds.value),
  }
  sendCommand({ action: 'create_squad', payload })
  handleClearSelection()
}

const handleDeployIndividuals = () => {
  if (selectedUnitIds.value.size === 0) return
  const payload: DeployIndividualsPayload = {
    unit_ids: Array.from(selectedUnitIds.value),
  }
  sendCommand({ action: 'deploy_individuals', payload })
  handleClearSelection()
}

const handleDeploySquad = (squadId: string) => {
  const payload: DeploySquadPayload = { squad_id: squadId }
  sendCommand({ action: 'deploy_squad', payload })
}
</script>

<template>
  <div class="bg-slate-900 text-white min-h-screen font-sans">
    <header class="py-6">
      <h1 class="text-4xl font-bold text-center tracking-wider">Tower Defense: DataLand</h1>
    </header>

    <main class="container mx-auto p-4">
      <!-- Loading / Disconnected States -->
      <div v-if="!isConnected" class="text-center text-2xl text-yellow-400 animate-pulse">
        <h2>Connecting to Command Center...</h2>
      </div>
      <div v-else-if="!gameState" class="text-center text-2xl text-slate-400">
        <h2>Awaiting initial battlefield report...</h2>
      </div>

      <!-- Main Game Layout (2-Column on Large Screens) -->
      <div v-else class="flex flex-col lg:flex-row gap-4">
        <!-- === LEFT COLUMN: Main Interaction Panels === -->
        <div class="flex flex-col gap-4 flex-grow">
          <!-- Battlefield Panel -->
          <Battlefield
            :game-state="gameState"
            :selected-units="selectedUnitIds"
            @toggle-selection="handleToggleSelection"
          />

          <!-- Saved Squads Panel -->
          <SquadsPanel :game-state="gameState" @deploy-squad="handleDeploySquad" />

          <!-- Scoreboard and Actions (in a 2-column grid) -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Scoreboard :game-state="gameState" />
            <ControlPanel
              :selection-size="selectedUnitIds.size"
              @create-squad="handleCreateSquad"
              @deploy-individuals="handleDeployIndividuals"
              @clear-selection="handleClearSelection"
            />
          </div>

          <!-- Weapon Status Panel -->
          <WeaponStatus :game-state="gameState" />
        </div>

        <!-- === RIGHT COLUMN: Defense Corridor Visualization (Alone) === -->
        <div class="w-full lg:w-96 flex-shrink-0">
          <DefenseCorridor :game-state="gameState" />
        </div>
      </div>
    </main>
  </div>
</template>
