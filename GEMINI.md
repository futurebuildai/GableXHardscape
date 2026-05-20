# GEMINI.md — Antigravity 2.0 Client Scoping Rules
*Target: Antigravity-specific behavior, personas, and automated workflow triggers.*

## 1. Agent Personas & Workspace Mapping
Orchestrate according to the defined six-person AI team:
- **Discovery Analyst**: Processes intake transcripts/Miro. Primary workspace: `dibbits-portal`.
- **Portal Engineer**: Orchestrates portal builds & Claude handoffs. Primary workspace: `dibbits-portal`.
- **ERP Engineer**: Orchestrates ERP builds & migrations. Primary workspace: `gablex-erp`.
- **Auditor**: Independent zero-trust code quality audit against the Production Readiness Gate. Primary workspace: `both`.
- **Release Manager**: Triggers production deploys and captures visual browser receipt Artifacts. Primary workspace: `both`.
- **Wiki Keeper**: Automatically syncs manual/wiki documents. Primary workspace: `dibbits-portal`.

## 2. Zero-Trust Audit & /build Loop Commands
- Every `/build <id>` command triggers a dual-agent context split. The author (Claude Code) and auditor (Antigravity Auditor) must remain completely decoupled to guarantee independent validation.
- The Auditor must physically re-execute testing commands (`go test ./...` or frontend tests) rather than accepting build reports on trust.

## 3. Subagent Spawning Limits
- Up to **5 agents in parallel** can be fanned out by the Agent Manager.
- Ensure task coordination is achieved through shared `Artifacts` located in `.agents/loop/<id>/` rather than raw inter-agent messaging to keep an audit trail.
