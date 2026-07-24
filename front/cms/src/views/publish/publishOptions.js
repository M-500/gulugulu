export const publishTypes = [
  {
    key: 'video',
    label: '上传视频',
    shortLabel: '视频',
    accept: 'video/*',
    hint: '拖拽视频到此或点击上传',
    buttonText: '上传视频',
    description: '适合短视频、教程、探店与内容切片'
  },
  {
    key: 'imageText',
    label: '上传图文',
    shortLabel: '图文',
    accept: 'image/*',
    hint: '拖拽图片到此或点击上传',
    buttonText: '上传图文',
    description: '适合多图笔记、清单攻略与种草内容'
  },
  {
    key: 'mixed',
    label: '视频+图文',
    shortLabel: '混合',
    accept: 'video/*,image/*',
    hint: '拖拽视频和图片到此或点击上传',
    buttonText: '上传素材',
    description: '适合视频搭配图文补充的复合作品'
  }
]

export const defaultTopics = ['生活美学', '日常文案', '人生的意义', '每天都有值得记录的瞬间', '快乐瞬间', '我的生活碎片']
