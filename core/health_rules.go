package core

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// GetMaxAnnualImmunityLoss 返回每个难度单次年度结算允许的最大免疫损失。
// 疾病仍然需要治疗，但不会在一个回合内把健康角色直接打到死亡。
func GetMaxAnnualImmunityLoss(difficulty int) int {
	switch difficulty {
	case DifficultyEasy:
		return 8
	case DifficultyHard:
		return 16
	default:
		return 12
	}
}

// GetImmunityEvent 获取免疫力变化事件，免疫力事件（每年调用）
func GetImmunityEvent(difficulty int, user *User) (int, string) {
	// 变化量
	change := 0
	// 变化描述
	reason := ""

	// 基础变化
	baseChange := rand.Intn(8) - 5

	// 疾病影响
	diseaseImpact := 0
	for _, ud := range user.UDiseases {
		diseaseImpact += ud.UHealthImpact * ud.UDSeverity / 2
	}

	change = baseChange + diseaseImpact
	if change < 0 {
		change = int(math.Round(float64(change) * GetDifficultyConfig(difficulty).HealthBonus))
	}
	maxLoss := GetMaxAnnualImmunityLoss(difficulty)
	limited := change < -maxLoss
	if limited {
		change = -maxLoss
	}

	// 生成原因
	if diseaseImpact < 0 {
		reason = "🤒 疾病影响，免疫力" + strconv.Itoa(change)
		if limited {
			reason += fmt.Sprintf("（本年度损失已限制为最多%d点）", maxLoss)
		}
	} else if change > 0 {
		reason = "💪 身体恢复，免疫力+" + strconv.Itoa(change)
	} else if change < 0 {
		reason = "😫 身体不适，免疫力" + strconv.Itoa(change)
	} else {
		reason = "😊 身体状况稳定"
	}

	return change, reason
}

// GetHealthEmergencyStatus 判断当前是否必须先进行急诊，出现免疫力低或者严重疾病时
func GetHealthEmergencyStatus(user *User, difficulty ...int) *HealthEmergencyStatus {
	// 初始化急诊状态变量
	status := &HealthEmergencyStatus{
		Required:       false,
		Reasons:        []string{},
		SevereDiseases: []string{},
	}
	// 初始化总费用
	totalCost := 0
	// 如果免疫力低于最低阈值，需要急诊
	if user.UImmunity < ImmunityThreshold {
		status.Required = true
		status.Reasons = append(status.Reasons, "免疫力低于"+strconv.Itoa(ImmunityThreshold))
		// 根据低于多少点计算费用，每低于1点增加DiseaseSeriousCostPerStep费用
		totalCost += ImmunityEmergencyBaseCost + (ImmunityThreshold-user.UImmunity)*DiseaseSeriousCostPerStep
	}
	// 遍历用户疾病，判断是否有需要急诊的疾病
	for _, disease := range user.UDiseases {
		if disease.UDSeverity < 5 {
			continue
		}
		status.Required = true
		// 记录需要急诊的疾病
		status.SevereDiseases = append(status.SevereDiseases, disease.UDName)
		status.Reasons = append(status.Reasons, disease.UDName+"需要急诊")
		totalCost += getHealthEmergencyDiseaseCost(disease)
	}
	// 如果需要急诊，计算总费用
	if status.Required {
		level := DifficultyNormal
		if len(difficulty) > 0 {
			level = difficulty[0]
		}
		costMultiplier := 1.0
		switch level {
		case DifficultyEasy:
			costMultiplier = 0.75
		case DifficultyHard:
			costMultiplier = 1.25
		}
		status.Cost = int(math.Ceil(float64(totalCost) * costMultiplier))
	}
	return status
}

// GetHealthEmergencyDiseaseCost 获取急诊费用，根据疾病类型和严重程度计算费用
func getHealthEmergencyDiseaseCost(disease UDiseaseInfo) int {
	cost := DiseaseSeriousBaseCost
	switch disease.UDType {
	case "green":
		cost = DiseaseSeriousBaseCost
	case "cyan":
		cost = DiseaseSeriousBaseCost + DiseaseSeriousCostPerStep
	case "yellow":
		cost = DiseaseSeriousBaseCost + DiseaseSeriousCostPerStep*2
	case "red":
		cost = DiseaseSeriousBaseCost + DiseaseSeriousCostPerStep*3
	}
	return cost
}

