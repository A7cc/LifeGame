package core

// =========== 公告结构体类型 ===========
// 公告信息
type Announce struct {
	AnnounceIns     []string `json:"announceins"`
	AnnounceOut     []string `json:"announceout"`
	AnnounceCompany []string `json:"announcecompany"`
	AnnounceGame    []string `json:"announcegame"`
	AnnounceHealthy []string `json:"announcehealthy"`
}

// =========== 用户和游戏机会结构体类型 ===========
type Opportunity struct {
	//打工次数
	OWNum int `json:"ownum"`
	// 玩小游戏次数
	OGNum int `json:"ognum"`
	// 约会次数
	OMNum int `json:"omnum"`
	// 逛街次数
	OSNum int `json:"osnum"`
	// 参加拍卖会次数
	OANum int `json:"oanum"`
}

// =========== 游戏可以持有物品最大数量 ===========
type MaxHoldNum struct {
	// 国内商品最大数
	MDHoldNum int `json:"mdholdnum"`
	// 国外商品最大数
	MFHoldNum int `json:"mfholdnum"`
	// 公司最大数
	MCHoldNum int `json:"mcholdnum"`
	// 古董持有最大数
	MAHoldNum int `json:"maholdnum"`
	// 最大名声
	MFAHoldNum int `json:"mfaholdnum"`
	// 最大免疫力
	MIHoldNum int `json:"miholdnum"`
	// 每回合打工最大次数
	MWRoundNum int `json:"mwroundnum"`
	// 每回合玩小游戏最大次数
	MGRoundNum int `json:"mgroundnum"`
	// 每回合约会最大次数
	MMRoundNum int `json:"mmroundnum"`
	// 每回合逛街最大次数
	MSRoundNum int `json:"msroundnum"`
	// 每回合参加拍卖会最大次数
	MARoundNum int `json:"maroundnum"`
}

// =========== 市场物资结构体类型 ===========
// 物资信息
type Item struct {
	// 物资ID
	IId int `json:"iid"`
	// 物资名字
	IName string `json:"iname"`
	// 物资平常价格
	IPrice int `json:"iprice"`
	// 物资最低价格
	IPrice_min int `json:"iprice_min"`
	// 物资最高价格
	IPrice_max int `json:"iprice_max"`
	// 物资可能发生的效果，[]map[int]string，第一个int是物资ID，第二个string是物资效果
	IEffects map[int]string `json:"ieffects"`
}

// 用户和游戏存放商品信息
type ItemInfo struct {
	// 商品名字
	IIName string `json:"iiname"`
	// 商品当前价格
	IIPrice int `json:"iiprice"`
	// 商品效果
	IIEffect string `json:"iieffect"`
	// 商品在当前游戏中是否显示
	IIDisplay bool `json:"iidisplay"`
}

// =========== 公司结构体类型 ===========
// 公司信息
type Company struct {
	// 公司ID
	CId int `json:"cid"`
	// 公司名称
	CName string `json:"cname"`
	// 公司基础价格（每手的）
	CPrice int `json:"cprice"`
	// 公司风险程度（1-10）
	CRisk int `json:"crisk"`
	// 公司盈利百分比（区间1-10），由于是百分比，需要除以100，获得的分红：一般是每手的价格*持有的数量（这是用户有的参数）
	CProfit int `json:"cprofit"`
	// 创业时间（1-5年）
	CTime int `json:"ctime"`
}

// 用户和游戏存放公司信息
type CompanyInfo struct {
	// 公司名字
	CIName string `json:"ciname"`
	// 公司当前价格（每手价格）
	CIPrice int `json:"ciprice"`
	// 公司风险程度
	CIRisk int `json:"cirisk"`
	// 公司盈利百分比
	CIProfit int `json:"ciprofit"`
	// 创业需要的时间
	CITime int `json:"citime"`
	// 公司是否破产,true是正常，false是破产
	CIStatus bool `json:"cistatus"`
}

// =========== 股票结构体类型 ===========
// 股票信息
type Stock struct {
	// 股票id
	SId int `json:"sid"`
	// 股票名称
	SName string `json:"sname"`
	// 股票价格
	SPrice int `json:"sprice"`
	// 股票风险程度，只有1-5个等级，越高风险越大
	SRisk int `json:"srisk"`
}

// 游戏存放股票信息
type StockInfo struct {
	// 股票id
	SIId int `json:"siid"`
	// 股票名称
	SIName string `json:"siname"`
	// 股票价格
	SIPrice int `json:"siprice"`
	// 股票风险程度
	SIRisk int `json:"sirisk"`
	// 股票历史
	SIHistory []int `json:"sihistory"`
	// K线历史，[开盘, 收盘, 最低, 最高]
	SIKlineHistory [][4]int `json:"siklinehistory"`
	// 涨跌状态
	SIStatus string `json:"sistatus"`
}

