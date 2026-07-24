import { defineStore } from 'pinia'

const DB_NAME = 'gulugulu_publish_cache'
const DB_VERSION = 1
const DRAFT_STORE = 'drafts'
const ASSET_STORE = 'assets'
const MAX_DRAFTS = 100

let dbPromise

function openDb() {
  if (dbPromise) {
    return dbPromise
  }

  dbPromise = new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION)

    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve(request.result)
    request.onupgradeneeded = () => {
      const db = request.result

      if (!db.objectStoreNames.contains(DRAFT_STORE)) {
        db.createObjectStore(DRAFT_STORE, { keyPath: 'id' })
      }

      if (!db.objectStoreNames.contains(ASSET_STORE)) {
        const assetStore = db.createObjectStore(ASSET_STORE, { keyPath: 'assetId' })
        assetStore.createIndex('draftId', 'draftId', { unique: false })
      }
    }
  })

  return dbPromise
}

function runStore(storeName, mode, callback) {
  return openDb().then((db) => new Promise((resolve, reject) => {
    const transaction = db.transaction(storeName, mode)
    const store = transaction.objectStore(storeName)
    const result = callback(store)

    transaction.oncomplete = () => resolve(result)
    transaction.onerror = () => reject(transaction.error)
    transaction.onabort = () => reject(transaction.error)
  }))
}

function requestToPromise(request) {
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result)
    request.onerror = () => reject(request.error)
  })
}

async function getAll(storeName) {
  const db = await openDb()

  return requestToPromise(db.transaction(storeName, 'readonly').objectStore(storeName).getAll())
}

async function getAssetsByDraftId(draftId) {
  const db = await openDb()
  const transaction = db.transaction(ASSET_STORE, 'readonly')
  const index = transaction.objectStore(ASSET_STORE).index('draftId')

  return requestToPromise(index.getAll(draftId))
}

function createId(prefix) {
  return `${prefix}_${Date.now()}_${Math.random().toString(16).slice(2)}`
}

