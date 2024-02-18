import { defineStore } from 'pinia'
import { getToken } from '@/utils/auth'
import { getUserInfo } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref({
    uuid: '',
    nickName: '',
    role: {}
  })

  const token = ref(getToken())

  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.data.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }

  return {
    userInfo,
    token,
    GetUserInfo,
  }
})
