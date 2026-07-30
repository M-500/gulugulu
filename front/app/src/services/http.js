const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function request(path, options = {}) {
  const headers = {
    Accept: 'application/json',
    ...(options.headers || {})
  }
  if (options.body && !headers['Content-Type']) {
    headers['Content-Type'] = 'application/json'
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers
  })

  const body = await response.json().catch(() => null)
  if (!response.ok) {
    throw new Error(body?.message || `请求失败：${response.status}`)
  }
  if (!body || body.code !== 0) {
    throw new Error(body?.message || '请求失败')
  }

  // 后端统一返回 code/message/data，前端只消费 data。
  return body.data
}
