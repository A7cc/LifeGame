package db

import (
	"encoding/json"
	"fmt"
	"strings"
)

// seedData 按稳定 ID 幂等写入当前版本的内置参考数据。
func seedData(db sqlExecutor) error {
	// 写入物资数据
	if err := seedItems(db); err != nil {
		return err
	}

	// 写入公司数据
	if err := seedCompanies(db); err != nil {
		return err
	}

	// 写入古董数据
	if err := seedAntiques(db); err != nil {
		return err
	}

	// 写入股票数据
	if err := seedStocks(db); err != nil {
		return err
	}

	// 写入股票新闻数据
	if err := seedStockNews(db); err != nil {
		return err
	}

	// 写入房产数据
	if err := seedHouses(db); err != nil {
		return err
	}

	// 写入车辆数据
	if err := seedCars(db); err != nil {
		return err
	}

	// 写入初始约会对象数据
	if err := seedDatingInfo(db); err != nil {
		return err
	}

	// 写入银行任务数据
	if err := seedBankTasks(db); err != nil {
		return err
	}

	// 写入疾病数据
	if err := seedDiseases(db); err != nil {
		return err
	}

	// 写入药品数据
	if err := seedTreats(db); err != nil {
		return err
	}

	// 写入医院数据
	if err := seedHospitals(db); err != nil {
		return err
	}

	// 写入小游戏数据
	if err := seedMiniGames(db); err != nil {
		return err
	}

	return nil
}