// 用户存放股票信息
type UserStockInfo struct {
	// 股票名称
	USName string `json:"usname"`
	// 用户成本价格，当前价格不能存放在用户数据里，应该是通过id获取
	USPrice_init int `json:"usprice_init"`
	// 持有数量
	USNum int `json:"usnum"`
	// 累计盈亏
	USProfit int `json:"usprofit"`
}

// ========== 古董结构体类型 ===========
// 古董信息
type Antique struct {
	// 古董ID
	AId int `json:"AId"`
	// 古董名称
	AName string `json:"aname"`
	// 古董基础价格
	APrice int `json:"aprice"`
	// 古董稀有度，越稀有价格越高，稀有度数值（1~10）
	AMaterial int `json:"amaterial"`
	// 古董图片
	AImg string `json:"aimg"`
	// 古董描述
	ADesc string `json:"adesc"`
	// 拍卖会等级
	ALevel int `json:"alevel"`
}

// 用户和游戏存放古董信息
type AntiqueInfo struct {
	// 古董ID，这里的id在没有拍卖前是古董的id，拍卖后是用户持有古董的id
	AIId int `json:"aiid"`
	// 古董名字
	AIName string `json:"ainame"`
	// 古董当前价格
	AIPrice int `json:"aiprice"`
	// 古董真伪情况（1是真，2是假，其他都是没有鉴定）
	AIDisplay int `json:"aiidisplay"`
	// 稀有度
	AIMaterial int `json:"aiamaterial"`
	// 完好程度
	AICondition int `json:"aiacondition"`
	// 最高价格
	AIPrice_max int `json:"aiprice_max"`
	// 用户持有时间
	AITime int `json:"aiatime"`
	// 图片
	AIImg string `json:"aiimg"`
	// 古董描述
	AIDesc string `json:"aidesc"`
	// 古董行等级
	AILevel int `json:"ailevel"`
}

// ========== 小游戏结构体类型 ===========
type MiniGame struct {
	// 游戏ID
	MGId int `json:"mgid"`
	// 游戏名称
	MGName string `json:"mgname"`
	// 中文名
	MGCName string `json:"mgcname"`
	// ICON
	MGIcon string `json:"mgicon"`
	// 游戏描述
	MGDesc string `json:"mgdesc"`
	// 游戏类型
	MGType string `json:"mgtype"`
	// 游戏难度：0=简单, 1=中等, 2=困难
	MGDifficulty map[int]SubMiniGame `json:"mgdifficulty"`
}

// 小游戏子信息
type SubMiniGame struct {
	// 游玩一次游戏需要的钱
	SMGNeed int `json:"smgneed"`
	// 需要达到的目标分数或者值
	SMGTarget int `json:"smgtarget"`
	// 游戏获得的奖励，map[string]int，key是奖励类型，value是奖励数量
	SMGReward map[string]int `json:"smgreward"`
	// 最少游玩游戏时间
	SMGMinRunTime int `json:"smgminruntime"`
}

// 小游戏的记录
type MiniGameRecord struct {
	// 游戏类型
	MGRType string `json:"mgrtype"`
	// 游玩次数
	PlayCount int `json:"playcount"`
	// 胜利次数
	WinCount int `json:"wincount"`
}

// 小游戏的Session
type MiniGameSession struct {
	// 服务端会话标识，用于幂等取消并防止旧窗口取消新游戏。
	MGSID string `json:"-"`
	// 游戏名字
	MGSName string `json:"mgsname"`
	// 游戏子信息
	MGSSubInfo SubMiniGame `json:"mgssubinfo"`
	// 游戏难度
	MGSLevel int `json:"mgslevel"`
	// 游戏类型
	MGSType string `json:"mgstype"`
	// Session 开始时间
	MGSStartTime int64 `json:"mgsstarttime"`
	// 本局由后端确认并扣除的下注金额
	MGSWager int `json:"mgswager"`
	// 本局下注选项（例如赛马编号、轮盘颜色）
	MGSChoice string `json:"mgschoice"`
	// 本局购买数量（目前用于彩票）
	MGSQuantity int `json:"mgsquantity"`
	// 本局总成本（报名费、下注或彩票总价）
	MGSCost int `json:"mgscost"`
	// 后端预先生成的固定奖金（目前用于彩票）
	MGSPayout int `json:"mgspayout"`
	// 后端生成的彩票奖面
	MGSTickets []int `json:"mgstickets,omitempty"`
	// 后端权威结果。纯随机游戏结算时忽略前端提交的结果值。
	MGSOutcome int `json:"-"`
	// 不返回给前端的游戏秘密（例如猜数字答案）。
	MGSSecret int `json:"-"`
	// 由后端生成、可用于前端动画的公开局面。
	MGSRound map[string]interface{} `json:"-"`
	// 二十一点由后端维护的牌堆和双方手牌。
	MGSDeck        []string `json:"-"`
	MGSPlayerCards []string `json:"-"`
	MGSDealerCards []string `json:"-"`
	MGSResolved    bool     `json:"-"`
}

