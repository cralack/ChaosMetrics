<template>
  <div>
    <!--    Home-->
    <!--    <pre>{{ getToken }}</pre>-->
    <!--    <el-button-->
    <!--      @click="goLogin"-->
    <!--    >Login</el-button>-->
    <!--    <el-button-->
    <!--      @click="removeToken"-->
    <!--    >Clear Cookie</el-button>-->

    <el-container class="min-h-screen bg-dark-200">
      <el-aside class="w-1/5">
        <el-container class="min-h-screen ">
          <!--          logo-->
          <el-header
            class="allCenter
           h-50"
          >
            <div class="logo">
              <el-image
                :src="logoImg"
                :fit="'scale-down'"
              />
              <div class="allCenter text-gray-200">
                <span>Chaos</span>
                <span>Mertics</span>
              </div>
            </div>
          </el-header>

          <!--          sidebar menu-->
          <el-main
            class="allCenter
           h-full"
          >
            <el-menu
              active-text-color="#ffd04b"
              background-color="transparent"
              text-color="#d1d5db"
              default-active="1"
              :router="true"
              @open="handleOpen"
              @close="handleClose"
            >

              <el-menu-item index="1">
                <span>Classic</span>
              </el-menu-item>
              <el-menu-item index="2">
                <span>ARAM</span>
              </el-menu-item>
              <el-menu-item index="3">
                <span>Champions</span>
              </el-menu-item>
              <el-menu-item index="4">
                <span>Items</span>
              </el-menu-item>
            </el-menu>
          </el-main>
        </el-container>
      </el-aside>

      <el-container
        class="min-h-screen
      bg-dark-100  w-full"
      >
        <!--        banner-->
        <el-header
          class="flex items-center justify-end"
        >
          <div>
            <el-button
              v-if="!isLogin"
              plain
              class="style-button"
              @click="login"
            >登陆
              <el-icon
                class="el-icon--right"
                color="#d1d5db"
                size="large"
              >
                <user />
              </el-icon>
            </el-button>

            <el-dropdown
              v-if="isLogin"
              trigger="click"
            >
              <div class="flex items-center mr-3">
                <el-avatar
                  src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
                  class="mr-3"
                />
              </div>

              <template #dropdown>
                <el-dropdown-menu
                  v-if="isLogin"
                >
                  <div class="allCenter text-base">{{ username }}</div>
                  <el-dropdown-item divided>
                    <el-icon><Setting /></el-icon> 个人中心
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <el-icon><StarFilled /></el-icon>收藏夹
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <el-icon><Message /></el-icon>信息
                  </el-dropdown-item>

                  <el-dropdown-item
                    divided
                    @click="logout"
                  >
                    <el-icon>  <SwitchButton /></el-icon> 登出
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!--        main-->
        <el-main
          class="allCenter
         bg-blue-gray-500"
        > {{ isLogin }}
        </el-main></el-container>
      <el-container />
    </el-container>
  </div>

</template>

<script setup>
import { useRouter } from 'vue-router'
import { removeToken } from '@/utils/auth'
const router = useRouter()
const username = ref('snoop')
const isLogin = ref(false)

const goLogin = () => {
  router.push('/login')
}

const login = () => {
  isLogin.value = true
  // goLogin()
}

const logout = () => {
  removeToken()
  isLogin.value = false
}

import logoImg from '@/assets/logo.png'
import { ref } from 'vue'

const handleOpen = (key, keyPath) => {
  console.log(key, keyPath)
}
const handleClose = (key, keyPath) => {
  console.log(key, keyPath)
}

</script>

<style>
.allCenter{
  @apply flex items-center justify-center;
}
.el-menu{
  border: 0!important;
}

.style-button{
  &.is-plain {
    background-color:transparent;
    color:#d1d5db;
    border: transparent;
    &:hover {
      color: #ffd04b;
    }
  }
}

</style>
