<template>
	<el-menu data-testid="game-menu" :default-active="activeMenu" @select="handleSelect" :collapse="false">
    <el-menu-item data-testid="menu-domestic" index="menu-domestic">📈 国内市场</el-menu-item>
    <el-menu-item data-testid="menu-foreign" index="menu-foreign">🌍 国外市场</el-menu-item>
    <el-menu-item data-testid="menu-company" index="menu-company">🚀 创业</el-menu-item>
    <el-menu-item data-testid="menu-stock" index="menu-stock">📊 股市</el-menu-item>
    <el-menu-item data-testid="menu-antique" index="menu-antique">🖼️ 古玩</el-menu-item>
    <el-menu-item data-testid="menu-entertainment" index="menu-entertainment">🎬 娱乐</el-menu-item>
    <el-menu-item data-testid="menu-bank" index="menu-bank">🏦 银行</el-menu-item>
    <el-menu-item data-testid="menu-hospital" index="menu-hospital">🏥 医院</el-menu-item>
    <el-menu-item data-testid="menu-dating" index="menu-dating">💖 约会</el-menu-item>
    <el-menu-item data-testid="menu-carshop" index="menu-carshop">🚗 汽车店</el-menu-item>
    <el-menu-item data-testid="menu-realestate" index="menu-realestate">🏢 售楼部</el-menu-item>
    <el-menu-item data-testid="menu-profile" index="menu-profile">👤 个人属性</el-menu-item>
  </el-menu>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'

// 导入 Vue Router 的路由实例（用于跳转）
const router = useRouter()

// 当前选中的菜单项（默认选中 menu-domestic，即国内市场）
const activeMenu = ref('menu-domestic')

// 处理菜单点击事件的函数，参数 index 是被点击的菜单项 index
const handleSelect = (index) => {
  // 更新当前选中菜单，用于高亮菜单项
  activeMenu.value = index
  // 跳转到对应的路由路径，比如 index 是 menu-foreign，则跳转到 /menu-foreign
  router.push(`/${index}`)
}

</script>


<style scoped>
.el-menu {
	--el-menu-bg-color: var(--panel-main-color);
	--el-menu-text-color: var(--font-color);
	--el-menu-hover-bg-color: var(--select-color);
	--el-menu-active-color: var(--primary-color);
	height: 100%;
	border-right: 0;
}

/* 菜单项 hover 高亮 */
.el-menu-item {
  transition: background 0.3s ease;           /* 背景色变化动画 */
}

.el-menu-item:hover {
  background-color: var(--select-color);         /* 根据主题变化 */
}

/* 选中菜单项样式微调 */
.el-menu-item.is-active {
  background-color: var(--select-color);         /* 根据主题变化 */
	box-shadow: inset 3px 0 0 var(--primary-color);
  font-weight: bold;                            /* 字体加粗 */
}
</style>