// ApplyHealthEmergencyTreatment 执行急诊治疗，清除严重疾病并把免疫力拉回安全线。
func ApplyHealthEmergencyTreatment(user *User) []UDiseaseInfo {
	// 初始化
	curedDiseases := []UDiseaseInfo{}
	// 获取用户的存在的严重疾病
	for id, disease := range user.UDiseases {
		if disease.UDSeverity < 5 {
			continue
		}
		curedDiseases = append(curedDiseases, disease)
		delete(user.UDiseases, id)
	}

	if user.UImmunity < DefaultImmunityEmergencyThreshold {
		user.UImmunity = DefaultImmunityEmergencyThreshold
	} else {
		user.UImmunity = CalcImmunity(user.UImmunity + 10)
	}
	user.UCriticalHealthYears = 0

	return curedDiseases
}

// AdvanceCriticalHealthYear 每年只调用一次。低免疫前两年只预警，第三年仍未
// 恢复到安全线才结束游戏，让玩家拥有完整的就医和筹款窗口。
func AdvanceCriticalHealthYear(user *User) (gameOver bool, rescueYearsLeft int) {
	if user == nil {
		return false, DefaultCriticalHealthGraceYears
	}
	if user.UImmunity >= ImmunityThreshold {
		user.UCriticalHealthYears = 0
		return false, DefaultCriticalHealthGraceYears
	}
	if user.UImmunity < DefaultMinimumSurvivableImmunity {
		user.UImmunity = DefaultMinimumSurvivableImmunity
	}
	user.UCriticalHealthYears++
	rescueYearsLeft = DefaultCriticalHealthGraceYears - user.UCriticalHealthYears
	if rescueYearsLeft < 0 {
		rescueYearsLeft = 0
	}
	return user.UCriticalHealthYears > DefaultCriticalHealthGraceYears, rescueYearsLeft
}

// ResetCriticalHealthIfRecovered 在治疗或其他活动把免疫力恢复到安全线时清除危机记录。
func ResetCriticalHealthIfRecovered(user *User) {
	if user != nil && user.UImmunity >= ImmunityThreshold {
		user.UCriticalHealthYears = 0
	}
}

// CalculateDiseaseChance 计算每年新增疾病概率。健康青年约为15%，高龄或
// 低免疫玩家逐步升高，最高55%。
func CalculateDiseaseChance(user *User) int {
	if user == nil {
		return 0
	}
	immunity := user.UImmunity
	if immunity < 0 {
		immunity = 0
	} else if immunity > MaxImmunity {
		immunity = MaxImmunity
	}
	chance := 15 + (MaxImmunity-immunity)/4
	if user.UAge > 30 {
		chance += (user.UAge - 30) / 8
	}
	if chance > 55 {
		chance = 55
	}
	return chance
}

// GetInitialDiseaseSeverity 控制新病初始严重度，禁止新生成疾病直接达到5级急诊。
func GetInitialDiseaseSeverity(diseaseType string) int {
	switch diseaseType {
	case "cyan":
		return rand.Intn(2) + 1
	case "yellow":
		return 2
	case "red":
		return 3
	default:
		return 1
	}
}

// 疾病生成与处理，生成疾病（每年调用）
func GenerateDisease(user *User) DiseaseInfo {
	// 计算生病概率
	// 免疫力越高，生病概率越低
	// 年龄越大，生病概率越高
	chance := CalculateDiseaseChance(user)

	// 随机判断是否生病
	if rand.Intn(100) >= chance {
		return DiseaseInfo{}
	}

	// 根据免疫力决定疾病类型
	var possibleDiseases []DiseaseInfo
	immunity := user.UImmunity
	age := user.UAge

	for _, d := range GetDiseases() {
		// 已患疾病不会在同一年再次被随机抽中并额外升级；病情只按持续时间升级。
		if _, active := user.UDiseases[d.DId]; active {
			continue
		}
		switch d.DType {
		case "green":
			// 小病（感冒、发烧等），小病谁都会得
			possibleDiseases = append(possibleDiseases, d)
		case "cyan":
			// 中等疾病（肠胃炎、扭伤等），免疫力<60或年龄>50
			if immunity < 60 || age > 50 {
				possibleDiseases = append(possibleDiseases, d)
			}
		case "yellow":
			// 严重疾病（骨折、肺炎等），免疫力<40且年龄>40
			if immunity < 40 && age > 40 {
				possibleDiseases = append(possibleDiseases, d)
			}
		case "red":
			// 危急疾病（心脏病、中风等），免疫力<20且年龄>60
			if immunity < 20 && age > 60 {
				possibleDiseases = append(possibleDiseases, d)
			}
		}
	}

	if len(possibleDiseases) == 0 {
		return DiseaseInfo{}
	}

	// 随机选择一种疾病
	disease := possibleDiseases[rand.Intn(len(possibleDiseases))]

	return disease
}

