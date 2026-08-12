# LifeGame 二次开发手册

> 适用对象：后续开发者、维护者和协助开发的 AI  
> 当前基线：Wails 2.10.2、Go 1.25、Vue 3、SQLite  
> 文档更新时间：2026-08-12

## 1. 手册目的与最高原则

本手册用于让第一次接触 LifeGame 的开发者或 AI 快速理解项目，能够在不破坏现有玩法和资源约定的前提下继续优化或添加功能。

项目所有者是产品和玩法的最终决策者。优化内容、功能范围、数值方向和视觉风格由项目所有者提出或确认。开发者或 AI 应当：

1. 先阅读根目录 `README.md` 和本手册，再检查与需求直接相关的当前代码。
2. 以项目所有者本次明确提出的需求作为开发范围，不根据个人偏好擅自改玩法。
3. 可以发现问题、解释影响并提出方案；未经确认，不扩大为新的游戏规则或大规模重构。
4. 修改时保持现有功能可运行，并为新增核心规则补充测试。
5. 不添加旧数据库、旧配置、旧图片目录或旧存档兼容代码。本项目只维护当前格式。
6. 不直接覆盖、删除或重建用户真实的 `~/.lifegame`，除非项目所有者明确要求这样做。
7. 修改架构、目录约定、数据格式或关键规则后，同步更新本手册。

## 2. 项目概览

LifeGame 是单机桌面人生模拟游戏。Wails 将 Go 后端和 Vue 前端打包为一个桌面可执行程序。

主要系统包括：

- 年度推进、难度和结局评价
- 国内/国外市场交易
- 股票行情和创业公司
- 银行存取款、贷款和任务
- 古董拍卖、鉴定和出售
- 房产、车辆购买和出售
- 免疫力、疾病、药品与医院
- 男女各 50 位约会对象、相遇地点、约会、送礼、互动、换装、结婚、分手和离婚
- 工作、竞技、棋类、休闲和赌博类小游戏
- 当前格式存档、读取和删除
- 明暗主题及桌面端场景化图片预览
- 低占用程序化背景音乐及本地静音偏好

### 2.1 运行架构

```text
Vue 组件/对话框
       │
       ▼
Wails 生成的 frontend/wailsjs 绑定
       │
       ▼
services.App 公共方法（业务编排、校验、响应）
       │
       ├── core（状态模型、纯规则、数值计算）
       └── internal/db（SQLite 表、默认数据、读取缓存）
```

界面不能直接访问 SQLite。所有会改变游戏状态的操作都应经过 `services.App`，再将后端返回的快照写入 Pinia。

### 2.2 启动流程

```text
main.go
  ├── 确定 ~/.lifegame/images 与 ~/.lifegame/audio
  ├── 只补齐缺失的内置媒体
  ├── 启动 Wails 和外部资源优先处理器
  └── services.App.Startup
        ├── 读取/创建 config.yaml
        ├── 刷新 core 运行时配置值
        ├── 打开/创建 lifegame.db
        ├── 建立当前表和索引
        └── 把数据库参考数据载入 core 缓存
```

前端 `App.vue` 先读取启动状态。初始化完成后，玩家可以创建新游戏或读取存档；`GameMain.vue` 负责游戏主框架，具体系统由路由页面承载。

## 3. 目录与职责

