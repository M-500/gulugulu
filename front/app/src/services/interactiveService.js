import { request } from '@/services/http'

export function setResourceLike({ resourceType, resourceId, liked, token }) {
  return request(`/api/v1/interactions/${encodeURIComponent(resourceType)}/${encodeURIComponent(resourceId)}/like`, {
    method: liked ? 'POST' : 'DELETE',
    headers: {
      Authorization: `Bearer ${token}`
    }
  })
}
