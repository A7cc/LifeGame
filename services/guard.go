package services

// requireStartup 阻止初始化失败后的业务调用继续使用半初始化状态。
// 测试中直接构造的 App 没有启动状态，保留该方式以便验证纯业务规则。
func (a *App) requireStartup() M {
	if a == nil {
		return M{"code": -1, "msg": "应用未初始化"}
	}
	status := a.startupStatus
	if status.Ready || (status.Stage == "" && status.Error == "") {
		return nil
	}
	if status.Error != "" {
		return M{"code": -1, "msg": status.Error}
	}
	return M{"code": -1, "msg": "应用尚未完成初始化：" + status.Stage}
}

// 统一检查游戏是否开始
func (a *App) requireGame() M {
	if errResp := a.requireStartup(); errResp != nil {
		return errResp
	}
	if a == nil || a.Gameinfo == nil || a.Userinfo == nil {
		return M{
			"code": -1,
			"msg":  "游戏未开始",
		}
	}
	return nil
}
