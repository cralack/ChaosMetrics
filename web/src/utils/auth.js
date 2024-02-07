import { useCookies } from '@vueuse/integrations/useCookies'

const TokenKey = 'x-token'
const cookie = useCookies()

export function getToken() {
  return cookie.get(TokenKey)
}

export function setToken(token) {
  return cookie.set(TokenKey, token)
}

export function removeToken() {
  return cookie.remove(TokenKey)
}
