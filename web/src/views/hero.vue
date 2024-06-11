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
            v-for="(view, index) in perkViews"
            :key="index"
          >
            <el-row>
              <el-col
                :span="12"
                class="pri"
              >
                <div
                  v-for="(priItem,subIndex) in view.pri"
                  :key="subIndex"
                  class="rune"
                >
                  <el-tooltip
                    effect="dark"
                    placement="bottom"
                    :disabled="priItem.id % 100 === 0"
                  >
                    <template #content>
                      <div
                        class="tooltip-content"
                      >
                        <el-text type="warning">{{ priItem.name }}</el-text>
                        <el-divider class="my-divider" />
                        <el-text>{{ priItem.description }}</el-text>
                      </div>
                    </template>
                    <el-image
                      :src="getPerkImageUrl(priItem)"
                      :alt="priItem.name"
                      class="rune-icon"
                      :class="{ 'larger-icon': subIndex === 1 }"
                    />
                  </el-tooltip>
                </div>
              </el-col>
              <el-col
                :span="6"
                class="pick-rate"
              >
                <el-text
                  size="large"
                  type="warning"
                >登场: {{ (view.picks / heroData.total_played * 100).toFixed(1) }}%
                </el-text>
              </el-col>
            </el-row>
            <el-row
              class="flex items-center"
            >
              <el-col
                class="sub"
                :span="8"
              >
                <div
                  v-for="subItem in view.sub"
                  :key="subItem.id"
                  class="rune"
                >
                  <el-tooltip
                    effect="dark"
                    placement="bottom"
                    :disabled="subItem.id % 100 === 0"
                  >
                    <template #content>
                      <div class="tooltip-content">
                        <el-text type="warning">{{ subItem.name }}</el-text>
                        <el-divider class="my-divider" />
                        <el-text>{{ subItem.description }}</el-text>
                      </div>
                    </template>
                    <el-image
                      :src="getPerkImageUrl(subItem)"
                      :alt="subItem.name"
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
                  v-for="statItem in view.stat"
                  :key="statItem.id"
                >
                  <el-image
                    class="stat-icon"
                    :src="getStatImageUrl(statItem.id)"
                  />
                </div>
              </el-col>
              <el-col
                :span="6"
                class="win-rate"
              >
                <el-text
                  size="large"
                  type="warning"
                >
                  胜场: {{ (view.wins / view.picks * 100).toFixed(1) }}%
                </el-text>
              </el-col>
            </el-row>
            <el-divider
              v-if="index !== perkViews.length - 1"
              class="my-divider"
            />
          </div>
        </div>
        <div class="spell-container">
          <el-row :gutter="20">
            <el-space
              :size="'default'"
              :spacer="spacer"
            >
              <el-col
                v-for="spell in matchedSpells"
                :key="spell.id"
                :span="24"
                class="spell-col"
              >
                <div
                  v-for="detail in spell.details"
                  :key="detail.id"
                  class="flex items-center"
                >
                  <el-tooltip>
                    <template #content>
                      <div class="tooltip-content">
                        <el-text>{{ detail.name }}</el-text>
                        <el-divider class="my-divider" />
                        <el-text> {{ detail.description }}</el-text>
                      </div>
                    </template>
                    <el-image
                      class="mx-1"
                      style="width: 40px"
                      fit="cover"
                      :src="getSpellImageUrl(detail)"
                    />
                  </el-tooltip>
                </div>
                <div class="mx-2">
                  <el-text type="warning">登场率:{{ (spell.picks / heroData.total_played * 100).toFixed(1) }}%</el-text>
                  <br>
                  <el-text type="warning">胜率:{{ (spell.wins / spell.picks * 100).toFixed(1) }}%</el-text>
                </div>

              </el-col>
            </el-space>
          </el-row>
        </div>
        <div class="item-container">
          <el-row
            v-for="itemList in matchedItemsTri"
            :key="itemList.id"
            class="w-62 bg-blue-gray-500"
          >
            <div
              v-for="item in itemList.details"
              :key="item.id"
            >
              <el-tooltip
                placement="top"
                effect="dark"
                :visible-arrow="false"
              >
                <template #content>
                  <div
                    class="w-120"
                    style="white-space: pre-wrap"
                  >
                    <el-text
                      size="large"
                      type="warning"
                    >{{ item.name }}
                    </el-text><br>
                    <el-text
                      size="small"
                    >价格：{{ item.total_gold }}({{ item.base_gold }})
                    </el-text>
                    <el-divider class="my-divider" />
                    <el-text> {{ item.description }}</el-text>
                    <el-divider class="my-divider" />
                  </div>
                  <ItemTree
                    :item="item"
                    :items="items"
                  />
                </template>
                <el-image
                  class="mx-1"
                  style="width: 40px"
                  fit="cover"
                  :src="getItemsImageUrl(item)"
                />
              </el-tooltip>
            </div>
            <div class="mx-2">
              <el-text type="warning">登场率:{{ (itemList.picks/ heroData.total_played * 100).toFixed(1) }}%</el-text><br>
              <el-text type="warning">胜率:{{ (itemList.wins/itemList.picks*100).toFixed(1) }}%</el-text>
            </div>
          </el-row>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { computed, onMounted, ref, h } from 'vue'
