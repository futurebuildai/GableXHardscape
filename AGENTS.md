# AGENTS.md — Global Dibbits Handoff Rules & Guardrails
*Target: Commited to both repository roots for cross-tool synergy (Claude Code, Cursor, and Antigravity).*

## 1. Project Context & Engagement
- **Client Project**: Dibbits
- **Purpose**: Epicor BisTrack → GableXHardscape migration (Ontario Construction Act compliance, custom Hardscape module).
- **Core Architecture**:
  1. `dibbits-portal` (`futurebuildai/dibbits`): React + Vite + Tailwind CSS partner portal. Contains the docs/wiki/ + project manual.
  2. `gablex-erp` (`futurebuildai/GableXHardscape`): Go 1.25 + Chi router + pgx + Lit 3 client-side. PostgreSQL 16 database.

## 2. Security Boundaries & Guardrails
- **Miro Board Data Scoping**: The Miro board (`mcp.miro.com`) is **client-shared**. Under no circumstances may an agent write internal technical preparation notes, pricing discussions, or raw hours negotiations directly to the Miro board. ONLY client-safe triage rows, scope boards, and clear milestone update summaries are allowed.
- **Credential Safety**: Never commit credentials, raw tokens, or Railway keys. Use the `.env` boundaries scoped in each repository configuration.

## 3. Deployment & Branching Models
- **Single-Trunk Production Deploy**: We run on a strict single-trunk branching architecture. All production releases are triggered automatically when commits land on the `production` branch.
- **Handoff Loop**: Any code updates must clear the zero-trust audit checks in `.agents/loop/` before they are approved for deployment.
