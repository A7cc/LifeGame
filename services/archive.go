package services

import (
	"LifeGame/core"
	"LifeGame/internal/db"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const currentSaveVersion = 6

// SaveInfo 是存档列表中稳定、可生成前端类型的元数据。
type SaveInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	GameYear    int    `json:"game_year"`
	SaveVersion int    `json:"save_version"`
}

type SaveGameResponse struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	SaveID      int64  `json:"saveId,omitempty"`
	SaveVersion int    `json:"saveversion,omitempty"`
}

type LoadGameResponse struct {
	Code         int                   `json:"code"`
	Msg          string                `json:"msg"`
	Gameinfo     *core.Game            `json:"gameinfo,omitempty"`
	Userinfo     *core.User            `json:"userinfo,omitempty"`
	Announce     *core.Announce        `json:"announce,omitempty"`
	Difficulty   core.DifficultyConfig `json:"difficulty,omitempty"`
	StockEpoch   string                `json:"stockepoch,omitempty"`
	StockVersion uint64                `json:"stockversion,omitempty"`
	SaveVersion  int                   `json:"saveversion,omitempty"`
}

type ListSavesResponse struct {
	Code  int        `json:"code"`
	Msg   string     `json:"msg"`
	Saves []SaveInfo `json:"saves"`
}

// SaveGame 保存游戏，并明确写入当前存档格式版本。
func (a *App) SaveGame(name string) SaveGameResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if a.Gameinfo == nil || a.Userinfo == nil {
		return SaveGameResponse{Code: -1, Msg: "没有正在进行的游戏"}
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = a.Userinfo.UName + "的存档"
	}
	if utf8.RuneCountInString(name) > 30 {
		return SaveGameResponse{Code: -1, Msg: "存档名称过长"}
	}

	database := db.GetDB()
	if database == nil {
		return SaveGameResponse{Code: -1, Msg: "数据库未初始化"}
	}

	gameData, err := json.Marshal(a.Gameinfo)
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "游戏数据序列化失败: " + err.Error()}
	}
	announce := a.currentAnnounce()
	announceData, err := json.Marshal(announce)
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "公告数据序列化失败: " + err.Error()}
	}
	userData, err := json.Marshal(a.Userinfo)
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "用户数据序列化失败: " + err.Error()}
	}

	tx, err := database.Begin()
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "创建存档事务失败: " + err.Error()}
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO saves (name, created_at, game_year, save_version) VALUES (?, ?, ?, ?)",
		name, time.Now().Format("2006-01-02 15:04:05"), a.Userinfo.UAge, currentSaveVersion,
	)
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "保存失败: " + err.Error()}
	}
	saveID, err := result.LastInsertId()
	if err != nil {
		return SaveGameResponse{Code: -1, Msg: "获取存档编号失败: " + err.Error()}
	}

	if _, err = tx.Exec("INSERT INTO game_saves (save_id, game_data) VALUES (?, ?)", saveID, string(gameData)); err != nil {
		return SaveGameResponse{Code: -1, Msg: "游戏存档保存失败: " + err.Error()}
	}
	if _, err = tx.Exec("INSERT INTO user_saves (save_id, user_data) VALUES (?, ?)", saveID, string(userData)); err != nil {
		return SaveGameResponse{Code: -1, Msg: "用户存档保存失败: " + err.Error()}
	}
	if _, err = tx.Exec("INSERT INTO announce_saves (save_id, announce_data) VALUES (?, ?)", saveID, string(announceData)); err != nil {
		return SaveGameResponse{Code: -1, Msg: "公告存档保存失败: " + err.Error()}
	}
	if err := tx.Commit(); err != nil {
		return SaveGameResponse{Code: -1, Msg: "提交存档失败: " + err.Error()}
	}

	return SaveGameResponse{
		Code:        200,
		Msg:         "保存成功",
		SaveID:      saveID,
		SaveVersion: currentSaveVersion,
	}
}

