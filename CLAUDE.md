# CLAUDE.md — GableXHardscape

> **Fork of `futurebuildai/GableLBM-main`** (re-forked 2026-05-19). Hardscape
> vertical layered with migrations 070-074. Everything below 070 is upstream's
> conventions, unchanged.
>
> **Initial customer:** Dibbits (Trenton + Kingston, ON).
> Engagement repo: https://github.com/futurebuildai/dibbits
> Upstream:        https://github.com/futurebuildai/GableLBM-main

The upstream's full operator guide is preserved verbatim at
[`CLAUDE.upstream.md`](./CLAUDE.upstream.md) — start there for stack details,
conventions, pre-flight checks, and gotchas.

## Fork-specific gotchas

### `products.species` and `products.grade` are renamed

After migration `070_hardscape_product_model.sql`, the columns are
`manufacturer` and `collection` respectively. Upstream Go code that still
references `species` / `grade` will not compile — those references need
renaming when porting forward.

### UoM enum has extra values

`PLT`, `TON`, `LYR`, `PC`, `CYD` are added to `uom_type` by migration 071.
Code that switches on `uom_type` must handle them or default-case cleanly.

### Configurator rule data is hardscape-shaped

Migration 072 deletes upstream's lumber configurator seed and inserts
Techo-Bloc / Permacon / Belgard / Unilock rules. The configurator engine
itself is unchanged; only the seed differs.

### Lien notice table (`lien_notices`) is Dibbits-shaped today

Migration 073 ships an Ontario Construction Act-flavoured schema
(preservation_deadline, holdback_amount, supply_date). Other hardscape
customers in other jurisdictions may need a more generic shape later.

### BiTrack sync tables exist; the Go module is pending

Migration 074 creates `sync_jobs` + discrepancy tables. The
`internal/bistrack/` Go module that drives the bi-directional sync is
**not in this commit** — it lands in a follow-up once the Dibbits credentials
+ BiTrack webhook URLs are settled. The schema can be applied now without
runtime impact.

### Module path is `github.com/futurebuildai/gablexhardscape`

All Go imports use this path. The upstream's `github.com/gablelbm/gable` is
mass-renamed across 127 files. When cherry-picking new upstream changes,
re-run the rename on any new files.
