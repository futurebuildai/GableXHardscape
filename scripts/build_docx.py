#!/usr/bin/env python3
"""
build_docx.py — GableLBM Wiki Docx Compiler
Compiles docs/wiki/Architecture.md into a Word document (.docx)
at docs/wiki/dist/Architecture.docx using only Python stdlib
(produces a valid Open XML / .docx file without python-docx dependency).

Usage: python3 scripts/build_docx.py
Requires: Python 3.10+; no external deps.
"""

import re
import zipfile
import io
import xml.etree.ElementTree as ET
from pathlib import Path
from datetime import datetime, timezone

REPO_ROOT = Path(__file__).parent.parent
WIKI_DIR = REPO_ROOT / "docs" / "wiki"
DIST_DIR = WIKI_DIR / "compiled"
SOURCE = WIKI_DIR / "Architecture.md"
OUTPUT = DIST_DIR / "Architecture.docx"

# ─── Minimal Open XML (OOXML) templates ────────────────────────────────────

CONTENT_TYPES = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
</Types>"""

RELS = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>"""

WORD_RELS = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>"""

STYLES = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:style w:type="paragraph" w:styleId="Normal" w:default="1">
    <w:name w:val="Normal"/>
    <w:pPr><w:spacing w:after="160"/></w:pPr>
    <w:rPr><w:sz w:val="22"/><w:szCs w:val="22"/>
      <w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading1">
    <w:name w:val="heading 1"/>
    <w:pPr><w:outlineLvl w:val="0"/><w:spacing w:before="240" w:after="120"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="40"/><w:szCs w:val="40"/>
      <w:color w:val="00CC82"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading2">
    <w:name w:val="heading 2"/>
    <w:pPr><w:outlineLvl w:val="1"/><w:spacing w:before="200" w:after="80"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="28"/><w:szCs w:val="28"/>
      <w:color w:val="00CC82"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading3">
    <w:name w:val="heading 3"/>
    <w:pPr><w:outlineLvl w:val="2"/><w:spacing w:before="160" w:after="60"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="24"/><w:szCs w:val="24"/>
      <w:color w:val="38BDF8"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Code">
    <w:name w:val="Code"/>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="1A1B26"/></w:pPr>
    <w:rPr>
      <w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>
      <w:sz w:val="18"/><w:szCs w:val="18"/>
      <w:color w:val="A9B1D6"/>
    </w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="TableHeader">
    <w:name w:val="TableHeader"/>
    <w:rPr><w:b/><w:sz w:val="20"/><w:color w:val="94A3B8"/></w:rPr>
  </w:style>
</w:styles>"""

NS = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
W = f"{{{NS}}}"


def esc(text: str) -> str:
    """Escape XML special chars."""
    return (text
            .replace("&", "&amp;")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
            .replace('"', "&quot;"))


def make_run(text: str, bold: bool = False, code: bool = False, color: str = None) -> str:
    """Return a <w:r> run XML fragment."""
    rpr_parts = []
    if bold:
        rpr_parts.append("<w:b/>")
    if code:
        rpr_parts.append('<w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>')
        rpr_parts.append("<w:sz w:val=\"18\"/>")
    if color:
        rpr_parts.append(f'<w:color w:val="{color}"/>')

    rpr = f"<w:rPr>{''.join(rpr_parts)}</w:rPr>" if rpr_parts else ""
    return f"<w:r>{rpr}<w:t xml:space=\"preserve\">{esc(text)}</w:t></w:r>"


def inline_runs(text: str) -> str:
    """Convert inline markdown (bold, code) to OOXML runs."""
    out = []
    # Split on **bold** and `code`
    pattern = re.compile(r"\*\*(.+?)\*\*|`([^`]+)`")
    last = 0
    for m in pattern.finditer(text):
        if m.start() > last:
            out.append(make_run(text[last:m.start()]))
        if m.group(1):
            out.append(make_run(m.group(1), bold=True))
        elif m.group(2):
            out.append(make_run(m.group(2), code=True))
        last = m.end()
    if last < len(text):
        out.append(make_run(text[last:]))
    return "".join(out)


def para(content: str, style: str = "Normal") -> str:
    return f'<w:p><w:pPr><w:pStyle w:val="{style}"/></w:pPr>{content}</w:p>'