import { useGameStore } from '@/store/game'
import { useUserStore } from '@/store/user'
import { getHeroData, getHeroDetail, getPerks, getSpells, getItems } from '@/api/game'
import { ElDivider } from 'element-plus'
import ItemTree from '@/views/itemTree.vue'

const userStore = useUserStore()
const gameStore = useGameStore()
const heroName = ref('')
const mode = ref('')
const loc = ref('')
const version = ref()

const hero = ref()
const heroData = ref()

const perks = ref([])
const perkWinRates = ref({})
const perkViews = ref([])

const spells = ref([])
const matchedSpells = ref([])

const items = ref([])
// 示例使用
// const matchedItemsFir = ref([])
const matchedItemsTri = ref([])
// const matchedItemsSho = ref([])
// const matchedItemsOth = ref([])

onMounted(async() => {
  const route = useRoute()
  await gameStore.setVersions()

  heroName.value = route.query.hero
  mode.value = route.query.mode
  loc.value = route.query.loc
  version.value = gameStore.gameversion[0]

  await setData()
  if (perks.value.length > 0 && heroData.value) {
    processPerkWinRates()
  }
  getMatchedSpells()
  perkViews.value = perkWinRates.value.map(perk => setPerkView(perk))
  // matchedItemsFir.value = getMatchedItems(heroData.value.item.fir, items.value)
  matchedItemsTri.value = getMatchedItems(heroData.value.item.tri, items.value)
  // matchedItemsSho.value = getMatchedItems(heroData.value.item.Sho, items.value)
  // matchedItemsOth.value = getMatchedItems(heroData.value.item.Oth, items.value)
})

const setData = async() => {
  try {
    const [resHeroDetail, resHeroData, resPerks, resSpells, resItems] = await Promise.all([
      getHeroDetail(heroName.value, version.value, userStore.lang),
      getHeroData(heroName.value, loc.value, mode.value, version.value),
      getPerks(version.value, userStore.lang),
      getSpells(version.value, userStore.lang),
      getItems(version.value, userStore.lang, mode.value),
    ])

    if (resHeroDetail.code === 1) {
      hero.value = resHeroDetail.data
    } else {
      console.error('Failed to fetch hero details', resHeroDetail.message)
    }

    if (resHeroData.code === 1) {
      heroData.value = resHeroData.data
    } else {
      console.error('Failed to fetch hero data', resHeroData.message)
    }

    if (resPerks.code === 1) {
      perks.value = resPerks.data
    } else {
      console.error('Failed to fetch perks', resPerks.message)
    }

    if (resSpells.code === 1) {
      spells.value = resSpells.data
    } else {
      console.error('Failed to fetch perks', resSpells.message)
    }

    if (resItems.code === 1) {
      items.value = resItems.data
    } else {
      console.error('Failed to fetch perks', resItems.message)
    }
  } catch (error) {
    console.error('Error fetching data:', error)
  }
}

const getSkillShortId = (skillId) => {
  return skillId.replace(hero.value.id, '')
}

// 获取技能图片 URL 的方法
const getSpellImageUrl = (spell) => {
  return `src/assets/datadragon/spell/${spell.image.full}`
}

