<template>
  <div>
    Home
  </div>
  <div>
    <el-button @click="addCount">
      {{ count }}
    </el-button>
  </div>
</template>

<script setup>
import { ref } from "vue"

let count = ref(1)

function addCount() {
  count.value++
  console.log(count)
}
</script>