| 路径 | 职责 | 修改时注意 |
| --- | --- | --- |
| `main.go` | Wails 窗口、嵌入前端、外部图片/音频优先服务 | 不绕过路径安全检查；保持外部媒体 `no-store` |
| `user_assets.go` | 把发布包图片和音频补齐到用户目录 | 已存在文件视为用户内容，不能覆盖 |
| `core/` | 数据结构、游戏状态、纯规则、配置和数值计算 | 核心公式优先放这里并写 Go 单元测试 |
| `internal/db/db.go` | 当前数据库表和索引 | 不写历史迁移分支 |
| `internal/db/seeds.go` | 新数据库的默认参考数据 | 稳定 ID，不要只改运行时数据库 |
| `internal/db/queries.go` | 从数据库读取并转换为 `core` 类型 | SQL 字段顺序与结构体扫描必须同步 |
| `internal/db/dating_profiles.go` | 100 位约会对象的名字、国籍和简介 | 与 `seeds.go` 的 ID 一一对应 |
| `services/` | Wails 后端接口、业务编排和类型化响应 | 公共状态修改必须持有 `stateMu` |
| `frontend/src/stores/game.js` | Pinia 全局游戏状态和统一同步入口 | 不在多个组件复制同一份可变状态 |
| `frontend/src/components/Menu/` | 各主系统页面 | 大功能拆出对话框、工具函数或 composable |
| `frontend/src/components/Dialog/` | 场景、购买、存档等对话框 | 保持房车/约会预览交互风格一致 |
| `frontend/src/components/GameList/` | 小游戏界面 | 结算以服务端会话和规则为准 |
| `frontend/src/utils/` | 约会场景、礼物、服装等展示映射 | 与后端规则保持同名、同路径 |
| `frontend/src/composables/useBackgroundMusic.js` | 无缝背景音乐、自动播放解锁和静音偏好 | 保持单 AudioBufferSource、低音量并关闭 AudioContext |
| `frontend/public/images/` | 打入发布包的默认图片源 | 使用当前目录结构，不创建 `v2` 等平行目录 |
| `frontend/public/audio/` | 打入发布包的默认循环音乐 | 默认文件名为 `lifegame-theme.wav` |
| `frontend/public/vendor/` | 生产包本地 Vue/Element Plus ESM | Vite 已将这两个包 external，不能随意删除 |
| `frontend/wailsjs/` | Wails 自动生成的前端绑定 | 后端公开 API 变化后重新生成，不手工维护 |
| `frontend/tests/` | 前端状态和资源契约测试 | 资源路径或布局规则变化时同步更新 |
| `frontend/e2e/` | 连接真实后端的 Playwright 测试 | 必须继续使用临时家目录 |
| `scripts/process_dating_character_sheets.py` | 约会对象角色图处理脚本 | 项目要求保留 |
| `scripts/native_gui_smoke.py` | Linux 原生桌面冒烟测试 | 依赖已构建程序和 AT-SPI |

## 4. 数据源：先判断应该改哪里

同一个界面可能同时使用配置、数据库、Go 规则、前端映射和图片。开发前必须先确定唯一数据源。

| 内容 | 默认来源 | 运行时来源 | 示例 |
| --- | --- | --- | --- |
| 全局可调参数 | `core/config.go` | `~/.lifegame/config.yaml` | 年龄、市场数量、贷款利率、最大免疫力 |
| 参考目录数据 | `internal/db/seeds.go` | `~/.lifegame/lifegame.db` | 商品、公司、股票、房车、疾病、小游戏 |
| 约会人物资料 | `dating_profiles.go` + `seeds.go` | 数据库 `datings` 表 | 名字、职业、礼物、地点、解锁条件 |
| 核心规则与公式 | `core/*.go` | 编译进程序 | 资产、健康、解锁、关系等级、市场波动 |
| 业务操作 | `services/*.go` | 当前内存游戏状态 | 买卖、约会、结婚、结算、存档 |
| 前端状态 | `stores/game.js` | 当前 WebView 内存 | 玩家、游戏、公告、股票时钟 |
| 展示映射 | `frontend/src/utils/` | 编译进前端 | 场景分组、服装路径、礼物场景 |
| 默认图片 | `frontend/public/images/` | `~/.lifegame/images/` 优先 | 人物、房车、古董、场景 |
| 默认音乐 | `frontend/public/audio/` | `~/.lifegame/audio/` 优先 | `lifegame-theme.wav` |

### 4.1 为什么修改默认代码后游戏可能不变化

程序尊重用户已有数据，因此不同内容有不同的更新方式：

- 修改 `core/config.go` 的默认值，只影响以后新创建的 `config.yaml`；已有配置仍按用户文件读取。
- 修改 `internal/db/seeds.go`，只影响以后新创建的数据库；已有数据库不会重新灌入默认数据。
- 修改 `frontend/public/images/`，只改变发布包内置副本；已有 `~/.lifegame/images/` 同名文件不会被覆盖。
- 修改 `frontend/public/audio/`，只改变发布包内置音乐；已有 `~/.lifegame/audio/` 同名文件不会被覆盖。
- 修改 Go 规则或 Vue 代码，需要重新构建/运行新程序。

开发测试默认数据时，应使用临时 `HOME`，或者在确认备份后只删除测试目录中的对应配置、数据库或图片。不要为了看到新数据而直接删除真实存档。

