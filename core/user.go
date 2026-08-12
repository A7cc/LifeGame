package core

import (
	"errors"
	"fmt"
)

// 用户持有的物资信息
type UItemInfo struct {
	// 成本价格
	UIICostPrice int `json:"uicostprice"`
	// 物资数量
	UIINum int `json:"uitemnum"`
}

// 用户持有的公司信息
type UCompanyInfo struct {
	// 公司名字
	UCompanyName string `json:"ucompanyname"`
	// 持有的公司时间
	UCompanyHoldTime int `json:"ucompanyholdtime"`
	// 成本价
	UCompanyCostPrice int `json:"ucompanycostprice"`
	// 持有数量（一手价格）
	UCompanyNum int `json:"ucompanynum"`
}

// 用户当前疾病
type UDiseaseInfo struct {
	// 疾病名字
	UDName string `json:"udname"`
	// 疾病类型
	UDType string `json:"udtype"`
	// 症状描述
	USymptoms string `json:"usymptoms"`
	// 健康影响（每回合）
	UHealthImpact int `json:"uhealthimpact"`
	// 超过天数后疾病升级
	UUpgradeDays int `json:"uupgradedays"`
	// 可用治疗列表
	UTreatments []string `json:"utreatments"`
	// 严重程度倍数（1-5），1为最轻，5为最重，5表示必须打针或者手术
	UDSeverity int `json:"udseverity"`
	// 已持续时间
	UDTime int `json:"udtime"`
}

// 用户信息
type User struct {
	// 用户ID
	Uid int `json:"uid"`
	// 名字
	UName string `json:"uname"`
	// 性别
	USex bool `json:"usex"`
	// 岁数
	UAge int `json:"uage"`
	// 免疫力
	UImmunity int `json:"uimmunity"`
	// 连续处于低免疫危急状态的年度数；达到安全线后清零。
	UCriticalHealthYears int `json:"ucriticalhealthyears"`
	// 当前疾病列表
	UDiseases map[int]UDiseaseInfo `json:"udiseases"`
	// 名声
	UFame int `json:"ufame"`
	// 存放打工、游戏、约会次数、逛街等参数
	UOpportunity Opportunity `json:"uopportunity"`
	// 净资产=现金+存款+物资等资产-贷款
	UAssets int `json:"uassets"`
	// 现金
	UCash int `json:"ucash"`
	// 存款
	UBank int `json:"ubank"`
	// 贷款金额
	ULoan int `json:"uloan"`
	// 贷款逾期次数
	ULoanOverdue int `json:"uloanoverdue"`
	// 国内物资，map[int]UItemInfo，第一个int是物资ID，第二个物资信息
	UItemins map[int]UItemInfo `json:"uitemins"`
	// 国外物资，map[int]UItemInfo，第一个int是物资ID，第二个物资信息
	UItemout map[int]UItemInfo `json:"uitemout"`
	// 古董，最多只能有10个古董
	UAntique []AntiqueInfo `json:"uantique"`
	// 存放全部公司，但是手上最多只能有创业3个公司，所以想要设置持有数量，，map[int]UItemInfo，第一个int是公司ID，第二个用户持有的公司信息
	UCompany map[int]UCompanyInfo `json:"ucompany"`
	// 股票
	UStock map[int]UserStockInfo `json:"ustock"`
	// 记录股票当天盈利
	UStockProfit int `json:"ustockprofit"`
	// 约会关系信息，map[约会对象ID]约会信息
	UDating map[int]UserDatingInfo `json:"udating"`
	// 当前唯一配偶的约会对象 ID；0 表示未婚。
	UMarriedDatingID int `json:"umarrieddatingid"`
	// 车子
	UCar map[int]bool `json:"ucar"`
	// 房子
	UHouse map[int]bool `json:"uhouse"`
	// 小游戏总累计统计，key 为小游戏 ID，空字符串汇总所有小游戏
	UMiniGameRecords map[string]MiniGameRecord `json:"uminigamerecords"`
}

