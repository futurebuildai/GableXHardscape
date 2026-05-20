# deploy.md — Deployment and Release Integrity Rules

## 1. Release Orchestration
- All deploys are executed by the **Release Manager** agent.
- High-risk deploys require explicit human confirmation.
- Deploys must use the `dibbits-deploy` skill rules to execute backoffs/retries during Railway or Cloudflare routing incidents.

## 2. Post-Deployment Verification
- Every successful deployment must capture a visual browser recording Artifact of the live site at `dibbits.gablelbm.com`.
- Verify TLS handshake status and page structure to ensure no custom-domain failures are active.
