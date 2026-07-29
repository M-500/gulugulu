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
  const contentType = presign.headers?.['Content-Type'] || file.type || 'application/octet-stream'
  const response = await axios({
    method: presign.method || 'PUT',
    url: presign.uploadUrl,
    headers: {
      'Content-Type': contentType
    },
    data: file,
    responseType: 'blob',
    validateStatus: (status) => status >= 200 && status < 300
  })

  return response.data
}
