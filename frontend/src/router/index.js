/**
 * ============================================================
 * LifeGame 人生模拟器 - 路由配置文件
 * ============================================================
 *
 * 本文件定义了应用的路由结构，包含：
 * - 国内市场、国外市场、创业、股市等 10 个主要页面
 * - 使用 Vue Router 的 history 模式
 *
 * 路由命名规范：
 * - 路由路径：/menu-{功能名}
 * - 组件名：Menu{功能名}
 *
 * ============================================================
 */

// ==================== Vue Router 核心函数 ====================
import { createRouter, createWebHistory } from 'vue-router'

// ==================== 页面组件懒加载 ====================
const MenuProfile = () => import('@/src/components/Menu/MenuProfile.vue')                  // 个人属性
const MenuDomesticMarket = () => import('@/src/components/Menu/MenuDomesticMarket.vue')    // 国内市场
const MenuForeignMarket = () => import('@/src/components/Menu/MenuForeignMarket.vue')      // 国外市场
const MenuCompany = () => import('@/src/components/Menu/MenuCompany.vue')                  // 创业
const MenuStock = () => import('@/src/components/Menu/MenuStock.vue')                      // 股市
const MenuAntique = () => import('@/src/components/Menu/MenuAntique.vue')                  // 古玩
const MenuBank = () => import('@/src/components/Menu/MenuBank.vue')                        // 银行
const MenuHospital = () => import('@/src/components/Menu/MenuHospital.vue')                // 医院
const MenuRealEstate = () => import('@/src/components/Menu/MenuRealEstate.vue')            // 售楼部
const MenuDating = () => import('@/src/components/Menu/MenuDating.vue')                    // 约会
const MenuEntertainment = () => import('@/src/components/Menu/MenuEntertainment.vue')      // 娱乐
const MenuCarshop = () => import('@/src/components/Menu/MenuCarshop.vue')                  // 汽车店

// ==================== 路由配置 ====================
/**
 * 路由配置数组
 *
 * 路由结构说明：
 * - 默认路由重定向到国内市场
 * - 所有路由使用懒加载方式引入组件
 * - 路由路径采用 kebab-case 命名规范
 */
const routes = [
  // 默认路由：重定向到国内市场
  { path: '/', redirect: '/menu-domestic' },

  // ==================== 个人模块 ====================
  { path: '/menu-profile', component: MenuProfile },            // 个人属性

  // ==================== 市场模块 ====================
  { path: '/menu-domestic', component: MenuDomesticMarket },  // 国内市场
  { path: '/menu-foreign', component: MenuForeignMarket },    // 国外市场

  // ==================== 投资模块 ====================
  { path: '/menu-company', component: MenuCompany },          // 创业
  { path: '/menu-stock', component: MenuStock },              // 股市
  { path: '/menu-antique', component: MenuAntique },          // 古玩

  // ==================== 服务模块 ====================
  { path: '/menu-entertainment', component: MenuEntertainment }, // 娱乐
  { path: '/menu-bank', component: MenuBank },                // 银行
  { path: '/menu-hospital', component: MenuHospital },        // 医院

  // ==================== 生活模块 ====================
  { path: '/menu-dating', component: MenuDating },            // 约会
  { path: '/menu-realestate', component: MenuRealEstate },    // 售楼部
  { path: '/menu-carshop', component: MenuCarshop },          // 汽车店
]

// ==================== 路由实例创建 ====================
/**
 * 创建路由实例
 * 使用 HTML5 History 模式，URL 更加美观（无 # 号）
 */
const router = createRouter({
  history: createWebHistory(),
  routes
})

// ==================== 导出路由实例 ====================
export default router