## 5. `~/.lifegame` 运行时目录

程序通过 `os.UserHomeDir()` 获取家目录，所以实际位置由启动进程的 `HOME` 决定。Linux 常见位置为：

```text
/home/<用户名>/.lifegame
```

当前结构：

```text
.lifegame/
├── config.yaml
├── lifegame.db
├── audio/
│   └── lifegame-theme.wav
└── images/
    ├── antiqueinfo/
    ├── carinfo/
    │   ├── cars/
    │   └── car-moments/
    ├── datinginfo/
    │   ├── dating-careers/
    │   ├── dating-scenes/
    │   ├── dating-moments/
    │   └── dating-partner/
    │       ├── female/
    │       └── male/
    └── houseinfo/
        ├── houses/
        └── house-moments/
```

### 5.1 配置

`config.yaml` 只负责 `core.Config` 已定义的字段。增加配置项时必须同时完成：

1. 给 `core.Config` 或子结构增加 YAML 字段。
2. 给 `DefaultConfig` 设置当前默认值。
3. 在 `UpdateConfigValues` 中应用到运行时变量，或直接提供 getter。
4. 给缺失、非法、边界值增加 `core/config_test.go` 测试。
5. 明确现有配置不兼容即可，不添加旧字段别名或迁移分支。

### 5.2 数据库

当前数据库包含：

- 参考数据：`items`、`companies`、`antiques`、`stocks`、`stock_news`、`houses`、`cars`、`datings`、`bank_tasks`、`diseases`、`treats`、`hospitals`、`minigames`
- 存档数据：`saves`、`game_saves`、`user_saves`、`announce_saves`

新数据库会建立当前表、索引并写入默认数据。已有数据库只要求符合当前结构，不进行历史结构转换，也不会自动用 `seeds.go` 覆盖用户修改的参考数据。

直接编辑数据库时：

1. 完全退出游戏，确保 SQLite 连接已关闭。
2. 备份 `lifegame.db`。
3. 修改字段时保持类型、JSON 格式、ID 和资源路径有效。
4. 启动游戏重新加载数据库缓存。

### 5.3 存档版本

当前存档格式版本在 `services/archive.go` 的 `currentSaveVersion` 中维护。加载时要求版本完全相等，不转换旧存档。

如果 `core.Game`、`core.User`、`core.Announce` 的持久化结构发生不兼容变化：

1. 增加 `currentSaveVersion`。
2. 更新存档测试和提示。
3. 不编写历史存档转换器，除非项目所有者改变“不兼容历史版本”的决定。

## 6. 后端开发约定

### 6.1 分层

- 可复用、可测试的公式和状态规则放入 `core`。
- SQLite 表、默认目录数据和读取放入 `internal/db`。
- 一次用户操作的校验、扣费、次数消耗、状态更新和响应组装放入 `services`。
- 前端只决定展示和输入，不自行复制关键成功率、费用或结算公式。

### 6.2 并发与状态

`services.App` 保存当前 `Gameinfo`、`Userinfo`、公告、小游戏会话、约会互动会话和股票时钟。股票定时器与玩家操作可能同时进入后端。

新增公共业务方法必须遵循现有模式：

```go
func (a *App) NewOperation(...) SomeResponse {
    a.stateMu.Lock()
    defer a.stateMu.Unlock()

    if errResp := a.requireGame(); errResp != nil {
        return SomeResponse{Code: -1, Msg: responseMessage(errResp)}
    }
    // 校验并更新状态
}
```

返回 `Userinfo` 或 `Gameinfo` 时使用 `userSnapshot()` / `gameSnapshot()`，不要把正在变化的内部指针直接交给 Wails 序列化。

### 6.3 响应类型

新接口应在 `services/response.go` 或对应业务文件中定义明确响应结构，并包含稳定的 JSON tag。优先使用：

```go
type SomeResponse struct {
    Code     int        `json:"code"`
    Msg      string     `json:"msg"`
    Userinfo *core.User `json:"userinfo,omitempty"`
}
```

`H`/`M` 是仍在使用的旧式动态响应边界。新增功能不要继续扩大动态 `map[string]interface{}`；只有逐步重构既有接口时才处理它们。

后端公开方法、参数或返回类型变化后，运行一次不带 `-skipbindings` 的 Wails `dev` 或 `build`，并检查 `frontend/wailsjs/go/services/App.*` 与 `models.ts`。