// seedItems 写入物资数据
func seedItems(db sqlExecutor) error {
	// 国内物资
	domesticItems := []struct {
		ID       int
		Name     string
		Price    int
		PriceMin int
		PriceMax int
		Effects  map[int]string
		Region   string
	}{
		{1, "二手汽车", 70000, 20000, 120000, map[int]string{
			0: "⚖️ 二手汽车价格维持在合理区间",
			1: "📈 油价持续走低，二手汽车购入需求上涨",
			2: "🪂 环保政策加码，二手汽车需求显著下滑",
		}, "domestic"},
		{2, "电视机", 5000, 1500, 8000, map[int]string{
			0: "⚖️ 电视机目前销售平稳",
			1: "🔥 世界杯赛事刺激电视机热销",
			2: "📉 随着智能投影普及，电视机销量受到冲击",
		}, "domestic"},
		{3, "洗衣机", 2200, 1500, 3500, map[int]string{
			0: "⚖️ 洗衣机供需关系稳定",
			1: "📈 核心部件价格上涨推高洗衣机售价",
			2: "🪂 促销季来临，洗衣机纷纷降价清仓",
		}, "domestic"},
		{4, "二手手机", 4200, 1500, 10000, map[int]string{
			0: "⚖️ 当前市场对二手手机接受度良好",
			1: "🔥 绿色消费观念推动二手手机升温",
			2: "📉 新机上市拉低二手手机价格",
		}, "domestic"},
		{5, "电饭煲", 300, 150, 600, map[int]string{
			0: "⚖️ 电饭煲售价无明显波动",
			1: "🔥 年末促销活动让电饭煲销售翻番",
			2: "🪂 电饭煲市场饱和，价格趋于走低",
		}, "domestic"},
		{6, "羽绒服", 1000, 500, 1500, map[int]string{
			0: "⚖️ 羽绒服价格呈季节性波动",
			1: "🔥 强寒潮袭来引发羽绒服抢购潮",
			2: "📉 气温回暖过快，羽绒服大量滞销",
		}, "domestic"},
		{7, "矿泉水", 4, 1, 10, map[int]string{
			0: "⚖️ 矿泉水价格稳定",
			1: "🔥 夏季高温天气，矿泉水销量激增，价格也随着上涨",
			2: "📉 水源污染事件导致矿泉水价格下跌",
		}, "domestic"},
		{8, "方便面", 10, 5, 20, map[int]string{
			0: "⚖️ 方便面销量波动不大",
			1: "📈 推出新品口味带动方便面销售增长",
			2: "🪂 食品安全事件打击方便面销量",
		}, "domestic"},
		{9, "洗发水", 25, 15, 50, map[int]string{
			0: "⚖️ 洗发水价格总体平稳",
			1: "📈 品牌促销活动助推洗发水销量",
			2: "📉 老款洗发水因更新换代出现库存积压",
		}, "domestic"},
		{10, "自行车", 1000, 100, 1800, map[int]string{
			0: "⚖️ 自行车销售情况保持正常",
			1: "📈 绿色出行政策提升自行车需求",
			2: "📉 共享单车普及拉低自行车价格",
		}, "domestic"},
		{11, "纸巾", 30, 5, 80, map[int]string{
			0: "⚖️ 纸巾市场稳定",
			1: "🔥 纸巾在促销活动中销量猛增",
			2: "📉 物流延迟导致纸巾供给紧张",
		}, "domestic"},
		{12, "进口香水", 350, 300, 400, map[int]string{
			0: "⚖️ 进口香水作为高端消费品表现稳定",
			1: "🔥 明星代言效应带动进口香水热销",
			2: "🪂 经济放缓抑制了进口香水消费",
		}, "domestic"},
		{13, "鸡蛋", 10, 8, 15, map[int]string{
			0: "⚖️ 鸡蛋价格处于正常水平",
			1: "📈 养殖规模扩大使鸡蛋供需紧张",
			2: "📉 禽流感爆发引发鸡蛋价格回落",
		}, "domestic"},
		{14, "大米", 25, 20, 35, map[int]string{
			0: "⚖️ 大米供应充足，价格稳定",
			1: "📈 自然灾害影响大米产量上涨",
			2: "🪂 政府储备粮投放压低大米市场价",
		}, "domestic"},
	}

	// 国外物资（ID 从 15 开始，避免与国内物资冲突）
	foreignItems := []struct {
		ID       int
		Name     string
		Price    int
		PriceMin int
		PriceMax int
		Effects  map[int]string
		Region   string
	}{
		{15, "挪威三文鱼", 150, 50, 500, map[int]string{
			0: "⚖️ 挪威三文鱼价格整体维持在稳定区间",
			1: "📈 健康饮食趋势兴起带动挪威三文鱼价格上涨",
			2: "🪂 由于渔业限制，挪威三文鱼价格呈下跌趋势",
		}, "foreign"},
		{16, "法国葡萄酒", 1500, 200, 7000, map[int]string{
			0: "⚖️ 法国葡萄酒市场维持平稳需求",
			1: "📈 法国葡萄酒在年份酒推出后迎来价格上涨",
			2: "🪂 经济放缓导致法国葡萄酒价格略有下滑",
		}, "foreign"},
		{17, "意大利橄榄油", 100, 10, 600, map[int]string{
			0: "⚖️ 意大利橄榄油目前处于价格稳定阶段",
			1: "🔥 随着品质提升，意大利橄榄油价格走高",
			2: "📉 受到本地生产商冲击，意大利橄榄油价格下跌",
		}, "foreign"},
		{18, "日本绿茶", 100, 10, 400, map[int]string{
			0: "⚖️ 目前日本绿茶价格平稳，市场反应良好",
			1: "📈 健康消费趋势推动日本绿茶价格上扬",
			2: "🪂 天气不利影响日本绿茶产量，价格出现下行",
		}, "foreign"},
		{19, "瑞士手表", 10000, 1000, 100000, map[int]string{
			0: "⚖️ 瑞士手表作为奢侈品，其价格保持稳定",
			1: "📈 限量款热销推动瑞士手表价格大幅上涨",
			2: "🪂 经济低迷拖累瑞士手表市场价格下滑",
		}, "foreign"},
		{20, "德国汽车", 200000, 100000, 1000000, map[int]string{
			0: "⚖️ 德国汽车价格维持在相对稳定水平",
			1: "📈 新车款上市提升了德国汽车整体价格",
			2: "🪂 环保法规施压下，德国汽车价格有所下降",
		}, "foreign"},
		{21, "西班牙火腿", 100, 10, 300, map[int]string{
			0: "⚖️ 目前西班牙火腿价格稳定，市场接受度高",
			1: "📈 健康潮流促使西班牙火腿价格逐步上涨",
			2: "📉 进口政策变化导致西班牙火腿价格下行",
		}, "foreign"},
		{22, "澳大利亚羊毛", 100, 20, 700, map[int]string{
			0: "⚖️ 澳大利亚羊毛价格总体保持稳定",
			1: "🔥 天气影响供应，澳大利亚羊毛价格上涨明显",
			2: "🪂 受替代品竞争影响，澳大利亚羊毛价格下跌",
		}, "foreign"},
		{23, "美国牛肉", 150, 100, 500, map[int]string{
			0: "⚖️ 美国牛肉在国际市场上价格较为稳定",
			1: "📈 进口量上升拉动美国牛肉价格上涨",
			2: "🪂 人们开始提倡吃羊肉，使得美国牛肉价格有所下降",
		}, "foreign"},
		{24, "哥伦比亚咖啡", 100, 50, 300, map[int]string{
			0: "⚖️ 哥伦比亚咖啡维持稳定价格区间",
			1: "🔥 需求增长促使哥伦比亚咖啡价格提升",
			2: "📉 气候变化压制产量，使哥伦比亚咖啡价格下滑",
		}, "foreign"},
		{25, "比利时巧克力", 50, 10, 300, map[int]string{
			0: "⚖️ 比利时巧克力目前价格稳定，节假日销量影响不大",
			1: "🔥 新推出系列带动比利时巧克力价格上涨",
			2: "🪂 市场竞争激烈导致比利时巧克力价格走低",
		}, "foreign"},
		{26, "荷兰郁金香球茎", 120, 50, 900, map[int]string{
			0: "⚖️ 荷兰郁金香球茎价格保持稳定，供需均衡",
			1: "📈 进口需求增长，荷兰郁金香球茎价格上扬",
			2: "📉 积存积压使荷兰郁金香球茎价格下滑",
		}, "foreign"},
		{27, "加拿大枫糖浆", 200, 45, 800, map[int]string{
			0: "⚖️ 加拿大枫糖浆维持平稳售价，市场接受良好",
			1: "🔥 近日促销带动加拿大枫糖浆价格小幅上涨",
			2: "📉 替代品增多导致加拿大枫糖浆价格下降",
		}, "foreign"},
		{28, "巴西咖啡豆", 100, 30, 550, map[int]string{
			0: "⚖️ 巴西咖啡豆供应稳定，价格基本不变",
			1: "🔥 全球需求复苏推高了巴西咖啡豆价格",
			2: "📉 进口减少拖累巴西咖啡豆价格下行",
		}, "foreign"},
	}

	// 插入国内物资
	for _, item := range domesticItems {
		effectsJSON, _ := json.Marshal(item.Effects)
		_, err := db.Exec(`
			INSERT INTO items (id, name, price, price_min, price_max, effects, region)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, price_min=excluded.price_min,
				price_max=excluded.price_max, effects=excluded.effects, region=excluded.region
		`, item.ID, item.Name, item.Price, item.PriceMin, item.PriceMax, string(effectsJSON), item.Region)
		if err != nil {
			return err
		}
	}

	// 插入国外物资
	for _, item := range foreignItems {
		effectsJSON, _ := json.Marshal(item.Effects)
		_, err := db.Exec(`
			INSERT INTO items (id, name, price, price_min, price_max, effects, region)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, price_min=excluded.price_min,
				price_max=excluded.price_max, effects=excluded.effects, region=excluded.region
		`, item.ID, item.Name, item.Price, item.PriceMin, item.PriceMax, string(effectsJSON), item.Region)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedCompanies 写入公司数据
func seedCompanies(db sqlExecutor) error {
	companies := []struct {
		ID     int
		Name   string
		Price  int
		Risk   int
		Profit int
		Time   int
	}{
		{1, "星链智能公司", 120, 6, 8, 3},
		{2, "阳光新能源公司", 95, 6, 6, 2},
		{3, "飞驰物流公司", 60, 7, 5, 1},
		{4, "康瑞生物公司", 200, 4, 8, 4},
		{5, "远洋地产公司", 80, 8, 4, 2},
		{6, "恒星资本公司", 170, 5, 8, 3},
		{7, "钢毅制造公司", 110, 5, 7, 2},
		{8, "蓝天航运公司", 150, 8, 5, 4},
		{9, "星视传媒公司", 75, 6, 5, 1},
		{10, "绿野农科公司", 65, 4, 6, 1},
		{11, "光影游戏公司", 90, 6, 6, 2},
		{12, "雷动电车公司", 180, 7, 8, 3},
		{13, "碧海环保公司", 85, 5, 7, 2},
		{14, "博学教育公司", 70, 4, 5, 1},
		{15, "链信科技公司", 135, 8, 7, 4},
		{16, "智汇硬件公司", 120, 6, 7, 3},
		{17, "快购电商公司", 80, 5, 5, 2},
		{18, "安康医疗公司", 190, 4, 8, 4},
		{19, "绿筑建设公司", 85, 6, 5, 2},
		{20, "未来机器人公司", 155, 7, 8, 3},
		{21, "味道小吃公司", 30, 4, 5, 1},
		{22, "乐居装修公司", 40, 5, 5, 1},
		{23, "幸福宠物店公司", 35, 4, 6, 1},
		{24, "趣游旅行社公司", 50, 5, 5, 2},
		{25, "闪购便利店公司", 25, 3, 4, 1},
		{26, "潮流服饰公司", 45, 5, 5, 2},
		{27, "甜心咖啡公司", 40, 4, 5, 1},
		{28, "健康健身馆公司", 55, 5, 6, 1},
		{29, "安心保洁公司", 30, 3, 4, 1},
		{30, "开心托儿所公司", 50, 4, 5, 1},
	}

	for _, c := range companies {
		_, err := db.Exec(`
			INSERT INTO companies (id, name, price, risk, profit, time)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, risk=excluded.risk,
				profit=excluded.profit, time=excluded.time
		`, c.ID, c.Name, c.Price, c.Risk, c.Profit, c.Time)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedAntiques 写入古董数据
func seedAntiques(db sqlExecutor) error {
	antiques := []struct {
		ID       int
		Name     string
		Price    int
		Material int
		Img      string
		Desc     string
		Level    int
	}{
		{1, "青花瓷盘", 5800, 2, "/images/antiqueinfo/01-blue-white-porcelain-plate.webp", "明代青花瓷盘，釉色均匀，图案细腻，保存较完整。", 0},
		{2, "宋代龙泉窑碗", 24000, 5, "/images/antiqueinfo/02-song-longquan-celadon-bowl.webp", "宋代龙泉窑作品，釉面莹润，造型古朴，器型标准。", 1},
		{3, "清代鼻烟壶", 9200, 3, "/images/antiqueinfo/03-qing-snuff-bottle.webp", "清代早期手工雕刻鼻烟壶，图案精致，具有较高艺术价值。", 0},
		{4, "战国青铜剑", 148000, 7, "/images/antiqueinfo/04-warring-states-bronze-sword.webp", "战国时期铸造青铜剑，剑身刻有铭文，品相极佳，历史价值极高。", 2},
		{5, "唐三彩骏马", 395000, 8, "/images/antiqueinfo/05-tang-sancai-horse.webp", "唐代三彩陶马，釉色鲜明，造型生动，是唐代陶艺的代表之作。", 3},
		{6, "汉代玉璧", 630000, 9, "/images/antiqueinfo/06-han-jade-bi.webp", "汉代和田玉雕刻玉璧，温润透亮，制作精湛，为贵族礼仪之用。", 4},
		{7, "宋代青瓷瓶", 15800, 4, "/images/antiqueinfo/07-song-celadon-vase.webp", "宋代龙泉窑青瓷瓶，胎质细腻，釉色温润，造型古朴典雅。", 1},
		{8, "清代粉彩碗", 32000, 5, "/images/antiqueinfo/08-qing-famille-rose-bowl.webp", "清乾隆时期粉彩器皿，色彩丰富，画工精美，品相完好。", 1},
		{9, "元代青花瓷罐", 110000, 7, "/images/antiqueinfo/09-yuan-blue-white-porcelain-jar.webp", "元代青花罐，胎体厚重，画面生动，是元青花的重要代表。", 2},
		{10, "战国编钟残片", 9700, 3, "/images/antiqueinfo/10-warring-states-bell-fragment.webp", "战国时期编钟碎片，铭文清晰，具有重要考古研究价值。", 0},
		{11, "清代珐琅彩盘", 420000, 9, "/images/antiqueinfo/11-qing-enamel-plate.webp", "清代宫廷珐琅彩盘，纹饰华丽，器形规整，极为罕见。", 3},
		{12, "明代青花龙纹罐", 750000, 10, "/images/antiqueinfo/12-ming-blue-white-dragon-jar.webp", "明永乐时期龙纹青花罐，图案威武，烧造技术登峰造极。", 4},
		{13, "元代青花鱼藻纹碗", 198000, 7, "/images/antiqueinfo/13-yuan-blue-white-fish-algae-bowl.webp", "元青花鱼藻纹碗，主题鲜明，象征吉祥，工艺水平高超。", 2},
		{14, "春秋时期青铜钺", 465000, 8, "/images/antiqueinfo/14-spring-autumn-bronze-yue.webp", "春秋时期青铜兵器，铸造精良，是权力和威严的象征。", 3},
		{15, "清宫御用玉带扣", 990000, 10, "/images/antiqueinfo/15-qing-imperial-jade-belt-buckle.webp", "清代宫廷御用白玉带扣，雕工繁复，玉质温润，存世稀少。", 4},
		{16, "宋代汝窑杯", 88000, 7, "/images/antiqueinfo/16-song-ru-ware-cup.webp", "宋代汝窑青瓷杯，釉色温润，纹饰独特，极具收藏价值。", 2},
		{17, "唐代金饰", 250000, 8, "/images/antiqueinfo/17-tang-gold-ornament.webp", "唐代金饰，工艺精美，金质纯净，适合高端收藏。", 3},
		{18, "明代铜香炉", 35000, 6, "/images/antiqueinfo/18-ming-bronze-incense-burner.webp", "明代铜香炉，造型古朴，表面雕刻精美，具有独特的艺术价值。", 1},
		{19, "清代大元宝", 400000, 9, "/images/antiqueinfo/19-qing-sycee.webp", "清代光绪年间大元宝，历史悠久，铸工精良，市场稀缺。", 3},
		{20, "清代宫廷木雕屏风", 750000, 9, "/images/antiqueinfo/20-qing-palace-wood-screen.webp", "清代宫廷木雕屏风，雕工繁复，木质坚硬，具有较高艺术和历史价值。", 4},
		{21, "明代红釉瓷器", 98000, 6, "/images/antiqueinfo/21-ming-red-glazed-porcelain.webp", "明代红釉瓷器，色泽鲜艳，保存完好，具备较高的艺术收藏价值。", 2},
		{22, "战国铜车马", 550000, 10, "/images/antiqueinfo/22-warring-states-bronze-chariot.webp", "战国时期青铜车马，工艺复杂，保存完好，是重要的历史遗物。", 4},
		{23, "清代象牙雕刻", 320000, 8, "/images/antiqueinfo/23-qing-ivory-carving.webp", "清代象牙雕刻工艺品，雕刻细致，图案生动，象征富贵吉祥。", 3},
		{24, "唐代铜钱", 7500, 4, "/images/antiqueinfo/24-tang-copper-coin.webp", "唐代铜钱，保存完好，具有一定的历史价值，适合初学者收藏。", 0},
		{25, "明代白玉杯", 150000, 8, "/images/antiqueinfo/25-ming-white-jade-cup.webp", "明代白玉杯，玉质细腻，器型优雅，极具高端收藏价值。", 3},
		{26, "元代青花大盘", 120000, 7, "/images/antiqueinfo/26-yuan-blue-white-large-plate.webp", "元代青花大盘，画面生动，艺术价值极高，制作工艺精湛。", 2},
		{27, "清代珐琅铜壶", 220000, 7, "/images/antiqueinfo/27-qing-enamel-copper-pot.webp", "清代珐琅铜壶，彩绘精美，工艺绝佳，收藏价值极高。", 3},
		{28, "明代五福鼎", 510000, 9, "/images/antiqueinfo/28-ming-five-blessings-cauldron.webp", "明代五福鼎，铜器雕刻精细，五福象征，历史文化价值极高。", 4},
		{29, "宋代釉里红花瓶", 105000, 7, "/images/antiqueinfo/29-song-underglaze-red-vase.webp", "宋代釉里红花瓶，工艺精美，色彩鲜艳，是宋代陶瓷的典型代表。", 2},
		{30, "清代金箔大盘", 850000, 10, "/images/antiqueinfo/30-qing-gold-foil-large-plate.webp", "清代金箔大盘，金属材质高贵，工艺精细，极为罕见，适合私人收藏。", 4},
		{31, "战国玉佩", 19000, 6, "/images/antiqueinfo/31-warring-states-jade-pendant.webp", "战国时期和田玉玉佩，纹路细致，玉质温润，存世稀有。", 1},
		{32, "清代翡翠挂坠", 75000, 7, "/images/antiqueinfo/32-qing-jadeite-pendant.webp", "清代翡翠挂坠，雕刻精美，玉质翠绿，适合收藏投资。", 2},
		{33, "元代铜鼓", 200000, 8, "/images/antiqueinfo/33-yuan-bronze-drum.webp", "元代铜鼓，音质清晰，外形美观，具有较高的历史与文化价值。", 3},
		{34, "明代紫砂壶", 128000, 6, "/images/antiqueinfo/34-ming-zisha-teapot.webp", "明代紫砂壶，壶体沉稳，手工雕刻精美，茶文化爱好者必备。", 2},
		{35, "清代缂丝屏风", 650000, 9, "/images/antiqueinfo/35-qing-kesi-screen.webp", "清代宫廷缂丝屏风，丝线细致，工艺精湛，贵族气息浓厚。", 4},
		{36, "唐代金铜佛像", 90000, 8, "/images/antiqueinfo/36-tang-gilt-bronze-buddha.webp", "唐代金铜佛像，庄重典雅，铜质精美，具备很高的宗教和历史价值。", 3},
		{37, "宋代青瓷碗", 95000, 7, "/images/antiqueinfo/37-song-celadon-bowl.webp", "宋代青瓷碗，青釉质感温润，工艺精致，是宋代陶瓷的典型之作。", 2},
		{38, "明代桃花心木柜", 220000, 9, "/images/antiqueinfo/38-ming-mahogany-cabinet.webp", "明代桃花心木柜，木质坚实，雕工精美，具有极高的家具收藏价值。", 3},
		{39, "清代螺钿屏风", 540000, 10, "/images/antiqueinfo/39-qing-mother-of-pearl-screen.webp", "清代螺钿屏风，使用珍贵螺钿装饰，工艺精细，极具收藏价值。", 4},
		{40, "元代瓷器大盘", 150000, 8, "/images/antiqueinfo/40-yuan-porcelain-large-plate.webp", "元代青花瓷器大盘，造型简洁，图案生动，兼具实用与艺术价值。", 3},
		{41, "清代大理石碑文", 350000, 9, "/images/antiqueinfo/41-qing-marble-inscription.webp", "清代大理石碑文，刻有历史重要事件，文笔优雅，收藏价值极高。", 4},
		{42, "清代玉雕花瓶", 180000, 8, "/images/antiqueinfo/42-qing-jade-vase.webp", "清代玉雕花瓶，玉质温润，雕工精细，花纹自然，适合高端收藏。", 3},
		{43, "战国王室青铜剑", 520000, 10, "/images/antiqueinfo/43-warring-states-royal-bronze-sword.webp", "战国王室规格青铜剑，保存完好，历史意义重大，堪称国家级文物。", 4},
		{44, "明代龙纹瓷碗", 95000, 7, "/images/antiqueinfo/44-ming-dragon-porcelain-bowl.webp", "明代龙纹瓷碗，龙纹精美，釉色鲜亮，瓷质细腻，是收藏爱好者的心头好。", 2},
		{45, "元代青花盘", 145000, 8, "/images/antiqueinfo/45-yuan-blue-white-plate.webp", "元代青花盘，图案简洁，瓷质坚硬，具有很高的艺术和文化价值。", 3},
		{46, "清代金箔屏风", 780000, 10, "/images/antiqueinfo/46-qing-gold-foil-screen.webp", "清代金箔屏风，金色光泽夺目，镶嵌金箔，工艺精湛，具有极高的收藏价值。", 4},
		{47, "战国玉带钩", 34000, 6, "/images/antiqueinfo/47-warring-states-jade-belt-hook.webp", "战国时期玉带钩，玉质洁净，雕工精巧，是战国文化的代表性饰品。", 1},
		{48, "宋代青花瓷碗", 15000, 6, "/images/antiqueinfo/48-song-blue-white-porcelain-bowl.webp", "宋代青花瓷碗，青花纹饰古朴，釉面光滑，保存完好，适合新手收藏。", 0},
		{49, "唐代铜镜", 30000, 5, "/images/antiqueinfo/49-tang-bronze-mirror.webp", "唐代铜镜，镜面清晰，外观精致，适合喜欢铜器收藏的爱好者。", 1},
		{50, "明代黄花梨书柜", 550000, 10, "/images/antiqueinfo/50-ming-huanghuali-bookcase.webp", "明代黄花梨书柜，黄花梨木纹美，工艺精湛，是高端家具收藏中的极品。", 4},
	}

	for _, a := range antiques {
		_, err := db.Exec(`
			INSERT INTO antiques (id, name, price, material, img, desc, level)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, material=excluded.material,
				img=excluded.img, desc=excluded.desc, level=excluded.level
		`, a.ID, a.Name, a.Price, a.Material, a.Img, a.Desc, a.Level)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedStocks 写入股票数据
func seedStocks(db sqlExecutor) error {
	stocks := []struct {
		ID    int
		Name  string
		Price int
		Risk  int
	}{
		{1, "星辰科技", 285, 5},
		{2, "安康医疗", 98, 2},
		{3, "中联电网", 150, 1},
		{4, "极光生物", 78, 4},
		{5, "流光娱乐", 65, 3},
		{6, "震元重工", 125, 2},
		{7, "九天能源", 92, 3},
		{8, "未来通链", 45, 5},
		{9, "太和地产", 110, 4},
		{10, "凌峰数码", 205, 3},
		{11, "东海物流", 88, 2},
		{12, "云杉教育", 56, 1},
		{13, "万象百货", 138, 2},
		{14, "火种引擎", 302, 4},
		{15, "灵脉互娱", 73, 5},
	}

	for _, s := range stocks {
		_, err := db.Exec(`
			INSERT INTO stocks (id, name, price, risk)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, risk=excluded.risk
		`, s.ID, s.Name, s.Price, s.Risk)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedStockNews 写入股票新闻数据
func seedStockNews(db sqlExecutor) error {
	news := []string{
		"XXXX发布新款智能手表，市场反应热烈，股价上升",
		"XXXX季度营收稳步增长，股价上涨",
		"XXXX宣布进军新能源市场，股价上涨",
		"XXXX收购初创公司扩展业务，股价上扬",
		"XXXX电动车销量持续增长，股价上涨",
		"XXXX与国际品牌合作，市场看好，股价上升",
		"XXXX成功打开海外市场，股价上涨",
		"XXXX云服务用户数增加，股价上扬",
		"XXXX完成供应链优化，利润增长，股价上涨",
		"XXXX宣布新品发布计划，市场预期良好，股价上涨",
		"XXXX签订长期供货协议，股价稳步上升",
		"XXXX物流系统升级成功，成本下降，股价上涨",
		"XXXX获批新技术专利，前景被看好，股价上涨",
		"XXXX品牌影响力提升，消费者信心增强，股价上扬",
		"XXXX市场份额提升，股价上涨",
		"XXXX与高校合作推动科研进展，股价上升",
		"XXXX加大研发投入，技术领先，股价上涨",
		"XXXX用户满意度提升，客户留存率增长，股价上涨",
		"XXXX新产品获奖，品牌声誉增强，股价上扬",
		"XXXX通过绿色认证，获得环保基金关注，股价上涨",
		"XXXX公司债务水平下降，财务状况改善，股价反弹",
		"XXXX官网访问量创新高，投资者信心增强，股价上涨",
		"XXXX股东大会通过新增长战略，市场看好，股价上涨",
		"XXXX与政府达成补贴协议，业务利好，股价上扬",
		"XXXX推出节能技术方案，环保资本关注，股价上涨",
		"XXXX季度利润超预期，股价大涨",
		"XXXX推出AI芯片产品，市场追捧，股价飙升",
		"XXXX研发新药成功通过三期临床，股价飙升",
		"XXXX获得巨额订单，营收激增，股价大涨",
		"XXXX中标国家重大工程项目，股价飙升",
		"XXXX被纳入国际指数基金，资本流入推动股价大涨",
		"XXXX与科技巨头深度合作，股价飙升",
		"XXXX财报净利润翻倍，投资者热情高涨，股价大涨",
		"XXXX新能源车月销破纪录，股价飙升",
		"XXXX连续四季度强劲增长，股价大涨",
		"XXXX技术突破颠覆行业，资本看好，股价飙升",
		"XXXX宣布股票回购计划，股价大涨",
		"XXXX高管大规模增持，释放强烈信号，股价飙升",
		"XXXX收购海外巨头公司，扩张成功，股价大涨",
		"XXXX市值突破千亿关口，投资者热情高涨，股价飙升",
		"XXXXIPO表现强劲，市场追捧，股价大涨",
		"XXXX成为独角兽企业，市值重估，股价飙升",
		"XXXX子公司上市消息公布，市场激动，股价大涨",
		"XXXX行业地位稳固，分析师一致看多，股价飙升",
		"XXXX中期业绩惊艳市场，股价大涨",
		"XXXX新能源产品获国际奖项，品牌声望高涨，股价飙升",
		"XXXX签署百亿合同，业务井喷，股价大涨",
		"XXXX全球扩张成功，市占率翻倍，股价飙升",
		"XXXX全年营收破纪录，股价大涨",
		"XXXX获头部机构评级上调，股价飙升",
		"XXXX营收不及预期，股价下跌",
		"XXXX供应链中断影响生产，股价下降",
		"XXXX遭遇管理层变动，市场不安，股价下挫",
		"XXXX季报盈利下降，股价下跌",
		"XXXX调低全年预期，投资者信心动摇，股价下跌",
		"XXXX原材料价格上涨压缩利润，股价下挫",
		"XXXX在华市场销售乏力，股价下降",
		"XXXX遭遇产品召回事件，消费者信心下降，股价下跌",
		"XXXX财报显示成本上涨，利润受压，股价回调",
		"XXXX因政策变动受限，业务缩水，股价下跌",
		"XXXX被竞争对手超越，市场份额缩小，股价下跌",
		"XXXX主力产品销售下降，市场反应平淡，股价下挫",
		"XXXX遭遇专利纠纷，法律风险增加，股价下跌",
		"XXXX合作方违约，项目受阻，股价下降",
		"XXXX外部环境不利影响营收，股价下跌",
		"XXXX高管离职引发市场担忧，股价下跌",
		"XXXX遭遇舆论危机，声誉受损，股价下跌",
		"XXXX信用评级被下调，融资成本上升，股价下降",
		"XXXX新产品销售表现平淡，股价下跌",
		"XXXX遭遇黑客攻击，数据泄露，股价下跌",
		"XXXX订单量不及预期，股价回调",
		"XXXX投资失败导致亏损，股价下跌",
		"XXXX市场扩张未达标，股价下跌",
		"XXXX国际业务增长放缓，股价下降",
		"XXXX分析师下调评级，股价下挫",
		"XXXX被曝财务造假，股价暴跌",
		"XXXX遭遇重大安全事故，声誉受损，股价暴跌",
		"XXXX高管被调查，市场恐慌，股价暴跌",
		"XXXX季度亏损巨大，投资者抛售，股价暴跌",
		"XXXX核心技术被取代，竞争力削弱，股价暴跌",
		"XXXX股东大量减持，市场信心崩溃，股价暴跌",
		"XXXX债务违约风险上升，股价暴跌",
		"XXXX卷入商业贿赂案，股价暴跌",
		"XXXX出口业务被封锁，营收锐减，股价暴跌",
		"XXXX遭遇系统性欺诈指控，股价暴跌",
		"XXXX突发重大丑闻，公众抵制，股价暴跌",
		"XXXX涉嫌内幕交易被调查，股价暴跌",
		"XXXX技术项目失败，重大损失，股价暴跌",
		"XXXX资金链断裂传闻四起，市场恐慌，股价暴跌",
		"XXXX因业绩造假遭停牌，股价暴跌",
		"XXXX丧失关键客户订单，业务重创，股价暴跌",
		"XXXX面临大规模诉讼，法律风险剧增，股价暴跌",
		"XXXX子公司破产清算，投资者信心全无，股价暴跌",
		"XXXX突遭国际制裁，市场剧烈反应，股价暴跌",
		"XXXX遭遇重大黑天鹅事件，股价暴跌",
	}

	for i, n := range news {
		_, err := db.Exec(`
			INSERT INTO stock_news (id, content)
			VALUES (?, ?)
			ON CONFLICT(id) DO UPDATE SET content=excluded.content
		`, i+1, n)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedHouses 写入房产数据
func seedHouses(db sqlExecutor) error {
	houses := []struct {
		ID       int
		Name     string
		Price    int
		PriceMax int
		PriceMin int
		Health   int
		Fame     int
		Img      string
	}{
		{1, "单身公寓", 300000, 450000, 200000, 2, 1, "/images/houseinfo/houses/01-studio-apartment.webp"},
		{2, "一室一厅", 500000, 700000, 350000, 3, 2, "/images/houseinfo/houses/02-one-bedroom-apartment.webp"},
		{3, "两室一厅", 800000, 1200000, 600000, 5, 4, "/images/houseinfo/houses/03-two-bedroom-apartment.webp"},
		{4, "三室两厅", 1500000, 2200000, 1000000, 8, 7, "/images/houseinfo/houses/04-three-bedroom-apartment.webp"},
		{5, "四室两厅", 2500000, 3500000, 1800000, 10, 12, "/images/houseinfo/houses/05-four-bedroom-apartment.webp"},
		{6, "复式公寓", 2000000, 3000000, 1500000, 9, 10, "/images/houseinfo/houses/06-duplex-apartment.webp"},
		{7, "顶层复式", 3500000, 5000000, 2500000, 15, 18, "/images/houseinfo/houses/07-penthouse-duplex.webp"},
		{8, "花园洋房", 4500000, 6500000, 3000000, 18, 22, "/images/houseinfo/houses/08-garden-residence.webp"},
		{9, "联排别墅", 5500000, 8000000, 4000000, 20, 28, "/images/houseinfo/houses/09-townhouse.webp"},
		{10, "独栋别墅", 10000000, 15000000, 7000000, 28, 38, "/images/houseinfo/houses/10-detached-villa.webp"},
		{11, "精装公寓", 600000, 900000, 400000, 4, 3, "/images/houseinfo/houses/11-furnished-apartment.webp"},
		{12, "学区房小户型", 1200000, 1800000, 900000, 6, 8, "/images/houseinfo/houses/12-small-school-district-home.webp"},
		{13, "学区房大户型", 2500000, 3800000, 1800000, 10, 14, "/images/houseinfo/houses/13-large-school-district-home.webp"},
		{14, "市中心公寓", 1800000, 2600000, 1300000, 7, 10, "/images/houseinfo/houses/14-downtown-apartment.webp"},
		{15, "江景房", 3000000, 4500000, 2200000, 14, 20, "/images/houseinfo/houses/15-river-view-home.webp"},
		{16, "湖景房", 2800000, 4000000, 2000000, 13, 18, "/images/houseinfo/houses/16-lake-view-home.webp"},
		{17, "山景别墅", 6000000, 9000000, 4500000, 22, 32, "/images/houseinfo/houses/17-mountain-villa.webp"},
		{18, "海景别墅", 12000000, 18000000, 8000000, 30, 45, "/images/houseinfo/houses/18-seaside-villa.webp"},
		{19, "豪华平层", 4000000, 6000000, 3000000, 16, 25, "/images/houseinfo/houses/19-luxury-flat.webp"},
		{20, "大平层", 5000000, 7500000, 3500000, 18, 28, "/images/houseinfo/houses/20-large-flat.webp"},
		{21, "商务公寓", 700000, 1000000, 500000, 4, 4, "/images/houseinfo/houses/21-business-apartment.webp"},
		{22, "酒店式公寓", 900000, 1300000, 650000, 5, 5, "/images/houseinfo/houses/22-serviced-apartment.webp"},
		{23, "SOHO公寓", 550000, 800000, 380000, 4, 3, "/images/houseinfo/houses/23-soho-apartment.webp"},
		{24, "LOFT公寓", 650000, 950000, 450000, 4, 4, "/images/houseinfo/houses/24-loft-apartment.webp"},
		{25, "青年社区", 450000, 650000, 300000, 3, 2, "/images/houseinfo/houses/25-youth-community.webp"},
		{26, "老年公寓", 800000, 1200000, 550000, 6, 3, "/images/houseinfo/houses/26-senior-apartment.webp"},
		{27, "养老社区", 2000000, 3000000, 1500000, 10, 8, "/images/houseinfo/houses/27-retirement-community.webp"},
		{28, "经济适用房", 350000, 500000, 250000, 2, 1, "/images/houseinfo/houses/28-affordable-housing.webp"},
		{29, "限价房", 600000, 850000, 400000, 4, 2, "/images/houseinfo/houses/29-price-capped-housing.webp"},
		{30, "公租房", 200000, 300000, 150000, 1, 0, "/images/houseinfo/houses/30-public-rental-housing.webp"},
		{31, "保障性住房", 400000, 600000, 280000, 3, 1, "/images/houseinfo/houses/31-social-housing.webp"},
		{32, "人才公寓", 500000, 750000, 350000, 3, 3, "/images/houseinfo/houses/32-talent-apartment.webp"},
		{33, "高端豪宅", 8000000, 12000000, 5500000, 25, 42, "/images/houseinfo/houses/33-luxury-mansion.webp"},
		{34, "庄园别墅", 20000000, 30000000, 15000000, 38, 60, "/images/houseinfo/houses/34-manor-villa.webp"},
		{35, "城堡别墅", 50000000, 80000000, 35000000, 50, 100, "/images/houseinfo/houses/35-castle-villa.webp"},
		{36, "中式合院", 7000000, 10000000, 5000000, 24, 38, "/images/houseinfo/houses/36-chinese-courtyard-residence.webp"},
		{37, "四合院", 15000000, 25000000, 10000000, 35, 55, "/images/houseinfo/houses/37-traditional-courtyard.webp"},
		{38, "欧式别墅", 6000000, 9000000, 4500000, 22, 35, "/images/houseinfo/houses/38-european-villa.webp"},
		{39, "美式别墅", 5500000, 8000000, 4000000, 20, 32, "/images/houseinfo/houses/39-american-villa.webp"},
		{40, "现代简约别墅", 4500000, 6500000, 3200000, 17, 26, "/images/houseinfo/houses/40-modern-minimalist-villa.webp"},
		{41, "北欧风情别墅", 5000000, 7500000, 3800000, 19, 30, "/images/houseinfo/houses/41-nordic-villa.webp"},
		{42, "东南亚风情别墅", 4800000, 7000000, 3600000, 18, 28, "/images/houseinfo/houses/42-southeast-asian-villa.webp"},
		{43, "日式庭院", 4000000, 5800000, 3000000, 16, 25, "/images/houseinfo/houses/43-japanese-courtyard.webp"},
		{44, "法式城堡", 18000000, 25000000, 12000000, 36, 58, "/images/houseinfo/houses/44-french-chateau.webp"},
		{45, "英式庄园", 12000000, 18000000, 8500000, 30, 48, "/images/houseinfo/houses/45-english-manor.webp"},
		{46, "意大利式别墅", 6500000, 9500000, 5000000, 23, 35, "/images/houseinfo/houses/46-italian-villa.webp"},
		{47, "西班牙式别墅", 5800000, 8500000, 4500000, 21, 32, "/images/houseinfo/houses/47-spanish-villa.webp"},
		{48, "高尔夫别墅", 9000000, 13000000, 6500000, 27, 44, "/images/houseinfo/houses/48-golf-villa.webp"},
		{49, "温泉别墅", 7500000, 11000000, 5500000, 32, 40, "/images/houseinfo/houses/49-hot-spring-villa.webp"},
		{50, "私人岛屿豪宅", 50000000, 100000000, 30000000, 50, 100, "/images/houseinfo/houses/50-private-island-mansion.webp"},
	}

	for _, h := range houses {
		_, err := db.Exec(`
			INSERT INTO houses (id, name, price, price_max, price_min, health, fame, img)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, price_max=excluded.price_max,
				price_min=excluded.price_min, health=excluded.health,
				fame=excluded.fame, img=excluded.img
		`, h.ID, h.Name, h.Price, h.PriceMax, h.PriceMin, h.Health, h.Fame, h.Img)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedCars 写入车辆数据
func seedCars(db sqlExecutor) error {
	cars := []struct {
		ID       int
		Name     string
		Price    int
		PriceMax int
		PriceMin int
		Health   int
		Fame     int
		Img      string
	}{
		{1, "经济型轿车", 80000, 120000, 50000, 3, 2, "/images/carinfo/cars/01-economy-sedan.webp"},
		{2, "紧凑型轿车", 120000, 180000, 80000, 4, 4, "/images/carinfo/cars/02-compact-sedan.webp"},
		{3, "中型轿车", 200000, 300000, 150000, 5, 7, "/images/carinfo/cars/03-midsize-sedan.webp"},
		{4, "中大型轿车", 350000, 500000, 250000, 7, 12, "/images/carinfo/cars/04-executive-sedan.webp"},
		{5, "豪华轿车", 600000, 900000, 450000, 10, 22, "/images/carinfo/cars/05-luxury-sedan.webp"},
		{6, "紧凑型SUV", 150000, 220000, 100000, 5, 6, "/images/carinfo/cars/06-compact-suv.webp"},
		{7, "中型SUV", 250000, 380000, 180000, 6, 10, "/images/carinfo/cars/07-midsize-suv.webp"},
		{8, "中大型SUV", 400000, 600000, 300000, 8, 15, "/images/carinfo/cars/08-mid-large-suv.webp"},
		{9, "全尺寸SUV", 550000, 800000, 400000, 9, 20, "/images/carinfo/cars/09-full-size-suv.webp"},
		{10, "豪华SUV", 800000, 1200000, 600000, 12, 30, "/images/carinfo/cars/10-luxury-suv.webp"},
		{11, "跑车", 1000000, 1500000, 700000, 14, 40, "/images/carinfo/cars/11-sports-car.webp"},
		{12, "超级跑车", 3000000, 5000000, 2000000, 18, 80, "/images/carinfo/cars/12-supercar.webp"},
		{13, "商务MPV", 300000, 450000, 220000, 6, 10, "/images/carinfo/cars/13-business-mpv.webp"},
		{14, "家用MPV", 200000, 300000, 150000, 5, 6, "/images/carinfo/cars/14-family-mpv.webp"},
		{15, "豪华MPV", 500000, 750000, 380000, 9, 18, "/images/carinfo/cars/15-luxury-mpv.webp"},
		{16, "纯电动车", 250000, 400000, 180000, 5, 8, "/images/carinfo/cars/16-battery-electric-car.webp"},
		{17, "增程式电动车", 280000, 420000, 200000, 6, 9, "/images/carinfo/cars/17-range-extended-ev.webp"},
		{18, "插电混动", 220000, 350000, 160000, 5, 7, "/images/carinfo/cars/18-plug-in-hybrid.webp"},
		{19, "油电混动", 240000, 360000, 170000, 5, 8, "/images/carinfo/cars/19-hybrid-car.webp"},
		{20, "氢燃料电池车", 400000, 600000, 300000, 7, 15, "/images/carinfo/cars/20-hydrogen-fuel-cell-car.webp"},
		{21, "微型车", 50000, 80000, 35000, 2, 1, "/images/carinfo/cars/21-microcar.webp"},
		{22, "小型车", 70000, 100000, 50000, 2, 2, "/images/carinfo/cars/22-small-car.webp"},
		{23, "两厢车", 90000, 130000, 65000, 3, 3, "/images/carinfo/cars/23-hatchback.webp"},
		{24, "掀背车", 110000, 160000, 80000, 3, 4, "/images/carinfo/cars/24-liftback.webp"},
		{25, "旅行车", 180000, 280000, 130000, 5, 8, "/images/carinfo/cars/25-station-wagon.webp"},
		{26, "跨界车", 160000, 240000, 120000, 4, 6, "/images/carinfo/cars/26-crossover.webp"},
		{27, "硬顶敞篷", 450000, 650000, 320000, 8, 18, "/images/carinfo/cars/27-hardtop-convertible.webp"},
		{28, "软顶敞篷", 380000, 550000, 280000, 7, 15, "/images/carinfo/cars/28-soft-top-convertible.webp"},
		{29, "皮卡", 150000, 220000, 100000, 5, 5, "/images/carinfo/cars/29-pickup-truck.webp"},
		{30, "轻型货车", 120000, 180000, 80000, 4, 2, "/images/carinfo/cars/30-light-truck.webp"},
		{31, "面包车", 60000, 90000, 40000, 2, 1, "/images/carinfo/cars/31-passenger-van.webp"},
		{32, "厢式货车", 100000, 150000, 70000, 3, 1, "/images/carinfo/cars/32-cargo-van.webp"},
		{33, "房车", 800000, 1200000, 500000, 12, 28, "/images/carinfo/cars/33-motorhome.webp"},
		{34, "豪华房车", 2000000, 3500000, 1500000, 16, 50, "/images/carinfo/cars/34-luxury-motorhome.webp"},
		{35, "校车", 250000, 380000, 180000, 5, 2, "/images/carinfo/cars/35-school-bus.webp"},
		{36, "公交车", 400000, 600000, 280000, 6, 3, "/images/carinfo/cars/36-city-bus.webp"},
		{37, "观光车", 150000, 220000, 100000, 4, 4, "/images/carinfo/cars/37-sightseeing-cart.webp"},
		{38, "高尔夫球车", 80000, 120000, 50000, 2, 3, "/images/carinfo/cars/38-golf-cart.webp"},
		{39, "三轮摩托车", 30000, 50000, 20000, 1, 1, "/images/carinfo/cars/39-three-wheel-motorcycle.webp"},
		{40, "摩托车", 50000, 80000, 30000, 2, 2, "/images/carinfo/cars/40-motorcycle.webp"},
		{41, "电动摩托车", 15000, 25000, 10000, 1, 1, "/images/carinfo/cars/41-electric-motorcycle.webp"},
		{42, "电动自行车", 3000, 5000, 2000, 1, 0, "/images/carinfo/cars/42-electric-bicycle.webp"},
		{43, "自行车", 1000, 2000, 500, 1, 0, "/images/carinfo/cars/43-bicycle.webp"},
		{44, "山地自行车", 3000, 6000, 1800, 1, 0, "/images/carinfo/cars/44-mountain-bike.webp"},
		{45, "公路自行车", 5000, 10000, 3000, 1, 1, "/images/carinfo/cars/45-road-bicycle.webp"},
		{46, "折叠自行车", 2000, 4000, 1200, 1, 0, "/images/carinfo/cars/46-folding-bicycle.webp"},
		{47, "电动滑板车", 2000, 3500, 1200, 1, 0, "/images/carinfo/cars/47-electric-scooter.webp"},
		{48, "平衡车", 1500, 3000, 800, 1, 0, "/images/carinfo/cars/48-self-balancing-scooter.webp"},
		{49, "老爷车", 500000, 800000, 350000, 8, 25, "/images/carinfo/cars/49-classic-car.webp"},
		{50, "收藏级古董车", 2000000, 5000000, 1200000, 15, 60, "/images/carinfo/cars/50-collector-antique-car.webp"},
	}

	for _, c := range cars {
		_, err := db.Exec(`
			INSERT INTO cars (id, name, price, price_max, price_min, health, fame, img)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, price=excluded.price, price_max=excluded.price_max,
				price_min=excluded.price_min, health=excluded.health,
				fame=excluded.fame, img=excluded.img
		`, c.ID, c.Name, c.Price, c.PriceMax, c.PriceMin, c.Health, c.Fame, c.Img)
		if err != nil {
			return err
		}
	}

	return nil
}

// datingMeetScenes 只配置需要玩家主动前往场景才能认识的对象。
// 男性对象沿用对应女性对象的条件和场景，保证两侧难度一致。
var datingMeetScenes = map[int]string{
	1:  "公园",
	7:  "健身房",
	16: "校园",
	20: "医院",
	30: "游乐园",
	40: "诊所",
	47: "剧院",
}

// seedDatingInfo 直接写入当前版本男女各 50 位约会对象。
func seedDatingInfo(db sqlExecutor) error {
	datings := []struct {
		ID             int
		Age            int
		Occupation     string
		Cost           int
		MeetConditions string
		Gifts          string
		Locations      string
	}{
		{1, 25, "大学生", 500, `[{"ctype":"age","cvalue":18},{"ctype":"random","cvalue":40}]`, `["鲜花","巧克力","玩偶"]`, `["公园","电影院","咖啡厅"]`},
		{2, 28, "职场白领", 1500, `[{"ctype":"bank","cvalue":50000},{"ctype":"work_count","cvalue":10}]`, `["香水","红酒","首饰"]`, `["西餐厅","艺术馆","音乐厅"]`},
		{3, 30, "名媛千金", 9200, `[{"ctype":"car","cvalue":8},{"ctype":"bank","cvalue":8000000}]`, `["钻石","名表","礼服"]`, `["私人会所","游艇","豪华酒店"]`},
		{4, 32, "上流名媛", 10000, `[{"ctype":"house","cvalue":20},{"ctype":"cash","cvalue":10000000}]`, `["豪宅","私人飞机","钻戒"]`, `["私人庄园","海岛","度假村"]`},
		{5, 26, "图书管理员", 1400, `[{"ctype":"fame","cvalue":20},{"ctype":"play_game","cvalue":5}]`, `["书籍","钢笔","书签"]`, `["图书馆","书店","茶馆"]`},
		{6, 25, "游戏主播", 1500, `[{"ctype":"play_game","cvalue":10},{"ctype":"win_game","cvalue":3}]`, `["游戏机","耳机","外设"]`, `["电竞馆","游戏展","咖啡厅"]`},
		{7, 26, "健身教练", 700, `[{"ctype":"immunity","cvalue":90},{"ctype":"random","cvalue":50}]`, `["运动装","蛋白粉","手环"]`, `["健身房","海边","户外"]`},
		{8, 28, "艺术家", 6500, `[{"ctype":"antique_rare","cvalue":1},{"ctype":"fame","cvalue":70}]`, `["画具","艺术品","画册"]`, `["画廊","展览","艺术区"]`},
		{9, 26, "网红主播", 1800, `[{"ctype":"fame","cvalue":60},{"ctype":"lottery_win","cvalue":1}]`, `["手机","化妆品","名牌包"]`, `["商场","网红地","直播间"]`},
		{10, 30, "医生", 1800, `[{"ctype":"immunity","cvalue":90},{"ctype":"fame","cvalue":20}]`, `["医学书","养生品","器械"]`, `["医院花园","温泉","餐厅"]`},
		{11, 29, "律师", 1900, `[{"ctype":"fame","cvalue":40},{"ctype":"car","cvalue":1}]`, `["钢笔","公文包","手表"]`, `["律所","商务餐厅","会所"]`},
		{12, 27, "作家", 8800, `[{"ctype":"fame","cvalue":85},{"ctype":"bank","cvalue":500000}]`, `["书籍","钢笔","咖啡"]`, `["书房","咖啡馆","沙龙"]`},
		{13, 26, "歌手", 1700, `[{"ctype":"fame","cvalue":40},{"ctype":"play_game","cvalue":5}]`, `["麦克风","吉他","乐谱"]`, `["KTV","音乐节","录音棚"]`},
		{14, 28, "演员", 6000, `[{"ctype":"fame","cvalue":75},{"ctype":"win_game","cvalue":5}]`, `["剧本","戏服","化妆品"]`, `["剧院","片场","首映礼"]`},
		{15, 31, "企业家", 3500, `[{"ctype":"company_founder","cvalue":1},{"ctype":"bank","cvalue":200000}]`, `["名表","商业书","礼盒"]`, `["高尔夫","酒会","会所"]`},
		{16, 25, "博士后", 700, `[{"ctype":"random","cvalue":60},{"ctype":"age","cvalue":26}]`, `["奶茶","零食","玩偶"]`, `["校园","游乐园","奶茶店"]`},
		{17, 30, "舞蹈老师", 5500, `[{"ctype":"immunity","cvalue":85},{"ctype":"fame","cvalue":60}]`, `["舞鞋","舞衣","音乐盒"]`, `["舞蹈室","剧院","舞台"]`},
		{18, 27, "消防员", 3000, `[{"ctype":"immunity","cvalue":80},{"ctype":"work_count","cvalue":60}]`, `["装备","护具","饮料"]`, `["户外","训练场","烧烤"]`},
		{19, 31, "珠宝设计师", 7000, `[{"ctype":"antique_rare","cvalue":2},{"ctype":"fame","cvalue":75}]`, `["宝石","珠宝","设计工具"]`, `["拍卖会","展览","珠宝店"]`},
		{20, 25, "护士", 900, `[{"ctype":"immunity","cvalue":85},{"ctype":"random","cvalue":40}]`, `["香薰","护手霜","礼盒"]`, `["医院","花园","温泉"]`},
		{21, 28, "钢琴老师", 6000, `[{"ctype":"fame","cvalue":70},{"ctype":"bank","cvalue":150000}]`, `["乐谱","钢琴模型","专辑"]`, `["音乐厅","琴房","沙龙"]`},
		{22, 25, "摄影师", 3000, `[{"ctype":"fame","cvalue":50},{"ctype":"antique_rare","cvalue":1}]`, `["镜头","摄影集","胶片"]`, `["街拍","摄影棚","展览"]`},
		{23, 32, "大学教授", 9000, `[{"ctype":"fame","cvalue":90},{"ctype":"bank","cvalue":800000}]`, `["期刊","书籍","钢笔"]`, `["大学","讲座","图书馆"]`},
		{24, 27, "空姐", 6500, `[{"ctype":"fame","cvalue":70},{"ctype":"cash","cvalue":200000}]`, `["丝巾","香水","行李箱"]`, `["机场","酒店","景点"]`},
		{25, 27, "家政主管", 2800, `[{"ctype":"house","cvalue":1},{"ctype":"work_count","cvalue":100}]`, `["厨具","家电","收纳"]`, `["家居城","厨房","超市"]`},
		{26, 29, "瑜伽教练", 3200, `[{"ctype":"immunity","cvalue":90},{"ctype":"play_game","cvalue":5,"ctarget":"keepfit"}]`, `["瑜伽垫","精油","音乐"]`, `["瑜伽馆","森林","静修中心"]`},
		{27, 26, "赌场荷官", 7200, `[{"ctype":"win_game","cvalue":10,"ctarget":"poker"},{"ctype":"fame","cvalue":60}]`, `["筹码","扑克牌","饰品"]`, `["赌场","VIP厅","夜场"]`},
		{28, 25, "模特", 9500, `[{"ctype":"fame","cvalue":80},{"ctype":"immunity","cvalue":85}]`, `["高跟鞋","时装","香水"]`, `["T台","摄影棚","时尚派对"]`},
		{29, 26, "程序员", 2800, `[{"ctype":"play_game","cvalue":20,"ctarget":"chess"},{"ctype":"work_count","cvalue":70}]`, `["机械键盘","咖啡","数码产品"]`, `["科技园","咖啡厅","展会"]`},
		{30, 25, "幼儿园老师", 800, `[{"ctype":"random","cvalue":50},{"ctype":"age","cvalue":27}]`, `["玩具","绘本","手工"]`, `["幼儿园","游乐园","书店"]`},
		{31, 28, "花艺师", 1800, `[{"ctype":"fame","cvalue":30},{"ctype":"item_own","cvalue":2}]`, `["鲜花","花瓶","园艺工具"]`, `["花店","花园","植物园"]`},
		{32, 27, "市场经理", 1800, `[{"ctype":"fame","cvalue":60},{"ctype":"company_founder","cvalue":1}]`, `["商务包","礼盒","钢笔"]`, `["商务餐厅","展会","高尔夫"]`},
		{33, 30, "厨师", 1800, `[{"ctype":"work_count","cvalue":7},{"ctype":"immunity","cvalue":70}]`, `["厨具","食材","料理书"]`, `["餐厅","美食节","夜市"]`},
		{34, 26, "导游", 1700, `[{"ctype":"fame","cvalue":35},{"ctype":"date_count","cvalue":5}]`, `["背包","相机","纪念品"]`, `["景点","博物馆","街区"]`},
		{35, 29, "导演", 6500, `[{"ctype":"fame","cvalue":70},{"ctype":"win_game","cvalue":5,"ctarget":"poker"}]`, `["剧本","设备","电影周边"]`, `["片场","电影节","影院"]`},
		{36, 25, "服装设计师", 3000, `[{"ctype":"fame","cvalue":50},{"ctype":"antique_rare","cvalue":1}]`, `["布料","设计工具","配饰"]`, `["时装周","工作室","买手店"]`},
		{37, 31, "科学家", 4000, `[{"ctype":"play_game","cvalue":25,"ctarget":"chess"},{"ctype":"bank","cvalue":120000}]`, `["实验设备","书籍","器材"]`, `["实验室","科技馆","大学"]`},
		{38, 25, "舞者", 1700, `[{"ctype":"immunity","cvalue":80},{"ctype":"play_game","cvalue":8,"ctarget":"party"}]`, `["舞鞋","舞衣","音乐"]`, `["舞蹈室","剧院","舞台"]`},
		{39, 25, "马戏团演员", 3000, `[{"ctype":"immunity","cvalue":85},{"ctype":"play_game","cvalue":10}]`, `["表演服","道具","鲜花"]`, `["马戏团","剧场","游乐园"]`},
		{40, 30, "心理医生", 1000, `[{"ctype":"immunity","cvalue":85},{"ctype":"random","cvalue":40}]`, `["心理书籍","香薰","茶"]`, `["诊所","书店","咖啡厅"]`},
		{41, 27, "小提琴家", 6500, `[{"ctype":"fame","cvalue":70},{"ctype":"play_game","cvalue":15,"ctarget":"concert"}]`, `["小提琴","乐谱","音乐专辑"]`, `["音乐厅","歌剧院","音乐节"]`},
		{42, 26, "中医师", 3000, `[{"ctype":"immunity","cvalue":90},{"ctype":"work_count","cvalue":50}]`, `["中药材","茶具","养生品"]`, `["中医馆","茶园","养生馆"]`},
		{43, 29, "金融分析师", 4200, `[{"ctype":"bank","cvalue":180000},{"ctype":"play_game","cvalue":20,"ctarget":"poker"}]`, `["财经书籍","手表","礼盒"]`, `["证券所","商务餐厅","会所"]`},
		{44, 25, "插画师", 3000, `[{"ctype":"fame","cvalue":40},{"ctype":"antique_rare","cvalue":1}]`, `["画笔","颜料","画册"]`, `["市集","画展","艺术区"]`},
		{45, 28, "飞行员", 7000, `[{"ctype":"immunity","cvalue":85},{"ctype":"car","cvalue":1}]`, `["飞行眼镜","模型飞机","装备"]`, `["机场","俱乐部","观景台"]`},
		{46, 27, "电台主持人", 3000, `[{"ctype":"fame","cvalue":45},{"ctype":"play_game","cvalue":10,"ctarget":"concert"}]`, `["耳机","麦克风","专辑"]`, `["电台","录音棚","音乐会"]`},
		{47, 30, "话剧演员", 900, `[{"ctype":"fame","cvalue":30},{"ctype":"random","cvalue":40}]`, `["剧本","戏服","道具"]`, `["剧院","排练厅","舞台"]`},
		{48, 25, "花店老板", 1500, `[{"ctype":"house","cvalue":1},{"ctype":"item_own","cvalue":2}]`, `["鲜花","盆栽","花瓶"]`, `["花店","花市","花园"]`},
		{49, 32, "奥运冠军", 7200, `[{"ctype":"immunity","cvalue":95},{"ctype":"fame","cvalue":90}]`, `["运动装备","奖牌","纪念品"]`, `["体育馆","训练场","赛事"]`},
		{50, 26, "婚礼策划师", 6000, `[{"ctype":"date_count","cvalue":5},{"ctype":"fame","cvalue":50}]`, `["婚礼杂志","鲜花","装饰品"]`, `["婚礼场地","花园","酒店"]`},
	}

	for _, d := range datings {
		for gender := 0; gender <= 1; gender++ {
			id := d.ID + gender*50
			profile, err := currentDatingProfile(id)
			if err != nil {
				return err
			}
			genderDirectory := "female"
			occupation := d.Occupation
			gifts := d.Gifts
			if gender == 1 {
				genderDirectory = "male"
				occupation = maleDatingTextReplacer.Replace(occupation)
				gifts = maleDatingTextReplacer.Replace(gifts)
			}
			image := fmt.Sprintf("/images/datinginfo/dating-partner/%s/%02d.webp", genderDirectory, d.ID)
			meetScene := datingMeetScenes[d.ID]
			_, err = db.Exec(`
				INSERT INTO datings(
					id, name, image, age, gender, nationality, occupation, description,
					cost, meet_conditions, gifts, locations, meet_scene
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET
					name=excluded.name, image=excluded.image, age=excluded.age,
					gender=excluded.gender, nationality=excluded.nationality,
					occupation=excluded.occupation, description=excluded.description,
					cost=excluded.cost, meet_conditions=excluded.meet_conditions,
					gifts=excluded.gifts, locations=excluded.locations,
					meet_scene=excluded.meet_scene
			`, id, profile.Name, image, d.Age, gender, profile.Nationality,
				occupation, profile.Description, d.Cost, d.MeetConditions, gifts,
				d.Locations, meetScene)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var maleDatingTextReplacer = strings.NewReplacer(
	"名媛千金", "青年投资人",
	"上流名媛", "商界精英",
	"空姐", "空乘",
	"高跟鞋", "皮鞋",
	"化妆品", "护肤品",
	"名牌包", "皮具",
	"丝巾", "领带",
)

// seedBankTasks 写入银行任务数据
func seedBankTasks(db sqlExecutor) error {
	tasks := []struct {
		ID          int
		Name        string
		Description string
		Type        string
		TargetValue int
		Reward      int
	}{
		{1, "存钱新手", "今年存钱达到5000元", "deposit", 5000, 100},
		{2, "储蓄达人", "今年存钱达到10000元", "deposit", 10000, 500},
		{3, "储蓄大户", "今年存钱达到50000元", "deposit", 50000, 2000},
		{4, "超级储蓄", "今年存钱达到100000元", "deposit", 100000, 5000},
		{5, "频繁存款", "今年存款达到5次", "depositcount", 5, 300},
		{6, "取钱新手", "今年取钱达到5000元", "withdraw", 5000, 200},
		{7, "取钱达人", "今年取钱达到10000元", "withdraw", 10000, 1000},
		{8, "取钱大户", "今年取钱达到50000元", "withdraw", 50000, 2000},
		{9, "超级取款", "今年取钱达到100000元", "withdraw", 100000, 5000},
		{10, "频繁取款", "今年取钱达到5次", "withdrawcount", 5, 300},
		{11, "贷款新手", "今年贷款达到10000元", "loan", 10000, 300},
		{12, "贷款达人", "今年贷款达到50000元", "loan", 50000, 1000},
		{13, "贷款大户", "今年贷款达到100000元", "loan", 100000, 2000},
		{14, "超级贷款", "今年贷款达到500000元", "loan", 500000, 5000},
		{15, "勤劳致富", "今年打工达到3次", "work", 3, 500},
		{16, "工作狂", "今年打工达到5次", "work", 5, 1000},
		{17, "工作达人", "今年打工达到10次", "work", 10, 2000},
	}

	for _, t := range tasks {
		_, err := db.Exec(`
			INSERT INTO bank_tasks (id, name, description, type, target_value, reward)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, description=excluded.description, type=excluded.type,
				target_value=excluded.target_value, reward=excluded.reward
		`, t.ID, t.Name, t.Description, t.Type, t.TargetValue, t.Reward)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedDiseases 写入疾病数据
func seedDiseases(db sqlExecutor) error {
	diseases := []struct {
		ID           int
		Name         string
		Type         string
		Symptoms     string
		HealthImpact int
		UpgradeDays  int
		Treatment    string
	}{
		// 小病
		{1, "感冒", "green", "流鼻涕、打喷嚏、轻微发热", -2, 1, `["感冒灵", "板蓝根", "打针"]`},
		{2, "发烧", "green", "体温升高、头痛、乏力", -3, 1, `["布洛芬", "退烧药", "退烧针剂", "打针"]`},
		{3, "咳嗽", "green", "干咳或有痰、喉咙痒", -1, 1, `["止咳糖浆", "川贝枇杷膏", "打针"]`},
		{4, "头痛", "green", "头部疼痛、注意力不集中", -1, 1, `["止痛片", "止痛药", "止痛针", "打针"]`},
		// 中等疾病
		{5, "肠胃炎", "cyan", "腹泻、腹痛、恶心呕吐", -5, 2, `["藿香正气水", "消炎药", "中医调理", "打针"]`},
		{6, "扭伤", "cyan", "关节肿胀、疼痛、活动受限", -4, 2, `["红花油", "膏药", "针灸"]`},
		{7, "腰肌劳损", "cyan", "腰酸背痛、久坐加重", -3, 2, `["膏药", "中医推拿", "针灸"]`},
		// 严重疾病
		{8, "骨折", "yellow", "剧烈疼痛、无法活动、肿胀", -10, 3, `["中医调理", "手术"]`},
		{9, "肺炎", "yellow", "高烧不退、咳嗽有痰、呼吸困难", -8, 3, `["抗生素", "打针"]`},
		{10, "阑尾炎", "yellow", "右下腹剧痛、发烧、恶心", -12, 3, `["中医调理", "手术"]`},
		// 危急疾病
		{11, "心脏病", "red", "胸闷、心悸、呼吸困难", -20, 5, `["速效救心丸", "手术"]`},
		{12, "中风", "red", "突然昏倒、肢体麻木、口歪眼斜", -25, 5, `["中医调理", "手术"]`},
		// ==========资产变化相关疾病==========
		// 经济压力类疾病
		{13, "焦虑症", "cyan", "心慌、失眠、过度担忧财务状况", -5, 3, `["安神补脑液", "中医调理", "心理咨询", "脱胎换骨"]`},
		{14, "失眠症", "green", "入睡困难、睡眠质量差、白天疲惫", -3, 3, `["褪黑素", "中医调理", "安神补脑液", "脱胎换骨"]`},
		{15, "抑郁症", "yellow", "情绪持续低落、对生活失去兴趣、社交回避", -10, 3, `["心理咨询", "脱胎换骨"]`},
		{16, "神经衰弱", "cyan", "易疲劳、注意力不集中、记忆力下降", -4, 3, `["安神补脑液", "中医调理", "脱胎换骨"]`},
		// 富贵病（资产大幅增长或高额资产）
		{17, "高血压", "cyan", "头晕、心悸、血压升高", -4, 5, `["降压药", "打针"]`},
		{18, "痛风", "cyan", "关节红肿疼痛、活动受限", -3, 5, `["痛风药" ,"打针"]`},
		{19, "糖尿病", "yellow", "多饮多尿、体重下降、容易疲劳", -6, 5, `["降糖药", "打针"]`},
		{20, "脂肪肝", "cyan", "右上腹不适、易疲劳、食欲不振", -3, 5, `["护肝片" ,"打针"]`},
	}

	for _, t := range diseases {
		_, err := db.Exec(`
			INSERT INTO diseases (id, name, type, symptoms, healthimpact, upgradedays, treatment)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, type=excluded.type, symptoms=excluded.symptoms,
				healthimpact=excluded.healthimpact, upgradedays=excluded.upgradedays,
				treatment=excluded.treatment
		`, t.ID, t.Name, t.Type, t.Symptoms, t.HealthImpact, t.UpgradeDays, t.Treatment)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedTreats 写入治疗数据
func seedTreats(db sqlExecutor) error {
	treats := []struct {
		ID          int
		Name        string
		Description string
		Price       int
		Heal        int
		SideEffect  int
		Source      string
	}{
		// 药店药品
		{1, "感冒灵", "治疗感冒，缓解流鼻涕打喷嚏", 500, 5, 0, "pharmacy"},
		{2, "板蓝根", "预防治疗感冒，清热解毒", 400, 3, 0, "pharmacy"},
		{3, "布洛芬", "退烧止痛，缓解发热头痛", 1000, 8, 0, "pharmacy"},
		{4, "止咳糖浆", "止咳化痰，润喉", 600, 4, 0, "pharmacy"},
		{5, "止痛片", "缓解各种疼痛", 300, 3, -1, "pharmacy"},
		{6, "红花油", "活血化瘀，消肿止痛", 900, 5, 0, "pharmacy"},
		{7, "川贝枇杷膏", "止咳化痰，润喉", 700, 5, 0, "pharmacy"},
		// 诊所药品
		{8, "退烧针剂", "快速退烧，诊所专用", 2500, 15, 0, "clinic"},
		{9, "消炎药", "抗感染，消炎", 1800, 10, -2, "clinic"},
		{10, "止痛针", "强效止痛", 3000, 12, 0, "clinic"},
		// 中医药品
		{11, "中医调理", "中医处方，治本调理", 3500, 20, 0, "chinese"},
		{12, "藿香正气水", "解表化湿，理气和中", 1200, 12, 0, "chinese"},
		{13, "膏药", "活血通络，消肿止痛", 2200, 15, 0, "chinese"},
		// 大医院药品
		{14, "抗生素", "强效抗感染", 5000, 25, -3, "big_hospital"},
		{15, "速效救心丸", "缓解心脏不适", 8000, 30, 0, "big_hospital"},
		{16, "止痛药", "强效处方止痛", 4500, 20, -5, "big_hospital"},
		// 资产相关疾病药品
		{17, "安神补脑液", "安神定志，改善睡眠", 1000, 8, 0, "pharmacy"},
		{18, "褪黑素", "改善睡眠质量", 800, 5, 0, "pharmacy"},
		{19, "降压药", "控制血压", 2200, 10, 0, "clinic"},
		{20, "痛风药", "缓解痛风疼痛", 2800, 12, 0, "clinic"},
		{21, "降糖药", "控制血糖", 5000, 8, 0, "big_hospital"},
		{22, "护肝片", "保护肝脏", 1800, 10, 0, "pharmacy"},
		// 普遍
		{23, "心理咨询", "缓解心理压力", 5000, 15, 0, "clinic"},
		{24, "中医推拿", "中医治疗，缓解肌肉紧张", 2500, 15, 0, "chinese"},
		{25, "脱胎换骨", "专治压力疾病", 50000, 50, 0, "chinese"},
	}
	for _, t := range treats {
		_, err := db.Exec(`
			INSERT INTO treats (id, name, description, price, heal, sideeffect, source)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, description=excluded.description, price=excluded.price,
				heal=excluded.heal, sideeffect=excluded.sideeffect, source=excluded.source
		`, t.ID, t.Name, t.Description, t.Price, t.Heal, t.SideEffect, t.Source)
		if err != nil {
			return err
		}
	}

	return nil

}

// seedHospitals 写入医院数据
func seedHospitals(db sqlExecutor) error {
	hospitals := []struct {
		ID          int
		Name        string
		Type        string
		Icon        string
		Description string
		Services    string
	}{
		{1, "药店", "pharmacy", "💊", "价格便宜，自行选购药品", `[{"hsid":1,"hsname":"买药","hstype":"medicine","hsprice":0,"hsimmunity":0,"hsdesc":"根据症状选择药品"}]`},
		{2, "社区诊所", "clinic", "🏥", "价格适中，有专业医生", `[{"hsid":1,"hsname":"买药","hstype":"medicine","hsprice":0,"hsimmunity":0,"hsdesc":"根据症状选择药品"},{"hsid":2,"hsname":"打针","hstype":"injection","hsprice":5000,"hsimmunity":25,"hsdesc":"专业护理打针"}]`},
		{3, "中医馆", "chinese", "🌿", "中药治本，针灸调理", `[{"hsid":1,"hsname":"中药","hstype":"medicine","hsprice":3000,"hsimmunity":25,"hsdesc":"恢复慢但一定能治好"},{"hsid":2,"hsname":"针灸","hstype":"acupuncture","hsprice":8000,"hsimmunity":30,"hsdesc":"针灸调理"},{"hsid":3,"hsname":"脱胎换骨","hstype":"bereborn","hsprice":50000,"hsimmunity":50,"hsdesc":"专治压力疾病"}]`},
		{4, "大医院", "big_hospital", "🏨", "价格较高，设备先进", `[{"hsid":1,"hsname":"买药","hstype":"medicine","hsprice":0,"hsimmunity":0,"hsdesc":"处方药品"},{"hsid":2,"hsname":"打针","hstype":"injection","hsprice":5000,"hsimmunity":25,"hsdesc":"专业护理打针"},{"hsid":3,"hsname":"手术","hstype":"surgery","hsprice":50000,"hsimmunity":50,"hsdesc":"治疗严重疾病"}]`},
	}

	for _, t := range hospitals {
		_, err := db.Exec(`
			INSERT INTO hospitals (id, name, type, icon, description, services)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, type=excluded.type, icon=excluded.icon,
				description=excluded.description, services=excluded.services
		`, t.ID, t.Name, t.Type, t.Icon, t.Description, t.Services)
		if err != nil {
			return err
		}
	}
	return nil
}

// seedMiniGames 写入小游戏数据
func seedMiniGames(db sqlExecutor) error {
	games := []struct {
		ID          int
		Name        string
		CNName      string
		Icon        string
		Description string
		Type        string
		Difficulty  string
	}{
		{1, "bank", "银行打工", "🏦", "数钱游戏", "work", `{"0":{"smgneed":0,"smgreward":{"totalmoney":1000},"smgtarget":50,"smgminruntime":0},"1":{"smgneed":0,"smgreward":{"totalmoney":1500},"smgtarget":40,"smgminruntime":0},"2":{"smgneed":0,"smgreward":{"totalmoney":2000},"smgtarget":20,"smgminruntime":0}}`},
		{2, "taxi", "出车接单", "🚕", "接送乘客", "work", `{"0":{"smgneed":0,"smgreward":{"money":200},"smgtarget":5,"smgminruntime":1},"1":{"smgneed":0,"smgreward":{"money":300},"smgtarget":5,"smgminruntime":1},"2":{"smgneed":0,"smgreward":{"money":400},"smgtarget":5,"smgminruntime":1}}`},
		{3, "rps", "猜拳游戏", "✊", "和电脑比猜拳", "casual", `{"0":{"smgneed":50,"smgreward":{"totalmoney":150},"smgtarget":0,"smgminruntime":0}}`},
		{4, "guess", "猜数字", "🔢", "猜1-100的数字", "casual", `{"0":{"smgneed":10,"smgreward":{"totalmoney":300},"smgtarget":0,"smgminruntime":0}}`},
		{5, "dice", "掷骰子", "🎲", "比大小游戏", "casual", `{"0":{"smgneed":100,"smgreward":{"totalmoney":200},"smgtarget":0,"smgminruntime":0}}`},
		{6, "slot", "老虎机", "🎰", "经典老虎机", "casual", `{"0":{"smgneed":200,"smgreward":{"totalmoney":8000},"smgtarget":0,"smgminruntime":0}}`},
		{7, "gobang", "五子棋", "⚫", "连成五子获胜", "board", `{"0":{"smgneed":300,"smgreward":{"totalmoney":600},"smgtarget":0,"smgminruntime":0}}`},
		{8, "jungle", "斗兽棋", "🦁", "兽类相克，渡河制胜", "board", `{"0":{"smgneed":400,"smgreward":{"totalmoney":800},"smgtarget":0,"smgminruntime":0}}`},
		{9, "go", "围棋", "⚪", "围地吃子，策略深奥", "board", `{"0":{"smgneed":500,"smgreward":{"totalmoney":1000},"smgtarget":0,"smgminruntime":0}}`},
		{10, "othello", "黑白棋", "⚫", "翻转棋子", "board", `{"0":{"smgneed":300,"smgreward":{"totalmoney":600},"smgtarget":0,"smgminruntime":0}}`},
		{11, "landbattle", "军旗", "🚩", "陆战棋，夺旗获胜", "board", `{"0":{"smgneed":400,"smgreward":{"totalmoney":800},"smgtarget":0,"smgminruntime":0}}`},
		{12, "chess", "国际象棋", "♟️", "策略对弈", "board", `{"0":{"smgneed":500,"smgreward":{"totalmoney":1000},"smgtarget":0,"smgminruntime":0}}`},
		{13, "fps", "FPS射击", "🎯", "反应速度对决", "competitive", `{"0":{"smgneed":200,"smgreward":{"totalmoney":500},"smgtarget":0,"smgminruntime":0}}`},
		{14, "moba", "MOBA对战", "⚔️", "团队竞技", "competitive", `{"0":{"smgneed":300,"smgreward":{"totalmoney":1000},"smgtarget":0,"smgminruntime":0}}`},
		{15, "racing", "赛车竞速", "🏎️", "速度较量", "competitive", `{"0":{"smgneed":400,"smgreward":{"totalmoney":1000},"smgtarget":0,"smgminruntime":0}}`},
		{16, "fighting", "格斗竞技", "👊", "近身格斗", "competitive", `{"0":{"smgneed":300,"smgreward":{"totalmoney":700},"smgtarget":0,"smgminruntime":0}}`},
		{17, "war", "战争策略", "🛡️", "指挥军队", "competitive", `{"0":{"smgneed":300,"smgreward":{"totalmoney":800},"smgtarget":0,"smgminruntime":0}}`},
		{18, "poker", "德州扑克", "🃏", "心理博弈，长期返奖率约94.5%", "gambling", `{"0":{"smgneed":2000,"smgreward":{"totalmoney":5400},"smgtarget":0,"smgminruntime":10}}`},
		{19, "horseracing", "赛马博彩", "🏇", "按名次和赔率返奖，无额外入场费", "gambling", `{"0":{"smgneed":0,"smgreward":{"totalmoney":70000},"smgtarget":0,"smgminruntime":10}}`},
		{20, "roulette", "轮盘赌", "🎡", "庄家优势5.26%，无额外入场费", "gambling", `{"0":{"smgneed":0,"smgreward":{"totalmoney":4000},"smgtarget":0,"smgminruntime":10}}`},
		{21, "baccarat", "百家乐", "🎴", "庄家优势约1%，无额外入场费", "gambling", `{"0":{"smgneed":0,"smgreward":{"totalmoney":45000},"smgtarget":0,"smgminruntime":10}}`},
		{22, "blackjack", "二十一点", "🃏", "追求21点，无额外入场费", "gambling", `{"0":{"smgneed":0,"smgreward":{"totalmoney":15000},"smgtarget":0,"smgminruntime":10}}`},
		{23, "lottery", "彩票刮刮乐", "🎫", "长期返奖率约95%", "gambling", `{"0":{"smgneed":100,"smgreward":{"totalmoney":15000},"smgtarget":0,"smgminruntime":0}}`},
	}

	for _, game := range games {
		_, err := db.Exec(`
			INSERT INTO minigames (id, name, cnname, icon, description, type, difficulty)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name, cnname=excluded.cnname, icon=excluded.icon,
				description=excluded.description, type=excluded.type,
				difficulty=excluded.difficulty
		`, game.ID, game.Name, game.CNName, game.Icon, game.Description, game.Type, game.Difficulty)
		if err != nil {
			return err
		}
	}

	return nil
}
