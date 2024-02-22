<template>
  <div class="m-4">
    <el-container>
      <el-header w-full>
        <div class="selector">
          <span>Location</span>
          <el-select
            v-model="loc"
            class="child"
            placeholder="LOCATION"
            :default-first-option="true"
          >
            <el-option
              v-for="item in options1"
              :key="item.value"
              :label="item.label"
              :value="item.value"
              :disabled="item.disabled"
            />
          </el-select>
          <span>Version</span>
          <el-select
            v-model="ver"
            placeholder="versions"
            :default-first-option="true"
          >
            <el-option
              v-for="(version,index) in gameStore.gameversion"
              :key="index"
              :label="version"
              :value="version"
            />
          </el-select>
          <span>Tier</span>
          <el-select
            v-model="tier"
            placeholder="TIER"
            :default-first-option="true"
          >
            <el-option
              v-for="item in options2"
              :key="item.value"
              :label="item.label"
              :value="item.value"
              :disabled="item.disabled"
            />
          </el-select>
          <span>Queue</span>
          <el-select
            v-model="que"
            placeholder="LOCATION"
            :default-first-option="true"
          >
            <el-option
              v-for="item in options3"
              :key="item.value"
              :label="item.label"
              :value="item.value"
              :disabled="item.disabled"
            />
          </el-select>
        </div>

      </el-header>
      <el-main>
        <span />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { LOCATION_OPTIONS, TIERS_OPTIONS, QUECODES_OPTIONS } from '@/options'
import { onMounted, ref } from 'vue'
import { useGameStore } from '@/store/game'

const gameStore = useGameStore()
const loc = ref('na1')
const ver = ref()
const tier = ref('platinum')
const que = ref('ranked_solo_5x5')

const options1 = LOCATION_OPTIONS
const options2 = TIERS_OPTIONS
const options3 = QUECODES_OPTIONS

onMounted(async() => {
  await gameStore.setVersions()
  ver.value = gameStore.gameversion[0]
})
</script>

<style scoped>
.selector{
  @apply flex gap-2;
}
.selector .child{
  @apply w-100;
}

</style>
