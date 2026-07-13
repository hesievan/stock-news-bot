package main

import (
	"os"
	"strconv"
	"strings"

	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/models"
)

// envOr 读环境变量，空则返回 fallback。
func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

// AppConfig 是从环境变量/Secrets 解析出来的运行配置。
type AppConfig struct {
	DocsDir string // 归档与站点输出目录，默认 docs

	// 推送通道（留空则跳过对应通道）
	DingPushEnable bool
	DingRobot      string
	FeishuPush     bool
	FeishuRobot    string
	FeishuSecret   string

	// 抓取控制
	CrawlTimeoutSec int64
	OnlyRedPush     bool // 仅推送红色（重要）电报
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		DocsDir:         envOr("DOCS_DIR", "docs"),
		DingPushEnable:  envBool("DING_PUSH_ENABLE", os.Getenv("DING_WEBHOOK") != ""),
		DingRobot:       envOr("DING_WEBHOOK", ""),
		FeishuPush:      envBool("FEISHU_PUSH_ENABLE", os.Getenv("FEISHU_WEBHOOK") != ""),
		FeishuRobot:     envOr("FEISHU_WEBHOOK", ""),
		FeishuSecret:    envOr("FEISHU_SECRET", ""),
		CrawlTimeoutSec: 30,
		OnlyRedPush:     envBool("ONLY_RED_PUSH", true),
	}
}

// InjectSettings 把环境变量里的推送配置写进本地临时 SQLite 的 settings 表，
// 这样原项目里 GetSettingConfig() 读出来的配置就带上了 webhook，
// NewDingDingAPI()/NewFeishuAPI() 才能真正发出去。
//
// 原项目配置完全走 DB（见 settings_api.go GetSettingConfig），所以这里用插入一行的方式注入。
// 抓取代码零修改。
func InjectSettings(cfg *AppConfig) {
	settings := &data.Settings{
		LocalPushEnable: false,
		DingPushEnable:  cfg.DingPushEnable && cfg.DingRobot != "",
		DingRobot:       cfg.DingRobot,
		FeishuPushEnable: cfg.FeishuPush && cfg.FeishuRobot != "",
		FeishuRobot:     cfg.FeishuRobot,
		FeishuSecret:    cfg.FeishuSecret,
		CrawlTimeOut:    cfg.CrawlTimeoutSec,
		BrowserPoolSize: 1,
	}

	// settings 表是单行配置（First 取第一条），先清空再插入，保证幂等。
	db.Dao.Where("1 = 1").Delete(&data.Settings{})
	if err := db.Dao.Create(settings).Error; err != nil {
		// 不致命：推送功能会失效，但抓取/归档继续。
		println("warn: 写入 settings 表失败:", err.Error())
	}
}

// MigrateNewsTables 建好抓取函数和配置注入会触碰的表。
// 资讯抓取函数"边抓边写库"，缺表会直接报错；这里按需建表即可。
func MigrateNewsTables() {
	db.Dao.AutoMigrate(
		&models.Telegraph{},
		&models.TelegraphTags{},
		&models.Tags{},
		&data.Settings{},
		&data.AIConfig{},
	)
}
