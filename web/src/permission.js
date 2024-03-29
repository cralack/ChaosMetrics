import router from '@/router'
import Nprogress from 'nprogress'
import { useUserStore } from '@/store/user'

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  Nprogress.start()
  if (userStore.isLogin && to.path === '/login') {
    return next({ path: from.path ? from.path : '/' })
  }

  next()
})

router.afterEach(() => {
  const userStore = useUserStore()
  if (userStore.isLogin) {
    userStore.GetUserInfo()
  }
  Nprogress.done()
})

router.onError(() => {
  Nprogress.remove()
})
