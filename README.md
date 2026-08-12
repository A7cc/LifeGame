# LifeGame

> 这个游戏算是完成我小时候的开发游戏梦吧！

LifeGame 是一款使用 Wails 2、Go 和 Vue 3 开发的桌面人生模拟游戏。玩家可以伴随可静音的内置背景音乐，进行市场交易、股票投资、创业、银行业务、古董收藏、购置房车、医疗、约会以及多种小游戏，并通过年度推进完成一局人生。

项目当前只维护一套代码、数据库、配置和媒体目录结构，不提供历史数据库、历史配置或旧资源路径的兼容与迁移。

## 技术栈

- 桌面容器：Wails 2.10.2
- 后端：Go 1.25、SQLite（`modernc.org/sqlite`）、YAML
- 前端：Vue 3、Vite 3、Pinia、Vue Router、Element Plus、ECharts
- 测试：Go Test、Node Test、Playwright、原生 AT-SPI 冒烟测试

## 环境要求

- Go 1.25+
- Node.js 20+，CI 使用 Node.js 22
- Linux 桌面构建需要 GTK3 和 WebKitGTK 4.1

Ubuntu/Debian 可安装：

```bash
sudo apt-get update
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.1-dev
```

## 快速开始

安装前端依赖：

```bash
cd frontend
npm ci
cd ..
```

启动开发模式：

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 dev -tags webkit2_41
```

构建 Linux 可执行程序：

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 build -tags webkit2_41
```

产物默认位于：

```text
build/bin/LifeGame
```

如果本次没有修改 Go 对外方法、参数或响应结构，可以使用更快的重复构建：

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 build -tags webkit2_41 -skipbindings
```

修改过 `services.App` 的公开方法或响应类型时，至少运行一次不带 `-skipbindings` 的 `dev` 或 `build`，更新 `frontend/wailsjs/` 绑定。

## 测试

后端测试和静态检查：

```bash
go test -count=1 ./...
go vet ./...
```

前端状态/契约测试和生产构建：

```bash
cd frontend
npm test
npm run build
```

首次运行浏览器 E2E 前安装 Chromium：

```bash
cd frontend
npm run test:e2e:install
npm run test:e2e
```

E2E 会连接真实 Wails 后端，并使用临时家目录，避免改动开发者真实的 `~/.lifegame`。原生 Linux 桌面冒烟测试需要先构建应用：

```bash
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 build -tags webkit2_41
cd frontend
npm run test:native
```

更详细的前端测试说明见 [`frontend/README.md`](frontend/README.md)。

## 用户数据与外部资源

程序以当前用户家目录下的 `.lifegame` 作为唯一运行时数据目录：

```text
~/.lifegame/
├── config.yaml       # 当前游戏参数
├── lifegame.db       # 当前参考数据和存档
├── audio/
│   └── lifegame-theme.wav # 可替换的循环背景音乐
└── images/           # 可由用户直接替换的运行时图片
    ├── datinginfo/
    │   ├── dating-careers/
    │   ├── dating-scenes/
    │   ├── dating-moments/
    │   └── dating-partner/
    │       ├── female/
    │       └── male/
    ├── antiqueinfo/
    ├── carinfo/
    │   ├── cars/
    │   └── car-moments/
    └── houseinfo/
        ├── houses/
        └── house-moments/
```

运行时规则：

- 首次启动会创建默认配置和当前数据库，并将内置图片、音乐分别释放到 `~/.lifegame/images/` 和 `~/.lifegame/audio/`。
- 后续启动优先读取 `.lifegame` 外部媒体；已有文件不会被内置资源覆盖，缺少的文件才会补齐。
- `/images/...`、`/audio/...` 与 `.lifegame` 中对应目录后面的相对路径完全一致。
- 配置和数据库在启动时加载，修改后必须完全退出并重新启动游戏。
- 编辑数据库前应先退出游戏并备份；程序不转换历史数据库或历史存档。
- 删除整个 `.lifegame` 会丢失配置、自定义图片和存档。再次启动虽然会生成默认内容，但不会恢复已删除的个人存档。

例如：

```text
界面地址：/images/datinginfo/dating-partner/female/11.webp
运行文件：~/.lifegame/images/datinginfo/dating-partner/female/11.webp
内置来源：frontend/public/images/datinginfo/dating-partner/female/11.webp

音乐地址：/audio/lifegame-theme.wav
运行文件：~/.lifegame/audio/lifegame-theme.wav
内置来源：frontend/public/audio/lifegame-theme.wav
```

更换背景音乐时，完全退出游戏，用同名的 PCM WAV 文件覆盖 `~/.lifegame/audio/lifegame-theme.wav` 后重新启动。音乐会使用 AudioBuffer 首尾无缝循环；自定义音乐自身也应剪辑成首尾自然衔接的循环段。

约会对象共有七种造型。基础目录是职业装，其余造型必须替换同名子目录文件：

```text
female/11.webp               # 职业装/默认形象
female/homewear/11.webp      # 居家装
female/qipao/11.webp         # 旗袍/国风
female/cosplay/11.webp       # Cosplay
female/swimwear/11.webp      # 泳装
female/sleepwear/11.webp     # 睡衣
female/romantic/11.webp      # 情趣睡衣
```

如果图片修改后界面没有变化，请依次确认：文件名和扩展名未改变、修改的是当前造型目录、游戏已完全退出、启动的是本项目最新构建的 `build/bin/LifeGame`，而不是其他目录中的旧程序。

## 项目结构

```text
LifeGame/
├── main.go                  # Wails 启动、窗口和外部媒体服务
├── user_assets.go           # 内置图片/音乐首次释放
├── core/                    # 游戏状态、纯规则、配置和数值计算
├── internal/db/             # SQLite 表结构、默认数据和查询
├── services/                # Wails 对外业务接口与类型化响应
├── frontend/
│   ├── src/components/      # 页面、对话框和小游戏组件
│   ├── src/stores/          # Pinia 游戏状态
│   ├── src/utils/           # 前端场景、礼物、换装等映射
│   ├── public/images/       # 发布包内置图片源
│   ├── public/audio/        # 发布包内置背景音乐
│   └── tests/               # 前端契约测试
├── scripts/                 # 图片处理和原生测试辅助脚本
└── wails.json               # Wails 项目配置
```

## 二次开发

继续优化、添加功能或交给 AI 修改前，请先完整阅读 [`DEVELOPMENT_GUIDE.md`](DEVELOPMENT_GUIDE.md)。该手册记录了当前架构、游戏规则边界、数据源、图片索引、扩展流程、验证要求以及 AI 协作约定。

项目功能和游戏设计以项目所有者提出的需求为准。开发者或 AI 可以指出问题和方案，但不应在没有需求确认的情况下自行改变玩法、数值、关系规则或视觉方向。

## 持续集成

GitHub Actions 会执行 Go 测试、`go vet`、前端契约测试、生产构建和真实后端 Playwright E2E。原生 AT-SPI 测试作为 Linux 桌面发布前检查保留。

## 更新

v1.0.0：最初版本
v1.1.0：完善了很多，包括约会对象、房车、音乐背景等，以及修复一些BUG

## License

MIT License