### 6.4 一次操作的完整结算

同一次操作中的以下内容应在持有锁期间完成：

1. 启动/游戏状态校验。
2. 参数、对象、次数、资金和关系条件校验。
3. 随机判定。
4. 扣费和机会次数消耗。
5. 所有属性、资产和关系状态更新。
6. 返回脱离内部状态的快照。

不要先返回成功再由前端补扣费用，也不要让按钮重复点击造成重复结算。

## 7. 前端开发约定

### 7.1 状态同步

全局状态以 `useGameStore()` 为入口：

- 创建游戏、推进年度和读取存档：`applyGameData()`
- 只返回玩家变化：`applyUserInfo()`
- 只返回游戏变化：`applyGameInfo()`
- 股票轻量响应：`applyStockUpdate()`

不要在页面保存一份长期独立的 `userInfo` 或 `gameInfo` 副本。对话框可以保存本次展示快照，但结算成功后必须通过 store 的统一入口更新。

### 7.2 组件边界

- 主系统页面放在 `components/Menu/`。
- 可独立打开、关闭或复用的交互放在 `components/Dialog/`。
- 独立小游戏放在 `components/GameList/`。
- 多组件共享的状态流程放在 composable。
- 纯路径、文本和展示映射放在 `frontend/src/utils/`。

如果一个页面同时承担列表、详情、视频式场景、互动、换装和所有请求，应优先按业务边界拆分，而不是继续把逻辑堆入一个组件。

### 7.3 定时器与生命周期

页面级定时器必须在卸载时清除。可使用 `useCleanupTasks()` 管理 timeout/interval。

股票行情是全局游戏进程的一部分，当前由 `GameMain` 上的 `useStockTicker()` 管理，而不是由股市路由页面管理。这样切换页面不会停止行情，也不会创建多个轮询器。

新增定时功能时必须考虑：

- 页面反复进入后是否重复启动。
- 请求是否可能重叠。
- 旧响应是否会覆盖新状态。
- 存档读取、新游戏和推进年度是否需要重置计时状态。
- 组件卸载和应用退出时是否清理。

背景音乐使用 `frontend/public/audio/lifegame-theme.wav`，由 `scripts/generate_background_music.go` 生成，不依赖第三方版权音乐。播放时一次解码为 AudioBuffer，再由单个 BufferSource 设置 `loopStart/loopEnd` 连续循环，避免 HTML 媒体元素每轮重新起播产生停顿。浏览器或 WebView 禁止自动播放时，在玩家第一次点击、按键或点击“播放”按钮后启动。静音偏好保存在 `localStorage`，应用卸载时必须停止音源并关闭 AudioContext。

### 7.4 样式与主题

颜色优先使用 `frontend/src/style-variables.css` 的主题变量。新增界面需要同时检查明暗主题、1024×768 最小窗口、内容滚动和对话框溢出。

房车和约会对象的大图预览采用统一的场景化金边框架、右上角纵向小图标操作、左下角主要信息、右下角状态信息。除非项目所有者提出新设计，不要让同类界面重新分裂成多套布局。

## 8. 当前关键游戏规则

本节用于防止二次开发无意中破坏核心不变量。具体数值仍以代码为准。

### 8.1 年度与结束条件

- `NextTime()` 推进一个游戏年度，刷新市场、公司、房车价格、银行结算、健康和疾病。
- 每年重置打工、小游戏、关系行动、逛街和拍卖等机会次数。
- 净资产不大于 0 或达到最大年龄会结束游戏。
- 净资产由现金、存款、物资、股票、公司、古董、房产和车辆价值减贷款组成。

### 8.2 健康

- 免疫力不会因一次结算归零而立即死亡，会进入连续低免疫抢救期。
- 恢复到安全线会重置危急年度。
- 疾病会随年度发展，医院操作和急诊规则由 `core/health_rules.go` 与 `services/health.go` 共同完成。
- 难度会影响健康负面变化、公司破产率和古董假货概率。

### 8.3 股票和小游戏

- 股票行情每个游戏年度最多更新固定次数，服务端使用 epoch/version 防止旧响应覆盖新行情。
- 小游戏开始、操作、下注、结束和取消以服务端会话为准。
- 游戏记录参与约会对象解锁和银行任务，因此新增小游戏时必须正确写入 ID、类型、游玩次数和胜利次数。

