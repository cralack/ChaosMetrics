<template>
  <div>
    <el-result
      icon="info"
      title="Verifing register result"
    >
      <template #extra>
        <el-button
          type="primary"
          @click="$router.push('/')"
        >Home</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { onMounted } from 'vue'
import { verifyRegister } from '@/api/user'

onMounted(async() => {
  const route = useRoute()
  const token = route.query.token
  if (token) {
    try {
      const res = await verifyRegister(token)
      console.log(res)
    } catch (e) {
      console.log('failed')
    }
  } else {
    console.log('lack token')
  }
})
</script>