// 初始化User
// uname: 用户名，usex: 性别，uage: 年龄，ucash: 现金，ufame: 名声
func NewUser(uname string, usex bool, maxitemnum, uage, ucash int, CompanyData []Company) *User {
	return &User{
		// id
		Uid: 1,
		// 名字
		UName: uname,
		// 性别
		USex: usex,
		// 年龄
		UAge: uage,
		// 全部资产
		UAssets: ucash,
		// 现金
		UCash: ucash,
		// 存款
		UBank: 0,
		// 名声
		UFame: 0,
		// 免疫力
		UImmunity: 80,
		// 疾病
		UDiseases: make(map[int]UDiseaseInfo),
		// 国内物资，map[int]UItemInfo，第一个int是物资ID，第二个物资信息
		UItemins: make(map[int]UItemInfo),
		// 国外物资，map[int]UItemInfo，第一个int是物资ID，第二个物资信息
		UItemout: make(map[int]UItemInfo),
		// 古董，最多只能有10个古董
		UAntique: []AntiqueInfo{},
		// 股票
		UStock: make(map[int]UserStockInfo),
		// 公司
		UCompany: make(map[int]UCompanyInfo),
		// 每回合打工、游戏、约会次数、逛街等参数
		UOpportunity: Opportunity{
			// 打工次数
			OWNum: 0,
			// 游戏次数
			OGNum: 0,
			// 约会次数
			OMNum: 0,
			// 逛街次数
			OSNum: 0,
			// 购买次数
			OANum: 0,
		},
		// 贷款金额
		ULoan: 0,
		// 贷款逾期次数
		ULoanOverdue: 0,
		// 房子
		UHouse: make(map[int]bool),
		// 车子
		UCar: make(map[int]bool),
		// 约会关系信息
		UDating: make(map[int]UserDatingInfo),
		// 小游戏总的累计统计，key 为小游戏 ID，空字符串汇总所有小游戏
		UMiniGameRecords: make(map[string]MiniGameRecord),
	}
}

// 刷新用户资产并校验状态
// 返回值依次为：当前资产、刷新前资产、错误
func (u *User) RefreshAndValidateUserState(gameinfo *Game) (int, int, error) {
	// 获取之前的用户资产
	previousAssets := u.UAssets
	// 获取现在的总资产
	u.UAssets = CalculateUserAssets(u, gameinfo)
	// 检测资产
	if u.UAssets <= 0 {
		return u.UAssets, previousAssets, errors.New("资产为0，游戏结束！")
	}
	// 免疫力归零不再因一次结算立即死亡，而是进入连续低免疫抢救期。
	if u.UImmunity <= 0 {
		u.UImmunity = DefaultMinimumSurvivableImmunity
	} else if u.UImmunity > MaxImmunity {
		u.UImmunity = MaxImmunity
	}
	ResetCriticalHealthIfRecovered(u)
	// 检测年龄
	if u.UAge >= gameinfo.GTime {
		return u.UAssets, previousAssets, errors.New("年龄为100，游戏结束！")
	}
	// 校验名声，名声不能小于0，大于500
	if u.UFame <= 0 {
		u.UFame = 0
	} else if u.UFame >= MaxFame {
		u.UFame = MaxFame
	}

	return u.UAssets, previousAssets, nil
}

// 用户达到某个资产里程碑，增加名声
func (u *User) ApplyAssetFameMilestones(previousAssets, currentAssets int) []string {

	milestones := []struct {
		assets int
		fame   int
	}{
		{assets: 100000, fame: 1},
		{assets: 1000000, fame: 2},
		{assets: 10000000, fame: 3},
		{assets: 100000000, fame: 5},
	}

	announces := []string{}
	for _, milestone := range milestones {
		if previousAssets >= milestone.assets || currentAssets < milestone.assets {
			continue
		}
		u.UFame = CalcFame(u.UFame + milestone.fame)
		announces = append(announces, fmt.Sprintf("🌟 资产突破 %d 元，名声+%d", milestone.assets, milestone.fame))
	}
	return announces
}
