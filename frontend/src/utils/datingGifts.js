const giftSceneRules = [
  {
    pattern: /钻石|钻戒|珠宝|宝石|名表|豪宅|私人飞机|礼服|名牌包|首饰|手表|香水/,
    icon: '✨', label: '精致礼遇', image: '/images/datinginfo/dating-moments/date-sunset.webp',
  },
  {
    pattern: /鲜花|盆栽|花瓶|园艺|香薰|精油/,
    icon: '🌸', label: '自然心意', image: '/images/datinginfo/dating-scenes/nature-travel.webp',
  },
  {
    pattern: /巧克力|红酒|奶茶|零食|咖啡|茶|食材|饮料|养生品|中药材/,
    icon: '🍫', label: '味觉惊喜', image: '/images/datinginfo/dating-scenes/dining-city.webp',
  },
  {
    pattern: /画|乐谱|音乐|专辑|剧本|戏服|舞衣|舞鞋|电影|婚礼杂志|摄影集|胶片/,
    icon: '🎨', label: '艺术心意', image: '/images/datinginfo/dating-scenes/arts-culture.webp',
  },
  {
    pattern: /游戏|耳机|手机|相机|镜头|设备|器材|装备|护具|运动|瑜伽|键盘|模型|道具/,
    icon: '🎮', label: '兴趣好礼', image: '/images/datinginfo/dating-scenes/modern-activity.webp',
  },
]

const fallbackGiftScene = {
  icon: '🎁', label: '用心挑选', image: '/images/datinginfo/dating-scenes/study-career.webp',
}

export const getDatingGiftScene = gift => (
  giftSceneRules.find(rule => rule.pattern.test(gift || '')) || fallbackGiftScene
)