function formatDate(timestamp) {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(new Date(timestamp)).replace(/\//g, '-')
}

function fileKind(file) {
  if (file.type.startsWith('video/')) {
    return 'video'
  }

  if (file.type.startsWith('image/')) {
    return 'image'
  }

  return 'file'
}

function normalizeAsset(asset) {
  const localUrl = asset.blob ? URL.createObjectURL(asset.blob) : ''

  return {
    ...asset,
    url: localUrl || asset.previewUrl || '',
    localUrl,
    previewUrl: asset.previewUrl || ''
  }
}

function normalizeDraft(meta, assets = []) {
  return {
    ...meta,
    savedText: formatDate(meta.updatedAt),
    assets: assets.map(normalizeAsset)
  }
}

function toAssetMeta(asset) {
  return {
    assetId: asset.assetId,
    draftId: asset.draftId,
    name: asset.name,
    type: asset.type,
    size: asset.size,
    kind: asset.kind,
    order: asset.order,
    mediaId: asset.mediaId,
    bucket: asset.bucket,
    objectKey: asset.objectKey,
    previewUrl: asset.previewUrl,
    uploadStatus: asset.uploadStatus
  }
}

function toDraftForStore(draft) {
  return {
    id: draft.id,
    type: draft.type,
    title: draft.title,
    body: draft.body,
    collection: draft.collection,
    visibility: draft.visibility,
    allowComment: draft.allowComment,
    original: draft.original,
    scheduled: draft.scheduled,
    // Pinia deeply wraps arrays and their entries with Vue proxies. IndexedDB
    // uses the structured clone algorithm and cannot persist those proxies, so
    // keep the storage boundary strictly plain-data-only.
    topics: Array.from(draft.topics || [], (topic) => String(topic)),
    createdAt: draft.createdAt,
    updatedAt: draft.updatedAt,
    coverAssetId: draft.coverAssetId,
    assetMetas: Array.from(draft.assetMetas || [], toAssetMeta)
  }
}

export const usePublishDraftStore = defineStore('publishDrafts', {
  state: () => ({
    drafts: [],
    currentDraft: null,
    loaded: false
  }),
  getters: {
    draftCountByType: (state) => (type) => state.drafts.filter((draft) => draft.type === type).length,
    sortedDrafts: (state) => [...state.drafts].sort((a, b) => b.updatedAt - a.updatedAt),
    videoDrafts: (state) => state.drafts.filter((draft) => draft.type === 'video' || draft.type === 'mixed')
  },
  actions: {
    async loadDrafts() {
      const drafts = await getAll(DRAFT_STORE)
      this.drafts = drafts
        .map((draft) => ({
          ...draft,
          savedText: formatDate(draft.updatedAt)
        }))
        .sort((a, b) => b.updatedAt - a.updatedAt)
      this.loaded = true
    },
    async createDraft(type, files, uploadedAssets = []) {
      const now = Date.now()
      const id = createId('draft')
      const assets = Array.from(files).map((file, index) => ({
        assetId: createId('asset'),
        draftId: id,
        name: file.name,
        type: file.type,
        size: file.size,
        kind: fileKind(file),
        order: index,
        blob: file,
        mediaId: uploadedAssets[index]?.mediaId || 0,
        bucket: uploadedAssets[index]?.bucket || '',
        objectKey: uploadedAssets[index]?.objectKey || '',
        previewUrl: uploadedAssets[index]?.previewUrl || '',
        uploadStatus: uploadedAssets[index]?.status || 'local'
      }))

      const draft = {
        id,
        type,
        title: '',
        body: '',
        collection: '',
        visibility: 'public',
        allowComment: true,
        original: false,
        scheduled: false,
        topics: [],
        createdAt: now,
        updatedAt: now,
        coverAssetId: assets[0]?.assetId || '',
        assetMetas: assets.map(toAssetMeta)
      }

      await runStore(DRAFT_STORE, 'readwrite', (store) => store.put(draft))
      await runStore(ASSET_STORE, 'readwrite', (store) => {
        assets.forEach((asset) => store.put(asset))
      })

      await this.loadDrafts()
      await this.openDraft(id)
      await this.trimDrafts()

      return id
    },
    async openDraft(id) {
      const db = await openDb()
      const draft = await requestToPromise(db.transaction(DRAFT_STORE, 'readonly').objectStore(DRAFT_STORE).get(id))

      if (!draft) {
        return null
      }

      const assets = await getAssetsByDraftId(id)
      this.releaseCurrentUrls()
      this.currentDraft = normalizeDraft(draft, assets.sort((a, b) => a.order - b.order))

      return this.currentDraft
    },
    async saveCurrentDraft(patch = {}) {
      if (!this.currentDraft) {
        return
      }

      const now = Date.now()
      const updated = {
        ...this.currentDraft,
        ...patch,
        updatedAt: now
      }
      const assets = updated.assets

      await runStore(DRAFT_STORE, 'readwrite', (store) => store.put(toDraftForStore(updated)))
      this.currentDraft = {
        ...updated,
        savedText: formatDate(now),
        assets
      }
      await this.loadDrafts()
    },
    async appendImages(files, uploadedAssets = []) {
      if (!this.currentDraft || this.currentDraft.type !== 'imageText') {
        return
      }

      const remaining = Math.max(0, 19 - this.currentDraft.assets.length)
      const selectedFiles = Array.from(files).filter((file) => file.type.startsWith('image/')).slice(0, remaining)
      const startOrder = this.currentDraft.assets.length
      const assets = selectedFiles.map((file, index) => ({
        assetId: createId('asset'),
        draftId: this.currentDraft.id,
        name: file.name,
        type: file.type,
        size: file.size,
        kind: 'image',
        order: startOrder + index,
        blob: file,
        mediaId: uploadedAssets[index]?.mediaId || 0,
        bucket: uploadedAssets[index]?.bucket || '',
        objectKey: uploadedAssets[index]?.objectKey || '',
        previewUrl: uploadedAssets[index]?.previewUrl || '',
        uploadStatus: uploadedAssets[index]?.status || 'local'
      }))

      if (!assets.length) {
        return
      }

      await runStore(ASSET_STORE, 'readwrite', (store) => {
        assets.forEach((asset) => store.put(asset))
      })

      const normalizedAssets = assets.map(normalizeAsset)
      const nextAssets = [...this.currentDraft.assets, ...normalizedAssets]
      await this.saveCurrentDraft({
        assets: nextAssets,
        assetMetas: nextAssets.map(toAssetMeta)
      })
    },
    async deleteCurrentImage(assetId) {
      if (!this.currentDraft || this.currentDraft.type !== 'imageText') {
        return
      }

      const target = this.currentDraft.assets.find((asset) => asset.assetId === assetId)
      const nextAssets = this.currentDraft.assets
        .filter((asset) => asset.assetId !== assetId)
        .map((asset, index) => ({ ...asset, order: index }))

      if (!target) {
        return
      }

      if (target.localUrl) {
        URL.revokeObjectURL(target.localUrl)
      }

      await runStore(ASSET_STORE, 'readwrite', (store) => store.delete(assetId))
      await runStore(ASSET_STORE, 'readwrite', (store) => {
        nextAssets.forEach((asset) => store.put({ ...toAssetMeta(asset), blob: asset.blob }))
      })
      await this.saveCurrentDraft({
        assets: nextAssets,
        assetMetas: nextAssets.map(toAssetMeta),
        coverAssetId: this.currentDraft.coverAssetId === assetId ? (nextAssets[0]?.assetId || '') : this.currentDraft.coverAssetId
      })
    },
    async deleteDraft(id) {
      const assets = await getAssetsByDraftId(id)

      await runStore(ASSET_STORE, 'readwrite', (store) => {
        assets.forEach((asset) => store.delete(asset.assetId))
      })
      await runStore(DRAFT_STORE, 'readwrite', (store) => store.delete(id))

      if (this.currentDraft?.id === id) {
        this.releaseCurrentUrls()
        this.currentDraft = null
      }

      await this.loadDrafts()
    },
    async trimDrafts() {
      const drafts = this.sortedDrafts

      if (drafts.length <= MAX_DRAFTS) {
        return
      }

      const expiredDrafts = drafts.slice(MAX_DRAFTS)
      await Promise.all(expiredDrafts.map((draft) => this.deleteDraft(draft.id)))
    },
    releaseCurrentUrls() {
      this.currentDraft?.assets?.forEach((asset) => {
        if (asset.localUrl) {
          URL.revokeObjectURL(asset.localUrl)
        }
      })
    }
  }
})
