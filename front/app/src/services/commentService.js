import { request } from '@/services/http'

export async function getWorkComments(workId, { page = 1, pageSize = 10 } = {}) {
  const query = new URLSearchParams({ page: String(page), pageSize: String(pageSize) })
  const data = await request(`/app/v1/works/${encodeURIComponent(workId)}/comments?${query}`)
  return mapCommentPage(data)
}

export async function getCommentReplies(commentId, { page = 1, pageSize = 10 } = {}) {
  const query = new URLSearchParams({ page: String(page), pageSize: String(pageSize) })
  const data = await request(`/app/v1/comments/${encodeURIComponent(commentId)}/replies?${query}`)
  return mapCommentPage(data)
}

export async function createComment(workId, payload, token) {
  const data = await request(`/api/v1/works/${encodeURIComponent(workId)}/comments`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify(payload)
  })
  return mapComment(data.comment)
}

export async function uploadCommentImage(file, token) {
  const presign = await request('/api/v1/media/upload/presign', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify({ resourceType: 'image', fileName: file.name, contentType: file.type })
  })
  const uploadResponse = await fetch(presign.uploadUrl, {
    method: presign.method || 'PUT',
    headers: { ...(presign.headers || {}), 'Content-Type': file.type || presign.headers?.['Content-Type'] },
    body: file
  })
  if (!uploadResponse.ok) throw new Error('评论图片上传失败')
  const completed = await request('/api/v1/media/upload/complete', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify({ mediaId: presign.mediaId, objectKey: presign.objectKey })
  })
  return { mediaId: completed.mediaId, previewUrl: completed.previewUrl }
}

function mapCommentPage(data) {
  return {
    total: Number(data.total || 0), page: Number(data.page || 1),
    pageSize: Number(data.pageSize || 10), hasMore: Boolean(data.hasMore),
    list: (data.list || []).map(mapComment)
  }
}

export function mapComment(item) {
  return {
    id: item.commentId,
    workId: item.workId,
    rootId: item.rootCommentId,
    parentId: item.parentCommentId,
    replyToUserId: item.replyToUserId,
    replyToName: item.replyToName || '',
    author: {
      id: item.author?.userId || 0,
      nickname: item.author?.nickName || '咕噜用户',
      avatar: item.author?.avatarUrl || ''
    },
    content: item.content || '',
    image: item.imageUrl || '',
    meta: formatCommentTime(item.createdAt),
    likes: Number(item.like?.count || 0),
    liked: Boolean(item.like?.liked),
    replyCount: Number(item.replyCount || 0),
    replyHasMore: Boolean(item.replyHasMore),
    replies: (item.replies || []).map(mapComment)
  }
}

function formatCommentTime(value) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const elapsed = Date.now() - date.getTime()
  if (elapsed < 60_000) return '刚刚'
  if (elapsed < 3_600_000) return `${Math.floor(elapsed / 60_000)}分钟前`
  if (elapsed < 86_400_000) return `${Math.floor(elapsed / 3_600_000)}小时前`
  return date.toLocaleDateString('zh-CN')
}
