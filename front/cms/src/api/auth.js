import http from './http'

export function getCaptcha() {
  return http.get('/na/v1/captcha')
}

export function login(data) {
  return http.post('/na/v1/user/login', {
    email: data.email,
    password: data.password,
    captchaId: data.captchaId,
    captchaCode: data.captchaCode
  })
}

export function register(data) {
  return http.post('/na/v1/user/register', {
    email: data.email,
    nickName: data.nickName,
    password: data.password,
    verificationCode: data.verificationCode
  })
}
