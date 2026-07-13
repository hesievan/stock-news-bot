package site

// indexHTML 是 GitHub Pages 首页模板。纯静态、无依赖、内联 JS。
// 前端逻辑：fetch manifest.json 拿日期列表 → 渲染侧栏 → 点击某天加载 data/某天.json。
// 设计目标：单文件、零构建、可直接被 GitHub Pages 托管。
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>股票资讯归档</title>
<style>
  * { box-sizing: border-box; }
  body { margin:0; font-family: -apple-system, "PingFang SC", "Microsoft YaHei", sans-serif; background:#0d1117; color:#c9d1d9; }
  .layout { display:flex; min-height:100vh; }
  .sidebar { width:240px; background:#161b22; border-right:1px solid #30363d; padding:16px; overflow-y:auto; }
  .sidebar h1 { font-size:15px; color:#58a6ff; margin:0 0 12px; }
  .sidebar ol { list-style:none; padding:0; margin:0; }
  .sidebar li { padding:8px 10px; border-radius:6px; cursor:pointer; font-size:13px; color:#8b949e; }
  .sidebar li:hover { background:#21262d; color:#c9d1d9; }
  .sidebar li.active { background:#1f6feb; color:#fff; }
  .main { flex:1; padding:24px; overflow-y:auto; }
  .meta { color:#8b949e; font-size:13px; margin-bottom:16px; }
  .item { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:14px 16px; margin-bottom:12px; }
  .item.red { border-left:3px solid #f85149; }
  .item.normal { border-left:3px solid #30363d; }
  .item-head { display:flex; gap:8px; align-items:center; font-size:12px; color:#8b949e; margin-bottom:6px; }
  .badge { background:#21262d; padding:2px 6px; border-radius:4px; }
  .badge.red { background:#3d1418; color:#f85149; }
  .badge.sent { background:#1f2a1f; color:#7ee787; }
  .title { font-weight:600; color:#e6edf3; margin:2px 0 6px; font-size:15px; }
  .content { font-size:14px; line-height:1.6; white-space:pre-wrap; color:#c9d1d9; }
  .url { margin-top:6px; }
  .url a { color:#58a6ff; font-size:12px; text-decoration:none; }
  .empty { color:#8b949e; text-align:center; padding:60px; }
  .stats { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:12px 16px; margin-bottom:16px; font-size:13px; color:#8b949e; }
  .stats b { color:#e6edf3; }
</style>
</head>
<body>
<div class="layout">
  <aside class="sidebar">
    <h1>📅 资讯归档</h1>
    <ol id="days"><li class="empty">加载中…</li></ol>
  </aside>
  <main class="main">
    <div id="stats" class="stats"></div>
    <div id="news"></div>
  </main>
</div>
<script>
const $days = document.getElementById('days');
const $news = document.getElementById('news');
const $stats = document.getElementById('stats');
let activeDate = null;

async function loadManifest() {
  try {
    const r = await fetch('manifest.json?t=' + Date.now());
    const dates = await r.json();
    if (!dates.length) {
      $days.innerHTML = '<li class="empty">暂无归档</li>';
      return;
    }
    $days.innerHTML = dates.map(d =>
      '<li data-date="' + d + '">' + d + '</li>').join('');
    [...$days.children].forEach(li => {
      li.onclick = () => loadDay(li.dataset.date);
    });
    loadDay(dates[0]);
  } catch (e) {
    $days.innerHTML = '<li class="empty">加载失败</li>';
  }
}

async function loadDay(date) {
  activeDate = date;
  [...$days.children].forEach(li =>
    li.classList.toggle('active', li.dataset.date === date));
  $news.innerHTML = '<div class="empty">加载中…</div>';
  $stats.textContent = '';
  try {
    const r = await fetch('data/' + date + '.json?t=' + Date.now());
    const arc = await r.json();
    const red = arc.items.filter(i => i.isRed).length;
    $stats.innerHTML = '📅 <b>' + arc.date + '</b> · 共 <b>' + arc.items.length +
      '</b> 条 · 🔴 重要 <b>' + red + '</b> 条 · 抓取于 ' + arc.fetched_at;
    if (!arc.items.length) {
      $news.innerHTML = '<div class="empty">当天无资讯</div>';
      return;
    }
    $news.innerHTML = arc.items.map(renderItem).join('');
  } catch (e) {
    $news.innerHTML = '<div class="empty">加载失败</div>';
  }
}

function renderItem(i) {
  const title = i.title || (i.content || '').slice(0, 40);
  const cls = i.isRed ? 'red' : 'normal';
  const redBadge = i.isRed ? '<span class="badge red">重要</span>' : '';
  const sentBadge = i.sentiment ? '<span class="badge sent">' + i.sentiment + '</span>' : '';
  const url = i.url ? '<div class="url"><a href="' + i.url + '" target="_blank">查看原文 →</a></div>' : '';
  const content = i.content && i.content !== title ?
    '<div class="content">' + escapeHtml(i.content) + '</div>' : '';
  return '<div class="item ' + cls + '">' +
    '<div class="item-head">' +
      '<span class="badge">' + escapeHtml(i.source || '') + '</span>' +
      '<span>' + escapeHtml(i.time || '') + '</span>' +
      redBadge + sentBadge +
    '</div>' +
    '<div class="title">' + escapeHtml(title) + '</div>' +
    content + url +
  '</div>';
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, c => ({
    '&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'
  }[c]));
}

loadManifest();
</script>
</body>
</html>
`
