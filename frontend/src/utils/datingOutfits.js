export const DATING_OUTFITS = [
  { key: 'career', label: '职业装', icon: '💼', minRank: 1 },
  { key: 'homewear', label: '居家装', icon: '🏠', minRank: 1 },
  { key: 'qipao', label: '旗袍/国风', icon: '🌸', minRank: 2 },
  { key: 'cosplay', label: 'Cosplay', icon: '🎭', minRank: 3 },
  { key: 'swimwear', label: '泳装', icon: '🌊', minRank: 3 },
  { key: 'sleepwear', label: '睡衣', icon: '🌙', minRank: 4 },
  { key: 'romantic', label: '情趣睡衣', icon: '🌹', minRank: 5, adultOnly: true },
]

const STATUS_RANKS = {
  '陌生人': 0,
  '朋友': 1,
  '暧昧中': 2,
  '交往中': 3,
  '恋人': 4,
  '专属恋人': 5,
  '爱人': 6,
  '已婚': 7,
}

export const datingRelationshipRank = status => STATUS_RANKS[status] ?? 0

export const availableDatingOutfits = (status, playerAge, partnerAge) => {
  const rank = datingRelationshipRank(status)
  const bothAdults = Number(playerAge || 0) >= 18 && Number(partnerAge || 0) >= 18
  return DATING_OUTFITS.filter(outfit => (
    rank >= outfit.minRank && (!outfit.adultOnly || bothAdults)
  ))
}

export const findDatingOutfit = key => DATING_OUTFITS.find(outfit => outfit.key === key)

export const datingOutfitImage = (image, outfit) => {
  if (!image || outfit === 'career' || !findDatingOutfit(outfit)) return image || ''
  const slash = image.lastIndexOf('/')
  if (slash < 0) return image
  return `${image.slice(0, slash)}/${outfit}/${image.slice(slash + 1)}`
}
