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
    videoPlaylist: normalizePlaylistUrls(result.videoPlaylist),
    assets: (result.assets || []).map((asset) => ({
      ...asset,
      url: normalizeObjectStorageUrl(asset.url)
    }))
  })
}

export function auditWork(workId, decision, reason = '') {
  return http.post(`/api/v1/works/${workId}/audit`, { decision, reason })
}

function normalizeReviewUrls(item) {
  return {
    ...item,
    coverUrl: normalizeObjectStorageUrl(item.coverUrl)
  }
}

function normalizePlaylistUrls(value) {
  if (!value) {
    return ''
  }
  return value.split('\n').map((line) => {
    const trimmed = line.trim()
    return trimmed && !trimmed.startsWith('#') ? normalizeObjectStorageUrl(trimmed) : line
  }).join('\n')
}

function normalizeObjectStorageUrl(value) {
  if (!value) {
    return ''
  }
  const url = new URL(value)
  const pageHost = window.location.hostname
  if (!['localhost', '127.0.0.1'].includes(pageHost)
    && ['localhost', '127.0.0.1'].includes(url.hostname)) {
    url.hostname = pageHost
  }
  return url.toString()
}
