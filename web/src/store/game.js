import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getGameVersion } from '@/api/common'

export const useGameStore = defineStore('game', () => {
  const gameversion = ref()

  const setVersions = async() => {
    const res = await getGameVersion()
    if (res.code === 1) {
      gameversion.value = res.data
    }
  }

  return {
    gameversion,
    setVersions,
  }
})