### 8.4 约会系统

当前约会系统的重要规则：

- 数据库共有 100 位对象：女性 ID 1–50，男性 ID 51–100。
- 玩家当前只能看到与玩家性别相反的一侧，共 50 位对象。
- 男性 51–100 的人物图片仍使用 `male/01.webp`–`male/50.webp`，不要用 `male/51.webp`。
- 普通对象满足全部认识条件后自动解锁；配置 `meet_scene` 的对象必须主动前往对应场景认识。
- 当前每侧需要主动前往场景的基础编号是 1、7、16、20、30、40、47。
- 每次约会展示 3 个文字地点，保证至少包含一个对方喜欢的地点。
- 喜欢地点成功率高、收益中等；非偏好地点成功率低，但成功收益更高。
- 每次送礼给出 3 个选项且恰好一个是偏好礼物。偏好礼物容易送出、普通收益；非偏好礼物很难送出、成功时高收益。
- 一次成功约会只允许第一次互动产生属性结算；按钮可以继续点击用于观看，但不能重复获得效果。
- 关系状态依次包括：陌生人、朋友、暧昧中、交往中、恋人、专属恋人、爱人、已婚、前任。
- 可以同时维持多段较低等级关系，但只能有一位“爱人”级承诺对象和一位配偶。
- 结婚需要累计约会 20 次、好感度至少 90，并支付婚礼费用；只能和一位对象结婚。
- 未婚关系使用分手，婚姻使用离婚。结束后该对象成为前任，玩家可以和其他对象发展新关系。
- 婚后可以继续约会，并有配偶洗澡互动。

关系门槛、唯一关系约束和解锁条件以 `core/dating_rules.go`、`services/dating.go` 为准。前端只负责呈现可用操作。

### 8.5 约会对象换装权限

七种造型不是随机选择，由玩家在关系允许时指定：

| 造型 | 资源目录 | 最低关系 |
| --- | --- | --- |
| 职业装 | 人物性别基础目录 | 朋友 |
| 居家装 | `homewear/` | 朋友 |
| 旗袍/国风 | `qipao/` | 暧昧中 |
| Cosplay | `cosplay/` | 交往中 |
| 泳装 | `swimwear/` | 交往中 |
| 睡衣 | `sleepwear/` | 恋人 |
| 情趣睡衣 | `romantic/` | 专属恋人，且玩家和对象均已成年 |

换装权限定义在 `frontend/src/utils/datingOutfits.js`，后端互动还会校验对应规则。修改关系门槛时必须同时更新前后端测试。

## 9. 媒体资源与索引

### 9.1 两层资源

```text
frontend/public/images/...   发布包默认资源
~/.lifegame/images/...       用户运行时资源，优先级更高
frontend/public/audio/...    发布包默认音乐
~/.lifegame/audio/...        用户运行时音乐，优先级更高
```

Vite 构建时把 `frontend/public` 媒体原样复制到 `frontend/dist`，Go 再把 `frontend/dist` 嵌入可执行程序。启动时 `user_assets.go` 只释放缺失文件。请求 `/images/...` 或 `/audio/...` 时，`main.go` 先查找 `.lifegame` 对应目录，不存在才读取内置资源。

因此复制整个项目或单独移动可执行程序后图片仍能显示：默认图片已经嵌入可执行程序；同时，每个用户仍可以在 `.lifegame` 中替换同一路径图片。

### 9.2 命名规则

- 使用 WebP。
- 文件名使用小写英文和连字符，目录名保持当前英文结构。
- 有数据库 ID 的资产以两位 ID 开头，如 `01-economy-sedan.webp`。
- 不创建 `v2`、`new`、`old` 等平行资源目录。
- 数据库和前端引用必须使用 `/images/...` URL，不能写开发机器绝对路径。
- 同类资源的画布比例、主体大小和留白参考同目录现有图片。
- 房屋、车辆和需要镶嵌场景的人物图保持透明背景；职业环境、约会地点和交付场景使用完整背景图。
- 资源删除前不能只依赖简单文本搜索；还要检查数据库路径、动态拼接规则和契约测试。

### 9.3 约会人物图片索引

每个性别有 50 位人物、每人 7 张，共 350 张：

