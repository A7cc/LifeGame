export const DATING_CAREER_SCENE_COUNT = 50

const careerIcons = [
  [/歌手|演员|导演|舞蹈|舞者|钢琴|小提琴|马戏|话剧/, '🎭'],
  [/学生|图书|博士|教授|老师/, '📚'],
  [/白领|名媛|投资|商界|律师|企业家|经理|金融/, '💼'],
  [/主播|程序员|科学家|电台/, '💻'],
  [/医生|护士|中医|心理|瑜伽/, '🩺'],
  [/教练|冠军/, '🏅'],
  [/艺术|设计|摄影|插画|花艺|花店/, '🎨'],
  [/消防/, '🚒'],
  [/空姐|空乘|飞行员/, '✈️'],
  [/厨师|家政/, '🏠'],
  [/导游/, '🧭'],
  [/荷官/, '♠️'],
  [/模特/, '✨'],
  [/婚礼/, '💐'],
]

const sceneAccents = [
  '#70b7d6', '#7895ff', '#d38cff', '#d4a267', '#7cc6a2',
  '#8b7cff', '#ef7e72', '#db9b55', '#ef8db2', '#6bc1bc',
]

const normalizeCareerSceneId = datingId => {
  const numericId = Number(datingId)
  if (!Number.isInteger(numericId) || numericId < 1) return 1
  return ((numericId - 1) % DATING_CAREER_SCENE_COUNT) + 1
}

export const getDatingCareerScene = dating => {
  const sceneId = normalizeCareerSceneId(dating?.did)
  const occupation = String(dating?.doccup || '职业环境')
  const icon = careerIcons.find(([pattern]) => pattern.test(occupation))?.[1] || '💼'
  const filename = String(sceneId).padStart(2, '0')

  return {
    id: sceneId,
    icon,
    label: `${occupation}环境`,
    occupation,
    image: `/images/datinginfo/dating-careers/${filename}.webp`,
    accent: sceneAccents[(sceneId - 1) % sceneAccents.length],
  }
}