// ========== 银行结构体类型 ===========
// 银行任务统计（记录今年的操作）
type BankTaskStats struct {
	// 本年度相对年初的银行净流入；存款为正，取款为负。
	NetBankFlow int `json:"netbankflow"`
	// 今年存钱总额
	DepositAmount int `json:"depositamount"`
	// 今年存钱次数
	DepositCount int `json:"depositcount"`
	// 今年取钱总额
	WithdrawAmount int `json:"withdrawamount"`
	// 今年取钱次数
	WithdrawCount int `json:"withdrawcount"`
	// 今年贷款总额
	LoanAmount int `json:"loanamount"`
	// 今年打工次数
	WorkCount int `json:"workcount"`
	// 今年已完成的任务ID列表
	CompletedTasks []int `json:"completedtasks"`
	// 今年已领取奖励的任务ID列表
	ClaimedTasks []int `json:"claimedtasks"`
	// 今年的任务列表（随机抽取的3个任务）
	CurrentTasks []BankTask `json:"currenttasks"`
}

// 银行任务配置
type BankTask struct {
	// 任务ID
	TaskId int `json:"taskid"`
	// 任务名称
	TaskName string `json:"taskname"`
	// 任务描述
	TaskDesc string `json:"taskdesc"`
	// 任务类型：deposit=存钱, withdraw=取钱, withdrawcount=取钱次数, loan=贷款, work=打工
	TaskType string `json:"tasktype"`
	// 目标值
	TargetValue int `json:"targetvalue"`
	// 奖励金额
	Reward int `json:"reward"`
}

// ========== 疾病系统结构体类型 ===========

// 疾病定义
type DiseaseInfo struct {
	// 疾病ID
	DId int `json:"did"`
	// 疾病名称
	DName string `json:"dname"`
	// 疾病类型
	DType string `json:"dtype"`
	// 症状描述
	DSymptoms string `json:"dsymptoms"`
	// 健康影响（每回合）
	DHealthImpact int `json:"dhealthimpact"`
	// 超过天数后疾病升级
	DUpgradeDays int `json:"dupgradedays"`
	// 可用治疗列表
	DTreatments []string `json:"dTreatments"`
}

// 治疗信息
type TreatInfo struct {
	// 治疗ID
	TId int `json:"tid"`
	// 治疗名称
	TName string `json:"tname"`
	// 治疗描述
	TDesc string `json:"tdesc"`
	// 价格
	TPrice int `json:"tprice"`
	// 恢复健康值
	THeal int `json:"theal"`
	// 副作用
	TSideEffect int `json:"tsideeffect"`
	// 来源
	TSource string `json:"tsource"`
}

// 医院卡片信息
type HospitalInfo struct {
	HId          int                   `json:"hid"`
	HType        string                `json:"htype"`
	HName        string                `json:"hname"`
	HIcon        string                `json:"hicon"`
	HDescription string                `json:"hdescription"`
	HServices    []HospitalServiceInfo `json:"hservices"`
}

// 医院服务
type HospitalServiceInfo struct {
	// 服务ID
	HSId int `json:"hsid"`
	// 服务名称
	HSName string `json:"hsname"`
	// 服务类型（medicine/injection/acupuncture/surgery）
	HSType string `json:"hstype"`
	// 价格
	HSPrice int `json:"hsprice"`
	// 恢复健康值
	HSImmunity int `json:"hsimmunity"`
	// 描述
	HSDesc string `json:"hsdesc"`
}

// 健康急诊状态
type HealthEmergencyStatus struct {
	// 是否必须先进行急诊
	Required bool `json:"required"`
	// 急诊原因
	Reasons []string `json:"reasons"`
	// 急诊费用
	Cost int `json:"cost"`
	// 触发急诊的严重疾病
	SevereDiseases []string `json:"severediseases"`
}