// GenerateAssetDisease 根据资产变化生成疾病
// currentAssets: 当前资产, prevAssets: 上一年资产
func GenerateAssetDisease(user *User, currentAssets, prevAssets int) (DiseaseInfo, string) {
	// 避免除零错误
	if prevAssets <= 0 {
		prevAssets = 1
	}

	// 计算资产变化百分比
	changePercent := float64(currentAssets-prevAssets) / float64(prevAssets) * 100

	// ==========资产大幅下跌（经济压力）==========
	// 资产下跌超过50% - 可能得抑郁症
	if changePercent <= -50 {
		// 30%概率得抑郁症
		if rand.Intn(100) < 30 {
			for _, disease := range GetDiseases() {
				if disease.DId == 15 {
					return disease, "📉 资产腰斩让你陷入深深的自责和抑郁..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// 资产下跌30%-50% - 可能得焦虑症
	if changePercent <= -30 && changePercent > -50 {
		// 35%概率得焦虑症
		if rand.Intn(100) < 35 {
			for _, disease := range GetDiseases() {
				if disease.DId == 13 {
					return disease, "😰 资产大幅缩水让你整夜睡不着觉..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// 资产下跌20%-30% - 可能得失眠症
	if changePercent <= -20 && changePercent > -30 {
		// 25%概率得失眠症
		if rand.Intn(100) < 25 {
			for _, disease := range GetDiseases() {
				if disease.DId == 14 {
					return disease, "😫 资产缩水让你开始失眠..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// ==========资产大幅增长或高额资产（富贵病）==========
	// 资产超过500万 - 可能得富贵病
	if currentAssets >= 5000000 {
		// 年龄越大越容易得富贵病
		ageFactor := user.UAge / 10
		chance := 10 + ageFactor // 基础10% + 年龄因素
		if chance > 30 {
			chance = 30
		}

		if rand.Intn(100) < chance {
			// 随机选择一种富贵病
			wealthDiseases := []int{17, 18, 19, 20} // 高血压、痛风、糖尿病、脂肪肝
			selectedDisease := wealthDiseases[rand.Intn(len(wealthDiseases))]

			var msg string
			switch selectedDisease {
			case 17:
				msg = "🍽️ 大鱼大肉的生活让你得了高血压..."
			case 18:
				msg = "🍺 应酬太多，痛风找上门了..."
			case 19:
				msg = "🍰 甜食吃太多，血糖开始失控..."
			case 20:
				msg = "🥩 油腻饮食让肝脏不堪重负..."
			}
			for _, disease := range GetDiseases() {
				if disease.DId == selectedDisease {
					return disease, msg
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// ==========资产暴涨相关疾病==========
	// 资产暴涨超过150% - 可能得焦虑症（担心守不住财富）
	if changePercent >= 150 {
		// 25%概率得焦虑症
		if rand.Intn(100) < 25 {
			for _, disease := range GetDiseases() {
				if disease.DId == 13 {
					return disease, "😰 资产暴涨让你患得患失，总担心守不住财富..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// 资产暴涨100%-150% - 可能得神经衰弱
	if changePercent >= 100 && changePercent < 150 {
		// 20%概率得神经衰弱
		if rand.Intn(100) < 20 {
			for _, disease := range GetDiseases() {
				if disease.DId == 16 {
					return disease, "🤯 资产暴涨让你兴奋过度，有点神经衰弱了..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	// 资产暴涨50%-100% - 可能得失眠症（兴奋睡不着）
	if changePercent >= 50 && changePercent < 100 {
		// 15%概率得失眠症
		if rand.Intn(100) < 15 {
			for _, disease := range GetDiseases() {
				if disease.DId == 14 {
					return disease, "😆 资产大涨让你兴奋得睡不着觉..."
				}
			}
			return DiseaseInfo{}, ""
		}
	}

	return DiseaseInfo{}, ""
}
