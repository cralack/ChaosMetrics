import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login, logout, getUserInfo } from '@/api/user'
import { ElLoading } from 'element-plus'
import { useRouter } from 'vue-router'

export const useUserStore = defineStore('user', () => {
  const router = useRouter()
  const loadingInstance = ref(null)
  const isLogin = computed(() => token.value !== '')
  const userInfo = ref({
    uuid: '',
    nickName: '',
    email: '',
    role: {},
  })

  const token = ref(getToken() || '')
  const settkn = (val) => {
    token.value = val
  }

  const LoginIn = async(loginInfo) => {
    loadingInstance.value = ElLoading.service({
      fullscreen: false,
      text: '登录中，请稍候...',
    })
    try {
      const res = await login(loginInfo)
      if (res.code === 0) {
        settkn(res.data.token)
        setToken(res.data.token)
        console.log(res.data.token)
        await router.push('/')
      }
    } catch (e) {
      loadingInstance.value.close()
    }
    loadingInstance.value.close()
  }

  const LoginOut = async() => {
    const res = await logout()
    if (res.code === 0) {
      // await ClearStorage()
      console.log(res)
    }
  }

  const ClearStorage = async() => {
    token.value = ''
    sessionStorage.clear()
    localStorage.clear()
    removeToken()
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }
  const setUserInfo = (val) => {
    userInfo.value = val
  }
  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.data.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }

  return {
    isLogin,
    userInfo,
    token,
    LoginIn,
    LoginOut,
    ClearStorage,
    ResetUserInfo,
    GetUserInfo,
  }
})
