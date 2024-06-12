<template>
  <div class="container">
    <el-header class="selector">
      <span>Location</span>
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
      <span>Version</span>
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
    </el-header>
    <el-main class="table-container">
      <el-table
        :data="brief"
        :row-style="rowstyle"
      >
        <el-table-column
          type="index"
          label="排名"
          align="center"
          width="60"
        />

        <el-table-column
          prop="id"
          label="英雄"
          align="center"
          width="200"
          sortable
        >
          <template #default="{ row }">
            <div style="display: flex; align-items: center;">
              <div style="width: 36px; height: 48px;" />
              <router-link :to="`/herodetail?hero=${row.id}&loc=${loc}&mode=aram`">
                <el-image
                  :src="`src/assets/datadragon/champion/${row.id}.png`"
                  style="width: 48px; height: 48px;"
                />
              </router-link>
              <span style="margin-left: 10px">{{ row.id }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column
          prop="win_rate"
          label="胜率"
          align="center"
          sortable
        >
          <template #default="{ row }">
            {{ (row.win_rate * 100).toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column
          prop="pick_rate"
          label="登场率 "
          align="center"
          sortable
        >
          <template #default="{ row }">
            {{ (row.pick_rate * 100).toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column
          prop="avg_damage_dealt"
          label="场均伤害"
          align="center"
          sortable
        >
          <template #default="{ row }">
            {{ (row.avg_damage_dealt).toFixed(0) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="avg_dead_time"
          label="场均死亡"
          align="center"
          sortable
        >
          <template #default="{ row }">
            {{ (row.avg_dead_time).toFixed() }}秒
          </template>
        </el-table-column>

      </el-table>
    </el-main>
    <el-footer class="footer" />
  </div>
</template>

<script setup>
import { LOCATION_OPTIONS } from '@/options'
import { computed, onMounted, ref } from 'vue'
import { useGameStore } from '@/store/game'

const gameStore = useGameStore()
const loc = ref('na1')
const ver = ref()

const brief = computed(() => gameStore.aramBrief)

const options1 = LOCATION_OPTIONS

onMounted(async() => {
  await gameStore.setVersions()
  ver.value = gameStore.gameversion[0]
  await gameStore.setARAMBrief(loc.value, ver.value)
})

const handleSelectChange = async() => {
  await gameStore.setARAMBrief(loc.value, ver.value)
  console.log(gameStore.aramBrief)
}

const rowstyle = ({ rowIndex }) => {
  if (rowIndex % 2 === 0) {
    return {
      backgroundColor: '#9ca3af',
    }
  }
}

</script>

<style scoped>
.container {
  @apply w-full h-screen ;
}

.selector {
  @apply flex items-center justify-start gap-4 text-gray-300;
}

.selector .child {
  @apply w-30;
}

.table-container {
  @apply w-screen-md;
}

.footer {
  @apply h-40;
}

.el-table {
  --el-table-border-color: #d1d5db;
  --el-table-border: transparent;
  --el-table-text-color: #f9fafb;
  --el-table-header-text-color: #d1d5db;
  --el-table-row-hover-bg-color: transparent;
  --el-table-current-row-bg-color: transparent;
  --el-table-header-bg-color: transparent;
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-expanded-cell-bg-color: transparent;
}

</style>
