const image = (id, width = 120, height = 120) => `https://picsum.photos/id/${id}/${width}/${height}`

// 后端通知接口接入前保留完整页面结构，后续只需将这里替换为接口返回数据。
export const notificationMock = {
  comments: [
    {
      id: 1,
      actor: 'ouou',
      avatar: image(64),
      action: '回复了你的评论',
      time: '07-09',
      content: '有共鸣的等于是同性，没共鸣的等于是异性，你的观点也挺有意思的。',
      quote: '因为我不是 gay 啊，和同性共鸣的概率比异性共鸣高太多',
      thumbnail: image(24)
    },
    {
      id: 2,
      actor: '没错我不是乔',
      avatar: image(65),
      action: '回复了你的评论',
      time: '06-11',
      content: '@核颜悦色文玩店 这里哦',
      quote: '帅啊。在哪里买的，怎么买，我也想入一对',
      thumbnail: image(106)
    },
    {
      id: 3,
      actor: '漂藏',
      avatar: image(91),
      action: '回复了你的评论',
      time: '06-11',
      content: '还可以',
      quote: '老师，我这对怎么样',
      thumbnail: image(175)
    },
    {
      id: 4,
      actor: '老A的AI研究所',
      avatar: image(180),
      action: '回复了你的评论',
      time: '2025-08-14',
      content: '该评论已删除',
      thumbnail: image(20)
    },
    {
      id: 5,
      actor: '人机飓风',
      avatar: image(237),
      action: '回复了你的评论',
      time: '2025-06-15',
      content: 'mac pro？🥺🥺',
      quote: '反正我的电脑在同时打开开发工具的时候，频繁跟我说内存不足。',
      thumbnail: image(48)
    }
  ],
  likes: [
    { id: 11, actor: '小红薯6A5B6EA7', avatar: image(1027), action: '赞了你的笔记', time: '昨天 13:40', thumbnail: image(342) },
    { id: 12, actor: '小红薯6A6903D0', avatar: image(1025), action: '赞了你的头像', time: '昨天 13:40', unsupported: true },
    { id: 13, actor: 'momo', avatar: image(433), action: '赞了你的评论', time: '07-10', quote: '因为我不是 gay 啊，和同性共鸣的概率比异性共鸣高太多', thumbnail: image(24) },
    { id: 14, actor: '天生不爱上班', avatar: image(593), action: '赞了你的评论', time: '07-07', quote: '这不就压抑上了吗', thumbnail: image(175) },
    { id: 15, actor: '花果山小余', avatar: image(823), action: '赞了你的评论', time: '07-07', quote: '这不就压抑上了吗', thumbnail: image(175) },
    { id: 16, actor: '托斯托耶夫吃鸡', avatar: image(1005), action: '收藏了你的笔记', time: '06-22', thumbnail: image(106) }
  ],
  follows: [
    { id: 21, actor: '托斯托耶夫吃鸡', avatar: image(1005), action: '开始关注你了', time: '05-28', relation: 'follow-back' },
    { id: 22, actor: '珊珊来迟', avatar: image(433), action: '开始关注你了', time: '04-10', relation: 'mutual' },
    { id: 23, actor: '铭', avatar: image(823), action: '开始关注你了', time: '2024-09-06', relation: 'none' },
    { id: 24, actor: '巴芒之道', avatar: image(237), action: 'Ta关注了你，期待你的回关', time: '2024-09-05', relation: 'follow-back' },
    { id: 25, actor: '晓陈', avatar: image(91), action: 'Ta关注了你，期待你的回关', time: '2024-09-05', relation: 'follow-back' },
    { id: 26, actor: '123', avatar: image(48), action: '开始关注你了', time: '2024-09-05', relation: 'follow-back' },
    { id: 27, actor: 'wwweerrr', avatar: image(180), action: 'Ta关注了你，期待你的回关', time: '2024-09-05', relation: 'follow-back' }
  ]
}
