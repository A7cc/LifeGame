package core

// 计算用户名声，在什么级别
func CalcReputationLevel(reputation int) (int, string, int, int) {
	// 退出的参数：等级，描述，股票额度，公司数
	// 股票额度表示该回合股票盈利最大，根据名气，名气普通：10w，中级：50w，高级：100w，豪华500w，私人：1000w
	switch {
	case reputation >= 0 && reputation <= 70:
		// 普通
		return 0, "普通", 100000, 1
	case reputation > 70 && reputation <= 110:
		// 中等
		return 1, "中等", 500000, 2
	case reputation > 110 && reputation <= 130:
		// 高级
		return 2, "高级", 1000000, 3
	case reputation > 130 && reputation <= 145:
		// 豪华
		return 3, "豪华", 5000000, 4
	case reputation > 145 && reputation <= 150:
		// 私人
		return 4, "私人", 10000000, 5
	default:
		// 老赖
		return -1, "老赖", 0, 0
	}
}
