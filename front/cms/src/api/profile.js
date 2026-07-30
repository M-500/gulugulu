import http from './http'

export async function getCurrentUser() {
  return normalizeProfile(await http.get('/api/v1/users/me'))
}

export async function updateCurrentUser(data) {
  return normalizeProfile(await http.put('/api/v1/users/me', data))
}

export async function uploadCurrentUserAvatar(file) {
  const formData = new FormData()
  formData.append('avatar', file, file.name || 'avatar.jpg')
  return normalizeProfile(await http.post('/api/v1/users/me/avatar', formData, {
    timeout: 30000
  }))
}

function normalizeProfile(profile) {
  return {
    ...profile,
    avatarUrl: normalizeObjectStorageUrl(profile?.avatarUrl)
  }
}

function normalizeObjectStorageUrl(value) {
  if (!value) return ''
  const url = new URL(value)
  const proxyBuckets = ['/gulugulu-media/', '/gulugulu-public/']
  if (['localhost', '127.0.0.1'].includes(url.hostname)
    && proxyBuckets.some((bucketPath) => url.pathname.startsWith(bucketPath))) {
    return `${url.pathname}${url.search}`
  }
  return url.toString()
}
