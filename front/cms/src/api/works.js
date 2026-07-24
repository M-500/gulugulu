import http from './http'

export function createWork(payload, cover, idempotencyKey) {
  const formData = new FormData()
  formData.append('payload', JSON.stringify(payload))
  formData.append('cover', cover, cover.name || 'cover.jpg')

  return http.post('/api/v1/works', formData, {
    headers: {
      'Idempotency-Key': idempotencyKey
    },
    timeout: 30000
  })
}

export function getWorkStatus(workId) {
  return http.get(`/api/v1/works/${workId}/status`)
}

export function getWork(workId) {
  return http.get(`/api/v1/works/${workId}`)
}