const getItemsImageUrl = (item) => {
  return `src/assets/datadragon/item/${item.image}`
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
  5001: 'StatModsHealthPlusIcon.png',
  5002: 'StatModsArmorIcon.png',
  5003: 'StatModsMagicResIcon.png',
  5005: 'StatModsAttackSpeedIcon.png',
  5007: 'StatModsCDRScalingIcon.png',
  5008: 'StatModsAdaptiveForceIcon.png',
  5010: 'StatModsMovementSpeedIcon.png',
  5011: 'StatModsHealthScalingIcon.png',
  5013: 'StatModsTenacityIcon.png'
}

const getStatImageUrl = (id) => {
  return `src/assets/datadragon/perk-images/StatMods/${statModIcons[id]}`
}

const processPerkWinRates = () => {
  const winData = heroData.value.perk

  let winDataArray = Object.entries(winData).map(([key, wp]) => {
    const [priWithLabel, subWithLabel, statWithLabel] = key.split(' ')
    const pri = priWithLabel.replace('pri:', '').split(',').map(Number)
    const sub = subWithLabel.replace('sub:', '').split(',').map(Number)
    const stat = statWithLabel.replace('stat:', '').split(',').map(Number)
    return { pri, sub, stat, wp }
  })

  winDataArray.sort((a, b) => b.wp.wins - a.wp.wins)
  winDataArray = winDataArray.slice(0, 3)
  perkWinRates.value = winDataArray.map(perkData => {
    return {
      pri: perkData.pri,
      sub: perkData.sub,
      stats: perkData.stat,
      wins: perkData.wp.wins,
      picks: perkData.wp.picks,
    }
  })
}

const getRunesDetails = (ids) => {
  const details = []
  ids.forEach(id => {
    perks.value.forEach(style => {
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
    picks: perk.picks,
    wins: perk.wins,
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

const getMatchedSpells = () => {
  matchedSpells.value = Object.entries(heroData.value.spell).map(([key, stats]) => {
    const spellIds = key.split(',').map(id => id.trim())
    const spellDetails = spellIds.map(spellId => spells.value.find(spell => spell.key === spellId))
    return {
      spellIds: spellIds,
      details: spellDetails.filter(Boolean),
      picks: stats.picks,
      wins: stats.wins
    }
  })
    .sort((a, b) => b.wins - a.wins)
    .slice(0, 2)
}

const getMatchedItems = (itemArry, items) => {
  return Object.entries(itemArry).map(([key, stats]) => {
    const itemIds = key.split(',').map(id => id.trim())
    const itemDetails = itemIds.map(itemId => items.find(item => item.id === itemId)).filter(Boolean)
    return {
      details: itemDetails,
      picks: stats.picks,
      wins: stats.wins
    }
  })
    .sort((a, b) => b.picks - a.picks)
    .slice(0, 3)
}

const spacer = h(ElDivider, { direction: 'vertical' })
</script>

<style scoped>
.container {
  @apply flex w-full;
}

.left-container {
  @apply w-60;
}

.main-container {
  @apply w-xl;
}

.skills-container {
  @apply flex items-center mt-1;
}

.talent-container {
  @apply min-h-xs bg-gray-600;
}

.talent-container .pri {
  @apply mx-4 flex items-end ;
}

.talent-container .sub {
  @apply ml-4;
}

.talent-container .stat {
  @apply flex items-center;
}

.talent-container .pick-rate {
  @apply text-center flex items-start h-12 mt-4 ml-22;
}

.talent-container .win-rate {
  @apply text-center flex items-center ml-26 mb-6;
}

.spell-container {
  @apply p-4 mt-2 bg-gray-600 flex jusity-center items-center ;
}

.spell-col {
  @apply flex my-1 mx-2 ;
}

.rune {
  @apply inline-flex flex-col items-center m-2;
}

.skill-icon {
  @apply flex items-center mx-1 w-6 border-2 border-black;
  border-radius: 4px;
}

.rune-icon {
  @apply w-8 mb-1;
}

.larger-icon {
  @apply w-12;
}

.stat-icon {
  @apply w-6 h-6;
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

.item-container {
  @apply p-4 mt-2 bg-gray-600 flex flex-col items-start;
}

.text-bg {
  @apply text-start text-gray-700 mx-1;
}

</style>
