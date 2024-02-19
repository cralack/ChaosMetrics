<template>
  <div>
    <el-row class="login-container">
      <el-col
        :md="16"
        class="left"
      >
        <div>
          <div>Chaos Metrics</div>
          <div>League of Legends Personal Data Analysis Project</div>
        </div>
      </el-col>

      <el-col
        :md="8"
        class="right"
      >
        <div>
          <h2 class="title">Welcome</h2>
          <el-divider class="line" />

          <el-form
            ref="loginForm"
            :model="form"
            :rules="rules"
            class="w-[250px]"
          >
            <el-form-item prop="username">
              <el-input
                v-model="form.username"
                class="w-[210px]"
                placeholder="username"
              >
                <template #prefix>
                  <el-icon class="el-input__icon">
                    <user />
                  </el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                show-password
                class="w-[210px]"
                placeholder="password"
              >
                <template #prefix>
                  <el-icon class="el-input__icon">
                    <lock />
                  </el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                class="w-[100px]"
                :loading="loading"
                @click="onSubmit"
              >
                Login
              </el-button>
              <el-button
                type="primary"
                class="w-[100px]"
              >
                Register
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>

import { ref, reactive } from 'vue'
import { useUserStore } from '@/store/user'

const login = useUserStore().LoginIn

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: 'User name cannot be empty', trigger: 'blur' },
    { min: 5, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Password cannot be empty', trigger: 'blur' },
    { min: 6, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ]
}

const loginForm = ref(null)
const loading = ref(false)

const onSubmit = () => {
  loginForm.value.validate((valid) => {
    if (!valid) {
      return false
    }
    loading.value = true
    console.log(loading)
    login(form)
    loading.value = false
  })
}
</script>

<style scoped>
.login-container {
  @apply min-h-screen bg-dark-200;
}

.login-container .left,
.login-container .right {
  @apply flex items-center justify-center;
}

.left > div > div:first-child {
  @apply font-bold text-5xl text-light-50 mb-4;
}

.left > div > div:nth-child(2) {
  @apply text-light-50 mb-4;
}

.login-container .right {
  @apply bg-light-50;
}

.right .title {
  @apply text-gray-500 text-3xl mb-4;
}

.right .line {
  @apply h-[1px] w-52 bg-gray-200;
}
</style>
