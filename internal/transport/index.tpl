<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width,initial-scale=1" />
<title>Meilisearch MCP v{{ .Version }}</title>
<meta name="description" content="Meilisearch MCP (Model Context Protocol) HTTP transport landing page" />
<link rel="icon" href="https://www.meilisearch.com/favicon.ico" />
<style>
:root { --bg:#ffffff; --fg:#1b1f23; --accent:#0047ff; --accent-glow:#8ab4ff; --border:#e2e8f0; --muted:#566170; --code-bg:#f5f7fa; --radius:12px; --brand:#ff5caa; --brand2:#7b16ff; }
html,body { margin:0; padding:0; font-family: system-ui,-apple-system,Segoe UI,Roboto,Ubuntu,Cantarell,'Fira Sans','Droid Sans','Helvetica Neue',Arial,sans-serif; background:var(--bg); color:var(--fg); -webkit-font-smoothing: antialiased; }
body { line-height:1.55; }
a { color:var(--accent); text-decoration:none; }
a:hover { text-decoration:underline; }
header { padding:clamp(1.8rem,4vw,3rem) clamp(1.2rem,5vw,4rem) 1rem; }
main { padding:0 clamp(1.2rem,5vw,4rem) 4rem; max-width:1100px; margin:0 auto; }
footer { border-top:1px solid var(--border); padding:2rem clamp(1.2rem,5vw,4rem); font-size:.75rem; color:var(--muted); }
h1 { font-size:clamp(2.2rem,4vw,3.4rem); line-height:1.05; letter-spacing:-.02em; margin:.2rem 0 1.2rem; background:linear-gradient(90deg,var(--brand),var(--brand2)); -webkit-background-clip: text; color:transparent; }
.badge { display:inline-block; background:var(--brand); color:#fff; font-size:.6em; padding:.35em .6em .3em; border-radius:8px; vertical-align:middle; font-weight:600; letter-spacing:.05em; }
h2 { margin-top:3.2rem; font-size:1.45rem; letter-spacing:-.01em; }
h3 { margin-top:2.4rem; font-size:1.05rem; text-transform:uppercase; letter-spacing:.08em; font-weight:600; color:var(--muted); }
p { max-width:70ch; }
pre,code { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace; font-size:.85rem; }
pre { position:relative; background:var(--code-bg); padding:1rem 1.2rem; border-radius:var(--radius); overflow:auto; box-shadow:0 2px 4px -2px rgba(0,0,0,.25); border:1px solid var(--border); }
code { background:var(--code-bg); padding:.15em .35em; border-radius:6px; }
.hero-lede { font-size:clamp(1.05rem,1.6vw,1.25rem); font-weight:500; max-width:62ch; }
ul { padding-left:1.1rem; }
li { margin:.3rem 0 .3rem; }
.btn-row { display:flex; gap:.75rem; margin:1.6rem 0 2rem; flex-wrap:wrap; }
button, .btn { cursor:pointer; border:1px solid var(--border); background:var(--code-bg); color:var(--fg); padding:.7rem 1.05rem; font-size:.78rem; font-weight:500; border-radius:10px; line-height:1; display:inline-flex; gap:.45rem; align-items:center; transition:.18s background,.18s border-color,.18s color; }
button:hover,.btn:hover { background:var(--accent); color:#fff; border-color:var(--accent); text-decoration:none; }
.toggle { position:fixed; top:.9rem; right:.9rem; z-index:10; }
.github-link { position:fixed; top:.9rem; right:4.4rem; z-index:15; display:inline-flex; align-items:center; justify-content:center; width:34px; height:34px; border:1px solid var(--border); border-radius:10px; background:var(--code-bg); color:var(--fg); }
.github-link:hover { background:var(--accent); color:#fff; border-color:var(--accent); text-decoration:none; }
:root.dark { --bg:#0b0d11; --fg:#f5f7fa; --accent:#5b8dff; --accent-glow:#b1cfff; --border:#242a33; --muted:#9aa7ba; --code-bg:#151a21; }
.dark pre,.dark code { box-shadow:0 2px 4px -2px rgba(0,0,0,.6); }
.fade-in { animation:fade .7s ease; }
@keyframes fade { from { opacity:0; transform:translateY(4px);} to { opacity:1; transform:translateY(0);} }
.notice { border-left:4px solid var(--accent); padding:.9rem 1rem; background:var(--code-bg); border-radius:0 var(--radius) var(--radius) 0; margin:1.6rem 0 2.2rem; }
.small { font-size:.7rem; letter-spacing:.03em; }
.mark { background:linear-gradient(90deg,var(--accent-glow),var(--accent)); -webkit-background-clip: text; color:transparent; font-weight:600; }
.copy-box { margin:1.4rem 0 2.2rem; }
.copy-row { display:flex; gap:.6rem; max-width:560px; }
.copy-row input { flex:1; padding:.65rem .8rem; border:1px solid var(--border); border-radius:10px; background:var(--code-bg); font-size:.8rem; font-family:inherit; color:var(--fg); }
.copy-row button { white-space:nowrap; }
.copy-row button.copied { background:#16a34a; border-color:#16a34a; color:#fff; }
.code-wrap { position:relative; margin:1.4rem 0 2.1rem; }
.copy-btn { position:absolute; top:.45rem; right:.45rem; background:var(--code-bg); border:1px solid var(--border); padding:.35rem .55rem; border-radius:8px; font-size:.65rem; line-height:1; display:inline-flex; align-items:center; gap:.35rem; }
.copy-btn svg { width:14px; height:14px; stroke:currentColor; }
.copy-btn:hover { background:var(--accent); color:#fff; border-color:var(--accent); }
.copy-btn.copied { background:#16a34a; border-color:#16a34a; color:#fff; }
.hl .tok-key { color:#b45309; font-weight:600; }
.hl .tok-str { color:#0369a1; }
.dark .hl .tok-str { color:#38bdf8; }
.hl .tok-num { color:#7c2d12; }
.dark .hl .tok-num { color:#fbbf24; }
.hl .tok-bool { color:#0f766e; font-weight:600; }
.dark .hl .tok-bool { color:#34d399; }
.hl .tok-null { color:#6366f1; font-style:italic; }
.hl .tok-sym { color:#64748b; }
/* new config grid styles */
.config-grid { display:grid; gap:1.6rem; margin:1.4rem 0 2.2rem; }
@media (min-width:900px){ .config-grid { grid-template-columns:1fr 1fr; align-items:start; } }
.config-card { background:var(--code-bg); border:1px solid var(--border); padding:1rem 1.1rem 1.3rem; border-radius:var(--radius); box-shadow:0 2px 4px -2px rgba(0,0,0,.25); }
.dark .config-card { box-shadow:0 2px 4px -2px rgba(0,0,0,.6); }
.config-card h3 { margin-top:0; }
.config-card pre { margin:0.75rem 0 0; }
</style>
<script>
const toggleTheme = () => { const root=document.documentElement; if(root.classList.contains('dark')){ root.classList.remove('dark'); localStorage.setItem('prefers-dark','0'); } else { root.classList.add('dark'); localStorage.setItem('prefers-dark','1'); }};
window.addEventListener('DOMContentLoaded',()=>{ if(localStorage.getItem('prefers-dark')==='1'){ document.documentElement.classList.add('dark'); } enhanceCodeBlocks(); });
async function ping() { const t0=performance.now(); try { await fetch("/healthz",{method:'HEAD'}); const ms=(performance.now()-t0).toFixed(1); document.getElementById('latency').textContent=ms+' ms'; } catch(e){ document.getElementById('latency').textContent='n/a'; } }
function copyURL(){ const inp=document.getElementById('mcp-url'); inp.select(); inp.setSelectionRange(0,99999); try { navigator.clipboard.writeText(inp.value); } catch(e){} flashButton(document.getElementById('url-copy-btn')); }
function copySSEURL(){ const inp=document.getElementById('sse-url'); inp.select(); inp.setSelectionRange(0,99999); try { navigator.clipboard.writeText(inp.value); } catch(e){} flashButton(document.getElementById('url-sse-copy-btn')); }
function flashButton(btn){ if(!btn) return; btn.classList.add('copied'); const orig=btn.dataset.label||btn.textContent; const span=btn.querySelector('span'); if(span){ span.textContent='Copied'; } else { btn.textContent='Copied'; } setTimeout(()=>{ if(span){ span.textContent=orig; } else { btn.textContent=orig; } btn.classList.remove('copied'); },1400); }
function enhanceCodeBlocks(){ document.querySelectorAll('pre code').forEach(code=>{ const wrapper=document.createElement('div'); wrapper.className='code-wrap'; const pre=code.parentElement; pre.parentElement.insertBefore(wrapper,pre); wrapper.appendChild(pre); const raw=code.textContent.trim(); const highlighted=highlightJSONLike(raw); code.innerHTML=highlighted; code.classList.add('hl'); const btn=document.createElement('button'); btn.type='button'; btn.className='copy-btn'; btn.dataset.label='Copy'; btn.innerHTML='<svg fill="none" stroke-width="1.6" stroke="currentColor" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg><span>Copy</span>'; btn.addEventListener('click',()=>{ try { navigator.clipboard.writeText(raw);} catch(e){} flashButton(btn); }); wrapper.appendChild(btn); }); }
function highlightJSONLike(src){ try { const obj=JSON.parse(src); let json=JSON.stringify(obj, null, 2); json=json.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); const keyStore=[]; json=json.replace(/("(?:\\"|[^"\\])*?")(?=:\s)/g,function(m){ var i=keyStore.length; keyStore.push(m); return '@@K'+i+'@@'; }); json=json.replace(/"(?:\\"|[^"\\])*?"/g,function(m){ return '<span class="tok-str">'+m+'</span>'; }); json=json.replace(/\b-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?\b/g,function(m){ return '<span class="tok-num">'+m+'</span>'; }); json=json.replace(/\b(true|false)\b/g,function(m){ return '<span class="tok-bool">'+m+'</span>'; }); json=json.replace(/\bnull\b/g,'<span class="tok-null">null</span>'); json=json.replace(/[{}\[\],]/g,function(m){ return '<span class="tok-sym">'+m+'</span>'; }); keyStore.forEach(function(k,i){ var escaped=k.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); json=json.replace('@@K'+i+'@@','<span class="tok-key">'+escaped+'</span>'); }); return json; } catch(e) { return src.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); } }
</script>
</head>
<body class="fade-in">
<a class="github-link" href="https://github.com/ja7ad/meilisearch-mcp" target="_blank" aria-label="GitHub Repository" rel="noopener">
  <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.009-.866-.014-1.7-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.53 2.341 1.088 2.91.833.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.563 9.563 0 0 1 12 6.844a9.56 9.56 0 0 1 2.504.337c1.909-1.296 2.748-1.026 2.748-1.026.546 1.378.203 2.397.1 2.65.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.31.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0 0 22 12.017C22 6.484 17.523 2 12 2Z" clip-rule="evenodd"/></svg>
</a>
<button class="toggle" onclick="toggleTheme()" aria-label="Toggle dark mode">🌗</button>
<header>
  <h1>Meilisearch MCP v{{ .Version }}</h1>
  <p class="hero-lede">This server speaks the <span class="mark">Model Context Protocol (MCP)</span> over simple HTTP streaming. Use it to connect AI tooling (desktop & cloud) with your Meilisearch data.</p>
  <div class="btn-row">
    <a class="btn" href="https://github.com/modelcontextprotocol" target="_blank" rel="noopener">MCP Spec ↗</a>
    <a class="btn" href="https://www.meilisearch.com" target="_blank" rel="noopener">Meilisearch ↗</a>
    <button onclick="ping()">Ping <span id="latency" class="small" style="opacity:.8">...</span></button>
  </div>
  <div class="copy-box">
    <h3 style="margin-top:0;">MCP Endpoint URL</h3>
    <div class="copy-row">
      <input id="mcp-url" value="https://meilisearch.javad.dev/mcp" readonly aria-label="MCP Endpoint URL" />
      <button id="url-copy-btn" onclick="copyURL()" data-label="Copy" class="btn" style="padding:.55rem .85rem;">
        <svg fill="none" stroke-width="1.6" stroke="currentColor" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
        <span>Copy</span>
      </button>
    </div>
  </div>
  {{ if .EnableSSE }}
  <div class="copy-box">
      <h3 style="margin-top:0;">SSE Endpoint URL</h3>
      <div class="copy-row">
        <input id="sse-url" value="https://meilisearch.javad.dev/sse" readonly aria-label="SSE Endpoint URL" />
        <button id="url-sse-copy-btn" onclick="copySSEURL()" data-label="Copy" class="btn" style="padding:.55rem .85rem;">
          <svg fill="none" stroke-width="1.6" stroke="currentColor" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
          <span>Copy</span>
        </button>
      </div>
    </div>
  {{ end }}
  <div class="notice small">Root path is human friendly; POST JSON MCP requests (newline-delimited streaming) to the same endpoint. Non-root paths also route here.</div>
</header>
<main>
  <section>
    <h2>Supported Clients & Configuration</h2>
    <p>Below are example configurations & notes for MCP-capable clients (install / enable as applicable):</p>
    <ul class="small" style="list-style:disc; padding-left:1.2rem;">
      <li><strong>Jan</strong> (desktop) – native MCP provider config.</li>
      <li><strong>Claude Desktop</strong> – add via tool integration (future MCP support).</li>
      <li><strong>Cursor / VS Code (Continue)</strong> – configure remote/stdio MCP backend.</li>
      <li><strong>Zed / JetBrains (plugins)</strong> – emerging MCP adopters.</li>
      <li><strong>Custom</strong> – use <code>mcp-remote</code> CLI or direct HTTP POST.</li>
    </ul>
    <div class="config-grid">
      <div class="config-card">
        <h3 style="margin-top:0;">Generic HTTP (mcp-remote)</h3>
        <pre><code>{
  "command": "npx",
  "args": ["-y","mcp-remote@latest","https://meilisearch.javad.dev/mcp", "--header","X-Meili-Instance: ${MEILISEARCH_INSTANCE}", "--header","X-Meili-APIKey: ${MEILISEARCH_API_KEY}"],
  "env": {"MEILISEARCH_INSTANCE": "http://localhost:7700", "MEILISEARCH_API_KEY": "masterKey"},
  "active": true
}</code></pre>
      </div>
      <div class="config-card">
        <h3 style="margin-top:0;">Local STDIO</h3>
        <pre><code>{
  "command": "/usr/bin/meilisearch-mcp",
  "args": ["serve", "stdio","--meili-host","http://localhost:7700","--meili-api-key","masterKey"],
  "env": {},
  "active": false
}</code></pre>
      </div>
    </div>
    <p class="small">Flip <code>active</code> flags to select transport. Prefer stdio locally; HTTP for remote/container usage.</p>
  </section>
  <section>
    <h2>Security Notes</h2>
    <ul>
      <li>Place behind TLS (reverse proxy) when exposed publicly.</li>
      <li>Forward required auth headers only; strip unknown ones.</li>
      <li>Enforce payload size limits & rate limiting at the proxy.</li>
    </ul>
  </section>
</main>
<footer>
  <div>Meilisearch MCP • Open Source • <span id="year"></span></div>
  <script>document.getElementById('year').textContent=new Date().getFullYear();</script>
</footer>
</body>
</html>