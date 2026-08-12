package services

import "LifeGame/core"

// BasicResponse 是只包含状态码和消息的稳定响应边界。
type BasicResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// StartupResponse 让前端在应用尚未就绪时也能获得稳定的启动状态类型。
type StartupResponse struct {
	Code   int              `json:"code"`
	Msg    string           `json:"msg"`
	Status AppStartupStatus `json:"status"`
}

// GameStateResponse 用于创建游戏和推进年份。可选字段使用指针，确保
// Wails 生成的前端声明能反映真实字段，而不是退化为 interface{}。
type GameStateResponse struct {
	Code         int                    `json:"code"`
	Msg          string                 `json:"msg"`
	Gameinfo     *core.Game             `json:"gameinfo,omitempty"`
	Userinfo     *core.User             `json:"userinfo,omitempty"`
	Announce     *core.Announce         `json:"announce,omitempty"`
	Difficulty   *core.DifficultyConfig `json:"difficulty,omitempty"`
	StockEpoch   string                 `json:"stockepoch,omitempty"`
	StockVersion uint64                 `json:"stockversion,omitempty"`
}

// StockUpdateResponse 是轻量行情响应，只同步股票相关字段，避免定时器
// 每次轮询都传输并覆盖完整游戏状态。
type StockUpdateResponse struct {
	Code         int              `json:"code"`
	Msg          string           `json:"msg"`
	Stocks       []core.StockInfo `json:"stocks"`
	News         []string         `json:"news"`
	Epoch        string           `json:"epoch"`
	Version      uint64           `json:"version"`
	UpdatedAt    int64            `json:"updatedAt"`
	Remaining    int              `json:"remaining"`
	MarketClosed bool             `json:"marketclosed"`
}

type EvaluationResponse struct {
	Code       int             `json:"code"`
	Msg        string          `json:"msg"`
	Evaluation *GameEvaluation `json:"evaluation,omitempty"`
}

type DatingListResponse struct {
	Code          int               `json:"code"`
	Msg           string            `json:"msg"`
	DatingList    []core.DatingInfo `json:"datinglist"`
	MeetingScenes []string          `json:"meetingscenes"`
	Userinfo      *core.User        `json:"userinfo,omitempty"`
}

type DatingSceneResponse struct {
	Code     int               `json:"code"`
	Msg      string            `json:"msg"`
	Scene    string            `json:"scene,omitempty"`
	Met      []core.DatingInfo `json:"met"`
	Userinfo *core.User        `json:"userinfo,omitempty"`
}

type DatingSceneOutcome struct {
	Location             string  `json:"location"`
	Category             string  `json:"category"`
	Label                string  `json:"label"`
	Event                string  `json:"event"`
	Effect               string  `json:"effect"`
	Preferred            bool    `json:"preferred"`
	PreferenceRateChange float64 `json:"preferenceratechange"`
	SuccessRate          float64 `json:"successrate"`
	RewardTier           string  `json:"rewardtier"`
	RewardMultiplier     int     `json:"rewardmultiplier"`
	Moment               string  `json:"moment,omitempty"`
	MomentLabel          string  `json:"momentlabel,omitempty"`
	MomentEvent          string  `json:"momentevent,omitempty"`
	MomentAffinityChange int     `json:"momentaffinitychange,omitempty"`
	FameChange           int     `json:"famechange"`
	HealthChange         int     `json:"healthchange"`
	AffinityChange       int     `json:"affinitychange"`
}

type DatingResultResponse struct {
	Code           int                  `json:"code"`
	Msg            string               `json:"msg"`
	Success        bool                 `json:"success"`
	FameChange     int                  `json:"famechange"`
	HealthChange   int                  `json:"healthchange"`
	AffinityChange int                  `json:"affinitychange"`
	Scene          *DatingSceneOutcome  `json:"scene,omitempty"`
	Datinginfo     *core.UserDatingInfo `json:"datinginfo,omitempty"`
	Userinfo       *core.User           `json:"userinfo,omitempty"`
}

type DatingRelationshipResponse struct {
	Code           int                  `json:"code"`
	Msg            string               `json:"msg"`
	Event          string               `json:"event,omitempty"`
	Gift           string               `json:"gift,omitempty"`
	GiftCost       int                  `json:"giftcost,omitempty"`
	AffinityChange int                  `json:"affinitychange,omitempty"`
	Outcome        string               `json:"outcome,omitempty"`
	Preferred      bool                 `json:"preferred"`
	Success        bool                 `json:"success"`
	SuccessRate    float64              `json:"successrate,omitempty"`
	Datinginfo     *core.UserDatingInfo `json:"datinginfo,omitempty"`
	Userinfo       *core.User           `json:"userinfo,omitempty"`
}

type DatingGiftOption struct {
	Name            string `json:"name"`
	Cost            int    `json:"cost"`
	PreferredEffect int    `json:"preferredeffect"`
	RiskyEffect     int    `json:"riskyeffect"`
}

type DatingGiftOptionsResponse struct {
	Code    int                `json:"code"`
	Msg     string             `json:"msg"`
	Options []DatingGiftOption `json:"options,omitempty"`
}

type DatingInteractionResponse struct {
	Code           int                  `json:"code"`
	Msg            string               `json:"msg"`
	Interaction    string               `json:"interaction,omitempty"`
	Outcome        string               `json:"outcome,omitempty"`
	Label          string               `json:"label,omitempty"`
	Event          string               `json:"event,omitempty"`
	Outfit         string               `json:"outfit,omitempty"`
	OutfitVariant  string               `json:"outfitvariant,omitempty"`
	OutfitImage    string               `json:"outfitimage,omitempty"`
	AffinityChange int                  `json:"affinitychange"`
	Datinginfo     *core.UserDatingInfo `json:"datinginfo,omitempty"`
	Userinfo       *core.User           `json:"userinfo,omitempty"`
}

type SpouseInteractionResponse struct {
	Code           int                  `json:"code"`
	Msg            string               `json:"msg"`
	Interaction    string               `json:"interaction,omitempty"`
	Cost           int                  `json:"cost,omitempty"`
	AffinityChange int                  `json:"affinitychange,omitempty"`
	HealthChange   int                  `json:"healthchange,omitempty"`
	Datinginfo     *core.UserDatingInfo `json:"datinginfo,omitempty"`
	Userinfo       *core.User           `json:"userinfo,omitempty"`
}

func gameStateError(message string) GameStateResponse {
	return GameStateResponse{Code: -1, Msg: message}
}

func responseMessage(response M) string {
	if message, ok := response["msg"].(string); ok && message != "" {
		return message
	}
	return "操作失败"
}
