package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DayArchive 是单日归档文件的结构，和 pusher.NewsItem 字段保持一致但解耦（避免循环 import）。
type DayArchive struct {
	Date      string     `json:"date"`
	FetchedAt string     `json:"fetched_at"`
	Items     []ArchItem `json:"items"`
}

// ArchItem 归档条目。用 map[string]any 的精简版，字段和 pusher.NewsItem 对齐。
type ArchItem struct {
	Time      string `json:"time"`
	Source    string `json:"source"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Url       string `json:"url"`
	Sentiment string `json:"sentiment"`
	IsRed     bool   `json:"isRed"`
}

// WriteDay 把当天归档写入 docs/data/YYYY-MM-DD.json。
// 同一天多次运行会覆盖当天文件（以最新抓取为准），历史日期文件不动。
func WriteDay(docsDir, date, fetchedAt string, items []ArchItem) (string, error) {
	dir := filepath.Join(docsDir, "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	arc := DayArchive{Date: date, FetchedAt: fetchedAt, Items: items}
	b, err := json.MarshalIndent(arc, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, date+".json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// RebuildManifest 扫描 docs/data/*.json，按日期倒序生成 manifest.json。
// 静态站前端通过它知道有哪些天可看。
func RebuildManifest(docsDir string) error {
	dir := filepath.Join(docsDir, "data")
	entries, err := os.ReadDir(dir)
	if err != nil {
		// 目录还没数据，写个空数组
		return os.WriteFile(filepath.Join(docsDir, "manifest.json"), []byte("[]"), 0o644)
	}
	var dates []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		dates = append(dates, strings.TrimSuffix(e.Name(), ".json"))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	b, err := json.Marshal(dates)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(docsDir, "manifest.json"), b, 0o644)
}

// EnsureIndex 确保 docs/index.html 存在。已存在则不覆盖（保留用户自定义）。
func EnsureIndex(docsDir string) error {
	path := filepath.Join(docsDir, "index.html")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(indexHTML), 0o644)
}

// Stats 用于生成首页摘要。
type Stats struct {
	TotalDays int
	TotalItems int
	LastDate   string
}

// ComputeStats 扫描 docs/data 计算简单统计。
func ComputeStats(docsDir string) Stats {
	st := Stats{}
	dir := filepath.Join(docsDir, "data")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return st
	}
	var dates []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		dates = append(dates, name)
		b, _ := os.ReadFile(filepath.Join(dir, e.Name()))
		var arc DayArchive
		if err := json.Unmarshal(b, &arc); err == nil {
			st.TotalItems += len(arc.Items)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	st.TotalDays = len(dates)
	if len(dates) > 0 {
		st.LastDate = dates[0]
	}
	return st
}

// SummaryMarkdown 生成一份今日摘要 markdown，可供推送或归档首页使用。
func SummaryMarkdown(date string, items []ArchItem) string {
	red := 0
	for _, i := range items {
		if i.IsRed {
			red++
		}
	}
	return fmt.Sprintf("## %s 资讯摘要\n\n共 %d 条，其中重要（红色）%d 条。\n",
		date, len(items), red)
}
