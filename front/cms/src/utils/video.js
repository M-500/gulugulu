import mpegts from 'mpegts.js'

const MPEG_TS_CONTENT_TYPES = new Set([
  'video/mp2t',
  'video/mpegts',
  'video/ts'
])

export function isMpegTsSource(source) {
  const name = String(source?.name || source?.fileName || '').toLowerCase()
  const contentType = String(source?.type || source?.contentType || '').toLowerCase().split(';', 1)[0]

  return name.endsWith('.ts') || MPEG_TS_CONTENT_TYPES.has(contentType)
}

export function isVideoFile(file) {
  return String(file?.type || '').toLowerCase().startsWith('video/') || isMpegTsSource(file)
}

export function getVideoUploadContentType(file) {
  if (isMpegTsSource(file)) {
    return 'video/mp2t'
  }
  return file?.type || 'application/octet-stream'
}

export function mountVideoSource(video, source, options = {}) {
  if (!video || !source?.url) {
    throw new Error('视频预览地址不可用')
  }

  resetVideoElement(video)

  if (!isMpegTsSource(source)) {
    video.src = source.url
    video.load()
    return {
      isMpegTs: false,
      destroy: () => resetVideoElement(video)
    }
  }

  if (!mpegts.isSupported()) {
    throw new Error('当前浏览器不支持 MPEG-TS 视频预览')
  }

  const player = mpegts.createPlayer({
    type: 'mpegts',
    isLive: false,
    url: source.url
  }, {
    enableWorker: true,
    lazyLoad: false,
    enableStashBuffer: true
  })

  if (options.onError) {
    player.on(mpegts.Events.ERROR, (errorType, errorDetail) => {
      options.onError(new Error(`MPEG-TS 视频加载失败：${errorDetail || errorType}`))
    })
  }

  player.attachMediaElement(video)
  player.load()

  return {
    isMpegTs: true,
    player,
    destroy() {
      try {
        player.pause()
        player.unload()
        player.detachMediaElement()
        player.destroy()
      } finally {
        resetVideoElement(video)
      }
    }
  }
}

export async function extractVideoFrame(blob, seconds = 0, sourceMeta = {}) {
  const video = document.createElement('video')
  const objectUrl = URL.createObjectURL(blob)
  let playback = null

  video.preload = 'auto'
  video.muted = true
  video.playsInline = true

  try {
    await waitUntilFrameReady(video, seconds, (onError) => {
      playback = mountVideoSource(video, {
        url: objectUrl,
        name: sourceMeta.name || '',
        type: sourceMeta.type || blob.type
      }, { onError })
    })

    if (!video.videoWidth || !video.videoHeight) {
      throw new Error('视频画面尺寸无效，无法截取封面')
    }

    const canvas = document.createElement('canvas')
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height)

    return await new Promise((resolve, reject) => {
      canvas.toBlob((result) => {
        result ? resolve(result) : reject(new Error('视频封面生成失败'))
      }, 'image/jpeg', 0.92)
    })
  } finally {
    playback?.destroy()
    URL.revokeObjectURL(objectUrl)
  }
}

function waitUntilFrameReady(video, seconds, mountSource) {
  return new Promise((resolve, reject) => {
    let settled = false
    let metadataHandled = false
    const timeout = window.setTimeout(() => finish(new Error('视频加载超时，封面截取失败')), 30000)

    const cleanup = () => {
      window.clearTimeout(timeout)
      video.removeEventListener('loadedmetadata', handleMetadata)
      video.removeEventListener('loadeddata', handleLoadedData)
      video.removeEventListener('seeked', handleSeeked)
      video.removeEventListener('error', handleMediaError)
    }
    const finish = (error) => {
      if (settled) return
      settled = true
      cleanup()
      error ? reject(error) : resolve()
    }
    const finishOnNextFrame = () => {
      window.requestAnimationFrame(() => window.requestAnimationFrame(() => finish()))
    }
    const handleLoadedData = () => {
      if (video.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA && Math.abs(video.currentTime) < 0.05) {
        finishOnNextFrame()
      }
    }
    const handleSeeked = () => finishOnNextFrame()
    const handleMediaError = () => finish(new Error('无法加载视频，封面截取失败'))
    const handleMetadata = () => {
      if (metadataHandled) return
      metadataHandled = true
      const duration = Number.isFinite(video.duration) ? video.duration : 0
      const requestedTime = Number.isFinite(Number(seconds)) ? Number(seconds) : 0
      const targetTime = duration > 0
        ? Math.min(Math.max(requestedTime, 0), Math.max(duration - 0.1, 0))
        : Math.max(requestedTime, 0)

      if (targetTime < 0.05) {
        if (video.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA) {
          finishOnNextFrame()
        }
        return
      }
      video.currentTime = targetTime
    }

    video.addEventListener('loadedmetadata', handleMetadata)
    video.addEventListener('loadeddata', handleLoadedData)
    video.addEventListener('seeked', handleSeeked)
    video.addEventListener('error', handleMediaError)

    try {
      mountSource((error) => finish(error))
    } catch (error) {
      finish(error)
    }
  })
}

function resetVideoElement(video) {
  video.pause()
  video.removeAttribute('src')
  video.load()
}
