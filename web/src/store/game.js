import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getARAMChampionRankBrief, getClassicChampionRankBrief, getGameVersion } from '@/api/game'

export const useGameStore = defineStore('game', () => {
  const gameversion = ref()
  const classicBrief = ref([])
  const aramBrief = ref([])

  const setVersions = async() => {
    const res = await getGameVersion()
    if (res.code === 1) {
      gameversion.value = res.data
    }
  }

  const setARAMBrief = async(loc, version) => {
    const res = await getARAMChampionRankBrief(loc, version)
    if (res.code === 1) {
      aramBrief.value = res.data
    }
  }

  const setClassicBrief = async(loc, version) => {
    const res = await getClassicChampionRankBrief(loc, version)
    if (res.code === 1) {
      classicBrief.value = res.data
    }
  }

  return {
    gameversion,
    aramBrief,
    classicBrief,
    setVersions,
    setARAMBrief,
    setClassicBrief,
  }
})
