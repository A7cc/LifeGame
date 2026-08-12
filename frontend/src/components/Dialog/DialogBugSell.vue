<template>
    <el-dialog data-testid="trade-dialog" v-model="DialogVisible" :title="currentStr+'商品'" width="250px" align-center class="custom-dialog">
    <div class="dialog-body">
      <p>请输入{{ currentStr }}数量：</p>
      <el-input-number data-testid="trade-quantity" v-model="Count" :min="1" :max="9999" size="small" class="quantity-input"/>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button data-testid="trade-cancel" @click="cancel">取消</el-button>
        <el-button data-testid="trade-confirm" type="primary" @click="confirm">确认</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { BuyItem, SellItem } from "@/wailsjs/go/services/App.js";
import { useGameStore } from '@/src/stores/game'


// 用pinia存放数据
const gameStore = useGameStore()

// 打开弹窗控制
const DialogVisible = ref(false)
// 弹窗显示的文本购买出售
const currentStr = ref('')
// 弹窗购买出售的数量
const Count = ref(1)
// 弹窗的一些信息
const currentItem = ref(null)
// 设置地区
const regionStr = ref('')

// 打开弹窗设置
const useItem = (index, name, usestr, regionstr) => {
  currentItem.value = { id: Number(index), name: name }
  // 设置弹窗物品统计
  if ( name.slice(-2) == "公司" ){
    Count.value = 1000
  }else{
    Count.value = 1
  }
  // 弹窗隐藏开关
  DialogVisible.value = true
  currentStr.value = usestr
  // 设置地区
  regionStr.value = regionstr
}

// 点击确认后的参数
const confirm = async () => {
  // 选中后物品的信息
  const itemdata = currentItem.value
  // 如果itemdata为空，则直接返回
  if (!itemdata) return
  try{
    let data = null
    if (Count.value <= 0) {
      ElMessage.error('数量不能小于等于0')
      return
    }
    // 判断当前是购买还是出售
    if (currentStr.value === '购买') {
      data = await BuyItem(itemdata.id, Count.value, regionStr.value)
    }else{
      data = await SellItem(itemdata.id, Count.value, regionStr.value)      
    }
    if (data.code == 200) {
      gameStore.applyUserInfo(data.userinfo)
      // 弹窗提示
      let msg = `已${currentStr.value} ${Count.value} 个${itemdata.name}`
      if (currentStr.value != '购买') {
        msg = msg+data.msg
      }
      ElMessage.success(msg)
    } else {
      ElMessage.error(data.msg || `${currentStr.value}失败`);
    }
  }catch (err) {
      console.error('调用 BuyItem/SellItem 异常：', err)
      ElMessage.error(`${currentStr.value}失败`);
  }
  // 隐藏弹窗
  DialogVisible.value = false
}


// 取消
const cancel = () => {
  DialogVisible.value = false
}
// 通过 defineExpose 暴露方法给父组件
defineExpose({
    useItem
})
</script>

<style scoped>

/* ===== 游戏风格弹窗统一样式：购买 & 使用 ===== */
.custom-dialog .el-dialog {
  width: 340px !important;          /* 控制弹窗宽度更紧凑 */
  border-radius: 16px;              /* 整体圆角 */
}

/* 弹窗头部样式 */
.custom-dialog .el-dialog__header {
  background: var(--gradient-primary);  /* 蓝色渐变背景 */
  color: var(--font-color);                        /* 白色文字 */
  font-weight: bold;                  /* 加粗字体 */
  font-size: 16px;                    /* 标题字号 */
  padding: 10px 16px;                 /* 更小的头部高度 */
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;
}

/* 弹窗主体内容区域 */
.custom-dialog .el-dialog__body {
  padding: 16px;
  padding-bottom: 8px;               /* 缩小底部 padding */
}

/* 文本描述样式 */
.dialog-body p {
  margin-bottom: 8px;
  font-size: 14px;
  text-align: center;                /* 居中文本，更适合游戏弹窗风格 */
}

/* 数量输入框 */
.quantity-input {
  width: 180px;                      /* 不占整行，更精致 */
  margin: 0 auto 12px;               /* 居中并与按钮拉开距离 */
  display: block;
}

/* 弹窗底部按钮区域 */
.dialog-footer {
  display: flex;
  justify-content: center;          /* 居中按钮 */
  gap: 15px;                          /* 按钮间距 */
  margin-top: 4px;
  margin-bottom: 4px;
}
</style>
