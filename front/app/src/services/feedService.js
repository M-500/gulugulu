import { request } from '@/services/http'

export async function getHomeFeed({ page = 1, pageSize = 20, userId = 0 } = {}) {
  const query = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize)
  })
	if (Number(userId) > 0) query.set('userId', String(userId))
  const data = await request(`/app/v1/recommend/works?${query.toString()}`)

  return {
    page: data.page,
    pageSize: data.pageSize,
    hasMore: data.hasMore,
    list: data.list.map(mapRecommendWorkToNote)
  }
}

export function mapRecommendWorkToNote(item) {
  const authorName = item.author?.nickName || '咕噜用户'
  return {
    id: item.workId,
    authorId: item.author?.userId || 0,
    title: item.title || item.contentExcerpt || '未命名作品',
    author: authorName,
    avatar: item.author?.avatarUrl || '',
    image: item.coverUrl || '',
    likes: item.like?.text || String(item.like?.count || 0),
	liked: Boolean(item.like?.liked),
    isVideo: item.type === 'video',
    quote: item.coverUrl ? '' : item.contentExcerpt,
    tone: item.type === 'video' ? 'dark' : 'paper',
    workType: item.type
  }
}
