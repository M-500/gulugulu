import { request } from '@/services/http'

export function getCaptcha() {
  return request('/na/v1/captcha')
}

export function loginWithEmail(data) {
  return request('/na/v1/user/login', {
    method: 'POST',
    body: JSON.stringify({
      email: data.email,
      password: data.password,
      captchaId: data.captchaId,
      captchaCode: data.captchaCode
    })
  })
}

export function getUserInfo(token) {
  return request('/api/v1/user/info', {
    headers: {
      Authorization: `Bearer ${token}`
    }
  })
}
