import service from '@/api/axios'

export const login = (data) => {
  return service({
    url: '/user/login',
    method: 'post',
    data: data
  })
}
export const logout = () => {
  return service({
    url: '/user/logout',
    method: 'get',
  })
}

export const register = (data) => {
  return service({
    url: '/user/register',
    method: 'post',
    data: data
  })
}

export const verifyRegister = () => {
  return service({
    url: '/user/verify',
    method: 'get',
    params: {
      token
    }
  })
}

export const changePassword = (data) => {
  return service({
    url: '/user/changepasswd',
    method: 'post',
    data: data
  })
}

export const getUserInfo = () => {
  return service({
    url: '/user/info',
    method: 'get'
  })
}
