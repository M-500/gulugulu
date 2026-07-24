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

export async function getCreatorWorks(params) {
  const result = await http.get('/api/v1/creator/works', { params })

  return {
    ...result,
    list: (result.list || []).map((item) => ({
      ...item,
      coverUrl: normalizeObjectStorageUrl(item.coverUrl)
    }))
  }
}

export function updateCreatorWorkTitle(workId, title) {
  return http.put(`/api/v1/creator/works/${workId}/title`, { title })
}

export function updateCreatorWorkVisibility(workId, visibility) {
  return http.put(`/api/v1/creator/works/${workId}/visibility`, { visibility })
}

export function deleteCreatorWork(workId) {
  return http.delete(`/api/v1/creator/works/${workId}`)
}

function normalizeObjectStorageUrl(value) {
  if (!value) {
    return ''
  }
  const url = new URL(value)
  const pageHost = window.location.hostname
  const shouldUsePageHost = !['localhost', '127.0.0.1'].includes(pageHost)
    && ['localhost', '127.0.0.1'].includes(url.hostname)

  if (shouldUsePageHost) {
    url.hostname = pageHost
  }
  return url.toString()
}
