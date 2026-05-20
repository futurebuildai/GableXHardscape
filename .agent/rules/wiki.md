# wiki.md — Wiki Maintenance & Upstream Sync Rules

## 1. Documentation Alignment
- **Role**: The **Wiki Keeper** agent owns the execution of documentation generation.
- **Trigger**: Every documentation rebuild is triggered via the `/wiki-refresh` command.
- **Rebuild Operations**:
  - Run `build_html.py` and `build_docx.py` to compile formatted outputs.
  - Automatically update `refresh_inventory.py` to reflect newly identified hardscape models.

## 2. Upstream Tracking
- On a weekly schedule, the **ERP Engineer** must run upstream checks against the `GableLBM-main` repository branch to report new upstream commits and avoid rebase sync surprises.
