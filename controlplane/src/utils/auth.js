import Cookies from 'js-cookie'

const TokenKey = 'vue_admin_template_token'

export function getToken() {
  return localStorage.token
}

export function setToken(token) {
  localStorage.token = token
}

export function removeToken() {
  localStorage.removeItem('token')
}
