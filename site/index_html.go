package site

// indexHTML 是 GitHub Pages 首页模板。纯静态、无依赖、内联 JS。
// Mobile-first 设计：手机端日期列表在顶部可折叠抽屉，桌面端保持侧栏。
// 设计目标：单文件、零构建、可直接被 GitHub Pages 托管。
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<title>股票资讯归档</title>
<style>
  * { box-sizing: border-box; -webkit-tap-highlight-color: transparent; }
  body { margin:0; font-family: -apple-system, "PingFang SC", "Microsoft YaHei", Helvetica, sans-serif; background:#0d1117; color:#c9d1d9; }
  a { color:#58a6ff; text-decoration:none; }

  /* ===== 移动端优先 ===== */

  /* 顶部导航栏 */
  .navbar { position:sticky; top:0; z-index:100; background:#161b22; border-bottom:1px solid #30363d; padding:0 12px; display:flex; align-items:center; height:48px; gap:8px; }
  .navbar h1 { margin:0; font-size:15px; color:#58a6ff; flex-shrink:0; white-space:nowrap; }
  .navbar .date-picker { flex:1; text-align:center; }
  .navbar .date-picker select { background:#0d1117; color:#c9d1d9; border:1px solid #30363d; border-radius:6px; padding:6px 10px; font-size:14px; width:100%; max-width:200px; text-align:center; }
  .navbar .toggle-days { background:none; border:1px solid #30363d; border-radius:6px; color:#c9d1d9; font-size:18px; padding:4px 10px; cursor:pointer; flex-shrink:0; }
  .navbar .toggle-days:active { background:#21262d; }

  /* 日期抽屉（移动端） */
  .drawer { background:#161b22; border-bottom:1px solid #30363d; max-height:0; overflow-y:auto; transition:max-height 0.25s ease; -webkit-overflow-scrolling:touch; }
  .drawer.open { max-height:50vh; }
  .drawer ol { list-style:none; margin:0; padding:4px 0; display:flex; flex-wrap:wrap; gap:4px; justify-content:center; }
  .drawer li { padding:8px 14px; border-radius:6px; cursor:pointer; font-size:13px; color:#8b949e; white-space:nowrap; }
  .drawer li:active { background:#21262d; }
  .drawer li.active { background:#1f6feb; color:#fff; }

  /* 主内容区 */
  .main { padding:12px; }
  .stats { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:10px 12px; margin-bottom:12px; font-size:12px; color:#8b949e; line-height:1.5; }
  .stats b { color:#e6edf3; }

  .item { background:#161b22; border:1px solid #30363d; border-radius:8px; padding:12px; margin-bottom:10px; }
  .item.red { border-left:3px solid #f85149; }
  .item.normal { border-left:3px solid #30363d; }
  .item-head { display:flex; flex-wrap:wrap; gap:4px 6px; align-items:center; font-size:11px; color:#8b949e; margin-bottom:4px; }
  .badge { background:#21262d; padding:2px 6px; border-radius:4px; white-space:nowrap; }
  .badge.red { background:#3d1418; color:#f85149; }
  .badge.sent { background:#1f2a1f; color:#7ee787; }
  .title { font-weight:600; color:#e6edf3; margin:2px 0 4px; font-size:14px; line-height:1.4; }
  .content { font-size:13px; line-height:1.6; white-space:pre-wrap; color:#c9d1d9; word-break:break-word; }
  .url { margin-top:6px; }
  .url a { font-size:12px; color:#58a6ff; }
  .empty { color:#8b949e; text-align:center; padding:40px 16px; font-size:14px; }

  /* ===== 桌面端覆盖 ===== */
  @media (min-width: 768px) {
    .navbar { display:none; }
    .drawer { display:none; }
    .layout { display:flex; min-height:100vh; }
    .sidebar { width:220px; background:#161b22; border-right:1px solid #30363d; padding:16px; overflow-y:auto; position:sticky; top:0; height:100vh; }
    .sidebar h1 { font-size:15px; color:#58a6ff; margin:0 0 12px; }
    .sidebar ol { list-style:none; padding:0; margin:0; display:block; }
    .sidebar li { padding:8px 10px; border-radius:6px; cursor:pointer; font-size:13px; color:#8b949e; }
    .sidebar li:hover { background:#21262d; color:#c9d1d9; }
    .sidebar li.active { background:#1f6feb; color:#fff; }
    .main { flex:1; padding:24px; max-width:960px; }
    .stats { font-size:13px; }
    .title { font-size:15px; }
    .content { font-size:14px; }
  }
</style>
</head>
<body>

<!-- ===== 移动端导航栏（>=768px 隐藏） ===== -->
<nav class="navbar">
  <h1>📰 资讯</h1>
  <div class="date-picker">
    <select id="dateSelect" onchange="loadDay(this.value)"><option>加载中…</option></select>
  </div>
  <button class="toggle-days" id="toggleBtn" aria-label="切换日期列表">📅</button>
</nav>
<div class="drawer" id="drawer">
  <ol id="days"></ol>
</div>

<!-- ===== 桌面端侧栏（<768px 隐藏） ===== -->
<div class="layout">
  <aside class="sidebar">
    <h1>📅 资讯归档</h1>
    <ol id="daysDesktop"></ol>
  </aside>
  <main class="main">
    <div id="stats" class="stats"></div>
    <div id="news"></div>
  </main>
</div>

<script>
const $days = document.getElementById('days');
const $daysDesktop = document.getElementById('daysDesktop');
const $news = document.getElementById('news');
const $stats = document.getElementById('stats');
const $select = document.getElementById('dateSelect');
const $toggle = document.getElementById('toggleBtn');
const $drawer = document.getElementById('drawer');
let dates = [];
let activeDate = null;

// 打开/关闭日期抽屉（移动端）
$toggle.onclick = () => {
  $drawer.classList.toggle('open');
};
// 选了日期后关抽屉
function closeDrawer() {
  $drawer.classList.remove('open');
}

async function loadManifest() {
  try {
    const r = await fetch('manifest.json?t=' + Date.now());
    dates = await r.json();
    if (!dates.length) {
      $days.innerHTML = '<li class="empty" style="width:100%">暂无归档</li>';
      $daysDesktop.innerHTML = '<li class="empty">暂无归档</li>';
      $select.innerHTML = '<option>暂无归档</option>';
      return;
    }
    // 移动端抽屉
    $days.innerHTML = dates.map(d =>
      '<li data-date="' + d + '">' + d + '</li>').join('');
    // 桌面端侧栏
    $daysDesktop.innerHTML = dates.map(d =>
      '<li data-date="' + d + '">' + d + '</li>').join('');

    // 下拉选择器
    $select.innerHTML = dates.map(d =>
      '<option value="' + d + '">' + d + '</option>').join('');

    // 点击事件：两个列表
    const clickDay = (li, date) => {
      [...$days.children].forEach(el => el.classList.toggle('active', el === li));
      [...$daysDesktop.children].forEach(el => el.classList.toggle('active', el === li));
      loadDay(date);
      closeDrawer();
      $select.value = date;
    };
    [...$days.children].forEach(li => {
      li.onclick = () => clickDay(li, li.dataset.date);
    });
    [...$daysDesktop.children].forEach(li => {
      li.onclick = () => clickDay(li, li.dataset.date);
    });

    loadDay(dates[0]);
  } catch (e) {
    $days.innerHTML = '<li class="empty" style="width:100%">加载失败</li>';
    $daysDesktop.innerHTML = '<li class="empty">加载失败</li>';
    $select.innerHTML = '<option>加载失败</option>';
  }
}

async function loadDay(date) {
  if (!date) return;
  activeDate = date;
  // 高亮
  [...$days.children].forEach(li => li.classList.toggle('active', li.dataset.date === date));
  [...$daysDesktop.children].forEach(li => li.classList.toggle('active', li.dataset.date === date));
  $select.value = date;

  $news.innerHTML = '<div class="empty">加载中…</div>';
  $stats.textContent = '';
  try {
    const r = await fetch('data/' + date + '.json?t=' + Date.now());
    if (!r.ok) throw new Error('HTTP ' + r.status);
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
  const url = i.url ? '<div class="url"><a href="' + i.url + '" target="_blank" rel="noopener">查看原文 →</a></div>' : '';
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