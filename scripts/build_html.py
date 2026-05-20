#!/usr/bin/env python3
"""
build_html.py — GableLBM Wiki HTML Compiler
Compiles docs/wiki/Architecture.md into a styled standalone HTML file
at docs/wiki/dist/Architecture.html.

Usage: python3 scripts/build_html.py
Requires: Python 3.10+; no external deps (stdlib only).
"""

import re
import html
from pathlib import Path
from datetime import datetime, timezone

REPO_ROOT = Path(__file__).parent.parent
WIKI_DIR = REPO_ROOT / "docs" / "wiki"
DIST_DIR = WIKI_DIR / "compiled"
SOURCE = WIKI_DIR / "Architecture.md"
OUTPUT = DIST_DIR / "Architecture.html"

CSS = """
:root {
  --bg: #0A0B10;
  --card: #161821;
  --green: #00FFA3;
  --red: #F43F5E;
  --blue: #38BDF8;
  --text: #e2e8f0;
  --muted: #94a3b8;
  --border: rgba(255,255,255,0.1);
  --font-ui: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', monospace;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  background: var(--bg); color: var(--text);
  font-family: var(--font-ui); font-size: 14px; line-height: 1.7;
  padding: 0 0 80px;
}
header {
  background: var(--card); border-bottom: 1px solid var(--border);
  padding: 20px 40px; display: flex; align-items: center; gap: 16px;
}
header .logo { color: var(--green); font-weight: 700; font-size: 18px; letter-spacing: -0.5px; }
header .badge { background: rgba(0,255,163,0.1); color: var(--green);
  font-size: 11px; padding: 2px 8px; border-radius: 999px; font-family: var(--font-mono); }
header .ts { color: var(--muted); font-size: 12px; margin-left: auto; font-family: var(--font-mono); }
main { max-width: 1200px; margin: 40px auto; padding: 0 40px; }
h1 { font-size: 28px; font-weight: 700; color: #fff; margin-bottom: 8px; }
h2 { font-size: 18px; font-weight: 600; color: var(--green); margin: 40px 0 16px;
  padding-bottom: 8px; border-bottom: 1px solid var(--border); }
h3 { font-size: 15px; font-weight: 600; color: var(--blue); margin: 24px 0 12px; }
h4 { font-size: 13px; font-weight: 600; color: var(--muted); margin: 16px 0 8px; text-transform: uppercase; letter-spacing: 0.08em; }
p { margin: 8px 0; color: var(--text); }
a { color: var(--blue); text-decoration: none; }
a:hover { text-decoration: underline; }
table { width: 100%; border-collapse: collapse; margin: 16px 0; font-size: 13px; }
thead { background: rgba(255,255,255,0.04); }
th { padding: 10px 12px; text-align: left; color: var(--muted); font-weight: 600;
  border-bottom: 1px solid var(--border); font-size: 11px; text-transform: uppercase; letter-spacing: 0.06em; }
td { padding: 9px 12px; border-bottom: 1px solid rgba(255,255,255,0.05); vertical-align: top; }
tr:hover td { background: rgba(255,255,255,0.02); }
code { font-family: var(--font-mono); font-size: 12px;
  background: rgba(255,255,255,0.08); padding: 1px 5px; border-radius: 3px; }
pre { background: var(--card); border: 1px solid var(--border); border-radius: 6px;
  padding: 16px; overflow-x: auto; margin: 12px 0; }
pre code { background: none; padding: 0; font-size: 13px; }
blockquote { border-left: 3px solid var(--green); padding-left: 16px; color: var(--muted);
  margin: 16px 0; font-style: italic; }
ul, ol { padding-left: 20px; margin: 8px 0; }
li { margin: 4px 0; }
strong { color: #fff; font-weight: 600; }
.pill { display: inline-block; font-size: 11px; padding: 1px 7px; border-radius: 999px; font-weight: 600; }
.pill.green { background: rgba(0,255,163,0.1); color: var(--green); }
.pill.red { background: rgba(244,63,94,0.1); color: var(--red); }
.pill.blue { background: rgba(56,189,248,0.1); color: var(--blue); }
footer { text-align: center; color: var(--muted); font-size: 12px; margin-top: 60px;
  padding-top: 24px; border-top: 1px solid var(--border); font-family: var(--font-mono); }
"""

