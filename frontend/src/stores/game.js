/**
 * ============================================================
 * LifeGame 人生模拟器 - Pinia 状态管理仓库
 * ============================================================
 *
 * 本文件使用 Pinia 定义全局游戏状态，包含：
 * - 游戏信息（gameInfo）：游戏运行时的全局数据
 * - 用户信息（userInfo）：玩家角色的属性和资产
 * - 商品信息：国内/国际商品的持有情况
 * - 公告信息：游戏内的消息通知
 *
 * 使用方法：
 * import { useGameStore } from '@/src/stores/game'
 * const gameStore = useGameStore()
 *
 * ============================================================
 */

import { defineStore } from 'pinia'
// --------------------
// 游戏信息
// --------------------
const createInitialGameInfo = () => ({
    gid: 1,                    // 游戏 ID
    gname: '人生模拟器',       // 游戏名称
    gtime: 100,                // 游戏时间（年）
    gdifficulty: 1,            // 游戏难度
    giteminsinfo: null,        // 国内商品信息列表
    gitemoutinfo: null,        // 国际商品信息列表
    gmaxholdnum: null,          // 游戏可以最大数量
    gcompanyinfo: null,        // 公司信息列表
    gstocknews: [],            // 股市新闻动态
    gstockinfo: [],            // 股票信息列表
    gstockupdatecount: 0,      // 本年度已生成的行情次数
    gdatinginfo: null,            // 女友信息列表
    ghouseinfo: null,            // 房产信息列表
    gcarinfo: null,            // 汽车信息列表
    gdiseaseinfo: null,        // 疾病信息列表
    gfame: 150,                // 名声值
    gimmunity: 100,                // 免疫力
})

// --------------------
// 用户信息
// --------------------
const createInitialUserInfo = () => ({
    uid: 1,                    // 用户 ID
    uname: '',                 // 用户名
    usex: true,                // 性别：true=男, false=女
    uage: 18,                  // 年龄
    uassets: 0,                // 净资产
    ucash: 0,                  // 现金
    ubank: 0,                  // 银行存款
    uimmunity: 80,             // 免疫力
    ucriticalhealthyears: 0,   // 连续低免疫年度；恢复到安全线后清零
    udiseases: {},             // 疾病列表
    ufame: 0,                  // 名声值
    uitemins: {},              // 持有的国内商品
    uitemout: {},            // 持有的国际商品
    uantique: [],              // 持有的古董列表
    ucompany: {},            // 持有的公司
    ustock: {},              // 持有的股票
    uopportunity: null,      // 用户可打工、游戏、约会次数、逛街、参加拍卖会等次数
    uloan: 0,                  // 贷款金额
    uloanoverdue: 0,           // 逾期贷款
    uhouse: {},                // 持有的房产列表
    ucar: {},                  // 持有的汽车列表
    udating: {},                // 女友
    umarrieddatingid: 0,        // 当前配偶约会对象 ID，0 表示未婚
    uminigamerecords: {},    // 小游戏累计统计
    
})

const createInitialAnnouncements = () => ({
    announceins: ['欢迎来到人生模拟器！'],
    announceout: [],
    announcecompany: [],
    announcegame: [],
    announcehealthy: [],
})

