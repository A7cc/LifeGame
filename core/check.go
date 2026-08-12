package core

// 计算免疫力是否达到要求，返回数值
func CalcImmunity(immunity int) int {
	// 检测免疫力
	if immunity <= 0 {
		return 0
	} else if immunity > MaxImmunity {
		immunity = MaxImmunity
	}
	return immunity
}

// 计算名声是否达到要求，返回数值
func CalcFame(fame int) int {
	// 检测名声
	if fame < 0 {
		fame = 0
	} else if fame > MaxFame {
		fame = MaxFame
	}
	return fame
}

// 计算疾病严重性，返回数值
func CalcDisease(disease int) int {
	// 检测疾病严重性
	if disease < 1 {
		disease = 1
	} else if disease >= 5 {
		disease = 5
	}
	return disease
}