def md_to_html(md: str) -> str:
    """Minimal Markdown → HTML converter for the wiki compiler."""
    lines = md.split("\n")
    out = []
    in_table = False
    in_code = False
    in_list = False
    in_blockquote = False
    code_buf = []

    def flush_list():
        nonlocal in_list
        if in_list:
            out.append("</ul>")
            in_list = False

    def flush_blockquote():
        nonlocal in_blockquote
        if in_blockquote:
            out.append("</blockquote>")
            in_blockquote = False

    def inline(text: str) -> str:
        # code
        text = re.sub(r"`([^`]+)`", lambda m: f"<code>{html.escape(m.group(1))}</code>", text)
        # bold
        text = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", text)
        # italic
        text = re.sub(r"\*(.+?)\*", r"<em>\1</em>", text)
        # link
        text = re.sub(r"\[(.+?)\]\((.+?)\)", r'<a href="\2">\1</a>', text)
        # ✅ 🔧 emojis — keep as-is (already unicode)
        return text

    i = 0
    while i < len(lines):
        line = lines[i]

        # Fenced code block
        if line.strip().startswith("```"):
            if in_code:
                in_code = False
                code_html = html.escape("\n".join(code_buf))
                lang = ""
                out.append(f"<pre><code>{code_html}</code></pre>")
                code_buf = []
            else:
                flush_list()
                flush_blockquote()
                in_code = True
            i += 1
            continue

        if in_code:
            code_buf.append(line)
            i += 1
            continue

        # Table row detection
        if "|" in line and line.strip().startswith("|"):
            flush_list()
            flush_blockquote()
            if not in_table:
                out.append('<table>')
                in_table = True
                # Header row
                cols = [c.strip() for c in line.strip().strip("|").split("|")]
                out.append("<thead><tr>" + "".join(f"<th>{inline(c)}</th>" for c in cols) + "</tr></thead><tbody>")
            elif re.match(r"\|[-| :]+\|", line):
                pass  # separator row — skip
            else:
                cols = [c.strip() for c in line.strip().strip("|").split("|")]
                out.append("<tr>" + "".join(f"<td>{inline(c)}</td>" for c in cols) + "</tr>")
            i += 1
            continue
        elif in_table:
            out.append("</tbody></table>")
            in_table = False

        # Blockquote
        if line.startswith(">"):
            flush_list()
            if not in_blockquote:
                out.append("<blockquote>")
                in_blockquote = True
            out.append(f"<p>{inline(line[1:].strip())}</p>")
            i += 1
            continue
        else:
            flush_blockquote()

        # Headings
        if line.startswith("# "):
            flush_list()
            out.append(f"<h1>{inline(line[2:])}</h1>")
        elif line.startswith("## "):
            flush_list()
            out.append(f"<h2>{inline(line[3:])}</h2>")
        elif line.startswith("### "):
            flush_list()
            out.append(f"<h3>{inline(line[4:])}</h3>")
        elif line.startswith("#### "):
            flush_list()
            out.append(f"<h4>{inline(line[5:])}</h4>")
        # Horizontal rule
        elif line.strip() in ("---", "***", "___"):
            flush_list()
            out.append("<hr>")
        # List item
        elif line.strip().startswith("- ") or line.strip().startswith("* "):
            if not in_list:
                out.append("<ul>")
                in_list = True
            out.append(f"<li>{inline(line.strip()[2:])}</li>")
        # Ordered list
        elif re.match(r"^\d+\. ", line.strip()):
            if not in_list:
                out.append("<ul>")
                in_list = True
            content = re.sub(r"^\d+\. ", "", line.strip())
            out.append(f"<li>{inline(content)}</li>")
        # Empty line
        elif line.strip() == "":
            flush_list()
            out.append("")
        # Paragraph
        else:
            flush_list()
            out.append(f"<p>{inline(line)}</p>")

        i += 1

    flush_list()
    flush_blockquote()
    if in_table:
        out.append("</tbody></table>")

    return "\n".join(out)


def build():
    DIST_DIR.mkdir(parents=True, exist_ok=True)

    if not SOURCE.exists():
        print(f"❌ Source not found: {SOURCE}")
        raise SystemExit(1)

    md = SOURCE.read_text(encoding="utf-8")
    body = md_to_html(md)
    ts = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC")

    html_out = f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GableLBM — Architecture Wiki</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono&display=swap" rel="stylesheet">
  <style>{CSS}</style>
</head>
<body>
  <header>
    <span class="logo">GableLBM ERP</span>
    <span class="badge">Architecture Wiki</span>
    <span class="ts">Auto-generated · {ts}</span>
  </header>
  <main>
    {body}
  </main>
  <footer>
    GableLBM Architecture Wiki · Auto-compiled by wiki-refresh · {ts}
  </footer>
</body>
</html>"""

    OUTPUT.write_text(html_out, encoding="utf-8")
    size_kb = OUTPUT.stat().st_size / 1024
    print(f"✅ HTML compiled: {OUTPUT} ({size_kb:.1f} KB)")


if __name__ == "__main__":
    build()