const buildUserItemList = (userItems = {}, marketItems = {}) => {
    const items = []
    let count = 0

    for (const [itemId, userItem] of Object.entries(userItems || {})) {
        const marketItem = marketItems?.[itemId]
        if (!marketItem || Number(userItem?.uitemnum) === 0) {
            continue
        }

        const price = marketItem.iidisplay ? marketItem.iiprice : userItem.uicostprice
        items.push({
            id: Number(itemId),
            iiname: marketItem.iiname,
            iiprice: price,
            buyprice: userItem.uicostprice,
            num: userItem.uitemnum,
        })
        count += Number(userItem.uitemnum)
    }

    return { items, count }
}
export const useGameStore = defineStore('game', {
    state: () => ({
        gameInfo: createInitialGameInfo(),
        userInfo: createInitialUserInfo(),
        // --------------------
        // 国内商品信息
        // --------------------
        userItemInsData: [],           // 用户持有的国内商品详情
        userItemInsCount: 0,           // 用户持有的国内商品总数

        // --------------------
        // 国际商品信息
        // --------------------
        userItemOutData: [],           // 用户持有的国际商品详情
        userItemOutCount: 0,           // 用户持有的国际商品总数

        // --------------------
        // 其他信息
        // --------------------
        announcements: createInitialAnnouncements(),           // 公告信息
        difficulty: null,              // 游戏难度配置
		stockEpoch: '',
		stockVersion: 0,
		stockUpdatedAt: 0,
		stockRemaining: 20,
		stockMarketClosed: false,
    }),

    // ==================== 方法定义 ====================
    /**
     * actions: 定义修改状态的方法
     * 可以包含异步操作和业务逻辑
     */
    actions: {
        // 统一刷新持仓派生列表，避免组件手动调用旧逻辑。
        syncUserItems() {
            const domestic = buildUserItemList(this.userInfo.uitemins, this.gameInfo.giteminsinfo)
            const foreign = buildUserItemList(this.userInfo.uitemout, this.gameInfo.gitemoutinfo)

            this.userItemInsData = domestic.items
            this.userItemInsCount = domestic.count
            this.userItemOutData = foreign.items
            this.userItemOutCount = foreign.count
        },

        // 游戏初始化或年推进后，统一接收整包数据并同步派生状态。
        applyGameData(data = {}) {
            if (data.userinfo) {
                this.userInfo = data.userinfo
            }
            if (data.gameinfo) {
                const hasStockClock = Object.prototype.hasOwnProperty.call(data, 'stockepoch')
                const incomingEpoch = String(data.stockepoch || '')
                const incomingVersion = Number(data.stockversion || 0)
                const keepNewerStocks = hasStockClock &&
                    incomingEpoch === this.stockEpoch &&
                    incomingVersion < this.stockVersion

                this.gameInfo = keepNewerStocks
                    ? {
                        ...data.gameinfo,
                        gstockinfo: this.gameInfo.gstockinfo,
                        gstocknews: this.gameInfo.gstocknews,
                    }
                    : data.gameinfo
				this.stockRemaining = Math.max(0, 20 - Number(data.gameinfo.gstockupdatecount || 0))
				this.stockMarketClosed = this.stockRemaining === 0
            }
            if (data.announce) {
                this.announcements = data.announce
            }
            if (Object.prototype.hasOwnProperty.call(data, 'difficulty')) {
                this.difficulty = data.difficulty
            }
            if (Object.prototype.hasOwnProperty.call(data, 'stockepoch')) {
                const incomingEpoch = String(data.stockepoch || '')
                const incomingVersion = Number(data.stockversion || 0)
                if (incomingEpoch !== this.stockEpoch || incomingVersion >= this.stockVersion) {
                    this.stockEpoch = incomingEpoch
                    this.stockVersion = incomingVersion
                    this.stockUpdatedAt = 0
                }
            }
            this.syncUserItems()
        },

        // 局部用户变更也走统一入口，保证持仓派生数据同步刷新。
        applyUserInfo(userInfo) {
            if (!userInfo) return
            this.userInfo = userInfo
            this.syncUserItems()
        },

        // 局部游戏信息变更同样统一处理，避免页面漏掉持仓重算。
        applyGameInfo(gameInfo) {
            if (!gameInfo) return
            this.gameInfo = gameInfo
            this.syncUserItems()
        },

        // 股票轮询只替换行情和新闻，避免每五秒覆盖完整游戏状态。
        applyStockUpdate(update = {}) {
            const epoch = String(update.epoch || '')
            const version = Number(update.version || 0)
            if (!Array.isArray(update.stocks)) return false
            if (this.stockEpoch && epoch !== this.stockEpoch) return false
            if (epoch === this.stockEpoch && version <= this.stockVersion) return false

            this.gameInfo = {
                ...this.gameInfo,
                gstockinfo: update.stocks,
                gstocknews: Array.isArray(update.news) ? update.news : this.gameInfo.gstocknews,
            }
            this.stockEpoch = epoch
            this.stockVersion = version
            this.stockUpdatedAt = Number(update.updatedAt || 0)
			this.stockRemaining = Math.max(0, Number(update.remaining ?? this.stockRemaining))
			this.stockMarketClosed = Boolean(update.marketclosed) || this.stockRemaining === 0
            return true
        },

        setAnnouncements(announcements) {
            this.announcements = announcements || createInitialAnnouncements()
        },

        resetGameState() {
            this.gameInfo = createInitialGameInfo()
            this.userInfo = createInitialUserInfo()
            this.userItemInsData = []
            this.userItemInsCount = 0
            this.userItemOutData = []
            this.userItemOutCount = 0
            this.announcements = createInitialAnnouncements()
            this.difficulty = null
            this.stockEpoch = ''
            this.stockVersion = 0
            this.stockUpdatedAt = 0
			this.stockRemaining = 20
			this.stockMarketClosed = false
        },

        /**
         * 计算名声等级
         * 根据用户名声值返回对应的等级
         * 用于解锁不同级别的拍卖会
         *
         * @returns {number} 名声等级
         */
        calcReputationLevel() {
            if (this.userInfo.ufame >= 0 && this.userInfo.ufame <= 70) {
                return 0  // 普通
            } else if (this.userInfo.ufame > 70 && this.userInfo.ufame <= 110) {
                return 1  // 中等
            } else if (this.userInfo.ufame > 110 && this.userInfo.ufame <= 130) {
                return 2  // 高级
            } else if (this.userInfo.ufame > 130 && this.userInfo.ufame <= 145) {
                return 3  // 豪华
            } else if (this.userInfo.ufame > 145 && this.userInfo.ufame <= 150) {
                return 4  // 私人
            } else {
                return -1  // 老赖（超出范围）
            }
        }
    },
})
