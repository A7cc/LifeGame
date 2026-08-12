const sceneGroups = {
  dining: {
    label: '餐饮时光',
    description: '灯光、美食与轻松交谈',
    effect: '成功时好感额外 +2',
    image: '/images/datinginfo/dating-scenes/dining-city.webp',
    accent: '#f2a65a',
    locations: [
      '餐厅', '茶馆', '厨房', '咖啡馆', '咖啡厅', '美食节', '奶茶店',
      '商务餐厅', '烧烤', '西餐厅', '夜市',
    ],
  },
  culture: {
    label: '艺术约会',
    description: '在艺术与旋律里分享心情',
    effect: '成功时名声额外 +2',
    image: '/images/datinginfo/dating-scenes/arts-culture.webp',
    accent: '#d38cff',
    locations: [
      '博物馆', '电台', '电影节', '电影院', '歌剧院', '画廊', '画展', '剧场',
      '剧院', '录音棚', '排练厅', '片场', '琴房', '摄影棚', '首映礼',
      '舞蹈室', '舞台', '艺术馆', '艺术区', '音乐会', '音乐节', '音乐厅',
      '影院', '展览', 'KTV', 'T台', '时装周',
    ],
  },
  nature: {
    label: '自然漫游',
    description: '在风景与晚霞中慢慢靠近',
    effect: '成功时健康额外 +2',
    image: '/images/datinginfo/dating-scenes/nature-travel.webp',
    accent: '#76d39b',
    locations: [
      '茶园', '公园', '观景台', '海边', '海岛', '户外', '花店', '花市', '花园',
      '景点', '街拍', '森林', '医院花园', '植物园',
    ],
  },
  activity: {
    label: '活力同行',
    description: '一起体验新鲜有趣的活动',
    effect: '成功时健康、好感额外 +1',
    image: '/images/datinginfo/dating-scenes/modern-activity.webp',
    accent: '#7d8cff',
    locations: [
      '超市', '电竞馆', '赌场', '高尔夫', '家居城', '健身房', '街区', '俱乐部',
      '马戏团', '买手店', '赛事', '商场', '市集', '体育馆', '网红地', '训练场',
      '夜场', '游乐园', '游戏展', '展会', '直播间',
    ],
  },
  study: {
    label: '知性相伴',
    description: '在安静空间里了解彼此',
    effect: '成功时名声、好感额外 +1',
    image: '/images/datinginfo/dating-scenes/study-career.webp',
    accent: '#70b7d6',
    locations: [
      '大学', '工作室', '讲座', '科技馆', '科技园', '律所', '实验室', '书店',
      '书房', '图书馆', '校园', '医院', '幼儿园', '诊所', '证券所',
    ],
  },
  wellness: {
    label: '疗愈放松',
    description: '放慢节奏，享受安静陪伴',
    effect: '成功时健康额外 +3',
    image: '/images/datinginfo/dating-moments/bath-spa.webp',
    accent: '#70c9c4',
    locations: ['静修中心', '温泉', '养生馆', '瑜伽馆', '中医馆'],
  },
  luxury: {
    label: '浪漫礼遇',
    description: '用一场精心安排留下回忆',
    effect: '成功时名声 +1、好感额外 +2',
    image: '/images/datinginfo/dating-moments/date-sunset.webp',
    accent: '#ff9a7a',
    locations: [
      '度假村', '豪华酒店', '会所', '婚礼场地', '机场', '酒店', '酒会', '拍卖会',
      '沙龙', '时尚派对', '私人会所', '私人庄园', '游艇', '珠宝店', 'VIP厅',
    ],
  },
}

export const DATING_SCENE_GROUPS = Object.freeze(sceneGroups)

export const DATING_LOCATION_SCENES = Object.freeze(Object.fromEntries(
  Object.entries(sceneGroups).flatMap(([group, scene]) => (
    scene.locations.map(location => [location, group])
  )),
))

export const DATING_SCENE_OPTION_COUNT = 3

const shuffled = (values, random) => {
  const result = [...values]
  for (let index = result.length - 1; index > 0; index -= 1) {
    const swapIndex = Math.floor(random() * (index + 1))
    ;[result[index], result[swapIndex]] = [result[swapIndex], result[index]]
  }
  return result
}

// 每次约会保证至少有一个对方喜欢的地点，其余选项从非偏好地点中随机抽取。
export const createDatingSceneOptions = (
  preferredLocations,
  count = DATING_SCENE_OPTION_COUNT,
  random = Math.random,
) => {
  const allLocations = Object.keys(DATING_LOCATION_SCENES)
  const preferred = [...new Set(preferredLocations || [])]
    .filter(location => DATING_LOCATION_SCENES[location])
  const safeCount = Math.max(1, Math.min(Math.trunc(count) || DATING_SCENE_OPTION_COUNT, allLocations.length))

  if (!preferred.length) return shuffled(allLocations, random).slice(0, safeCount)

  const favorite = preferred[Math.floor(random() * preferred.length)]
  const preferredSet = new Set(preferred)
  const alternatives = allLocations.filter(location => !preferredSet.has(location))
  const otherFavorites = preferred.filter(location => location !== favorite)
  const fillers = [
    ...shuffled(alternatives, random),
    ...shuffled(otherFavorites, random),
  ].slice(0, safeCount - 1)
  return shuffled([
    favorite,
    ...fillers,
  ], random)
}

const iconRules = [
  [/咖啡|奶茶/, '☕'],
  [/餐厅|厨房|美食|烧烤|夜市/, '🍽️'],
  [/茶/, '🍵'],
  [/海边|海岛/, '🏝️'],
  [/花|植物/, '🌸'],
  [/公园|森林|户外/, '🌿'],
  [/医院|诊所|中医/, '🏥'],
  [/书|图书馆|大学|校园|讲座|幼儿园/, '📚'],
  [/电影|影院|片场|首映/, '🎬'],
  [/音乐|录音|琴房|KTV|电台/, '🎵'],
  [/剧|舞台|舞蹈|马戏/, '🎭'],
  [/画|艺术|博物馆|展览/, '🎨'],
  [/健身|训练|体育|赛事|高尔夫|瑜伽/, '🏅'],
  [/游戏|电竞|赌场|VIP/, '🎮'],
  [/机场/, '✈️'],
  [/游艇/, '⛵'],
  [/酒店|度假村/, '🏨'],
  [/商场|超市|市集|店|家居城|买手店/, '🛍️'],
  [/婚礼/, '💒'],
  [/温泉|养生|静修/, '♨️'],
  [/科技|实验室/, '🔬'],
  [/律所|证券所|商务|工作室/, '💼'],
  [/景点|观景台/, '🌇'],
  [/会所|俱乐部|酒会|沙龙|派对|夜场/, '🥂'],
]

export const getDatingLocationIcon = location => (
  iconRules.find(([pattern]) => pattern.test(location || ''))?.[1] || '💗'
)

export const getDatingScene = location => {
  const groupKey = DATING_LOCATION_SCENES[location] || 'dining'
  const group = sceneGroups[groupKey]
  return {
    group: groupKey,
    location: location || '约会地点',
    icon: getDatingLocationIcon(location),
    label: group.label,
    description: group.description,
    effect: group.effect,
    image: group.image,
    accent: group.accent,
  }
}
