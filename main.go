package main

import (
	"LifeGame/core"
	"LifeGame/services"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// 使用受版本控制的前端图标，避免干净仓库因 build/ 被忽略而无法编译。
//
//go:embed frontend/src/assets/images/icon-4096.png
var icon []byte

const fallbackAntiqueImage = "/images/antiqueinfo/default.webp"

func newAssetHandler(packagedAssets fs.FS) http.Handler {
	packagedServer := http.FileServer(http.FS(packagedAssets))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/images/antiqueinfo/") {
			if r.URL.Path != fallbackAntiqueImage && !packagedAssetExists(packagedAssets, r.URL.Path) {
				fallbackRequest := r.Clone(r.Context())
				fallbackRequest.URL.Path = fallbackAntiqueImage
				packagedServer.ServeHTTP(w, fallbackRequest)
				return
			}
		}
		packagedServer.ServeHTTP(w, r)
	})
}

type externalAssetRoute struct {
	urlPrefix string
	directory string
}

// externalAssetsFirst 在 Wails 尝试读取内置 Assets 之前检查用户资源。
// Wails 的 Handler 只会在内置文件不存在时调用，因此外部同名覆盖必须放在
// Middleware 中，不能仅放在 Handler 中。
func externalAssetsFirst(routes ...externalAssetRoute) assetserver.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				for _, route := range routes {
					if serveExternalAsset(w, r, route.directory, route.urlPrefix, r.URL.Path) {
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func newAssetServerOptions(packagedAssets fs.FS, imagesDir, audioDir string) *assetserver.Options {
	return &assetserver.Options{
		Assets:  packagedAssets,
		Handler: newAssetHandler(packagedAssets),
		Middleware: externalAssetsFirst(
			externalAssetRoute{urlPrefix: "/images/", directory: imagesDir},
			externalAssetRoute{urlPrefix: "/audio/", directory: audioDir},
		),
	}
}

func serveExternalAsset(w http.ResponseWriter, r *http.Request, assetDir, urlPrefix, requestPath string) bool {
	if assetDir == "" || urlPrefix == "" || !strings.HasPrefix(requestPath, urlPrefix) {
		return false
	}
	relativeURLPath := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(requestPath, urlPrefix)), "/")
	localPath := filepath.Join(assetDir, filepath.FromSlash(relativeURLPath))
	relativeLocalPath, err := filepath.Rel(assetDir, localPath)
	if err != nil || relativeLocalPath == "." || relativeLocalPath == ".." || strings.HasPrefix(relativeLocalPath, ".."+string(filepath.Separator)) {
		return false
	}
	if info, err := os.Stat(localPath); err != nil || info.IsDir() {
		return false
	}
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, localPath)
	return true
}

func packagedAssetExists(packagedAssets fs.FS, requestPath string) bool {
	assetPath := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(requestPath, "/")), "/")
	info, err := fs.Stat(packagedAssets, assetPath)
	return err == nil && !info.IsDir()
}

