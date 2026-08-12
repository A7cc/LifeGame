package services

import (
	"LifeGame/core"
	"context"
	cryptorand "crypto/rand"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type AppStartupStatus struct {
	Ready     bool   `json:"ready"`
	Stage     string `json:"stage"`
	Error     string `json:"error"`
	UpdatedAt int64  `json:"updatedAt"`
}

var runtimeIDFallback atomic.Uint64

// App 是 Wails 暴露给前端的应用桥接层。
//
// 前端只绑定这个对象，具体业务按文件拆分在同一个 package 中，避免
// 改动 Wails 生成的前端调用名。
type App struct {
	// stateMu 串行化所有会读取或修改游戏内存状态的 Wails 调用。
	// 前端定时任务（例如股票行情）和用户操作可能来自不同 goroutine，
	// 因此不能依赖页面侧的 updating 标记保证一致性。
	stateMu                  sync.Mutex
	ctx                      context.Context
	Gameinfo                 *core.Game
	Userinfo                 *core.User
	Announce                 *core.Announce
	MiniGameSession          *core.MiniGameSession
	pendingDatingInteraction *datingInteractionSession
	startupStatus            AppStartupStatus
	stockEpoch               string
	stockVersion             uint64
	stockUpdatedAt           int64
	randomRoll               func() float64
}

func NewApp() *App {
	return &App{startupStatus: AppStartupStatus{Stage: "正在启动", UpdatedAt: time.Now().UnixMilli()}}
}

func newRuntimeID() string {
	data := make([]byte, 16)
	if _, err := cryptorand.Read(data); err == nil {
		return fmt.Sprintf("%x", data)
	}
	return fmt.Sprintf("fallback-%d-%d", time.Now().UnixNano(), runtimeIDFallback.Add(1))
}

func (a *App) resetStockClock() {
	a.stockEpoch = newRuntimeID()
	a.stockVersion = 0
	a.stockUpdatedAt = 0
}

// userSnapshot/gameSnapshot 返回脱离 App 可变状态的深拷贝。
// Wails 会在方法返回后序列化响应；如果直接返回内部指针，下一次调用可能
// 在序列化期间修改同一对象，绕过 stateMu 的保护。
func (a *App) userSnapshot() *core.User {
	if a.Userinfo == nil {
		return nil
	}
	data, err := json.Marshal(a.Userinfo)
	if err != nil {
		copy := *a.Userinfo
		return &copy
	}
	var snapshot core.User
	if err := json.Unmarshal(data, &snapshot); err != nil {
		copy := *a.Userinfo
		return &copy
	}
	return &snapshot
}

func (a *App) gameSnapshot() *core.Game {
	if a.Gameinfo == nil {
		return nil
	}
	data, err := json.Marshal(a.Gameinfo)
	if err != nil {
		copy := *a.Gameinfo
		return &copy
	}
	var snapshot core.Game
	if err := json.Unmarshal(data, &snapshot); err != nil {
		copy := *a.Gameinfo
		return &copy
	}
	return &snapshot
}
