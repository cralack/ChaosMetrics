import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useCookies } from '@vueuse/integrations/useCookies'

const service = axios.create({
  baseURL: 'http://localhost:8080'
})
const successMessage = (str) => {
  ElMessage({
    showClose: true,
    message: str,
    type: 'success',
  })
}
const errorMessage = (str) => {
  ElMessage({
    showClose: true,
    message: str,
    type: 'error',
  })
}

service.interceptors.request.use(function(config) {
  const cookie = useCookies()
  const token = cookie.get('x-token')
  if (token) {
    config.headers.set('x-token', token)
  }
  return config
}, function(error) {
  return Promise.reject(error)
})

service.interceptors.response.use(function(response) {
  switch (response.data.code) {
    case 0:
      successMessage(response.data.msg)
      break
    case 4:
      errorMessage(response.data.msg)
      break
  }
  return response.data
}, function(error) {
  ElMessage({
    showClose: true,
    message: error.response.data.msg,
    type: 'error',
  })
  return Promise.reject(error)
})
export default service
