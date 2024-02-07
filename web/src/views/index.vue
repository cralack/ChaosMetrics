<template>
  <div>
    Home
    <pre>{{ getToken }}</pre>
    <el-button
      @click="goLogin"
    >Login</el-button>
    <el-button
      @click="removeToken"
    >Clear Cookie</el-button>
  </div>

</template>

<script setup>
import { useRouter } from 'vue-router'
import { removeToken, getToken } from '@/utils/auth'

const router = useRouter()

const goLogin = () => {
  router.push('/login')
}

</script>
