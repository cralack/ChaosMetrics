import axios from 'axios'
import { getToken } from '@/utils/auth'
import { successMessage, errorMessage } from '@/utils/message'

const service = axios.create({
  baseURL: 'http://localhost:8080'
})

service.interceptors.request.use(function(config) {
  const token = getToken()
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
  errorMessage(error.response.data.msg)
  return Promise.reject(error)
})
export default service
