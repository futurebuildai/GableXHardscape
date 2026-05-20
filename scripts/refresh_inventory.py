#!/usr/bin/env python3
"""
refresh_inventory.py — GableLBM Wiki Inventory Scanner
Parses Go module boundaries, migration count, and Lit component list
to generate an up-to-date snapshot for docs/wiki/Architecture.md.

Usage: python3 scripts/refresh_inventory.py
Output: prints summary; Architecture.md is updated by the wiki-refresh workflow.
"""

import os
import subprocess
import sys
from pathlib import Path
from datetime import datetime, timezone

REPO_ROOT = Path(__file__).parent.parent
DOCS_WIKI = REPO_ROOT / "docs" / "wiki"
BACKEND_INTERNAL = REPO_ROOT / "backend" / "internal"
MIGRATIONS = REPO_ROOT / "backend" / "migrations"
APP_PAGES = REPO_ROOT / "app" / "src" / "pages"
APP_COMPONENTS = REPO_ROOT / "app" / "src" / "components"


def count_migrations() -> int:
    """Count all .sql migration files."""
    return len(list(MIGRATIONS.glob("*.sql")))


def list_modules() -> list[dict]:
    """Find all internal packages with handler.go files."""
    modules = []
    for d in sorted(BACKEND_INTERNAL.iterdir()):
        if not d.is_dir():
            continue
        has_handler = (d / "handler.go").exists()
        modules.append({
            "name": d.name,
            "package": f"internal/{d.name}",
            "has_handler": has_handler,
        })
    return modules


def list_pages() -> list[str]:
    """Find all TypeScript page files."""
    pages = []
    for f in sorted(APP_PAGES.rglob("*.ts")):
        pages.append(str(f.relative_to(APP_PAGES)))
    return pages


def list_components() -> list[str]:
    """Find all TypeScript component files."""
    components = []
    for f in sorted(APP_COMPONENTS.rglob("*.ts")):
        components.append(str(f.relative_to(APP_COMPONENTS)))
    return components


def get_current_branch() -> str:
    """Get the current git branch name."""
    try:
        result = subprocess.run(
            ["git", "rev-parse", "--abbrev-ref", "HEAD"],
            capture_output=True, text=True, cwd=REPO_ROOT
        )
        return result.stdout.strip()
    except Exception:
        return "unknown"


def main():
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M UTC")
    branch = get_current_branch()
    migration_count = count_migrations()
    modules = list_modules()
    pages = list_pages()
    components = list_components()

    modules_with_handler = [m for m in modules if m["has_handler"]]
    modules_without = [m for m in modules if not m["has_handler"]]

    print("=" * 60)
    print("GableLBM Wiki Inventory Scanner")
    print("=" * 60)
    print(f"  Generated:      {timestamp}")
    print(f"  Branch:         {branch}")
    print(f"  Migrations:     {migration_count}")
    print(f"  Modules total:  {len(modules)}")
    print(f"    With handler: {len(modules_with_handler)}")
    print(f"    No handler:   {len(modules_without)}")
    print(f"  Pages:          {len(pages)}")
    print(f"  Components:     {len(components)}")
    print()
    print("Modules with HTTP handlers:")
    for m in modules_with_handler:
        print(f"  ✅  {m['package']}")
    print()
    print("Packages (no handler — shared/internal):")
    for m in modules_without:
        print(f"  —   {m['package']}")
    print()
    print("✅ Inventory scan complete.")
    print(f"   Run build_html.py and build_docx.py to compile outputs.")

    # Write a machine-readable snapshot for the build scripts to consume
    snapshot_path = DOCS_WIKI / ".inventory_snapshot.txt"
    DOCS_WIKI.mkdir(parents=True, exist_ok=True)
    with open(snapshot_path, "w") as f:
        f.write(f"generated_at={timestamp}\n")
        f.write(f"branch={branch}\n")
        f.write(f"migration_count={migration_count}\n")
        f.write(f"module_count={len(modules)}\n")
        f.write(f"handler_count={len(modules_with_handler)}\n")
        f.write(f"page_count={len(pages)}\n")
        f.write(f"component_count={len(components)}\n")

    print(f"\nSnapshot written: {snapshot_path}")


if __name__ == "__main__":
    main()