```text
dating-partner/female/01.webp
dating-partner/female/homewear/01.webp
dating-partner/female/qipao/01.webp
dating-partner/female/cosplay/01.webp
dating-partner/female/swimwear/01.webp
dating-partner/female/sleepwear/01.webp
dating-partner/female/romantic/01.webp
```

男性使用完全相同的子目录结构。`datingOutfitImage()` 根据基础 URL 动态插入造型目录，所以任何一个状态缺图都会只影响该状态。

职业环境图按人物在各自性别组内的基础编号索引：

```text
dating-careers/01.webp ... dating-careers/50.webp
```

### 9.4 房车、古董和场景

- 车辆：`carinfo/cars/01-...webp` 至 `50-...webp`
- 车辆交付场景：`carinfo/car-moments/car-showroom.webp`
- 房产：`houseinfo/houses/01-...webp` 至 `50-...webp`
- 房产交付场景：`houseinfo/house-moments/home-handover.webp`
- 古董：`antiqueinfo/01-...webp` 至 `50-...webp`，另有 `default.webp`
- 约会地点背景：`datinginfo/dating-scenes/`
- 约会短片式场景：`datinginfo/dating-moments/`

数据库中的 `img` 字段必须与文件名完全一致。

### 9.5 替换图片的正确步骤

只替换当前用户图片：

1. 完全退出游戏。
2. 找到数据库或界面 URL 对应的相对路径。
3. 用同名 WebP 覆盖 `~/.lifegame/images/` 下的文件。
4. 确认修改的是当前人物、当前性别和当前造型。
5. 启动最新构建程序验证。

修改以后所有新安装的默认图片：

1. 替换 `frontend/public/images/` 下的源文件。
2. 运行前端资源契约测试和生产构建。
3. 重新构建 Wails 程序。
4. 对已有用户，仍需手动替换外部同名文件，或在备份后删除该单个外部文件，让程序下次启动补齐。程序不会自动覆盖用户图片。

### 9.6 替换背景音乐

- 当前固定地址为 `/audio/lifegame-theme.wav`，外部文件为 `~/.lifegame/audio/lifegame-theme.wav`。
- 完全退出游戏后，用同名 PCM WAV 文件覆盖外部文件并重新启动即可切换。
- 也可以在游戏中先静音、替换文件、再点击开启；播放器会用 `no-store` 重新读取，但完全重启验证最可靠。
- AudioBuffer 可以消除播放器重新起播的空隙，但无法修复音乐文件本身首尾节奏、和声或波形不连续。自定义音乐必须先制作成自然循环段。
- 修改内置默认音乐后运行 `go run ./scripts/generate_background_music.go`（使用生成器时）或直接替换默认 WAV，再执行前端契约测试和 Wails 构建。

## 10. 常见扩展流程

### 10.1 修改游戏数值

1. 先确认数值属于配置、数据库目录数据还是核心公式。
2. 修改唯一数据源，不在 Vue 中复制后端公式。
3. 补充边界测试、概率测试或固定随机数测试。
4. 检查简单/普通/困难三种难度。
5. 检查资产、免疫力、名声、费用和机会次数是否保持合法范围。

### 10.2 添加后端操作

1. 在 `core` 写可复用规则和类型。
2. 在合适的 `services/*.go` 增加 `App` 公共方法。
3. 持有 `stateMu`，调用 `requireGame()`，完整校验和结算。
4. 返回显式响应结构和状态快照。
5. 写 Go 规则/服务测试。
6. 不带 `-skipbindings` 生成 Wails 绑定。
7. 前端通过生成的 `App.js` 调用并同步 Pinia。
8. 增加前端契约测试或 E2E。

### 10.3 添加主菜单功能

1. 在 `components/Menu/` 新建页面。
2. 在 `router/index.js` 添加懒加载路由。
3. 在 `GameMenu.vue` 添加入口。
4. 如需全局状态，扩展 Pinia；页面临时状态留在组件。
5. 如需后端数据，按上一节增加类型化 API。
6. 检查明暗主题、最小窗口和路由反复进入。

### 10.4 修改或添加约会对象

修改已有对象默认资料时，需要同步检查：

