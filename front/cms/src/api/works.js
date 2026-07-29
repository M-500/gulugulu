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

export async function getWork(workId) {
  const result = await http.get(`/api/v1/works/${workId}`)
  return {
    ...result,
    coverUrl: result.coverUrl,
    videoPlaylist: result.videoPlaylist,
    assets: (result.assets || []).map((asset) => ({
      ...asset,
      url: asset.url
    }))
  }
}

export async function getCreatorWorks(params) {
  const result = await http.get('/api/v1/creator/works', { params })

  return {
    ...result,
    list: (result.list || []).map((item) => ({
      ...item,
      coverUrl: item.coverUrl
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
