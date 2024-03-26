<template>
  <div
    v-if="detail"
    class="container"
  >
    <el-row class="flex">
      <el-col class="coltainer">
        <el-image
          class="hero-image"
          :src="heroImage"
          :fit="'contain'"
        />
        <el-col>
          <h1>{{ detail.name + " " + detail.title }}</h1>
        </el-col>
        <el-col>
          {{ hero.value }}
          <el-image :src="`src/assets/datadragon/spell/${detail.id}Q.png`" />
        </el-col>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { computed, onMounted } from 'vue'
import { ref } from 'vue'
import { useGameStore } from '@/store/game'

const gameStore = useGameStore()
const hero = ref('')
const mode = ref('')
const loc = ref('')
const version = ref()
const detail = computed(() => gameStore.detail)

onMounted(async() => {
  const route = useRoute()
  await gameStore.setVersions()

  hero.value = route.query.hero
  mode.value = route.query.mode
  loc.value = route.query.loc
  version.value = gameStore.gameversion[0]

  await gameStore.setHeroDetail(hero.value, loc.value, mode.value, version.value)
})

const heroImage = computed(() => {
  return 'src/assets/datadragon/champion_og/loading/' + hero.value + '_0.jpg'
})

</script>

<style scoped>
.container {
  @apply w-full h-screen;
}

.hero-image {
  @apply w-50;
}

</style>