def table_row(cells: list[str], header: bool = False) -> str:
    tcs = ""
    for c in cells:
        style = "TableHeader" if header else "Normal"
        cell_content = f'<w:p><w:pPr><w:pStyle w:val="{style}"/></w:pPr>{inline_runs(c)}</w:p>'
        tcs += f"<w:tc><w:tcPr><w:tcW w:w=\"0\" w:type=\"auto\"/></w:tcPr>{cell_content}</w:tc>"
    trpr = "<w:trPr><w:tblHeader/></w:trPr>" if header else ""
    return f"<w:tr>{trpr}{tcs}</w:tr>"


def md_to_ooxml(md: str) -> str:
    """Convert Markdown to OOXML body paragraphs."""
    lines = md.split("\n")
    body_parts = []
    in_code = False
    in_table = False
    code_buf = []

    i = 0
    while i < len(lines):
        line = lines[i]

        # Fenced code block
        if line.strip().startswith("```"):
            if in_code:
                in_code = False
                for cl in code_buf:
                    body_parts.append(para(make_run(cl, code=True), "Code"))
                code_buf = []
            else:
                in_code = True
            i += 1
            continue

        if in_code:
            code_buf.append(line)
            i += 1
            continue

        # Table
        if "|" in line and line.strip().startswith("|"):
            if not in_table:
                in_table = True
                body_parts.append("<w:tbl><w:tblPr><w:tblW w:w=\"0\" w:type=\"auto\"/></w:tblPr>")
                cols = [c.strip() for c in line.strip().strip("|").split("|")]
                body_parts.append(table_row(cols, header=True))
            elif re.match(r"\|[-| :]+\|", line):
                pass  # separator
            else:
                cols = [c.strip() for c in line.strip().strip("|").split("|")]
                body_parts.append(table_row(cols, header=False))
            i += 1
            continue
        elif in_table:
            body_parts.append("</w:tbl>")
            in_table = False

        # Headings
        if line.startswith("# "):
            body_parts.append(para(inline_runs(line[2:]), "Heading1"))
        elif line.startswith("## "):
            body_parts.append(para(inline_runs(line[3:]), "Heading2"))
        elif line.startswith("### "):
            body_parts.append(para(inline_runs(line[4:]), "Heading3"))
        elif line.strip() in ("---", "***"):
            body_parts.append('<w:p><w:pPr><w:pBdr><w:bottom w:val="single" w:sz="6" w:color="333344"/></w:pBdr></w:pPr></w:p>')
        elif line.strip().startswith("- ") or line.strip().startswith("* "):
            content = line.strip()[2:]
            body_parts.append(f'<w:p><w:pPr><w:pStyle w:val="Normal"/><w:numPr><w:ilvl w:val="0"/><w:numId w:val="1"/></w:numPr></w:pPr>{inline_runs(content)}</w:p>')
        elif re.match(r"^\d+\. ", line.strip()):
            content = re.sub(r"^\d+\. ", "", line.strip())
            body_parts.append(para(inline_runs(content), "Normal"))
        elif line.strip() == "":
            body_parts.append("<w:p/>")
        else:
            body_parts.append(para(inline_runs(line)))

        i += 1

    if in_table:
        body_parts.append("</w:tbl>")

    return "\n".join(body_parts)


def build():
    DIST_DIR.mkdir(parents=True, exist_ok=True)

    if not SOURCE.exists():
        print(f"❌ Source not found: {SOURCE}")
        raise SystemExit(1)

    md = SOURCE.read_text(encoding="utf-8")
    ts = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC")

    body_xml = md_to_ooxml(md)
    doc_xml = f"""<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    {body_xml}
    <w:p>
      <w:pPr><w:pStyle w:val="Normal"/></w:pPr>
      <w:r><w:rPr><w:color w:val="94A3B8"/><w:sz w:val="16"/></w:rPr>
        <w:t>Auto-compiled by wiki-refresh · {esc(ts)}</w:t>
      </w:r>
    </w:p>
    <w:sectPr>
      <w:pgSz w:w="12240" w:h="15840"/>
      <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"/>
    </w:sectPr>
  </w:body>
</w:document>"""

    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.writestr("[Content_Types].xml", CONTENT_TYPES)
        zf.writestr("_rels/.rels", RELS)
        zf.writestr("word/_rels/document.xml.rels", WORD_RELS)
        zf.writestr("word/styles.xml", STYLES)
        zf.writestr("word/document.xml", doc_xml)

    OUTPUT.write_bytes(buf.getvalue())
    size_kb = OUTPUT.stat().st_size / 1024
    print(f"✅ Docx compiled: {OUTPUT} ({size_kb:.1f} KB)")


if __name__ == "__main__":
    build()
