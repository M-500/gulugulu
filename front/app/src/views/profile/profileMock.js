const image = (id, width = 720, height = 960) => `https://picsum.photos/id/${id}/${width}/${height}`

export const profileUser = {
  userId: 'mock',
  nickname: '花花🌺_有闲版',
  avatar: image(64, 240, 240),
  redId: '63698788670',
  location: '浙江',
  bio: '不一样的穿搭风格\n闲情雅致😋\n2589330358@qq.com',
  following: 0,
  followers: 297,
  receives: 2863
}

export const profileNotes = [
  { id: 101, authorId: 'mock', title: '小挎包，黑短裤穿搭', author: profileUser.nickname, avatar: profileUser.avatar, image: image(1011, 680, 900), likes: '2237' },
  { id: 102, authorId: 'mock', title: '短裤，丝袜穿搭', author: profileUser.nickname, avatar: profileUser.avatar, image: image(1027, 680, 900), likes: '27' },
  { id: 103, authorId: 'mock', title: '这裙子谁穿谁仙！～～～', author: profileUser.nickname, avatar: profileUser.avatar, image: image(1060, 680, 900), likes: '120' },
  { id: 104, authorId: 'mock', title: '粉色露肩衫穿搭', author: profileUser.nickname, avatar: profileUser.avatar, image: image(325, 680, 900), likes: '6' },
  { id: 105, authorId: 'mock', title: '紫色穿搭', author: profileUser.nickname, avatar: profileUser.avatar, image: image(65, 680, 900), likes: '2' },
  { id: 106, authorId: 'mock', title: '镜子里的我', author: profileUser.nickname, avatar: profileUser.avatar, image: image(996, 680, 920), likes: '18' }
]
