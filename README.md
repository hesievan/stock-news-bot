# stock-news-bot 📰

> 从 [go-stock](https://github.com/ArvinLovegood/go-stock) 项目提取的**独立资讯采集工具**，部署在 GitHub Actions 上自动运行。
> 抓取财联社、华尔街见闻、新浪、TradingView 等多源资讯，归档 JSON + 推送通知 + 生成 GitHub Pages 静态站。

## ✨ 功能

| 功能 | 说明 |
|------|------|
| 📡 **多源抓取** | 财联社电报、华尔街见闻（全球/A股/美股）、新浪财经、TradingView 中文新闻 |
| 🗂️ **JSON 归档** | 每天一个 JSON 文件（`docs/data/YYYY-MM-DD.json`），commit 到仓库形成历史数据集 |
| 📢 **消息推送** | 支持钉钉机器人、飞书机器人，只推重要（红色）电报，避免刷屏 |
| 🌐 **GitHub Pages** | 自动生成暗色主题的资讯浏览站，按日期浏览历史归档 |
| ⏰ **自动运行** | GitHub Actions cron 定时触发，交易时段每 30 分钟一次 |

## 🚀 快速开始

### 1. Fork 本仓库

点击右上角 Fork 按钮。

### 2. 配置 Secrets（推送可选）

进入你 Fork 的仓库 → **Settings → Secrets and variables → Actions → New repository secret**

| Secret 名称 | 必填 | 说明 |
|-------------|------|------|
| `DING_WEBHOOK` | 否 | 钉钉自定义机器人 webhook 地址。留空则不推送钉钉 |
| `FEISHU_WEBHOOK` | 否 | 飞书自定义机器人 webhook 地址。留空则不推送飞书 |
| `FEISHU_SECRET` | 否 | 飞书机器人签名密钥（安全设置里开的）。不填则不签名 |

> **不配任何 Secret 也能用**——只是没有推送，归档和 Pages 照常。

### 3. 开启 GitHub Pages

仓库 → **Settings → Pages → Source** → 选择 **Deploy from a branch** → 分支选 **gh-pages**（首次 workflow 运行后自动创建）→ 目录 **/ (root)** → Save

### 4. 手动触发一次验证

仓库 → **Actions → 资讯采集 → Run workflow** → 等待运行完成。

完成后：
- `docs/data/` 下会出现当天的 JSON 归档
- 如果配了 Secrets，钉钉/飞书会收到一条资讯推送
- GitHub Pages 会更新（首次可能要等 1-2 分钟）

### 5. 定时已自动开启

workflow 里配了 cron schedule（周一到周五每 30 分钟），Push 后自动生效。

## 📁 仓库结构

```
├── main.go                       # 入口：编排抓取→归档→推送→生成站点
├── config.go                     # 从环境变量读配置，注入 settings 表
├── pusher/pusher.go              # 钉钉/飞书推送封装
├── site/site.go                  # JSON 归档写入 + manifest 重建
├── site/index_html.go            # GitHub Pages 首页模板（暗色主题）
├── backend/                      # 从 go-stock 原样拷贝，零修改
│   ├── data/                     #   资讯抓取核心（100+ 源文件）
│   ├── db/                       #   SQLite 初始化
│   ├── models/                   #   数据模型
│   ├── logger/                   #   日志
│   ├── machineid/                #   机器码（本项目未使用）
│   └── util/                     #   工具函数
├── docs/                         # GitHub Pages + JSON 归档（自动生成）
│   ├── index.html
│   ├── manifest.json
│   └── data/
│       └── 2026-07-13.json
├── .github/workflows/crawl.yml   # GitHub Actions 定时任务
└── go.mod / go.sum
```

## ⚙️ 环境变量（workflow 中使用）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_PATH` | `data/news.db` | 临时 SQLite 路径（CI 用 `/tmp/news.db`） |
| `DOCS_DIR` | `docs` | 归档输出目录 |
| `DING_WEBHOOK` | 空 | 钉钉 webhook，空=不推 |
| `FEISHU_WEBHOOK` | 空 | 飞书 webhook，空=不推 |
| `FEISHU_SECRET` | 空 | 飞书签名密钥 |
| `ONLY_RED_PUSH` | `true` | 只推送红色（重要）电报 |

## 🔄 同步上游

本项目 `backend/` 目录从 [go-stock](https://github.com/ArvinLovegood/go-stock) 整体拷贝。上游更新时：

```bash
# 1. 拉取最新 go-stock
cd /path/to/go-stock-dev
git pull

# 2. 覆盖 backend 目录
cp -R backend/ /path/to/stock-news-bot/backend/
cp go.mod go.sum /path/to/stock-news-bot/

# 3. 重新编译验证
cd /path/to/stock-news-bot
go build -o /dev/null .

# 4. 提交
git add backend/ go.mod go.sum
git commit -m "sync: 更新 backend 到最新 upstream"
git push
```

入口代码（`main.go`、`config.go`、`pusher/`、`site/`）与 backend 完全解耦，上游更新不会影响它们。

## 📄 JSON 归档格式

`docs/data/YYYY-MM-DD.json`：

```json
{
  "date": "2026-07-13",
  "fetched_at": "2026-07-13T08:30:00+08:00",
  "items": [
    {
      "time": "08:01:23",
      "source": "财联社电报",
      "title": "央行开展逆回购操作",
      "content": "央行今日开展1500亿元7天期逆回购操作...",
      "url": "https://www.cls.cn/telegraph/12345",
      "sentiment": "利好",
      "isRed": true
    }
  ]
}
```

## 🛡️ 许可

本项目使用的 `backend/` 代码来自 [go-stock](https://github.com/ArvinLovegood/go-stock)，遵循其原始许可证（AGPL-3.0）。新增的入口代码同样遵循 AGPL-3.0。
