import { ref, onMounted, onUnmounted } from 'vue'
import type { Command, GameState } from '@/types/game'

const WEBSOCKET_URL = `ws://${window.location.host}/ws`

export function useGameWebSocket() {
  const gameState = ref<GameState | null>(null)
  const isConnected = ref(false)

  let socket: WebSocket | null = null

  const sendCommand = (command: Command) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(command))
    } else {
      console.error('Cannot send command: WebSocket is not connected.')
    }
  }
  const connect = () => {
    socket = new WebSocket(WEBSOCKET_URL)

    socket.onopen = () => {
      console.log('Battlefield connection established.')
      isConnected.value = true
    }

    socket.onclose = () => {
      console.log('Battlefield connection lost.')
      isConnected.value = false
      gameState.value = null
    }
    socket.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnected.value = false
    }

    socket.onmessage = (event) => {
      try {
        const state = JSON.parse(event.data) as GameState
        gameState.value = state
      } catch (error) {
        console.error('Failed to parse game state:', error)
      }
    }
  }

  const disconnect = () => {
    if (socket) socket.close()
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return {
    gameState: gameState,
    isConnected: isConnected,
    sendCommand,
  }
}
