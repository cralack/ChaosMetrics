import service from '@/api/axios'

export const getGameVersion = () => {
  return service({
    url: '/gameversion',
    method: 'get'
  })
}
