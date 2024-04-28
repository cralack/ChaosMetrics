import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'
import { getToken, removeToken, setToken } from '@/utils/auth'
import { getUserInfo, login, logout } from '@/api/user'
import { useRouter } from 'vue-router'

export const useUserStore = defineStore('user', () => {
  const router = useRouter()
  const lang = ref('zh_CN')
  const isLogin = computed(() => token.value !== '')
  const userInfo = reactive({
    uuid: '',
    nickName: '',
    email: '',
    role: 0,
  })

  const check = computed(() => userInfo.value.role * 10)
  const token = ref(getToken() || '')
  const settkn = (val) => {
    token.value = val
  }

  const LoginIn = async(loginInfo) => {
    const res = await login(loginInfo)
    if (res.code === 0) {
      settkn(res.data.token)
      setToken(res.data.token)
      await router.push('/')
    }
  }

  const LoginOut = async() => {
    const res = await logout()
    if (res.code === 0) {
      await ClearStorage()
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
    userInfo.uuid = val.uuid
    userInfo.nickName = val.NickName
    userInfo.email = val.Email
    userInfo.role = val.Role
  }

  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.code === 1) {
      setUserInfo(res.data)
    } else {
      await ClearStorage
    }
    return res
  }

  return {
    check,
    lang,
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
