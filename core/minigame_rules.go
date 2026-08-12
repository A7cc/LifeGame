package core

import (
	"errors"
	"time"
)

// 查找小游戏是否存在，如果存在则返回小游戏配置
func FindMiniGameConfig(name string) (MiniGame, bool) {
	for _, mgame := range GetMiniGames() {
		if mgame.MGName == name {
			return mgame, true
		}
	}
	return MiniGame{}, false
}

// checkAndClearExpiredSession 检查并清理过期的 Session
func CheckAndClearExpiredSession(minigamesession *MiniGameSession) error {
	if minigamesession == nil {
		return errors.New("没有找到对应的小游戏")
	}

	startTime := time.Unix(minigamesession.MGSStartTime, 0)
	// Session 最少要超过 minigamesession.MGSSubInfo.SMGMinRunTime
	if time.Since(startTime) < time.Duration(minigamesession.MGSSubInfo.SMGMinRunTime)*time.Second {
		return errors.New("你的游戏时长太短了，视为作弊，请重新开始")
	}
	return nil
}
