package core

// 难度等级
const (
	DifficultyEasy   int = iota // 简单模式
	DifficultyNormal            // 普通模式（默认）
	DifficultyHard              // 困难模式
)

// 难度配置
type DifficultyConfig struct {
	Level        int     `json:"level"`
	Name         string  `json:"name"`
	InitMoney    int     `json:"initmoney"`    // 初始资金
	HealthBonus  float64 `json:"healthbonus"`  // 健康变化倍数
	BankruptRate float64 `json:"bankruptRate"` // 破产概率倍数
	AntiqueFake  float64 `json:"antiquefake"`  // 古董假货概率
	Description  string  `json:"description"`  // 描述
}

// 难度配置列表
var DifficultyConfigs = map[int]DifficultyConfig{
	DifficultyEasy: {
		Level:        DifficultyEasy,
		Name:         "简单",
		InitMoney:    1000000,
		HealthBonus:  0.5, // 健康负面影响减半
		BankruptRate: 0.5, // 破产概率减半
		AntiqueFake:  0.8, // 古董假货概率降低
		Description:  "初始资金100万，健康波动小，破产概率低，适合新手",
	},
	DifficultyNormal: {
		Level:        DifficultyNormal,
		Name:         "普通",
		InitMoney:    300000,
		HealthBonus:  1.0,
		BankruptRate: 1.0,
		AntiqueFake:  1.0,
		Description:  "初始资金30万，标准平衡，适合有经验的玩家",
	},
	DifficultyHard: {
		Level:        DifficultyHard,
		Name:         "困难",
		InitMoney:    50000,
		HealthBonus:  1.5, // 健康负面影响增加50%
		BankruptRate: 1.5, // 破产概率增加50%
		AntiqueFake:  1.2, // 古董假货概率增加
		Description:  "初始资金5万，健康波动大，破产概率高，极具挑战性",
	},
}

// GetDifficultyConfig 获取难度配置
func GetDifficultyConfig(level int) DifficultyConfig {
	if config, ok := DifficultyConfigs[level]; ok {
		return config
	}
	// 默认返回普通难度
	return DifficultyConfigs[DifficultyNormal]
}

// GetDifficultyName 获取难度名称
func GetDifficultyName(level int) string {
	config := GetDifficultyConfig(level)
	return config.Name
}
