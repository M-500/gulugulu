import axios from 'axios'

import http from './http'


// 获取上传预签名信息
export function createUploadPresign(data) {
  return http.post('/api/v1/media/upload/presign', {
    resourceType: data.resourceType,
    fileName: data.fileName,
    contentType: data.contentType
  })
}

// 检查是否完成上传
export function completeUpload(data) {
  return http.post('/api/v1/media/upload/complete', {
    mediaId: data.mediaId,
    objectKey: data.objectKey
  })
}

export async function uploadToObjectStorage(file, presign) {
  const response = await axios({
    method: presign.method || 'PUT',
    url: normalizeObjectStorageUrl(presign.uploadUrl),
    headers: {
      'Content-Type': 'application/octet-stream'
    },
    data: new Blob([file]),
    responseType: 'blob',
    validateStatus: (status) => status >= 200 && status < 300
  })

  return response.data
}

function normalizeObjectStorageUrl(uploadUrl) {
  const url = new URL(uploadUrl)
  const pageHost = window.location.hostname
  const shouldUsePageHost = !['localhost', '127.0.0.1'].includes(pageHost)
    && ['localhost', '127.0.0.1'].includes(url.hostname)

  if (shouldUsePageHost) {
    url.hostname = pageHost
  }

  return url.toString()
}