1. `internal/db/dating_profiles.go`：姓名、国籍、简介。
2. `internal/db/seeds.go`：年龄、职业、费用、认识条件、喜欢礼物、喜欢地点、相遇场景。
3. `frontend/public/images/datinginfo/dating-partner/<gender>/`：七套同编号人物图。
4. `dating-careers/<基础编号>.webp`：职业环境。
5. 约会场景和礼物是否已在前后端映射中识别。
6. 女性/男性 ID 和基础图片编号是否正确。
7. 解锁、约会、送礼、关系和换装契约测试。

如果要把人数扩展到每侧 50 人以上，不能只追加数据库行；还必须修改职业场景数量、基础 ID 归一化、前端资源契约、图片集合和相关测试。

### 10.5 添加约会地点或场景

后端成功率和收益规则位于 `services/dating_scenes.go`，前端文字选项与背景映射位于 `frontend/src/utils/datingScenes.js`。新增地点时二者都必须识别相同名称。

保持当前玩法：3 个文字选项、至少 1 个偏好地点，只有结算成功后进入短片式场景。

### 10.6 添加礼物

1. 在约会对象数据库偏好列表中按需加入名称。
2. 在 `services/dating_gifts.go` 配置价格档位、成功率和效果。
3. 在 `frontend/src/utils/datingGifts.js` 配置展示场景类别。
4. 保持选项页面不显示价格和效果，但后端仍应返回并结算真实成本。
5. 保持“偏好礼物稳妥普通收益、非偏好礼物高风险高收益”的核心路线。

### 10.7 添加换装或互动

1. 确定关系门槛、是否仅限成年和每次结算次数。
2. 同步修改 `datingOutfits.js`、`services/dating_interactions.go` 和相关响应类型。
3. 给女性和男性每位对象补齐对应目录资源。
4. 在 `DialogSpouseMoment.vue`/约会页面中增加按钮和场景反馈。
5. 保持一次成功约会只首次互动结算，后续点击只展示、不重复加属性。

### 10.8 添加房产、车辆或古董

1. 在 `internal/db/seeds.go` 使用新稳定 ID 添加数据。
2. 添加数据库 `img` 路径对应的 WebP。
3. 房产/车辆检查价格波动、健康/名声加成、购买、出售和总资产。
4. 古董检查拍卖等级、真假、稀有度、鉴定和默认图片回退。
5. 更新资源契约测试和相关页面测试。

### 10.9 添加小游戏

1. 在 `seedMiniGames` 增加稳定英文 ID、中文名、类型、难度、费用、目标和奖励。
2. 在 `components/GameList/` 实现组件。
3. 在娱乐页面注册组件映射。
4. 使用后端开始/操作/结束/取消会话，不让前端自报胜负绕过校验。
5. 确认最短运行时间、FPS/定时器、重复结算、赌注和次数消耗。
6. 确认记录能正确参与约会解锁与银行任务。

## 11. 测试与构建

### 11.1 日常快速验证

```bash
go test -count=1 ./...
go vet ./...

cd frontend
npm test
npm run build
```

### 11.2 真实后端 E2E

```bash
cd frontend
npm run test:e2e:install
npm run test:e2e
```

E2E 使用临时家目录。新增测试不能读取或修改开发者真实的 `~/.lifegame`。

