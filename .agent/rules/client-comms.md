# client-comms.md — Client Communication & Data Boundaries

## 1. The Internal/Client Boundary
- **Client-Visible Surface**: The Miro board (`mcp.miro.com`) is a shared client surface. Updates posted to it must represent finalized features, high-level status changes, and clear milestone progress.
- **Strictly Internal**: Under no circumstances should internal call preparation notes, pricing rationale, developer effort debates, or raw logs be written or synced to the Miro board. Keep all such documents inside the `.agents/` project configurations.

## 2. Status Traceability
- Ensure that every piece of client feedback synced via `/discovery-intake` has a clear tracking reference link inside the `Implementation/Mapping.md` page to guarantee billing traceability.
