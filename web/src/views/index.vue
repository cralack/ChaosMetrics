<template>
  <div>
    Home
    <pre>{{ cookie.getAll() }}</pre>
    <el-button
      @click="goLogin"
    >Login</el-button>
    <el-button
      @click="remove"
    >Clear Cookie</el-button>
  </div>

</template>

<script setup>
import { useRouter } from 'vue-router'
import { useCookies } from '@vueuse/integrations/useCookies'

const router = useRouter()
const cookie = useCookies()

const goLogin = () => {
  router.push('/login')
}
const remove = () => {
  cookie.remove('x-token')
}
</script>
