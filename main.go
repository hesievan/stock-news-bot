package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/models"
	"go-stock/pusher"
	"go-stock/site"
)

// 版本号，构建时可注入。
var version = "dev"

func main() {
	fmt.Printf("go-stock news bot %s 启动 @ %s\n", version, time.Now().Format(time.RFC3339))

	cfg := LoadConfig()

	// 1. 初始化临时 SQLite（轻量剥离：让抓取函数的 db.Dao.Create 正常工作）
	//    跑完即扔，不持久化。CI 上放 /tmp，本地放工作目录下的 data/。
	dbPath := envOr("DB_PATH", filepath.Join("data", "news.db"))
	_ = os.MkdirAll(filepath.Dir(dbPath), 0o755)
	db.Init(dbPath)
	MigrateNewsTables()
	InjectSettings(cfg)

	// 2. 抓取所有资讯源。每个源独立抓，失败不阻塞其他源。
	items := crawlAll(cfg)
	fmt.Printf("\n汇总：共抓到 %d 条电报\n", len(items))
	if len(items) == 0 {
		fmt.Println("无资讯，结束。")
		return
	}

	// 3. 按时间倒序排序（Telegraph 有 DataTime）
	sortByTime(items)

	// 4. 归档为 JSON（覆盖当天文件）
	date := time.Now().Format("2006-01-02")
	archItems := toArchItems(items)
	path, err := site.WriteDay(cfg.DocsDir, date, time.Now().Format(time.RFC3339), archItems)
	if err != nil {
		fmt.Printf("⚠️ 归档失败: %v\n", err)
	} else {
		fmt.Printf("✅ 已归档: %s（%d 条）\n", path, len(archItems))
	}

	// 5. 重建索引 + 确保首页存在
	if err := site.RebuildManifest(cfg.DocsDir); err != nil {
		fmt.Printf("⚠️ 生成 manifest 失败: %v\n", err)
	}
	if err := site.EnsureIndex(cfg.DocsDir); err != nil {
		fmt.Printf("⚠️ 生成 index.html 失败: %v\n", err)
	}
	st := site.ComputeStats(cfg.DocsDir)
	fmt.Printf("📊 站点统计：%d 天 / %d 条，最新 %s\n", st.TotalDays, st.TotalItems, st.LastDate)

	// 6. 推送（仅重要电报，避免刷屏）
	if cfg.DingRobot != "" || cfg.FeishuRobot != "" {
		pc := &pusher.PushConfig{
			DingRobot:   cfg.DingRobot,
			FeishuRobot: cfg.FeishuRobot,
			OnlyRed:     cfg.OnlyRedPush,
		}
		newsForPush := make([]pusher.NewsItem, len(archItems))
		for i, a := range archItems {
			newsForPush[i] = pusher.NewsItem(a)
		}
		res := pusher.Push(pc, newsForPush)
		if res.Ding != "" {
			fmt.Println("钉钉:", res.Ding)
		}
		if res.Feishu != "" {
			fmt.Println("飞书:", res.Feishu)
		}
	} else {
		fmt.Println("未配置推送 webhook，跳过推送。")
	}

	fmt.Println("\n完成。")
}

// crawlAll 并行抓取所有资讯源，汇总成 Telegraph 列表。
// 任一源失败：打印日志，返回空切片，不影响其他源。
func crawlAll(cfg *AppConfig) []*models.Telegraph {
	type sourceResult struct {
		name string
		data []*models.Telegraph
		err  error
	}

	// 定义各抓取任务（闭包捕获 cfg）
	tasks := []struct {
		name string
		run  func() []*models.Telegraph
	}{
		{"财联社电报", func() []*models.Telegraph {
			r := data.NewMarketNewsApi().TelegraphList(cfg.CrawlTimeoutSec)
			return sliceFromPtr(r)
		}},
		{"华尔街见闻-全球", func() []*models.Telegraph {
			r := data.NewWallstreetcnApi().GetLivesAsTelegraph("global-channel", 30)
			return sliceFromPtr(r)
		}},
		{"华尔街见闻-A股", func() []*models.Telegraph {
			r := data.NewWallstreetcnApi().GetLivesAsTelegraph("a-stock-channel", 20)
			return sliceFromPtr(r)
		}},
		{"华尔街见闻-美股", func() []*models.Telegraph {
			r := data.NewWallstreetcnApi().GetLivesAsTelegraph("us-stock-channel", 20)
			return sliceFromPtr(r)
		}},
		{"新浪财经", func() []*models.Telegraph {
			r := data.NewMarketNewsApi().GetSinaNews(uint(cfg.CrawlTimeoutSec))
			return sliceFromPtr(r)
		}},
		{"TradingView", func() []*models.Telegraph {
			r := data.NewMarketNewsApi().TradingViewNews()
			return sliceFromPtr(r)
		}},
	}

	// 并行执行
	ch := make(chan sourceResult, len(tasks))
	for _, t := range tasks {
		t := t // 捕获
		go func() {
			defer func() {
				if r := recover(); r != nil {
					ch <- sourceResult{name: t.name, err: fmt.Errorf("panic: %v", r)}
				}
			}()
			start := time.Now()
			got := t.run()
			ch <- sourceResult{name: t.name, data: got}
			fmt.Printf("  [%s] %d 条（%s）\n", t.name, len(got), time.Since(start).Truncate(time.Millisecond))
		}()
	}

	// 收集结果
	var all []*models.Telegraph
	for i := 0; i < len(tasks); i++ {
		res := <-ch
		if res.err != nil {
			fmt.Printf("  ⚠️ [%s] 失败: %v\n", res.name, res.err)
			continue
		}
		all = append(all, res.data...)
	}
	return all
}

// sliceFromPtr 把 *[]Telegraph 安全转成切片，nil 返回空切片。
func sliceFromPtr(p *[]models.Telegraph) []*models.Telegraph {
	if p == nil {
		return nil
	}
	out := make([]*models.Telegraph, 0, len(*p))
	for i := range *p {
		out = append(out, &(*p)[i])
	}
	return out
}

// sortByTime 按 DataTime 倒序（最新的在前）。
func sortByTime(items []*models.Telegraph) {
	// 简单插入排序，数据量不大，避免引入 sort 之外的复杂度
	for i := 1; i < len(items); i++ {
		for j := i; j > 0; j-- {
			if timeOf(items[j]).After(timeOf(items[j-1])) {
				items[j], items[j-1] = items[j-1], items[j]
			} else {
				break
			}
		}
	}
}

func timeOf(t *models.Telegraph) time.Time {
	if t != nil && t.DataTime != nil {
		return *t.DataTime
	}
	return time.Time{}
}

// toArchItems 把 Telegraph 转成 site 归档结构（两者字段对齐但解耦）。
func toArchItems(items []*models.Telegraph) []site.ArchItem {
	out := make([]site.ArchItem, 0, len(items))
	for _, t := range items {
		if t == nil {
			continue
		}
		out = append(out, site.ArchItem{
			Time:      t.Time,
			Source:    t.Source,
			Title:     t.Title,
			Content:   t.Content,
			Url:       t.Url,
			Sentiment: t.SentimentResult,
			IsRed:     t.IsRed,
		})
	}
	return out
}
