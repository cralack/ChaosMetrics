<template>
  <div class="flex h-full w-full items-center justify-center">
    <el-form
      ref="registerForm"
      :model="form"
      :rules="rules"
      class="w-full max-w-md"
    >
      <el-form-item
        label="User Name"
        prop="username"
      >
        <el-input
          v-model="form.username"
          placeholder="johnysmith"
        />
      </el-form-item>

      <el-form-item
        label="Nick Name"
        prop="nickname"
      >
        <el-input
          v-model="form.nickname"
          placeholder="John"
        />
      </el-form-item>

      <el-form-item
        label="Email"
        prop="email"
      >
        <el-input
          v-model="form.email"
          placeholder="xxx@abc.com"
        />
      </el-form-item>

      <el-form-item
        label="Password"
        prop="password"
      >
        <el-input
          v-model="form.password"
          show-password
          type="password"
          placeholder="******"
        />
      </el-form-item>

      <el-form-item
        label="Confirm"
        prop="checkpass"
      >
        <el-input
          v-model="form.checkpass"
          show-password
          type="password"
          placeholder="******"
        />
      </el-form-item>

      <div class="flex gap-4">
        <el-form-item class="w-120">
          <el-button
            type="info"
            class="w-full"
            @click="goHome"
          >
            Home
          </el-button>
        </el-form-item>
        <el-form-item class="w-120">
          <el-button
            type="primary"
            class="w-full"
            :loading="loading"
            @click="onSubmit"
          >Register
          </el-button>
        </el-form-item>
      </div>
    </el-form>

  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { register } from '@/api/user'
import { useRouter } from 'vue-router'

const registerForm = ref(null)
const loading = ref(false)
const router = useRouter()

const form = reactive({
  username: '',
  nickname: '',
  email: '',
  password: '',
  checkpass: '',
})

const onSubmit = () => {
  registerForm.value.validate((valid) => {
    if (!valid) {
      return false
    }
    loading.value = true
    register(form)
      .then(res => {
        if (res.code === 0) {
          router.push('/login')
        }
      })
      .finally(() => {
        loading.value = false
      })
  })
}
const goHome = () => {
  router.push('/')
}

const validatePass2 = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error("Two inputs don't match!"))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: 'User name cannot be empty', trigger: 'blur' },
    { min: 5, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: 'Nick name cannot be empty', trigger: 'blur' },
    { min: 5, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Password cannot be empty', trigger: 'blur' },
    { min: 6, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ],
  checkpass: [
    { required: true, validator: validatePass2, trigger: 'blur' },
    { min: 6, max: 12, message: 'Length should be 6 to 12', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Email cannot be empty', trigger: 'blur' },
    { type: 'email', message: 'Please input correct email address', trigger: 'blur' }
  ]
}

</script>

<style scoped>

</style>
