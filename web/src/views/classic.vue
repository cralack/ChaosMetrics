<template>
  <div class="min-h-screen">
    <el-row :gutter="20">
      <div class="selector">
        <!--          <span>Location</span>-->
        <el-select
          v-model="loc"
          class="child"
          placeholder="Location"
          :default-first-option="true"
          @change="handleSelectChange"
        >
          <el-option
            v-for="item in options1"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
          />
        </el-select>
        <!--          <span>Version</span>-->
        <el-select
          v-model="ver"
          class="child"
          placeholder="Versions"
          :default-first-option="true"
          @change="handleSelectChange"
        >
          <el-option
            v-for="(version,index) in gameStore.gameversion"
            :key="index"
            :label="version"
            :value="version"
          />
        </el-select>
        <!--          <span>Tier</span>-->
        <el-select
          v-model="tier"
          class="child"
          placeholder="Tier"
          :default-first-option="true"
          @change="handleSelectChange"
        >
          <el-option
            v-for="item in options2"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
          />
        </el-select>
        <!--          <span>Queue</span>-->
        <el-select
          v-model="que"
          class="child"
          placeholder="Que"
          :default-first-option="true"
          @change="handleSelectChange"
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
    </el-row>

    <el-scrollbar class="scrollbar-container">
      <el-row :gutter="20">
        <el-col
          v-for="champion in brief"
          :key="champion.id"
          :span="24"
        >
          <el-card class="card">
            <!--          <img-->
            <!--            :src="getImageUrl(champion.image.full)"-->
            <!--            class="image"-->
            <!--            alt="champion-image"-->
            <!--          >-->
            <div style="padding: 14px;">
              <h5>{{ champion.id }}</h5>
              <p>胜率: {{ (champion.win_rate * 100).toFixed(2) }}%</p>
              <p>选择率: {{ (champion.pick_rate * 100).toFixed(2) }}%</p>
              <p>Ban率: {{ (champion.ban_rate * 100).toFixed(2) }}%</p>
              <p>场均输出占比: {{ (champion.avg_damage_dealt) }}</p>
              <p>场均死亡时长: {{ (champion.avg_dead_time ) }}</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-scrollbar>
  </div>
</template>

<script setup>
import { LOCATION_OPTIONS, TIERS_OPTIONS, QUECODES_OPTIONS } from '@/options'
import { computed, onMounted, ref } from 'vue'
import { useGameStore } from '@/store/game'

const gameStore = useGameStore()
const loc = ref('na1')
const ver = ref()
const tier = ref('platinum')
const que = ref('ranked_solo_5x5')

const brief = computed(() => gameStore.classicBrief)

const options1 = LOCATION_OPTIONS
const options2 = TIERS_OPTIONS
const options3 = QUECODES_OPTIONS

onMounted(async() => {
  await gameStore.setVersions()
  ver.value = gameStore.gameversion[0]
  await gameStore.setClassicBrief(loc.value, '14.1.1')
})

const handleSelectChange = async() => {
  await gameStore.setClassicBrief(loc.value, ver.value)
  console.log(gameStore.classicBrief)
}

</script>

<style scoped>
.selector{
  @apply flex flex-row items-start justify-center gap-2 mb-10;
}
.selector .child{
  @apply w-40;
}

.scrollbar-container{
  @apply h-screen overflow-y-auto;
}
.el-card {
  @apply mb-4 ;
}

.card{
  @apply flex justify-between items-start p-4;
}

</style>