### 11.3 原生桌面验证

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 build -tags webkit2_41
cd frontend
npm run test:native
```

原生验证重点包括：启动、最小窗口、系统 WebView、图片服务、原生对话框和 AT-SPI 可访问性。

### 11.4 发布前手工流程

至少完整游玩以下链路：

1. 新建简单、普通、困难游戏。
2. 推进年度并触发市场、股票、公司、银行、健康和疾病。
3. 每类小游戏至少开始、胜利/失败、取消一次。
4. 买卖商品、股票、房产、车辆和古董。
5. 解锁普通约会对象和场景相遇对象。
6. 测试偏好/非偏好地点、礼物、首次互动、七套换装权限。
7. 测试恋爱、爱人、结婚、婚后互动、分手、离婚和新关系。
8. 保存、读取和删除当前版本存档。
9. 达成一种正常结局，并检查结束评价比例和布局。
10. 用 `~/.lifegame` 替换一张人物、房产、车辆图片和背景音乐，重启确认外部资源生效。

## 12. 常见故障定位

| 现象 | 优先检查 |
| --- | --- |
| 修改人物图但界面不变 | 当前造型子目录、同名 `.webp`、人物基础编号、完整重启、是否启动旧程序 |
| 修改 `frontend/public/images` 不生效 | `~/.lifegame` 已有同名文件会覆盖内置图；手动替换或只删除该外部文件 |
| 替换背景音乐不生效 | 确认路径为 `~/.lifegame/audio/lifegame-theme.wav`、格式为 PCM WAV、已重新启动或静音后重新开启 |
| 音乐循环处有停顿/跳变 | 使用 AudioBuffer 播放后仍跳变通常是文件首尾不连续；重新剪辑节拍、和声与波形边界 |
| 修改 `seeds.go` 后数据不变 | 已有数据库不会重新 seed；使用临时 HOME 创建新数据库，或按要求编辑当前数据库 |
| 修改默认配置后数值不变 | 已有 `config.yaml` 优先；检查实际运行 HOME 和配置内容 |
| 前端找不到新后端方法 | 重新运行不带 `-skipbindings` 的 Wails dev/build，检查导入路径 |
| 返回字段在前端是空值 | Go JSON tag、Wails 绑定、响应字段名和 Pinia 同步入口 |
| 切换页面后速度越来越快 | 重复 interval/RAF/事件监听，检查卸载清理和请求重叠 |
| 股票回退为旧行情 | epoch/version 或整包响应覆盖轻量股票响应 |
| 操作可以重复获得奖励 | 后端会话/幂等状态缺失，只在前端禁用按钮不够 |
| 程序启动失败 | `config.yaml` 格式、当前数据库结构、JSON 字段、文件权限和启动状态错误 |
| 测试污染真实存档 | 测试未设置临时 HOME；立即停止并检查 E2E 启动脚本 |
| 发布包缺图 | 数据库/动态路径与 `frontend/public/images` 不一致，运行资源契约测试 |

## 13. 已知技术债观察

以下内容仅作为后续需求定位线索，不代表开发者或 AI 可以自行实施：

- 部分旧服务仍返回 `H`/`M` 动态 map，可以在项目所有者要求时逐步改成类型化响应。
- `MenuDating.vue` 功能较多，可按列表、关系操作和场景控制继续拆分，但拆分时必须保持现有布局与交互。
- 约会场景在后端结算和前端展示各有映射，新增地点时容易漏改，可在明确重构需求下建立共享生成流程或更强契约测试。
- 图片占发布包主要体积。若继续优化包体积，应先测量各资源类别，再决定有损压缩质量、尺寸或按需资源包，不能直接删图。
- 原生桌面测试依赖 Linux 桌面环境，CI 的浏览器 E2E 不能完全替代发布前原生冒烟测试。

## 14. 开发完成标准

一项需求只有同时满足以下条件才算完成：

- 实现内容与项目所有者的需求一致，没有额外改变玩法。
- 后端校验、扣费、机会次数和属性结算完整且不可重复利用。
- 前端正常同步 Pinia，不出现旧响应覆盖新状态。
- 新增资源存在、路径正确、命名符合当前目录结构。
- 明暗主题和 1024×768 最小窗口可用。
- 相关 Go 测试、前端测试、构建通过；高影响功能完成 E2E 或原生验证。
- 没有破坏当前数据库和当前存档格式；如格式明确变化，已更新版本并接受不兼容。
- README 或本手册需要更新时已经同步。
- 交付说明列出修改文件、验证命令和仍需项目所有者决定的事项。

## 15. 给后续 AI 的工作入口

接到新需求后，按下面顺序工作：

1. 读取 `README.md` 和完整的 `DEVELOPMENT_GUIDE.md`。
2. 把项目所有者最新消息视为当前需求，提取明确的功能、界面、数值和资源要求。
3. 使用文件搜索定位现有入口，不根据文件名猜测实现。
4. 判断改动属于配置、数据库、核心规则、服务、前端状态、组件还是资源；必要时跨层同步。
5. 先理解现有测试和当前不变量，再修改。
6. 遇到会明显改变玩法的歧义，向项目所有者说明具体选项和影响；不要自行决定。
7. 不为旧结构增加兼容代码，不创建第二套数据源或 `v2` 资源目录。
8. 不覆盖用户真实 `.lifegame`；测试使用临时家目录。
9. 完成后运行与影响范围匹配的测试，并说明实际结果。
10. 如果本次改变了架构、规则、目录、命令或扩展流程，同步维护本手册。

项目所有者会提出下一步需要优化或添加的功能。本手册负责提供上下文和边界，不替代项目所有者的需求。
