<template>
  <div class="product-panel" data-testid="page-antique">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">古玩市场</div>
      <div class="panel-controls">
        <el-tag size="small">{{ gameStore.userInfo.uantique.length }}/{{ gameStore.gameInfo.gmaxholdnum.maholdnum }}</el-tag>
        <el-select v-model="selectedAntique" size="small" placeholder="选择古董" style="width: 150px">
          <el-option v-for="antique in gameStore.userInfo.uantique" :key="antique" :label="antique.ainame" :value="antique" />
        </el-select>
      </div>
    </div>

    <!-- 中间信息区域 -->
    <div class="antique-main">
      <!-- 用户操作 -->
      <div class="antiquestore">
        <div class="section-header">
          <span class="section-title">古董店</span>
        </div>
        <div class="antiquestore-content">
          <div class="action-item">
            <span>鉴定古董：</span>
            <el-button class="action-button" :disabled="selectedAntique == null" type="primary" @click="handleAppraiseAntique">鉴定</el-button>
          </div>
          <div class="action-item">
            <span>修复古董：</span>
            <el-button class="action-button" :disabled="selectedAntique == null" type="success" @click="handleRepairAntique">修复</el-button>
          </div>
          <div class="action-item">
            <span>出售古董：</span>
            <el-button class="action-button" :disabled="selectedAntique == null" type="danger" @click="handleSellAntique">出售</el-button>
          </div>
        </div>
      </div>
      <div class="log-panel">
          <div class="section-header">
            <span class="section-title">📋 操作记录</span>
            <div class="section-controls">
              <el-button size="small" text @click="optLogInfo = []">清空</el-button>
            </div>
          </div>
          <div class="section-content">
            <div v-if="optLogInfo.length === 0" class="section-empty">
              <span class="section-empty-icon">📋</span>
              <span class="section-empty-text">暂无记录</span>
            </div>
            <div v-else class="log-list">
              <div v-for="loginfo in optLogInfo" :key="loginfo.id" class="log-item">
                <span class="log-time">{{ loginfo.time }}</span>
                <span class="log-action" :class="loginfo.type">{{ loginfo.action }}</span>
                <span class="log-item-text">{{ loginfo.detail }}</span>
              </div>
            </div>
          </div>
      </div>
      <div class="uantiqueinfo">
        <!-- 标题 + 古董名称 一行左右对齐 -->
        <div class="section-header">
          <span class="section-title">当前古董信息</span>
        </div>

        <!-- 描述信息 -->
        <el-descriptions border :column="2" size="small">
          <el-descriptions-item label="古董名字">
            {{ selectedAntique?.ainame || '暂无' }}
          </el-descriptions-item>
          <el-descriptions-item label="拍卖价格">
            {{ formatPrice(selectedAntique?.aiprice || '0') }}
          </el-descriptions-item>
          <el-descriptions-item label="真伪程度">
            {{ authenticityText }}
          </el-descriptions-item>
          <el-descriptions-item label="最高价格">
            {{ formatPrice(selectedAntique?.aiprice_max || '0') }}
          </el-descriptions-item>
          <el-descriptions-item label="稀有度">
            {{ selectedAntique?.aiamaterial || '0' }}
          </el-descriptions-item>
          <el-descriptions-item label="完好程度">
            {{ selectedAntique?.aiacondition || '0' }}
          </el-descriptions-item>
          <el-descriptions-item label="持有时间">
            {{ selectedAntique?.aiatime || '0' }} 天
          </el-descriptions-item>
          <el-descriptions-item label="古董等级">
            {{ authenticityLevel }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- 拍卖行展示区 -->
    <div class="antique-bottom">
      <div class="chart">
        <!-- 头部 -->
        <div class="section-header">
          <span class="section-title">拍卖会</span>
          <div class="section-controls">
            <el-tag size="small">{{ gameStore.userInfo.uopportunity.oanum }}/{{ gameStore.gameInfo.gmaxholdnum.maroundnum }}</el-tag>
            <el-button-group size="small">
              <el-button size="small" :type="auctionlevel === 0 ? 'primary' : ''" :disabled="0 > gameStore.calcReputationLevel()" @click="changeAuctionLevel(0)">普通</el-button>
              <el-button size="small" :type="auctionlevel === 1 ? 'primary' : ''" :disabled="1 > gameStore.calcReputationLevel()" @click="changeAuctionLevel(1)">中等</el-button>
              <el-button size="small" :type="auctionlevel === 2 ? 'primary' : ''" :disabled="2 > gameStore.calcReputationLevel()" @click="changeAuctionLevel(2)">高级</el-button>
              <el-button size="small" :type="auctionlevel === 3 ? 'primary' : ''" :disabled="3 > gameStore.calcReputationLevel()" @click="changeAuctionLevel(3)">豪华</el-button>
              <el-button size="small" :type="auctionlevel === 4 ? 'primary' : ''" :disabled="4 > gameStore.calcReputationLevel()" @click="changeAuctionLevel(4)">私人</el-button>
            </el-button-group>
          </div>
        </div>
        <div class="auction-panel" v-if="currentAntique">
          <!-- 左侧：图片+名称 -->
          <div class="antique-left">
            <img class="antique-img" :src="currentAntique?.aiimg" alt="古董图片" @click="previewImage(currentAntique?.aiimg)" @error="handleImgError"/>
            <div class="antique-name">{{ currentAntique?.ainame || '暂无古董'}}</div>
          </div>

          <!-- 右侧：描述+出价信息 -->
          <div class="auction-right">
            <p class="antique-desc">描述：{{ currentAntique?.aidesc || '暂无描述' }}</p>
            <el-tag type="warning" class="auction-tag">当前价格：{{ currentHighestBid || 0 }} 元</el-tag>
            <el-tag type="danger" class="auction-tag">最高价格：{{ currentAntique?.aiprice_max || 0 }} 元</el-tag>
            <div class="auction-action-row">
              <el-input-number
                v-model="bidderPrice"
                :min="currentHighestBid + 100"
                :step="100"
                controls-position="right"
                class="auction-input"
              />
              <el-button type="primary" class="auction-bid-btn" :disabled="!isAuctionActive" @click="placePlayerBid">出价</el-button>
            </div>
            <div class="auction-action-row">
                <el-button type="success" class="auction-button" :disabled="isAuctionActive" @click="nextAuction">下一场竞拍</el-button>
                <el-button type="warning" class="auction-button" :disabled="isAuctionActive || isNextAuctionActive" @click="startAuction">开始竞拍</el-button>
            </div>
          </div>
        </div>
        <div class="auction-panel" v-else>
          <div class="auction-start">
          <el-button type="success" class="auction-button" :disabled="!isNextAuctionActive" @click="nextAuction">下一场竞拍</el-button>
          </div>
        </div>
      </div>
      <div class="log-panel">
        <div class="section-header">
          <span class="section-title">📋 拍卖动态</span>
          <div class="section-controls">
            <el-tag size="small" v-if="isAuctionActive">竞拍剩余 {{ remainingTime }} 秒</el-tag>
          <el-button size="small" text @click="auctionlog = []">清空</el-button>
          </div>
          
        </div>
        <div class="section-content">
          <div v-if="auctionlog.length === 0" class="section-empty">
            <span class="section-empty-icon">🏺</span>
            <span class="section-empty-text">暂无记录</span>
          </div>
          <div v-else class="log-list">
            <div v-for="loginfo in auctionlog" :key="loginfo.id" class="log-item">
              <span class="log-time">{{ loginfo.time }}</span>
              <span class="log-action" :class="loginfo.type">{{ loginfo.action }}</span>
              <span class="log-item-text">{{ loginfo.detail }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <!-- 放大图片 -->
  <el-image-viewer v-if="showPreview" :url-list="[previewSrc]" @close="showPreview = false"
/>

</template>

<script setup>
import { ref, onBeforeUnmount, onMounted, computed, watch } from 'vue';
import { useGameStore } from '@/src/stores/game'
import { GetAntique, AuctionEnd, OperationAntique  } from "@/wailsjs/go/services/App.js";
import { formatPrice } from '@/src/utils/format';

// 常量定义
const BID_INCREMENT = 100;
const BID_COOLDOWN = 500;
const NPC_BID_INTERVAL = 3000;
const DEFAULT_AUCTION_TIME = 5;
const NPC_BIDDERS_COUNT = 5;

// 游戏状态
// 获取游戏信息
const gameStore = useGameStore();
// 设置选定的古董
const selectedAntique = ref(null);
// 古董市场操作日志
const optLogInfo = ref([]);
let optlogId = 0;
// 拍卖行信息
const auctionlog = ref([]);
let alogId = 0
// 当前拍卖的古董
const currentAntique = ref(null);
// 参加竞拍的机器人
const allBids = ref([]);
// 赢家
const winner = ref('');
// 拍卖行等级
const auctionlevel = ref(0);
// 用户拍卖当前价格
const bidderPrice = ref(0);
// 当前拍卖行古董最高的价
const currentHighestBid = ref(0);
// 当场拍卖行倒计时，默认是5秒
const remainingTime = ref(DEFAULT_AUCTION_TIME);
// 显示下一轮拍卖按钮
const isAuctionActive = ref(false);
// 开始拍卖
const isNextAuctionActive = ref(false);
// 防止并发出价 bug，设置锁
let isBiddingLocked = false;
// 定时器
// npc出价的时间
let auctionTimer = null;
// 用户出价的时间
let bidTimer = null;
// 图片预览
const showPreview = ref(false);
// 图片url
const previewSrc = ref('');
// 开启预览图片
const previewImage = (src) => {
  previewSrc.value = src
  showPreview.value = true
}

// 工具函数：添加日志
const addLog = (detail, logtype, type="red") => {
  const now = new Date()
  let timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
  const actionMap = { green: '日志', red: '主要'}
  if (logtype === 'optlog') {
    optLogInfo.value.unshift({
      id: optlogId++,
      time: timeStr,
      type: type,
      action: actionMap[type] || type,
      detail: detail,
    });
  } else if (logtype === 'alog'){
    auctionlog.value.unshift({
      id: alogId++,
      time: timeStr,
      type: type,
      action: actionMap[type] || type,
      detail: detail,
    });
  } else {
    optLogInfo.value.unshift({
      id: optlogId++,
      time: timeStr,
      type: type,
      action: actionMap[type] || type,
      detail: detail,
    });
    auctionlog.value.unshift({
      id: alogId++,
      time: timeStr,
      type: type,
      action: actionMap[type] || type,
      detail: detail,
    });
  }
}

// 工具函数：更新当前选中的古董
const updateSelectedAntique = () => {
  if (!selectedAntique.value) return
  
  const updated = gameStore.userInfo.uantique.find(a => a.aiid === selectedAntique.value.aiid)
  if (updated) selectedAntique.value = updated
}

// 工具函数：处理古董操作
const handleAntiqueOperation = async (operationType, successMsgPrefix) => {
  try {
    const data = await OperationAntique(selectedAntique.value.aiid, operationType)
    if (data.code === 200) {
      gameStore.applyUserInfo(data.userinfo)
      updateSelectedAntique()
      addLog(data.topinfo, 'optlog')
      if (operationType === 3) selectedAntique.value = null // 出售后清空选择
    } else {
      addLog(`${successMsgPrefix}失败，原因是${data.msg}`, 'optlog')
    }
  } catch (e) {
    addLog(`${successMsgPrefix}失败，原因是${e}`, 'optlog')
  }
}
// 古董的验证信息
const authenticityText = computed(() => {
  const value = selectedAntique.value?.aiidisplay
  switch (value) {
    case 1:
      return '真品'
    case 2:
      return '赝品'
    default:
      return '未知'
  }
})
// 用户点击选择拍卖行等级
function changeAuctionLevel(level) {
  auctionlevel.value = level;
}
// 古董的拍卖行等级
const authenticityLevel = computed(() => {
  const value = selectedAntique.value?.ailevel
  switch (value) {
    case 0:
      return '普通'
    case 1:
      return '中等'
    case 2:
      return '高级'
    case 3:
      return '豪华'
    case 4:
      return '私人'
    default:
      return '未知'
  }
})
// 鉴定古董
const handleAppraiseAntique = () => handleAntiqueOperation(1, '鉴定')
// 修复古董
const handleRepairAntique = () => handleAntiqueOperation(2, '修复')

// 出售古董
const handleSellAntique = () => handleAntiqueOperation(3, '出售')
// NPC 竞拍者配置
const NPC_TYPES = ['cautious', 'aggressive', 'balanced']
const NPC_NAMES = ['李', '王', '张', '刘', '赵', '孙', '周', '吴', '杨', '黄', '陈', '邹', '杜', '贾', '唐', '钱', '任', '傅', '林', '罗']
const npcList = NPC_NAMES.map(name => ({
  name: `小${name}`,
  type: NPC_TYPES[Math.floor(Math.random() * NPC_TYPES.length)]
}))

// 随机选取 N 个不重复的 NPC
function selectBidders(count) {
  const shuffled = [...npcList].sort(() => Math.random() - 0.5)
  return shuffled.slice(0, count)
}

// 下一轮拍卖
async function nextAuction() {
  const data = await GetAntique(auctionlevel.value);
  if (data.code === 200) {
    // 将数目给用户
    gameStore.userInfo.uopportunity.oanum = data.oanum;
    // 给当前古董变量赋值
    currentAntique.value = data.currentAntique;
    // 重置用户当前最高出价
    bidderPrice.value = 0;
    // 清空竞拍结果
    auctionlog.value = [];
    alogId = 0
    // 清空胜出者
    winner.value = '';
    // 设置当前拍卖行古董最高的价
    currentHighestBid.value = currentAntique.value.aiprice;
    isNextAuctionActive.value = false;    // 禁用“下一场竞拍”
    isAuctionActive.value = false;        // 等待开始
  } else {
    addLog(`获取古董失败，原因是${data.msg}`, 'optlog');
  }
}

// 开始拍卖
function startAuction() {
  // 清空出价信息
  allBids.value = [];
  isAuctionActive.value = true;         // 拍卖开始，出价按钮可用
  isNextAuctionActive.value = false;    // 下一轮拍卖按钮不可用
  // 设置用户当前拍卖的价格
  bidderPrice.value = currentHighestBid.value + BID_INCREMENT;
  // 设置拍卖倒计时
  remainingTime.value = DEFAULT_AUCTION_TIME;
  // 设置5个npc玩家
  const selectedBidders = selectBidders(NPC_BIDDERS_COUNT);
  // 输出拍卖信息
  const startMsg = `拍卖会开始，古董名称：${currentAntique.value.ainame}！`;
  addLog(startMsg, 'all');
  // 输出npc信息
  addLog(`参与竞拍的NPC：${selectedBidders.map(b => b.name).join('，')}`, 'optlog', 'green');
  // 开启npc拍卖倒计时
  auctionTimer = setInterval(() => {
    if (remainingTime.value > 0) {
      remainingTime.value--;
    } else {
      // 用于停止由setInterval创建的定时器。定时器是通过setInterval函数设置的，它会在指定的时间间隔内重复执行一段代码或函数。当不再需要定时器时，可以使用clearInterval函数来停止它
      endAuction();
    }
  }, 1000);
  // NPC 出价
  // 每 3 秒钟执行一次这个函数，用于模拟竞拍过程
  bidTimer = setInterval(() => {
    selectedBidders.forEach((npc) => {
      const amount = generateBid(npc.type, currentHighestBid.value);
      if (amount > currentHighestBid.value) {
        allBids.value.push({ name: npc.name, amount });
        addLog(`${npc.name}出价：${amount} 元`, 'alog', 'green');
        currentHighestBid.value = amount;
      }
    });
    // 玩家出价 = 当前最高出价 + 100
    bidderPrice.value = currentHighestBid.value + BID_INCREMENT;
  }, NPC_BID_INTERVAL);
}
// 设置机器人拍卖的核心方法
function generateBid(type, currentBid) {
  if (currentBid >= ((currentAntique.value.aiprice_max-currentAntique.value.aiprice) * 0.5+currentAntique.value.aiprice)  && Math.random() > 0.5) {
    return currentBid;
  }
  if (currentBid >= ((currentAntique.value.aiprice_max-currentAntique.value.aiprice) * 0.85+currentAntique.value.aiprice)  && Math.random() > 0.25) {
    return currentBid;
  }

  if (currentBid >= currentAntique.value.aiprice_max && Math.random() > 0.05) {
    return currentBid;
  }

  // 根据类型生成出价范围
  const bidRange = {
    cautious: [100, 1000],
    balanced: [100, 5000],
    aggressive: [100, 10000],
  };

  // 随机生成出价
  const [minAdd, maxAdd] = bidRange[type] || [100, 1000];
  const increment = Math.floor(Math.random() * (maxAdd - minAdd + 1)) + minAdd;
  return currentBid + increment;
}
// 结束竞拍
async function endAuction() {
  // 结束竞拍
  // 用于停止由setInterval创建的定时器。定时器是通过setInterval函数设置的，它会在指定的时间间隔内重复执行一段代码或函数。当不再需要定时器时，可以使用clearInterval函数来停止它
  clearInterval(bidTimer);
  clearInterval(auctionTimer);
  const sorted = [...allBids.value].sort((a, b) => b.amount - a.amount);
  const top = sorted[0];
  if (!top) {
    addLog('本场无人出价，拍卖流拍。', 'all');
    isAuctionActive.value = false;
    isNextAuctionActive.value = true;
    return;
  }
  winner.value = top.name;
  // 输出拍卖结果
  const endMsg = `${currentAntique.value.ainame}拍卖会竞拍结束！`;
  addLog(endMsg, 'all');
  if (top.name === '你') {
    const successMsg = `恭喜你以 ${top.amount} 元竞拍成功！`;
    addLog(successMsg, 'all');
    // 将古董添加到玩家的仓库中,并计算钱
    const data = await AuctionEnd(top.amount, currentAntique.value.aiid);
    if (data.code === 200) {
      gameStore.applyUserInfo(data.userinfo);
      currentAntique.value = null;
    } else {
      const failMsg = `竞拍失败！${data.msg}`;
      addLog(failMsg, 'all');
    }
  } else {
    const winnerMsg = `${top.name} 出价 ${top.amount} 元赢得了拍卖！`;
    addLog(winnerMsg, 'all');
  }
  isAuctionActive.value = false;         // 拍卖结束，出价按钮不可用
  isNextAuctionActive.value = true;      // 下一轮按钮可用
}

// 玩家出价
function placePlayerBid() {
  // 如果冷却时间还没结束，不让玩家出价
  if (isBiddingLocked || !isAuctionActive.value) return;
  isBiddingLocked = true;

  setTimeout(() => {
    isBiddingLocked = false;
  }, BID_COOLDOWN);

  if (bidderPrice.value > currentHighestBid.value) {
    allBids.value.push({ name: '你', amount: bidderPrice.value });
    addLog(`你出价：${bidderPrice.value} 元`, 'alog', 'green');
    currentHighestBid.value = bidderPrice.value;
    // 每次出价后重置竞拍倒计时为5秒
    remainingTime.value = DEFAULT_AUCTION_TIME;
  } else {
    addLog(`你的出价必须高于当前最高出价（${currentHighestBid.value} 元）。`, 'alog', 'green');
  }
}

// 图片不存在返回默认图片
function handleImgError(event) {
  const target = event.target
  target.src = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAAEElEQVR42mP8z8BQDwADjAHi4VbOygAAAABJRU5ErkJggg=='
}
// 页面加载时启动竞拍
onMounted(async () => {
  // 默认选择第一个古董
  if (gameStore.userInfo.uantique.length > 0) {
    selectedAntique.value = gameStore.userInfo.uantique[0]
  }
  // 调整状态
  isAuctionActive.value = false;        // 拍卖未开始
  isNextAuctionActive.value = true;    // 不能进入下一轮
});
// 监听用户古董列表变化，自动更新当前选中的古董数据
watch(() => gameStore.userInfo.uantique, (newAntiques) => {
  if (selectedAntique.value && newAntiques && newAntiques.length > 0) {
    // 找到当前选中古董的 aiId
    const currentAiId = selectedAntique.value.aiid;
    // 从新的古董列表中找到对应的古董
    const updatedAntique = newAntiques.find(a => a.aiid === currentAiId);
    if (updatedAntique) {
      // 更新当前选中的古董数据
      selectedAntique.value = updatedAntique;
    } else if (newAntiques.length > 0) {
      // 如果原古董不存在了（比如被出售），选择第一个
      selectedAntique.value = newAntiques[0];
    } else {
      // 如果没有古董了，清空选择
      selectedAntique.value = null;
    }
  }
}, { deep: true });
// 在组件卸载之前被调用。这个钩子函数的主要作用是在组件从 DOM 中移除之前执行一些清理工作，例如移除事件监听器、停止定时器等
onBeforeUnmount(() => {
  clearInterval(auctionTimer);
  clearInterval(bidTimer);
});

// 监听用户年龄变化，当新年到来时强制结束拍卖
watch(() => gameStore.userInfo?.uage, (newAge, oldAge) => {
  // 当年龄增加时（新年），强制结束拍卖
  if (newAge && oldAge && newAge > oldAge) {
    // 清除定时器
    clearInterval(auctionTimer);
    clearInterval(bidTimer);
    // 重置拍卖状态
    isAuctionActive.value = false;
    isNextAuctionActive.value = true;
    // 清空当前拍卖古董
    currentAntique.value = null;
    addLog('⏰ 拍卖会时间结束，未完成的竞拍已取消', 'optlog');
  }
})
</script>

<style scoped>
/* 拍卖行展示区 */
.antique-bottom {
  height: 50%;
  display: flex;  /* 使子元素水平排列 */
  flex: 1; /* 占据剩余空间 */
  gap: 5px; /* 子元素之间有5px的间距 */
  margin-top: 5px;   /* 子元素之间有5px的间距 */
}

/* 拍卖会 */
.chart {
  flex: 2; /* 占据一半宽度 */
  border: 1px solid var(--border-color); /* 添加实线边框 */
  border-radius: 5px; /* 添加圆角边框 */
  background: var(--panel-color);
  padding: 10px; /* 添加内边距 */
  display: flex; /* 使子元素水平排列 */
  flex-direction: column; /* 子元素垂直排列 */
  overflow: hidden; /* 隐藏溢出内容 */
}

/* 拍卖会开始页面 */
.auction-start {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.uantiqueinfo :deep(.el-descriptions) {
  flex: 1;
  overflow: auto;
}

.antiquestore, .uantiqueinfo {
  border: 1px solid var(--border-color); /* 添加实线边框 */
  border-radius: 8px; /* 添加圆角边框 */
  background: var(--panel-color);
  padding: 10px; /* 添加内边距 */
  flex-direction: column; /* 子元素垂直排列 */
  overflow: hidden; /* 隐藏溢出内容 */
}
/* 面板 */
.antiquestore {
  flex: 1; /* 占据一半宽度 */
  display: flex; /* 使子元素水平排列 */
}
.uantiqueinfo:nth-child(3) {
  flex: 1 1 20%;   /* 调整宽度 */
}

/* 操作日志面板 - 宽度覆盖 */
.log-panel {
  width: 240px;
}

/* 底部信息区：交易记录 + 用户操作 */
.antique-main {
  display: flex;
  gap: 5px;
  margin-top: 5px;
  height: 170px;
}

.el-timeline {
  display: flex;
  flex: 1;
  flex-direction: column;
  height: 80%;
  /* 移除滚动条预留空间 */
  scrollbar-width: none;                         /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;                      /* IE/Edge 浏览器隐藏滚动条 */
}
.el-timeline::-webkit-scrollbar {
  display: none;                                 /* Chrome/Safari 隐藏滚动条 */
}
.custom-timeline-item {
  font-size: 12px; /* 调整为你想要的大小，例如 12px 或更小 */
  line-height: 1.0;
}

.auction-panel {
  width: 100%;
  flex: 1;
  overflow: auto;       /* 避免内容撑开 chart，高度超出时滚动 */
  display: flex;
  border-radius: 12px;
  gap: 5px;
  margin: 0 auto;
  max-height: 100%;     /* 限高 */
}

.antique-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center; /* 垂直居中 */
  text-align: center;
}

.antique-img {
  max-width: 100%;
  max-height: 90%;        /* 设置最大高度 */
  border-radius: 8px;
  cursor: pointer;
  object-fit: contain;
}

.antique-name {
  font-weight: bold;
  font-size: 16px;
}

.auction-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.antique-desc {
  font-size: 14px;
  color: var(--font-light);
  overflow-y: auto;
  height: 40px;
  padding: 5px;                  /* 内边距 */
  line-height: 21px;
  background-color: var(--panel-color);      /* 柔和背景 */
  border: 1px solid var(--border-color);         /* 边框 */
  border-radius: 6px;             /* 圆角 */
  white-space: pre-wrap;          /* 保持换行格式 */
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.05); /* 添加内阴影 */
}

.antique-desc::-webkit-scrollbar {
  width: 0px;
  background: transparent;
}

.auction-tag {
  width: 100%;
  height: 20px;
  line-height: 32px;
  box-sizing: border-box;
  padding: 0 12px; /* 和默认 input padding 接近 */
  display: flex;
  align-items: center;
}

.auction-input {
  height: 32px;
  width: 100%;
  max-width: 240px;
  flex: 2;
  box-sizing: border-box;
}

.auction-button {
  white-space: nowrap; 
}

.auction-action-row {
  display: flex;
  gap: 8px;
  width: 100%;
  align-items: center;
}
.auction-action-row .auction-input,
.auction-action-row .auction-button {
  flex: 1;
  margin: 0;
  flex-shrink: 0;
}

.auction-bid-btn {
  flex: 1;
  white-space: nowrap;
  padding: 0 12px;
}
.antiquestore-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
  width: 100%;               /* 撑满容器宽度 */
  padding: 8px 0;
  box-sizing: border-box;
}

.action-item {
  display: flex;
  justify-content: space-between;  /* 左右两端对齐 */
  align-items: center;
  width: 100%;                     /* 撑满每一行 */
  padding: 0 10px;
  box-sizing: border-box;
}
.action-item .action-button {
  width: 30%;
}
</style>
