<template>
  <div>
    <el-container
      v-if="hero"
      class="container"
    >
      <el-aside
        class="left-container"
      >
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
                  <h3>{{ spell.name + " " + "快捷键:" + getSkillShortId(spell.id) }}</h3>
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
            <div v-if="index%2===1">
              <el-icon>
                <SemiSelect />
              </el-icon>
            </div>
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
          <h3 class="text-bg">参战率: {{ (heroData?.avg_kp * 100 ?? 0).toFixed(2) }} </h3>
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
        <div class="talent-container">
          <div
            v-for="(view,index) in perkViews"
            :key="index"
          >
            <el-row>
              <el-col
                :span="11"
                class="pri"
              >
                <div
                  v-for="item in view.pri"
                  :key="item.id"
                  class="rune"
                >
                  <el-tooltip
                    effect="dark"
                    placement="bottom"
                    :disabled="item.id % 100 === 0"
                  >
                    <template #content>
                      <div
                        class="tooltip-content"
                      >
                        <h4>{{ item.name }}</h4>
                        <el-divider class="my-divider" />
                        <span>{{ item.description }}</span>
                      </div>
                    </template>
                    <el-image
                      :src="getPerkImageUrl(item)"
                      :alt="item.name"
                      class="rune-icon"
                    />
                  </el-tooltip>
                </div>
              </el-col>
            </el-row>
            <el-row
              class="flex items-center"
            >
              <el-col
                class="sub"
                :span="7"
              >
                <div
                  v-for="item in view.sub"
                  :key="item.id"
                  class="rune"
                >
                  <el-tooltip
                    effect="dark"
                    placement="bottom"
                  >
                    <template #content>
                      <div class="tooltip-content">
                        <h4>{{ item.name }}</h4>
                        <el-divider class="my-divider" />
                        <span>{{ item.description }}</span>
                      </div>
                    </template>
                    <el-image
                      :src="getPerkImageUrl(item)"
                      :alt="item.name"
                      class="rune-icon"
                    />
                  </el-tooltip>
                </div>
              </el-col>
              <el-col
                class="stat"
                :span="4"
              >
                <div
                  v-for="item in view.stat"
                  :key="item.id"
                >
                  <el-image
                    class="stat-icon"
                    :src="getStatImageUrl(item.id)"
                  />
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </el-main>
    </el-container>

  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import { useGameStore } from '@/store/game'
import { useUserStore } from '@/store/user'
import { getHeroData, getHeroDetail, getPerks } from '@/api/game'

const userStore = useUserStore()
const gameStore = useGameStore()
const heroName = ref('')
const mode = ref('')
const loc = ref('')
const version = ref()

const hero = ref()
const heroData = ref()

const perksData = ref([])
const perkWinRates = ref({})
const perkViews = ref([])

onMounted(async() => {
  const route = useRoute()
  await gameStore.setVersions()

  heroName.value = route.query.hero
  mode.value = route.query.mode
  loc.value = route.query.loc
  version.value = gameStore.gameversion[0]

  await setData()
  if (perksData.value.length > 0 && heroData.value) {
    processPerkWinRates()
  }

  perkViews.value = perkWinRates.value.map(perk => setPerkView(perk))

  console.log(perkViews.value)
  console.log(perkWinRates.value)
})

const setData = async() => {
  const res = await getHeroDetail(heroName.value, version.value, userStore.lang)
  if (res.code === 1) {
    hero.value = res.data
  }

  const res1 = await getHeroData(heroName.value, loc.value, mode.value, version.value)
  if (res1.code === 1) {
    heroData.value = res1.data
  }

  const res2 = await getPerks(version.value, userStore.lang)
  if (res2.code === 1) {
    perksData.value = res2.data
  }
}

const getSkillShortId = (skillId) => {
  return skillId.replace(hero.value.id, '')
}

// 获取技能图片 URL 的方法
const getSpellImageUrl = (spell) => {
  return `src/assets/datadragon/spell/${spell.image.full}`
}

// 获取被动技能图片 URL 的方法
const getPassiveImageUrl = (passive) => {
  if (heroName.value === 'Fiddlesticks') {
    return 'src/assets/datadragon/passive/FiddleSticksP.png'
  }
  return `src/assets/datadragon/passive/${passive.image.full}`
}

const getTagImageUrl = (tag) => {
  return `src/assets/tags/${tag}.png`
}

const getPerkImageUrl = (perk) => {
  return `src/assets/datadragon/${perk.icon}`
}

const statModIcons = {
  5001: 'StatModsHealthScalingIcon.png',
  5002: 'StatModsArmorIcon.png',
  5003: 'StatModsMagicResIcon.png',
  5005: 'StatModsAttackSpeedIcon.png',
  5007: 'StatModsCDRScalingIcon.png',
  5008: 'StatModsAdaptiveForceIcon.png',
  5010: 'StatModsMovementSpeedIcon.png',
  5011: 'StatModsHealthPlusIcon.png',
  5013: 'StatModsTenacityIcon.png'
}

const getStatImageUrl = (id) => {
  return `src/assets/datadragon/perk-images/StatMods/${statModIcons[id]}`
}

const processPerkWinRates = () => {
  const winData = heroData.value.perk

  let winDataArray = Object.entries(winData).map(([key, wins]) => {
    const [priWithLabel, subWithLabel, statWithLabel] = key.split(' ')
    const pri = priWithLabel.replace('pri:', '').split(',').map(Number)
    const sub = subWithLabel.replace('sub:', '').split(',').map(Number)
    const stat = statWithLabel.replace('stat:', '').split(',').map(Number)
    return { pri, sub, stat, wins }
  })

  winDataArray.sort((a, b) => b.wins - a.wins)
  winDataArray = winDataArray.slice(0, 3)
  perkWinRates.value = winDataArray.map(perkData => {
    return {
      pri: perkData.pri,
      sub: perkData.sub,
      stats: perkData.stat,
      wins: perkData.wins,
    }
  })
}

const getRunesDetails = (ids) => {
  const details = []
  ids.forEach(id => {
    perksData.value.forEach(style => {
      if (style.id === id) {
        // 添加顶级分类信息
        details.push({
          id: style.id,
          key: style.key,
          icon: style.icon,
          description: style.name
        })
      }
      style.slots.forEach(slot => {
        const rune = slot.runes.find(rune => rune.id === id)
        if (rune) {
          details.push({
            id: rune.id,
            key: rune.key,
            icon: rune.icon,
            name: rune.name,
            description: rune.shortDesc
          })
        }
      })
    })
  })
  return details
}

const setPerkView = (perk) => {
  return {
    win: perk.wins,
    pri: getRunesDetails(perk.pri),
    sub: getRunesDetails(perk.sub),
    stat: perk.stats.map(statId => {
      return {
        id: statId,
        description: `Stat detail for ${statId}`
      }
    })
  }
}

const heroImage = computed(() => {
  if (heroName.value === 'Fiddlesticks') {
    return 'src/assets/datadragon/champion_og/loading/' + 'FiddleSticks' + '_0.jpg'
  }
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

.talent-container {
  @apply min-h-xs;
}

.talent-container .pri{
  @apply  mx-4 mt-4;
}
.talent-container .sub{
  @apply ml-4;
}
.talent-container .stat{
  @apply flex items-center;
}
.rune {
  @apply inline-flex flex-col items-center m-2;
}

.stat-icon {
  @apply w-6 h-6;
}

.rune-icon {
  @apply w-8 mb-1;
}

.skill-icon {
  @apply flex items-center mx-1 w-6 border-2 border-black;
  border-radius: 4px;
}

.tag-icon {
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
