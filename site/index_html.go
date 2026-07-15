package site

// indexHTML 是 GitHub Pages 首页模板。纯静态、无依赖、内联 JS。
// 设计目标：阅读优先、移动端友好、单文件零构建。
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#0d1117">
<title>股票资讯归档</title>
<style>
  :root {
    --bg:#0d1117; --bg-card:#161b22; --bg-hover:#1c2128; --bg-inset:#010409;
    --border:#30363d; --border-muted:#21262d;
    --text:#e6edf3; --text-muted:#8b949e; --text-dim:#6e7681;
    --accent:#58a6ff; --green:#7ee787; --red:#f85149; --orange:#d29922; --purple:#bc8cff;
    --radius:10px;
  }
  * { box-sizing:border-box; -webkit-tap-highlight-color:transparent; }
  html { -webkit-text-size-adjust:100%; }
  body {
    margin:0; background:var(--bg); color:var(--text);
    font-family:-apple-system,BlinkMacSystemFont,"PingFang SC","Microsoft YaHei","Segoe UI",Helvetica,Arial,sans-serif;
    font-size:16px; line-height:1.7; -webkit-font-smoothing:antialiased;
  }
  a { color:var(--accent); text-decoration:none; }
  a:active { opacity:0.7; }

  /* ===== 顶部导航栏（移动端 + 桌面共用） ===== */
  .navbar {
    position:sticky; top:0; z-index:100;
    background:rgba(13,17,23,0.92); backdrop-filter:blur(12px); -webkit-backdrop-filter:blur(12px);
    border-bottom:1px solid var(--border);
    display:flex; align-items:center; gap:8px;
    padding:10px 12px; min-height:52px;
  }
  .navbar .brand { font-size:15px; font-weight:600; color:var(--accent); white-space:nowrap; flex-shrink:0; }
  .navbar select {
    flex:1; min-width:0;
    background:var(--bg-inset); color:var(--text); border:1px solid var(--border);
    border-radius:8px; padding:10px 12px; font-size:15px; cursor:pointer;
  }
  .icon-btn {
    background:var(--bg-card); border:1px solid var(--border); border-radius:8px;
    color:var(--text); padding:10px 14px; font-size:16px; cursor:pointer; flex-shrink:0;
    line-height:1; transition:background 0.15s; min-height:40px; min-width:40px;
  }
  .icon-btn:active { background:var(--bg-hover); }

  /* ===== 筛选条 ===== */
  .filters {
    position:sticky; top:56px; z-index:90;
    background:rgba(13,17,23,0.92); backdrop-filter:blur(8px); -webkit-backdrop-filter:blur(8px);
    border-bottom:1px solid var(--border-muted);
    display:flex; gap:6px; padding:8px 14px; overflow-x:auto; -webkit-overflow-scrolling:touch;
    scrollbar-width:none;
  }
  .filters::-webkit-scrollbar { display:none; }
  .chip {
    background:var(--bg-card); border:1px solid var(--border); border-radius:999px;
    padding:7px 14px; font-size:13px; color:var(--text-muted); white-space:nowrap; cursor:pointer;
    transition:all 0.15s; flex-shrink:0; min-height:36px; user-select:none;
  }
  .chip:active { background:var(--bg-hover); }
  .chip.active { background:var(--accent); color:#fff; border-color:var(--accent); }
  .chip .count { opacity:0.7; margin-left:3px; font-size:11px; }

  /* ===== 日期抽屉（移动端） ===== */
  .drawer {
    position:fixed; top:0; left:0; right:0; bottom:0; z-index:200;
    background:rgba(0,0,0,0.6); opacity:0; pointer-events:none; transition:opacity 0.2s;
  }
  .drawer.open { opacity:1; pointer-events:auto; }
  .drawer-panel {
    position:absolute; top:0; left:0; bottom:0; width:78%; max-width:300px;
    background:var(--bg-card); padding:16px;
    transform:translateX(-100%); transition:transform 0.25s ease;
    overflow-y:auto; -webkit-overflow-scrolling:touch;
  }
  .drawer.open .drawer-panel { transform:translateX(0); }
  .drawer-panel h2 { font-size:14px; color:var(--text-muted); margin:0 0 12px; font-weight:500; }
  .drawer-panel ol { list-style:none; margin:0; padding:0; }
  .drawer-panel li {
    padding:12px 14px; border-radius:8px; cursor:pointer; font-size:14px; color:var(--text-muted);
    display:flex; justify-content:space-between; align-items:center;
  }
  .drawer-panel li:active { background:var(--bg-hover); }
  .drawer-panel li.active { background:var(--accent); color:#fff; }

  /* ===== 主内容 ===== */
  .container { max-width:880px; margin:0 auto; padding:14px 14px 60px; }
  .stats {
    background:var(--bg-card); border:1px solid var(--border); border-radius:var(--radius);
    padding:12px 16px; margin-bottom:16px; font-size:13px; color:var(--text-muted); line-height:1.6;
  }
  .stats b { color:var(--text); }
  .stats .sep { color:var(--border); margin:0 6px; }

  /* ===== 资讯分组 ===== */
  .group { margin-bottom:24px; }
  .group-title {
    display:flex; align-items:center; gap:8px;
    font-size:13px; color:var(--text-muted); font-weight:600;
    margin:0 4px 10px; text-transform:uppercase; letter-spacing:0.5px;
  }
  .group-title .dot { width:8px; height:8px; border-radius:50%; flex-shrink:0; }
  .group-title .name { flex:1; }
  .group-title .cnt { font-size:11px; color:var(--text-dim); font-weight:400; }

  /* ===== 资讯卡片 ===== */
  .item {
    background:var(--bg-card); border:1px solid var(--border); border-radius:var(--radius);
    padding:14px 14px; margin-bottom:10px; transition:border-color 0.15s;
  }
  .item.red { border-left:4px solid var(--red); }
  .item.normal { border-left:4px solid var(--border-muted); }
  .item-head {
    display:flex; flex-wrap:wrap; align-items:center; gap:4px 8px;
    font-size:12px; color:var(--text-dim); margin-bottom:8px; line-height:1.8;
  }
  .item-head .time { font-variant-numeric:tabular-nums; white-space:nowrap; }
  .badge {
    display:inline-block; padding:2px 8px; border-radius:4px; font-size:11px; line-height:1.6;
    background:var(--bg-inset); color:var(--text-muted);
  }
  .badge.red { background:rgba(248,81,73,0.12); color:var(--red); font-weight:600; }
  .badge.up { background:rgba(63,185,80,0.12); color:var(--green); }
  .badge.down { background:rgba(248,81,73,0.12); color:var(--red); }
  .badge.neutral { background:rgba(139,148,158,0.12); color:var(--text-muted); }
  .title { font-size:15px; font-weight:600; color:var(--text); margin:0 0 6px; line-height:1.5; }
  .content { font-size:15px; line-height:1.7; color:var(--text); word-break:break-word; overflow-wrap:break-word; }
  .item .url { margin-top:10px; }
  .item .url a { font-size:13px; color:var(--accent); }

  .empty { color:var(--text-muted); text-align:center; padding:60px 20px; font-size:14px; }
  .loading { color:var(--text-muted); text-align:center; padding:40px 20px; font-size:14px; }

  /* ===== 桌面端优化 ===== */
  @media (min-width:768px) {
    body { font-size:16px; }
    .navbar { padding:10px 24px; }
    .navbar .brand { font-size:17px; }
    .filters { padding:10px 24px; gap:8px; }
    .container { padding:24px 24px 80px; }
    .item { padding:16px 20px; margin-bottom:12px; }
    .item:hover { border-color:var(--text-dim); }
    .title { font-size:16px; }
    .content { font-size:15px; }
    .icon-btn.hide-desktop { display:none; }
    .chip { padding:6px 16px; font-size:13px; min-height:auto; }
  }
  /* 超宽屏限制阅读宽度，避免行太长 */
  @media (min-width:1100px) {
    .container { max-width:920px; }
  }
</style>
</head>
<body>

<!-- ===== 顶部导航 ===== -->
<nav class="navbar">
  <button class="icon-btn hide-desktop" id="menuBtn" aria-label="历史日期">☰</button>
  <span class="brand">📰 资讯归档</span>
  <select id="dateSelect"><option>加载中…</option></select>
  <button class="icon-btn" id="refreshBtn" aria-label="刷新">↻</button>
</nav>

<!-- ===== 来源筛选条 ===== -->
<div class="filters" id="filters"></div>

<!-- ===== 日期抽屉（移动端） ===== -->
<div class="drawer" id="drawer">
  <div class="drawer-panel">
    <h2>📅 历史归档</h2>
    <ol id="dayList"></ol>
  </div>
</div>

<!-- ===== 主内容 ===== -->
<div class="container">
  <div id="stats" class="stats"></div>
  <div id="news"></div>
</div>

<script>
const $ = id => document.getElementById(id);
const $news = $('news'), $stats = $('stats'), $filters = $('filters');
const $select = $('dateSelect'), $dayList = $('dayList'), $drawer = $('drawer');
let dates = [], activeDate = null, activeSource = 'ALL', currentArc = null;

// 来源配色
const sourceColor = {
  '财联社电报':'#f85149', '华尔街见闻-全球7x24':'#58a6ff', '华尔街见闻-A股':'#bc8cff',
  '华尔街见闻-美股':'#7ee787', '新浪财经':'#d29922', '外媒':'#8b949e',
  '雪球热门股':'#79c0ff', '雪球热点事件':'#ff7b72', 'TradingView':'#3fb950'
};
const colorFor = s => sourceColor[s] || '#8b949e';

// 情感标签映射
const sentClass = s => {
  if(!s) return '';
  if(/涨|多|好|强|利/.test(s)) return 'up';
  if(/跌|空|差|弱|险/.test(s)) return 'down';
  return 'neutral';
};

// ===== UTC 时间转北京时间显示 =====
function fmtFetched(iso) {
  if(!iso) return '';
  try {
    const d = new Date(iso);
    // 转东八区
    const beijing = new Date(d.getTime() + (8*60 - d.getTimezoneOffset()*-1) * 0); // 由 toLocaleString 处理
    return d.toLocaleString('zh-CN', { timeZone:'Asia/Shanghai', hour12:false });
  } catch(e) { return iso; }
}

// ===== 抽屉开关 =====
$('menuBtn').onclick = () => $drawer.classList.add('open');
$drawer.onclick = e => { if(e.target === $drawer) $drawer.classList.remove('open'); };

// ===== 刷新 =====
$('refreshBtn').onclick = () => { if(activeDate) loadDay(activeDate, true); };

// ===== 加载日期清单 =====
async function loadManifest() {
  try {
    const r = await fetch('manifest.json?t=' + Date.now());
    dates = await r.json();
    if(!dates.length) {
      $select.innerHTML = '<option>暂无归档</option>';
      $dayList.innerHTML = '<li>暂无归档</li>';
      $news.innerHTML = '<div class="empty">还没有归档数据<br>等待第一次抓取完成。</div>';
      return;
    }
    $select.innerHTML = dates.map(d => '<option value="'+d+'">'+d+'</option>').join('');
    $dayList.innerHTML = dates.map(d => '<li data-date="'+d+'"><span>'+d+'</span></li>').join('');
    [...$dayList.children].forEach(li => li.onclick = () => {
      loadDay(li.dataset.date);
      $drawer.classList.remove('open');
    });
    loadDay(dates[0]);
  } catch(e) {
    $news.innerHTML = '<div class="empty">加载 manifest 失败</div>';
  }
}

// ===== 加载某天 =====
async function loadDay(date, isRefresh) {
  if(!date) return;
  activeDate = date;
  $select.value = date;
  [...$dayList.children].forEach(li => li.classList.toggle('active', li.dataset.date === date));
  $news.innerHTML = '<div class="loading">加载中…</div>';
  $stats.textContent = '';
  $filters.innerHTML = '';
  activeSource = 'ALL';
  try {
    const r = await fetch('data/' + date + '.json?t=' + Date.now());
    if(!r.ok) throw new Error('HTTP '+r.status);
    currentArc = await r.json();
    renderStats();
    renderFilters();
    renderNews();
  } catch(e) {
    $news.innerHTML = '<div class="empty">加载失败：'+e.message+'</div>';
  }
}

function renderStats() {
  const items = currentArc.items || [];
  const red = items.filter(i => i.isRed).length;
  $stats.innerHTML =
    '📅 <b>' + currentArc.date + '</b>' +
    '<span class="sep">·</span>共 <b>' + items.length + '</b> 条' +
    '<span class="sep">·</span>🔴 重要 <b>' + red + '</b> 条' +
    '<br><span style="font-size:12px;color:var(--text-dim)">抓取于 ' + fmtFetched(currentArc.fetched_at) + ' (北京时间)</span>';
}

// ===== 渲染来源筛选条 =====
function renderFilters() {
  const items = currentArc.items || [];
  const groups = {};
  items.forEach(i => { const s = i.source || '其他'; groups[s] = (groups[s]||0)+1; });
  const sources = Object.keys(groups).sort((a,b) => groups[b]-groups[a]);
  let html = '<div class="chip '+(activeSource==='ALL'?'active':'')+'" data-src="ALL">全部 <span class="count">'+items.length+'</span></div>';
  html += sources.map(s =>
    '<div class="chip '+(activeSource===s?'active':'')+'" data-src="'+escapeAttr(s)+'">'+
    '<span style="display:inline-block;width:7px;height:7px;border-radius:50%;background:'+colorFor(s)+';margin-right:4px;vertical-align:middle"></span>'+
    escapeHtml(s)+' <span class="count">'+groups[s]+'</span></div>'
  ).join('');
  $filters.innerHTML = html;
  [...$filters.children].forEach(c => c.onclick = () => {
    activeSource = c.dataset.src;
    renderFilters();
    renderNews();
  });
}

// ===== 渲染资讯列表（按来源分组，重要资讯置顶） =====
function renderNews() {
  const items = currentArc.items || [];
  let filtered = items;
  if(activeSource !== 'ALL') filtered = items.filter(i => (i.source||'其他') === activeSource);
  if(!filtered.length) {
    $news.innerHTML = '<div class="empty">该来源当天无资讯</div>';
    return;
  }
  // 重要（isRed）资讯置顶，组内保持原有时间倒序
  const sortRedFirst = arr => {
    const red = arr.filter(i => i.isRed);
    const normal = arr.filter(i => !i.isRed);
    return red.concat(normal);
  };
  // 按来源分组（每组内重要置顶）
  const groups = {};
  filtered.forEach(i => { const s = i.source || '其他'; (groups[s]=groups[s]||[]).push(i); });
  Object.keys(groups).forEach(s => { groups[s] = sortRedFirst(groups[s]); });
  // 组顺序：按条数倒序
  const order = Object.keys(groups).sort((a,b) => groups[b].length - groups[a].length);
  let html = '';
  order.forEach(src => {
    html += '<div class="group">';
    html += '<div class="group-title">'+
      '<span class="dot" style="background:'+colorFor(src)+'"></span>'+
      '<span class="name">'+escapeHtml(src)+'</span>'+
      '<span class="cnt">'+groups[src].length+' 条</span></div>';
    html += groups[src].map(renderItem).join('');
    html += '</div>';
  });
  $news.innerHTML = html;
}

function renderItem(i) {
  const title = i.title || (i.content || '').slice(0, 50);
  const cls = i.isRed ? 'red' : 'normal';
  const sc = sentClass(i.sentiment);
  const sentBadge = i.sentiment ? '<span class="badge '+sc+'">'+escapeHtml(i.sentiment)+'</span>' : '';
  const redBadge = i.isRed ? '<span class="badge red">重要</span>' : '';
  const url = i.url ? '<div class="url"><a href="'+escapeAttr(i.url)+'" target="_blank" rel="noopener">查看原文 →</a></div>' : '';
  // 内容：若与标题相同则不重复显示
  const showContent = i.content && i.content !== i.title && i.content.trim() !== title.trim();
  const content = showContent ? '<div class="content">'+escapeHtml(i.content)+'</div>' : '';
  return '<div class="item '+cls+'">' +
    '<div class="item-head">' +
      '<span class="time">'+escapeHtml(i.time||'')+'</span>' +
      redBadge + sentBadge +
    '</div>' +
    '<div class="title">'+escapeHtml(title)+'</div>' +
    content + url +
  '</div>';
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
}
function escapeAttr(s) {
  return String(s).replace(/"/g, '&quot;');
}

$select.onchange = () => loadDay($select.value);
loadManifest();
</script>
</body>
</html>
`