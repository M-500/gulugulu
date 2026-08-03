import { request } from '@/services/http'

export async function getAppWorkDetail(workId) {
  const data = await request(`/app/v1/works/${encodeURIComponent(workId)}`)
  return {
    id: data.workId,
    type: data.type,
    title: data.title || '未命名作品',
    content: data.content || '',
    cover: data.coverUrl || '',
    durationMs: data.durationMs || 0,
    videoPlaylist: data.videoPlaylist || '',
    publishedAt: data.publishedAt || '',
    authorId: data.author?.userId || 0,
    author: data.author?.nickName || '咕噜用户',
    avatar: data.author?.avatarUrl || '',
    likes: data.like?.text || String(data.like?.count || 0),
    favoriteCount: data.favoriteCount || 0,
    commentCount: data.commentCount || 0,
    shareCount: data.shareCount || 0,
    assets: data.assets || [],
    topics: data.topics || []
  }
}
