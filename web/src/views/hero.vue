<template>
  <div>
    <el-container
      v-if="hero"
      class="container"
    >
      <el-aside class="left-container">
        <el-tooltip
          placement="bottom"
          effect="dark"
          :visible-arrow="false"
        >
          <template #content>
            <div class="tooltip-content">
              <h3>
                {{ hero.blurb }}
              </h3>
            </div>
          </template>
          <el-image
            class="hero-image"
            :src="heroImage"
            :fit="'contain'"
          />
        </el-tooltip>
        <h3 class="text-bg">{{ hero.name + " " + hero.title }}</h3>
        <el-divider class="my-divider" />
        <div class="skills-container">
          <span class="text-bg">技能描述</span>
          <el-tooltip
            placement="top"
            effect="dark"
            :visible-arrow="false"
          >
            <template #content>
              <div
                class="tooltip-content"
              >
                <h3>{{ hero.passive.name }} P</h3>
                <el-divider class="my-divider" />
                <span>{{ hero.passive.description }}</span>
              </div>
            </template>
            <el-image
              class="skill-icon"
              :src="getPassiveImageUrl(hero.passive)"
            />
          </el-tooltip>
          <div
            v-for="(spell) in hero.spells"
            :key="spell.id"
          >
            <el-tooltip
              placement="top"
              effect="dark"
              :visible-arrow="false"
            >
              <template #content>
                <div class="tooltip-content">
                  <h3>{{ spell.name + " "+"快捷键:" + getSkillShortId(spell.id) }}</h3>
                  <el-divider class="my-divider" />
                  <span>技能消耗 ：{{ spell.costBurn }}</span><br>
                  <span>冷却时间(秒）：{{ spell.cooldownBurn }}</span><br>
                  <span>范围 ：{{ spell.rangeBurn }}</span><br><br>
                  <span>{{ spell.tooltip }}</span>
                </div>
              </template>
              <el-image
                class="skill-icon"
                :src="getSpellImageUrl(spell)"
              />
            </el-tooltip>

          </div>
        </div>
        <el-divider class="my-divider" />
        <div class="flex">
          <h3 class="text-bg">英雄难度: </h3>
          <div
            v-for="index in hero.info.difficulty"
            :key="index"
            class="mx-0.5"
          >
            <el-icon><SemiSelect /></el-icon>
          </div>
        </div>
        <div class="flex items-center">
          <h3 class="text-bg">英雄定位: </h3>
          <div
            v-for="tag in hero.tags"
            :key="tag"
            class="tag-icon"
          >
            <el-tooltip
              placement="bottom"
              effect="dark"
              :visible-arrow="false"
              :content="tag"
            >
              <el-image :src="getTagImageUrl(tag)" />
            </el-tooltip>
          </div>
        </div>
        <el-divider class="my-divider" />
        <div class="data-container">
          <h2 class="text-bg">主要数据-</h2>
          <h3 class="text-bg">登场率: {{ (heroData?.pick_rate * 100).toFixed(2) }}%</h3>
          <h3 class="text-bg">胜率:{{ (heroData?.win_rate * 100).toFixed(2) }}%</h3>
          <h3
            v-if="mode==='classic'"
            class="text-bg"
          >ban率:{{ (heroData?.ban_rate * 100).toFixed(2) }}%</h3>
          <el-divider class="my-divider" />
          <h2 class="text-bg">场均数据-</h2>
          <h3 class="text-bg">KDA: {{ (heroData?.avg_kda ?? 0).toFixed(2) }} </h3>
          <h3 class="text-bg">参战率: {{ (heroData?.avg_kp*100 ?? 0).toFixed(2) }} </h3>
          <h3 class="text-bg">伤害: {{ (heroData?.avg_damage_dealt ?? 0).toFixed(0) }} </h3>
          <h3 class="text-bg">承伤: {{ (heroData?.avg_damage_taken ?? 0).toFixed(0) }} </h3>
          <h3 class="text-bg">控制时长: {{ (heroData?.avg_time_ccing ?? 0).toFixed(0) }} </h3>
          <h3 class="text-bg">死亡时长: {{ (heroData?.avg_dead_time ?? 0).toFixed(0) }} </h3>
          <h3
            v-if="mode==='classic'"
            class="text-bg"
          >场均视野得分: {{ (heroData?.avg_vision_score ?? 0).toFixed(0) }} </h3>
        </div>
        <el-divider class="my-divider" />
      </el-aside>
      <el-main class="main-container">
        Main
      </el-main>
    </el-container>

  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import { useGameStore } from '@/store/game'
import { useUserStore } from '@/store/user'
import { getHeroData, getHeroDetail } from '@/api/game'

const userStore = useUserStore()
const gameStore = useGameStore()
const heroName = ref('')
const mode = ref('')
const loc = ref('')
const version = ref()

const hero = ref()
const heroData = ref()

onMounted(async() => {
  const route = useRoute()
  await gameStore.setVersions()

  heroName.value = route.query.hero
  mode.value = route.query.mode
  loc.value = route.query.loc
  version.value = gameStore.gameversion[0]

  await setHero()
  console.log(hero.value)
  console.log(heroData.value)
})

const setHero = async() => {
  const res = await getHeroDetail(heroName.value, version.value, userStore.lang)
  if (res.code === 1) {
    hero.value = res.data
  }

  const res1 = await getHeroData(heroName.value, loc.value, mode.value, version.value)
  if (res1.code === 1) {
    heroData.value = res1.data
  }
}

function getSkillShortId(skillId) {
  return skillId.replace(hero.value.id, '')
}

// 获取技能图片 URL 的方法
const getSpellImageUrl = (spell) => {
  return `src/assets/datadragon/spell/${spell.image.full}`
}

// 获取被动技能图片 URL 的方法
const getPassiveImageUrl = (passive) => {
  return `src/assets/datadragon/passive/${passive.image.full}`
}

const getTagImageUrl = (tag) => {
  return `src/assets/tags/${tag}.png`
}

const heroImage = computed(() => {
  return 'src/assets/datadragon/champion_og/loading/' + heroName.value + '_0.jpg'
})

</script>

<style scoped>
.container {
  @apply flex w-full;
}

.left-container {
  @apply w-60;
}

.main-container {
  @apply w-xl bg-gray-600;
}

.skills-container {
  @apply flex items-center mt-1;
}

.skill-icon {
  @apply flex items-center mx-1 w-6 border-2 border-black;
  border-radius: 4px;
}

.tag-icon{
  @apply w-6 mx-0.5;
}

.tooltip-content {
  @apply w-80;
}

.text-bg {
  @apply text-start text-gray-700 mx-1;
}

.hero-image {
  @apply w-full;
}

.my-divider {
  @apply my-2 mx-0;
}
</style>
