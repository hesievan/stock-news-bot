package pusher

import (
	"fmt"
	"strings"
	"time"

	"go-stock/backend/data"
	"go-stock/backend/models"
)

// PushResult 记录各通道发送结果，供主流程打日志。
type PushResult struct {
	Ding   string
	Feishu string
}

// NewsItem 是归档 JSON 的一条记录，和 site 渲染共用。
type NewsItem struct {
	Time      string `json:"time"`
	Source    string `json:"source"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Url       string `json:"url"`
	Sentiment string `json:"sentiment"`
	IsRed     bool   `json:"isRed"`
}

// FromTelegraph 把原项目 Telegraph 转成精简的归档结构。
func FromTelegraph(t *models.Telegraph) NewsItem {
	return NewsItem{
		Time:      t.Time,
		Source:    t.Source,
		Title:     t.Title,
		Content:   t.Content,
		Url:       t.Url,
		Sentiment: t.SentimentResult,
		IsRed:     t.IsRed,
	}
}

// PushConfig 只含推送需要的字段，和 main.AppConfig 解耦。
type PushConfig struct {
	DingRobot   string // 空=跳过钉钉
	FeishuRobot string // 空=跳过飞书
	OnlyRed     bool
}

// Push 按配置把重要电报推送到钉钉/飞书。
// onlyRed=true 时只推 isRed=true 的条目，避免刷屏。
// 任一通道失败不影响另一通道，也不影响归档。
func Push(cfg *PushConfig, items []NewsItem) PushResult {
	res := PushResult{}

	toSend := items
	if cfg.OnlyRed {
		filtered := make([]NewsItem, 0, len(items))
		for _, it := range items {
			if it.IsRed {
				filtered = append(filtered, it)
			}
		}
		toSend = filtered
	}

	if len(toSend) == 0 {
		return res
	}

	title := fmt.Sprintf("资讯速递 %s（%d条）", time.Now().Format("01-02 15:04"), len(toSend))
	text := BuildMarkdown(toSend)

	if cfg.DingRobot != "" {
		res.Ding = data.NewDingDingAPI().SendToDingDing(title, text)
	}
	if cfg.FeishuRobot != "" {
		res.Feishu = data.NewFeishuAPI().SendToFeishu(title, text)
	}
	return res
}

// BuildMarkdown 把多条电报拼成一段 markdown，钉钉/飞书通用。
func BuildMarkdown(items []NewsItem) string {
	var b strings.Builder
	for _, it := range items {
		head := "🔹"
		if it.IsRed {
			head = "🔴"
		}
		title := it.Title
		if title == "" {
			// 电报类往往标题为空，用内容前 40 字充当
			if len([]rune(it.Content)) > 40 {
				title = string([]rune(it.Content)[:40]) + "…"
			} else {
				title = it.Content
			}
		}
		b.WriteString(fmt.Sprintf("%s **%s** `%s` `%s`\n", head, title, it.Source, it.Time))
		if it.Content != "" && it.Content != title {
			b.WriteString(it.Content + "\n")
		}
		if it.Url != "" {
			b.WriteString(fmt.Sprintf("[查看原文](%s)\n", it.Url))
		}
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}
