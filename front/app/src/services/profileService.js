import { mapRecommendWorkToNote } from '@/services/feedService'
import { request } from '@/services/http'

export async function getAppUserProfile(userId) {
  const data = await request(`/app/v1/users/${encodeURIComponent(userId)}`)
  return {
    userId: data.userId,
    nickname: data.nickName || '咕噜用户',
    avatar: data.avatarUrl || '',
    bio: data.bio || '',
    ipAddress: data.ipAddress || '',
    sex: data.sex
  }
}

export async function getUserPublishedWorks(userId, { page = 1, pageSize = 20, token = '' } = {}) {
  const query = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize)
  })
  const data = await request(`/app/v1/users/${encodeURIComponent(userId)}/works?${query.toString()}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {}
  })

  return {
    total: data.total,
    page: data.page,
    pageSize: data.pageSize,
    hasMore: data.hasMore,
    list: data.list.map(mapRecommendWorkToNote)
  }
}