// ========== 约会系统结构体类型 ===========
// 遇见条件类型
type MeetCondition struct {
	// 条件类型：fame=名声, cash=现金, bank=存款, house=房子, car=车子,
	// play_game=玩过游戏, win_game=游戏获胜, work_count=打工次数,
	// age=年龄, date_count=约会次数, item_own=拥有物品,
	// stock_profit=股票盈利, lottery_win=彩票中奖, antique_rare=稀有古董,
	// company_founder=创业者, immunity=健康值
	CType string `json:"ctype"`
	// 条件值
	CValue int `json:"cvalue"`
	// 条件目标（可选，如游戏ID、物品ID等）
	CTarget string `json:"ctarget"`
}

// 约会对象信息
type DatingInfo struct {
	// 约会对象ID
	DId int `json:"did"`
	// 名字
	DName string `json:"dname"`
	// 图片
	DImage string `json:"dimage"`
	// 年龄
	DAge int `json:"dage"`
	// 性别：true=男，false=女
	DSex bool `json:"dsex"`
	// 国籍
	DNationality string `json:"dnationality"`
	// 职业
	DOccup string `json:"doccup"`
	// 性格/描述，分5个等级：非常容易、容易、中等、困难、非常困难
	DDesc string `json:"ddesc"`
	// 约会花费，范围：500-10000，价格区间：非常容易(500-1000)、(1000-2000)、(2000-5000)、(5000-8000)、(8000-10000)
	DCost int `json:"dcost"`
	// 遇见条件
	DMeetConditions []MeetCondition `json:"dmeetconditions"`
	// 喜欢的礼物类型
	DGifts []string `json:"dgifts"`
	// 喜欢的约会地点
	DLocations []string `json:"dlocations"`
	// 需要主动前往才能触发认识判定的场景；为空时按普通条件自动解锁。
	DMeetScene string `json:"dmeetscene"`
	// 是否已解锁（运行时计算，不存储）
	DUnlocked bool `json:"dunlocked"`
	// 好感度等级提示（运行时计算，不存储）
	DAffinityLevel string `json:"daffinitylevel"`
}

// 用户约会信息
type UserDatingInfo struct {
	// 约会对象ID
	DDatingId int `json:"ddatingid"`
	// 约会对象名字
	DName string `json:"dname"`
	// 好感度 0-100
	DAffinity int `json:"daffinity"`
	// 累计约会次数
	DCount int `json:"dcount"`
	// 累计赠送礼物次数。
	DGiftCount int `json:"dgiftcount"`
	// 关系状态：stranger=陌生人, friend=朋友, ambiguous=暧昧中, dating=交往中, lover=恋人, exclusive=专属, sweetheart=爱人, predecessor=前任
	DStatus string `json:"dstatus"`
}

// ========== 房子结构体类型 ===========
// 房子信息
type House struct {
	// 房子ID
	HId int `json:"hid"`
	// 房子名字
	HName string `json:"hname"`
	// 房子正常价格
	HPrice int `json:"hprice"`
	// 房子最高价格
	HPrice_max int `json:"hprice_max"`
	// 房子最低价格
	HPrice_min int `json:"hprice_min"`
	// 健康加成
	HHealth int `json:"hhealth"`
	// 名声加成
	HFame int `json:"hfame"`
	// 房子图片
	HImg string `json:"himg"`
}

// 游戏存放房子信息
type HouseInfo struct {
	// 房子id
	HIId int `json:"hiid"`
	// 房子名字
	HIName string `json:"hiname"`
	// 房子当前价格
	HIPrice int `json:"hiprice"`
	// 健康加成
	HIHealth int `json:"hihealth"`
	// 名声加成
	HIFame int `json:"hifame"`
	// 房子图片
	HIImg string `json:"hiimg"`
}

// ========== 车子结构体类型 ===========
// 车子信息
type Car struct {
	// 车子ID
	CId int `json:"cid"`
	// 车子名字
	CName string `json:"cname"`
	// 车子正常价格
	CPrice int `json:"cprice"`
	// 车子最高价格
	CPrice_max int `json:"cprice_max"`
	// 车子最低价格
	CPrice_min int `json:"cprice_min"`
	// 健康加成
	CHealth int `json:"chealth"`
	// 名声加成
	CFame int `json:"cfame"`
	// 车子图片
	CImg string `json:"cimg"`
}

// 游戏存放车子信息
type CarInfo struct {
	// 车子id
	CIId int `json:"ciid"`
	// 车子名字
	CIName string `json:"ciname"`
	// 车子当前价格
	CIPrice int `json:"ciprice"`
	// 健康加成
	CIHealth int `json:"cihealth"`
	// 名声加成
	CIFame int `json:"cifame"`
	// 车子图片
	CIImg string `json:"ciimg"`
}
