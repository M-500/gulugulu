const image = (id, width = 720, height = 960) => `https://picsum.photos/id/${id}/${width}/${height}`

export const homeFeed = [
  {
    id: 1,
    title: '夏日海边慢慢散步的下午',
    author: '姜时安',
    avatar: image(64, 80, 80),
    image: image(1011, 720, 540),
    likes: '982'
  },
  {
    id: 2,
    title: '城市周末咖啡路线，顺路看展',
    author: '野原新之嗷',
    avatar: image(1062, 80, 80),
    image: image(1060, 720, 720),
    likes: '2'
  },
  {
    id: 3,
    title: '最近读到的一段话',
    author: '云雀',
    avatar: image(1027, 80, 80),
    quote: '把日子过成自己喜欢的节奏，慢一点也没有关系。',
    likes: '1149',
    tone: 'dark'
  },
  {
    id: 4,
    title: '晚风里的街角灯光',
    author: 'Coco',
    avatar: image(65, 80, 80),
    image: image(1043, 720, 980),
    likes: '7',
    isVideo: true
  },
  {
    id: 5,
    title: '今日训练记录',
    author: 'Yu',
    avatar: image(1025, 80, 80),
    image: image(1005, 720, 1180),
    likes: '3538'
  },
  {
    id: 6,
    title: '给房间换一点新鲜颜色',
    author: '山楂丸子',
    avatar: image(823, 80, 80),
    image: image(101, 720, 900),
    likes: '2819'
  },
  {
    id: 7,
    title: '一页便签写下今天的关键词',
    author: '帅磊',
    avatar: image(1012, 80, 80),
    quote: '灵感\n收藏\n复盘',
    likes: '666',
    tone: 'note'
  },
  {
    id: 8,
    title: '雨后公园的清新空气',
    author: '金角大王',
    avatar: image(338, 80, 80),
    image: image(1020, 720, 1020),
    likes: '32'
  },
  {
    id: 9,
    title: '一顿认真做的早餐',
    author: '爱跳舞的猫',
    avatar: image(1024, 80, 80),
    image: image(292, 720, 840),
    likes: '1432'
  },
  {
    id: 10,
    title: '通勤路上的一点点观察',
    author: '泡桐树半',
    avatar: image(996, 80, 80),
    image: image(433, 720, 1080),
    likes: '1728'
  },
  {
    id: 11,
    title: '旅行前的行李清单',
    author: '呼兰',
    avatar: image(883, 80, 80),
    image: image(1040, 720, 900),
    likes: '1万',
    isVideo: true
  },
  {
    id: 12,
    title: '森林步道，适合周末出发',
    author: '开飞机的贝塔',
    avatar: image(91, 80, 80),
    image: image(1039, 720, 980),
    likes: '88'
  },
  {
    id: 13,
    title: '小城夜晚的热闹与安静',
    author: '阿梨',
    avatar: image(177, 80, 80),
    image: image(1047, 720, 600),
    likes: '420'
  },
  {
    id: 14,
    title: '长文：如何记录一个月的生活变化',
    author: '银杏叶',
    avatar: image(447, 80, 80),
    quote: '从饮食、睡眠、运动和情绪四个维度开始，记录越轻，越容易坚持。',
    likes: '731',
    tone: 'paper'
  }
]
