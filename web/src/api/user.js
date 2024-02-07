import service from '@/api/axios'

// @Summary 用户登录
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /base/login [post]
export const login = (data) => {
  return service({
    url: '/user/login',
    method: 'post',
    data: data
  })
}
