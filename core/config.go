package core

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GetDataDir 获取数据目录路径（家目录下的 .lifegame）
func GetDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataDir := filepath.Join(homeDir, ".lifegame")
	return dataDir, nil
}

// EnsureDataDir 确保数据目录存在
func EnsureDataDir() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}
	return dataDir, nil
}

// GetConfigPath 获取配置文件完整路径
func GetConfigPath() (string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "config.yaml"), nil
}

// GetDBPath 获取数据库文件完整路径
func GetDBPath() (string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "lifegame.db"), nil
}

// GetImagesDir 获取图片资源目录路径（与数据库同目录下的 images 子目录）
func GetImagesDir() (string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	imagesDir := filepath.Join(dataDir, "images")
	// 确保 images 目录存在
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return "", err
	}
	return imagesDir, nil
}

// GetAudioDir 获取可由用户替换的音频资源目录。
func GetAudioDir() (string, error) {
	dataDir, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	audioDir := filepath.Join(dataDir, "audio")
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return "", err
	}
	return audioDir, nil
}

// Config 游戏配置结构体
type Config struct {
	Game    GameConfig    `yaml:"game"`
	Market  MarketConfig  `yaml:"market"`
	Company CompanyConfig `yaml:"company"`
	Loan    LoanConfig    `yaml:"loan"`
}

// GameConfig 游戏基础配置
type GameConfig struct {
	Name          string `yaml:"name"`
	AgeInit       int    `yaml:"age_init"`
	AgeMax        int    `yaml:"age_max"`
	AntiqueMaxNum int    `yaml:"antique_max_num"`
}

// MarketConfig 市场配置
type MarketConfig struct {
	UpPrice       float64 `yaml:"up_price"`
	DownPrice     float64 `yaml:"down_price"`
	ShowItemNum   int     `yaml:"show_item_num"`
	ShowMarketNum int     `yaml:"show_market_num"`
	MaxItemNum    int     `yaml:"max_item_num"`
	MaxFame       int     `yaml:"max_fame"`
	MaxImmunity   int     `yaml:"max_immunity"`
}

// CompanyConfig 公司配置
type CompanyConfig struct {
	Fluct  float64 `yaml:"fluct"`
	MaxNum int     `yaml:"max_num"`
}

// LoanConfig 贷款配置
type LoanConfig struct {
	InterestRate   float64 `yaml:"interest_rate"`
	OverduePenalty float64 `yaml:"overdue_penalty"`
}

// DefaultConfig 默认配置（代码内置）
var DefaultConfig = &Config{
	Game: GameConfig{
		Name:          "人生模拟器",
		AgeInit:       UserAgeInit,
		AgeMax:        UserAgeMax,
		AntiqueMaxNum: UserAntiqueNum,
	},
	Market: MarketConfig{
		UpPrice:       1.2,
		DownPrice:     0.8,
		ShowItemNum:   10,
		ShowMarketNum: 10,
		MaxItemNum:    MaxItemNum,
		MaxFame:       MaxFame,
		MaxImmunity:   MaxImmunity,
	},
	Company: CompanyConfig{
		Fluct:  0.02,
		MaxNum: 3,
	},
	Loan: LoanConfig{
		InterestRate:   0.10,
		OverduePenalty: 0.20,
	},
}

// AppConfig 全局配置实例
var AppConfig *Config

// InitConfig 初始化配置
// 如果配置文件不存在，则创建默认配置文件
func InitConfig() error {
	// 获取配置文件路径
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		if err := SaveConfigToYAML(configPath, DefaultConfig); err != nil {
			return err
		}
	}

	// 加载配置文件
	return LoadConfigFromYAML(configPath)
}

// LoadConfigFromYAML 从 YAML 文件加载配置
func LoadConfigFromYAML(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	AppConfig = &Config{}
	return yaml.Unmarshal(data, AppConfig)
}

// SaveConfigToYAML 保存配置到 YAML 文件
func SaveConfigToYAML(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetGameName 获取游戏名称
func GetGameName() string {
	if AppConfig != nil {
		return AppConfig.Game.Name
	}
	return DefaultConfig.Game.Name
}

// GetAgeInit 获取初始年龄
func GetAgeInit() int {
	if AppConfig != nil {
		return AppConfig.Game.AgeInit
	}
	return DefaultConfig.Game.AgeInit
}

// GetAgeMax 获取最大年龄
func GetAgeMax() int {
	if AppConfig != nil {
		return AppConfig.Game.AgeMax
	}
	return DefaultConfig.Game.AgeMax
}

// GetAntiqueMaxNum 获取最大古董数量
func GetAntiqueMaxNum() int {
	if AppConfig != nil {
		return AppConfig.Game.AntiqueMaxNum
	}
	return DefaultConfig.Game.AntiqueMaxNum
}

// GetUpPrice 获取上涨浮动
func GetUpPrice() float64 {
	if AppConfig != nil {
		return AppConfig.Market.UpPrice
	}
	return DefaultConfig.Market.UpPrice
}

// GetDownPrice 获取下跌浮动
func GetDownPrice() float64 {
	if AppConfig != nil {
		return AppConfig.Market.DownPrice
	}
	return DefaultConfig.Market.DownPrice
}

// GetShowItemNum 获取显示物资数量
func GetShowItemNum() int {
	if AppConfig != nil {
		return AppConfig.Market.ShowItemNum
	}
	return DefaultConfig.Market.ShowItemNum
}

// GetShowMarketNum 获取市场物资数量
func GetShowMarketNum() int {
	if AppConfig != nil {
		return AppConfig.Market.ShowMarketNum
	}
	return DefaultConfig.Market.ShowMarketNum
}

// GetMaxItemNum 获取最大物资持有量
func GetMaxItemNum() int {
	if AppConfig != nil {
		return AppConfig.Market.MaxItemNum
	}
	return DefaultConfig.Market.MaxItemNum
}

// GetMaxFame 获取最大声望
func GetMaxFame() int {
	if AppConfig != nil {
		return AppConfig.Market.MaxFame
	}
	return DefaultConfig.Market.MaxFame
}

// GetMaxImmunity 获取最大免疫
func GetMaxImmunity() int {
	if AppConfig != nil {
		return AppConfig.Market.MaxImmunity
	}
	return DefaultConfig.Market.MaxImmunity
}

// GetCompanyFluct 获取公司浮动
func GetCompanyFluct() float64 {
	if AppConfig != nil {
		return AppConfig.Company.Fluct
	}
	return DefaultConfig.Company.Fluct
}

// GetCompanyNum 获取最大创业数量
func GetCompanyNum() int {
	if AppConfig != nil {
		return AppConfig.Company.MaxNum
	}
	return DefaultConfig.Company.MaxNum
}

// GetLoanInterestRate 获取贷款利率
func GetLoanInterestRate() float64 {
	if AppConfig != nil {
		return AppConfig.Loan.InterestRate
	}
	return DefaultConfig.Loan.InterestRate
}

// GetOverduePenalty 获取逾期惩罚
func GetOverduePenalty() float64 {
	if AppConfig != nil {
		return AppConfig.Loan.OverduePenalty
	}
	return DefaultConfig.Loan.OverduePenalty
}

// ConfigToJSON 将配置转换为 JSON（用于存档）
func ConfigToJSON() (string, error) {
	data, err := json.Marshal(AppConfig)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ConfigFromJSON 从 JSON 恢复配置
func ConfigFromJSON(jsonStr string) error {
	AppConfig = &Config{}
	return json.Unmarshal([]byte(jsonStr), AppConfig)
}
