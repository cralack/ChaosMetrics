<template>
  <div>
    <el-container class="min-h-screen bg-dark-200">
      <el-aside class="max-w-50">
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
                <span>ChaosMertics</span>
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
              default-active="aram"
              :router="true"
            >

              <el-menu-item index="aram">
                <span>ARAM</span>
              </el-menu-item>
              <el-menu-item index="herodetail">
                <span>Hero</span>
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
              v-if="!store.isLogin"
              plain
              class="style-button"
              @click="goLogin"
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
              v-if="store.isLogin"
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
                  v-if="store.isLogin"
                >
                  <div class="allCenter text-base mt-1 mb-1">{{ store.userInfo.nickName }}</div>
                  <el-dropdown-item divided>
                    <el-icon>
                      <Setting />
                    </el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <el-icon>
                      <StarFilled />
                    </el-icon>
                    收藏夹
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <el-icon>
                      <Message />
                    </el-icon>
                    信息
                  </el-dropdown-item>

                  <el-dropdown-item
                    divided
                    @click="store.LoginOut"
                  >
                    <el-icon>
                      <SwitchButton />
                    </el-icon>
                    登出
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
        >
          <div class="main-content">
            <RouterView />
          </div>
        </el-main>
      </el-container>
      <el-container />
    </el-container>
  </div>

</template>

<script setup>
import { useRouter } from 'vue-router'
import logoImg from '@/assets/logo_inv.png'
import { useUserStore } from '@/store/user'

const store = useUserStore()
const router = useRouter()

const goLogin = () => {
  router.push('/login')
}

// const goSwagger=()=>{
//   window.open('http://localhost:8080/swagger/index.html')
// }

</script>

<style>
.allCenter {
  @apply flex items-center justify-center;
}

.el-menu {
  border: 0 !important;
}

.style-button {
  &.is-plain {
    background-color: transparent;
    color: #d1d5db;
    border: transparent;

    &:hover {
      color: #ffd04b;
    }
  }
}

</style>