func main() {
	// 创建应用结构的实例
	app := services.NewApp()
	packagedAssets, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		println("Error:", err.Error())
		return
	}
	imagesDir, err := core.GetImagesDir()
	if err != nil {
		log.Printf("无法使用外部图片目录，将读取内置资源: %v", err)
		imagesDir = ""
	} else if extracted, extractErr := ensureExternalImages(packagedAssets, imagesDir); extractErr != nil {
		log.Printf("释放默认图片未完整完成，缺失文件将读取内置资源: %v", extractErr)
	} else if extracted > 0 {
		log.Printf("已释放 %d 个默认图片文件到 %s", extracted, imagesDir)
	}
	audioDir, err := core.GetAudioDir()
	if err != nil {
		log.Printf("无法使用外部音频目录，将读取内置资源: %v", err)
		audioDir = ""
	} else if extracted, extractErr := ensureExternalAudio(packagedAssets, audioDir); extractErr != nil {
		log.Printf("释放默认音频未完整完成，缺失文件将读取内置资源: %v", extractErr)
	} else if extracted > 0 {
		log.Printf("已释放 %d 个默认音频文件到 %s", extracted, audioDir)
	}

	// 创建带有选项的应用程序
	err = wails.Run(&options.App{
		Title: "LifeGame",
		// 宽
		Width: 1024,
		// 高
		Height: 768,
		// 最小宽
		MinWidth: 1024,
		// 最小高
		MinHeight: 768,
		// 禁止用户调整窗口大小，这里为 false 表示允许缩放
		DisableResize: false,
		// 是否全屏，设为 false
		Fullscreen: false,
		// 是否启动时隐藏窗口，设为 false
		StartHidden: false,
		// Windows/macOS 使用前端自定义窗口按钮；Linux 保留系统标题栏。
		Frameless: false,
		// 是否在关闭窗口时只是隐藏（不退出）
		HideWindowOnClose: false,
		// AssetServer为应用程序配置资产，前端资源设置
		AssetServer: newAssetServerOptions(packagedAssets, imagesDir, audioDir),
		// Wails 的菜单栏设置。它用来定义 桌面应用的原生菜单栏（比如 macOS 顶部菜单栏、Windows 应用窗口菜单）。当前设置是 nil，表示不显示任何菜单。
		Menu: nil,
		// 窗口背景颜色
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		// 创建窗口并即将开始加载前端资源时的回调
		OnStartup: app.Startup,
		// 应用程序即将关闭时的回调
		OnShutdown: app.Shutdown,
		// 当 DOM 准备就绪时调用
		OnDomReady: app.DomReady,
		// 窗口关闭前调用
		OnBeforeClose: app.BeforeClose,
		// 指定一个自定义的日志处理器，当设置为 nil，Wails 会使用默认的日志输出方式
		Logger: nil,
		// 日志等级
		LogLevel: logger.DEBUG,
		// 生产环境日志等级
		LogLevelProduction: logger.ERROR,
		// 用来指定应用程序的窗口在启动时的初始状态，例如最小化（Minimised）、最大化（Maximised）、正常（Normal）、全屏（Fullscreen）
		WindowStartState: options.Normal,
		// 在 Wails 中，Bind 是一个关键配置选项，通常用来绑定应用程序中的变量或对象，使得这些变量可以在前端（JavaScript 或 TypeScript）与后端（Go）之间进行交互。简单来说，Bind 使得前端能够调用 Go 后端的函数或获取 Go 后端的数据。
		Bind: []interface{}{
			// app是一个用于操作界面的后端对象，或者包含应用程序业务逻辑的对象，将 app 绑定到前端
			app,
		},
		// Mac平台特定选项
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				// 设置标题栏是否看起来是透明的。如果为 true，标题栏将显得透明。如果为 false，标题栏会是实心的背景色
				TitlebarAppearsTransparent: false,
				// 控制是否隐藏窗口的标题文本。如果为 true，标题文本会被隐藏
				HideTitle: false,
				// 设置是否隐藏整个标题栏。如果为 true，标题栏将被隐藏，通常这种配置与无边框窗口结合使用
				HideTitleBar: false,
				// 如果为 true，窗口内容区域会扩展到整个窗口，覆盖掉标题栏部分。这个选项通常用于无标题
				FullSizeContent: true,
				// 控制是否使用 macOS 的工具栏。如果为 true，窗口将显示一个工具栏
				UseToolbar: false,
				// 设置是否隐藏工具栏和窗口内容之间的分隔符。 true 表示隐藏分隔符，通常用于更简洁的界面设计
				HideToolbarSeparator: true,
			},
			// 这个选项设置应用程序的主题。 mac.DefaultAppearance 使应用程序遵循系统默认的外观配置（自动选择浅色模式或深色模式）。
			Appearance: mac.DefaultAppearance,
			// 这个选项控制 Webview 是否透明。如果设置为 true，Webview 内容区域将变得透明，适合在需要透明背景的窗口中使用。
			WebviewIsTransparent: true,
			// 这个选项使得窗口背景半透明。设置为 true 后，窗口会有一种半透明的效果，允许背景内容部分透过来。
			WindowIsTranslucent: false,
			// 允许窗口进入 macOS 全屏模式
			Preferences: &mac.Preferences{
				FullscreenEnabled: mac.Enabled,
			},
			// 这个配置定义了应用的 "关于" 信息
			About: &mac.AboutInfo{
				Title:   "LifeGame",
				Message: "https://github.com/a7cc",
				Icon:    icon,
			},
		},
		// Linux平台特定选项
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: "LifeGame",
		},
		// Windows平台特定选项
		Windows: &windows.Options{
			Theme: windows.SystemDefault,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