// LoadGame 只加载当前格式存档，不转换历史版本。
func (a *App) LoadGame(saveID int) LoadGameResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireStartup(); errResp != nil {
		return LoadGameResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if saveID <= 0 {
		return LoadGameResponse{Code: -1, Msg: "存档编号无效"}
	}
	database := db.GetDB()
	if database == nil {
		return LoadGameResponse{Code: -1, Msg: "数据库未初始化"}
	}

	var saveVersion int
	var gameDataStr, userDataStr, announceDataStr string
	err := database.QueryRow(`
		SELECT s.save_version, g.game_data, u.user_data, a.announce_data
		FROM saves s
		JOIN game_saves g ON g.save_id = s.id
		JOIN user_saves u ON u.save_id = s.id
		JOIN announce_saves a ON a.save_id = s.id
		WHERE s.id = ?
	`, saveID).Scan(&saveVersion, &gameDataStr, &userDataStr, &announceDataStr)
	if err != nil {
		return LoadGameResponse{Code: -1, Msg: "未找到完整存档: " + err.Error()}
	}
	if saveVersion != currentSaveVersion {
		return LoadGameResponse{
			Code: -1,
			Msg:  fmt.Sprintf("存档格式版本 %d 与当前版本 %d 不一致，请创建新存档", saveVersion, currentSaveVersion),
		}
	}

	loadedGame := core.Game{}
	if err := json.Unmarshal([]byte(gameDataStr), &loadedGame); err != nil {
		return LoadGameResponse{Code: -1, Msg: "游戏数据解析失败: " + err.Error()}
	}
	loadedUser := core.User{}
	if err := json.Unmarshal([]byte(userDataStr), &loadedUser); err != nil {
		return LoadGameResponse{Code: -1, Msg: "用户数据解析失败: " + err.Error()}
	}

	loadedAnnounce := core.Announce{}
	if err := json.Unmarshal([]byte(announceDataStr), &loadedAnnounce); err != nil {
		return LoadGameResponse{Code: -1, Msg: "公告数据解析失败: " + err.Error()}
	}

	a.Gameinfo = &loadedGame
	a.Userinfo = &loadedUser
	a.setAnnounce(loadedAnnounce)
	a.MiniGameSession = nil
	a.pendingDatingInteraction = nil
	a.resetStockClock()

	return LoadGameResponse{
		Code:         200,
		Msg:          "加载成功",
		Gameinfo:     a.gameSnapshot(),
		Userinfo:     a.userSnapshot(),
		Announce:     &loadedAnnounce,
		Difficulty:   core.GetDifficultyConfig(loadedGame.GDifficulty),
		StockEpoch:   a.stockEpoch,
		StockVersion: a.stockVersion,
		SaveVersion:  saveVersion,
	}
}

func (a *App) ListSaves() ListSavesResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireStartup(); errResp != nil {
		return ListSavesResponse{Code: -1, Msg: responseMessage(errResp), Saves: []SaveInfo{}}
	}
	database := db.GetDB()
	if database == nil {
		return ListSavesResponse{Code: -1, Msg: "数据库未初始化", Saves: []SaveInfo{}}
	}

	rows, err := database.Query(`
		SELECT id, name, created_at, game_year, save_version
		FROM saves ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		return ListSavesResponse{Code: -1, Msg: "查询存档失败: " + err.Error(), Saves: []SaveInfo{}}
	}
	defer rows.Close()

	saves := make([]SaveInfo, 0)
	for rows.Next() {
		var save SaveInfo
		if err := rows.Scan(&save.ID, &save.Name, &save.CreatedAt, &save.GameYear, &save.SaveVersion); err != nil {
			return ListSavesResponse{Code: -1, Msg: "读取存档失败: " + err.Error(), Saves: []SaveInfo{}}
		}
		saves = append(saves, save)
	}
	if err := rows.Err(); err != nil {
		return ListSavesResponse{Code: -1, Msg: "遍历存档失败: " + err.Error(), Saves: []SaveInfo{}}
	}
	return ListSavesResponse{Code: 200, Msg: "查询成功", Saves: saves}
}

func (a *App) DeleteSave(saveID int) BasicResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireStartup(); errResp != nil {
		return BasicResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if saveID <= 0 {
		return BasicResponse{Code: -1, Msg: "存档编号无效"}
	}
	database := db.GetDB()
	if database == nil {
		return BasicResponse{Code: -1, Msg: "数据库未初始化"}
	}

	tx, err := database.Begin()
	if err != nil {
		return BasicResponse{Code: -1, Msg: "创建删除事务失败: " + err.Error()}
	}
	defer tx.Rollback()

	result, err := tx.Exec("DELETE FROM saves WHERE id = ?", saveID)
	if err != nil {
		return BasicResponse{Code: -1, Msg: "删除存档元数据失败: " + err.Error()}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return BasicResponse{Code: -1, Msg: "确认删除结果失败: " + err.Error()}
	}
	if rowsAffected == 0 {
		return BasicResponse{Code: -1, Msg: fmt.Sprintf("存档 %d 不存在", saveID)}
	}
	if err := tx.Commit(); err != nil {
		return BasicResponse{Code: -1, Msg: "提交删除失败: " + err.Error()}
	}
	return BasicResponse{Code: 200, Msg: "删除成功"}
}
