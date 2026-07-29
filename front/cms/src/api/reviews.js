import http from './http'

export async function getAuditWorks(params) {
  const result = await http.get('/api/v1/audit/works', { params })

  return {
    ...result,
    list: (result.list || []).map(normalizeReviewUrls)
  }
}

export async function getAuditWorkDetail(workId) {
  const result = await http.get(`/api/v1/audit/works/${workId}`)
  return normalizeReviewUrls({
    ...result,
    videoPlaylist: result.videoPlaylist,
    assets: (result.assets || []).map((asset) => ({
      ...asset,
      url: asset.url
    }))
  })
}

export function auditWork(workId, decision, reason = '') {
  return http.post(`/api/v1/works/${workId}/audit`, { decision, reason })
}

function normalizeReviewUrls(item) {
  return {
    ...item,
    coverUrl: item.coverUrl
  }
}
