package services

import (
	"LifeGame/core"
	"LifeGame/internal/db"
	"context"
	"fmt"
	"log"
	"time"
)

// Startup 在应用程序启动时被调用。
func (a *App) Startup(ctx context.Context) {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	a.ctx = ctx
	if err := a.initializeLocked(); err != nil {
		log.Printf("应用初始化失败: %v", err)
	}
}

// initializeLocked 在持有 stateMu 时运行，启动和手动重试共用同一流程。
func (a *App) initializeLocked() error {
	a.startupStatus = AppStartupStatus{Stage: "正在加载配置", UpdatedAt: time.Now().UnixMilli()}
	if db.GetDB() != nil {
		if err := db.CloseDB(); err != nil {
			return a.failStartupLocked("关闭现有数据库连接", err)
		}
	}
	if err := core.InitConfig(); err != nil {
		return a.failStartupLocked("加载配置", err)
	}

	core.UpdateConfigValues()

	a.startupStatus = AppStartupStatus{Stage: "正在连接数据库", UpdatedAt: time.Now().UnixMilli()}
	dbPath, err := core.GetDBPath()
	if err != nil {
		return a.failStartupLocked("获取数据库路径", err)
	}

	if _, err = db.InitDB(dbPath); err != nil {
		return a.failStartupLocked("初始化数据库", err)
	}

	a.startupStatus = AppStartupStatus{Stage: "正在加载游戏数据", UpdatedAt: time.Now().UnixMilli()}
	if err := db.LoadAllGameData(db.GetDB()); err != nil {
		db.CloseDB()
		return a.failStartupLocked("加载游戏数据", err)
	}

	a.startupStatus = AppStartupStatus{Ready: true, Stage: "初始化完成", UpdatedAt: time.Now().UnixMilli()}
	return nil
}

func (a *App) failStartupLocked(stage string, err error) error {
	wrapped := fmt.Errorf("%s失败: %w", stage, err)
	a.startupStatus = AppStartupStatus{
		Ready:     false,
		Stage:     stage,
		Error:     wrapped.Error(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	return wrapped
}

func (a *App) GetStartupStatus() StartupResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()
	return StartupResponse{Code: 200, Status: a.startupStatus}
}

func (a *App) RetryStartup() StartupResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()
	if a.startupStatus.Ready {
		return StartupResponse{Code: 200, Msg: "应用已经初始化", Status: a.startupStatus}
	}
	if err := a.initializeLocked(); err != nil {
		return StartupResponse{Code: -1, Msg: err.Error(), Status: a.startupStatus}
	}
	return StartupResponse{Code: 200, Msg: "初始化成功", Status: a.startupStatus}
}

// Shutdown 在应用程序关闭时被调用。
func (a *App) Shutdown(ctx context.Context) {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()
	if err := db.CloseDB(); err != nil {
		log.Printf("关闭数据库失败: %v", err)
	}
}

func (a *App) DomReady(ctx context.Context) {
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}
