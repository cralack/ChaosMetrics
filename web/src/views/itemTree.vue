<template>
  <div class="item-tree">
    <el-image
      class="item-icon"
      :src="getItemImageUrl(item)"
      :alt="item.name"
    />
    <div
      v-if="item.from"
      class="sub-items"
    >
      <ItemTree
        v-for="subItemId in leftMidRight(item.from).left"
        :key="subItemId"
        :item="getItemById(subItemId)"
        :items="items"
      />
      <ItemTree
        v-for="subItemId in leftMidRight(item.from).mid"
        :key="subItemId"
        :item="getItemById(subItemId)"
        :items="items"
      />
      <ItemTree
        v-for="subItemId in leftMidRight(item.from).right"
        :key="subItemId"
        :item="getItemById(subItemId)"
        :items="items"
      />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  items: {
    type: Array,
    required: true
  }
})

// 获取物品图片 URL 的方法
const getItemImageUrl = (item) => {
  return `src/assets/datadragon/item/${item.image}` // 根据实际路径调整
}

// 根据 ID 获取物品的方法
const getItemById = (id) => {
  return props.items.find(i => i.id === id)
}

// 将子项分成左、中、右三部分
const leftMidRight = (from) => {
  const length = from.length
  return {
    left: length >= 2 ? [from[0]] : [],
    mid: length === 3 ? [from[1]] : length === 1 ? [from[0]] : [],
    right: length >= 2 ? [from[length - 1]] : []
  }
}
</script>

<style scoped>
.item-tree {
  @apply relative flex flex-col items-center mt-2.5;
}

.sub-items {
  @apply flex justify-center mt-2.5 gap-10;
}

.item-icon {
  @apply w-10 h-10 m-1.25;
}
</style>
