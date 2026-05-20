# GableLBM Architecture — Auto-Refreshed Inventory
> **Generated:** 2026-05-20 · Branch: `feature/a20-internal-permissions`
> **Scanner:** wiki-refresh workflow · Do not edit by hand — re-run `/wiki-refresh` to update.

---

## 1. System Principles

| Principle | Detail |
|---|---|
| **Modular Monolith** | Single deployment binary; internally strictly decoupled packages |
| **Zero-Trust Modules** | Modules never access another module's database tables directly |
| **Interface-Driven Interop** | All cross-module calls go through Go interface contracts |
| **Event-Driven Side Effects** | Write side-effects flow via NATS JetStream subjects (aspirational; see §4) |
| **Light DOM Frontend** | All Lit 3 components use `createRenderRoot() { return this; }` for Tailwind class inheritance |
| **Federated Governance** | Co-op partners can propose changes via the RFC / Governance portal |

---

## 2. Technology Stack

### Backend
| Layer | Technology | Version | Notes |
|---|---|---|---|
| Language | Go | 1.25+ | `go.mod` constraint |
| HTTP Router | Chi | v5 | Middleware-first; all routes at `/api/v1/*` |
| Database | PostgreSQL | 16+ | Extensions: `pgvector` (AI), `postgis` (GIS/delivery) |
| DB Driver | pgx | v5 | `database/sql`-compatible |
| Messaging | NATS JetStream | embedded | Container in `docker-compose.yml`; **not yet wired in Go code** — all inter-module writes are synchronous Go interface calls; NATS sections in this doc are forward-looking design |
| Auth | JWT + JWKS | — | Keycloak-compatible; `godotenv` config fallback |

### Frontend
| Layer | Technology | Version | Notes |
|---|---|---|---|
| Component Framework | Lit 3 Web Components | 3.x | Light DOM mandatory |
| Language | TypeScript | 5.9 | Strict mode |
| Bundler | Vite | 7 | Dev server + production build |
| Styling | Tailwind CSS | 3.4 | Design tokens; no Shadcn |
| Charts | Chart.js | 4 | Revenue trend, order status |
| Maps | Leaflet | — | Route / delivery maps |

### Infrastructure (Non-Production)
| Environment | Host | Branch | URL | Logical DB |
|---|---|---|---|---|
| Demo | Digital Ocean App Platform | `community` | https://demo.gablelbm.com | `gable_demo` |
| Staging | Digital Ocean App Platform | `staging` | https://staging.gablelbm.com | `gable_staging` |

Production (`master`) is not managed by this repository. App specs: `.do/app-demo.yaml`, `.do/app-staging.yaml`.

---

## 3. Backend Module Inventory

**Total modules with HTTP handlers: 35** · **Total migrations: 81**

| Module | Package | Handler | Responsibility |
|---|---|---|---|
| Account | `internal/account` | ✅ | Customer accounts, credit limits, account rules |
| AI | `internal/ai` | — | AI orchestration, embedding pipelines |
| Accounts Payable | `internal/ap` | ✅ | Vendor invoices, AP workflows, remittance |
| Bank Reconciliation | `internal/bankrecon` | ✅ | Bank statement import, reconciliation matching |
| Config | `internal/config` | — | Environment config, feature flags |
| Configurator | `internal/configurator` | ✅ | Millwork door/product configurator rules |
| CRM | `internal/crm` | ✅ | Activity logging, contact management |
| Customer | `internal/customer` | ✅ | Customer CRUD, branch access, statements |
| Dashboard | `internal/dashboard` | ✅ | KPIs, revenue trends, alert widgets |
| Delivery | `internal/delivery` | ✅ | Dispatch board, routes, driver app |
| Document | `internal/document` | ✅ | Document storage, lien notices, PDF generation |
| Domain | `internal/domain` | — | Shared domain types (no HTTP exposure) |
| EDI | `internal/edi` | — | Electronic data interchange adapters |
| General Ledger | `internal/gl` | ✅ | Double-entry GL, journal entries, fiscal periods, trial balance, financial statements |
| Governance | `internal/governance` | ✅ | RFC portal, AI governance engine, backlog orchestrator |
| Integrations | `internal/integrations` | ✅ | BisTrack, Spruce, DMSi sync adapters |
| Inventory | `internal/inventory` | ✅ | Stock quants, moves, cycle counts, UOM conversions |
| Invoice | `internal/invoice` | ✅ | AR invoices, payment terms, credit notes |
| Location | `internal/location` | ✅ | Branches, bins, staff roles + approval requests |
| Matching | `internal/matching` | ✅ | PO–invoice 3-way matching |
| Millwork | `internal/millwork` | ✅ | Millwork product configuration, cut lists |
| Notification | `internal/notification` | — | In-app and email notification dispatch |
| Order | `internal/order` | ✅ | Sales orders, line items, order lifecycle |
| Parsing | `internal/parsing` | ✅ | Material list OCR/AI parsing |
| Partner | `internal/partner` | ✅ | Co-op partner portal, federated catalog |
| Payment | `internal/payment` | ✅ | Payment processing, AR application, AR aging |
| PIM | `internal/pim` | ✅ | AI-powered product content, images, descriptions |
| Portal | `internal/portal` | ✅ | B2B dealer portal (ordering, team management) |
| POS | `internal/pos` | ✅ | Point-of-sale terminal, till management |
| Pricing | `internal/pricing` | ✅ | Price matrix, account rules, category tree, escalators |
| Product | `internal/product` | ✅ | Product catalog, SKUs, categories, hardscape model |
| Project | `internal/project` | ✅ | Project tracking, job costing |
| Purchase Order | `internal/purchase_order` | ✅ | PO creation, receiving, vendor management |
| Quote | `internal/quote` | ✅ | Quote builder, analytics, escalator rules |
| Reporting | `internal/reporting` | ✅ | Report builder, AR aging, saved reports, customer statements |
| Sales Team | `internal/salesteam` | ✅ | Sales team management, territory assignment |
| Tax | `internal/tax` | ✅ | HST/GST/PST tax codes, rates, CRA compliance |
| Tech Admin | `internal/techadmin` | ✅ | Platform admin, AI key management |
| Vendor | `internal/vendor` | ✅ | Vendor master, AP contacts, rebate programs |
| Vision | `internal/vision` | ✅ | Computer vision / barcode scanning |

---

## 4. Inter-Module Communication

### 4.1 Synchronous (Reads — Current)
Direct Go interface calls within the same process. No network hop.

```go
// Example: Sales checks stock availability
inventoryService.GetAvailability(sku, locationID) // → AvailabilityResult
```

### 4.2 Asynchronous (Writes — Forward-Looking)
NATS JetStream subjects. **Not yet wired in production code.** Current implementation uses synchronous calls for all inter-module writes.

| Subject | Publisher | Subscribers |
|---|---|---|
| `sales.order.confirmed` | order | inventory (reserve stock), logistics (pick ticket), invoice (credit check) |
| `ar.invoice.issued` | invoice | gl (auto-post: DR AR / CR Revenue / CR HST) |
| `ar.payment.received` | payment | gl (auto-post: DR Cash / CR AR) |
| `ap.vendor_invoice.approved` | ap | gl (auto-post: DR Inventory / DR HST / CR AP) |
| `ap.vendor_payment.issued` | ap | gl (auto-post: DR AP / CR Cash) |
| `pos.till.closed` | pos | gl (auto-post: DR Cash+Card / CR Revenue / CR HST) |
| `gl.period.reopen_requested` | gl | notification (admin bell → period approval) |

---

## 5. Database Schema Inventory

**Total migrations applied: 81** (through `079_staff_roles_and_approval_requests.sql`)

| Migration Range | Area | Key Tables Added |
|---|---|---|
| 001–020 | Core foundation | `products`, `uoms`, `locations`, `stock_quants`, `inventory_moves`, `customers`, `orders`, `invoices`, `payments`, `vendors`, `purchase_orders` |
| 021–040 | Sales & pricing | `quotes`, `price_levels`, `pricing_rules`, `price_matrix`, `account_rules`, `delivery_routes`, `dispatch` |
| 041–060 | Operations | `pos_sessions`, `tills`, `barcode_scans`, `crm_activities`, `projects`, `lien_notices` |
| 061–070 | Multi-branch | `branch_users`, `orders.branch_id`, `quotes.branch_id`, `invoices.branch_id`, `po.branch_id`, `pos.branch_id`, `customers.branch_id` |
| 071–077 | LBM expansion | `hardscape_products`, `uom_expansion` (lineal ft, sq ft, bd ft, ton), `hardscape_configurator_rules`, `bistrack_sync` |
| 075 | Dibbits seed | Trenton + Kingston branch data, Dibbits product categories |
| 025 | General Ledger | `gl_accounts` (COA seed: 16 LBM accounts), `gl_fiscal_periods`, `gl_journal_entries`, `gl_journal_lines` |
| 078 | Replenishment | `replenishment_settings` (planned: GL gap closure — see A-01 execution plan) |
| 079 | Staff Roles | `staff_roles`, `approval_requests` (A-20: internal permission model) |

> **Next migration:** `078_gl_gap_closure.sql` — adds `gl_period_overrides`, `gl_audit_log`, covering index (per A-01 execution plan)

---

## 6. API Surface

**Base URL:** `/api/v1/`
**Auth:** `Authorization: Bearer <JWT>` (JWKS-validated)
**Format:** JSON request/response; standard error envelope `{ "error": { "code": "...", "message": "..." } }`

| Module | Route Prefix | Notable Endpoints |
|---|---|---|
| Account | `/accounts` | CRUD, credit limit, rules |
| AP | `/ap` | Vendor invoices, payment runs |
| Bank Recon | `/bankrecon` | Statement upload, match, reconcile |
| Configurator | `/configurator` | Product configuration rules |
| CRM | `/crm` | Activity log, contact timeline |
| Customer | `/customers` | CRUD, statement, branch access |
| Dashboard | `/dashboard` | KPIs, revenue trend, alerts |
| Delivery | `/delivery` | Routes, dispatch board, driver app |
| General Ledger | `/gl` | Accounts, journal entries, fiscal periods, trial balance, balance sheet, P&L, audit log |
| Governance | `/governance` | RFCs, voting, backlog |
| Integrations | `/integrations` | BisTrack/Spruce/DMSi sync |
| Inventory | `/inventory` | Stock, moves, cycle counts |
| Invoice | `/invoices` | AR invoices, credit notes |
| Location | `/locations` | Branches, bins, staff roles, approvals |
| Matching | `/matching` | 3-way PO match |
| Millwork | `/millwork` | Millwork config, cut lists |
| Order | `/orders` | Sales order lifecycle |
| Parsing | `/parsing` | Material list AI parse |
| Partner | `/partner` | Federated catalog, co-op portal |
| Payment | `/payments` | AR payment application |
| PIM | `/pim` | AI product content |
| Portal | `/portal` | B2B ordering, team, cart |
| POS | `/pos` | Terminal, till session, closeout |
| Pricing | `/pricing` | Matrix, rules, category tree |
| Product | `/products` | Catalog, SKU, hardscape |
| Project | `/projects` | Job tracking, costing |
| Purchase Order | `/purchase-orders` | PO lifecycle, receiving |
| Quote | `/quotes` | Builder, analytics, escalators |
| Reporting | `/reports` | AR aging, report builder, statements |
| Sales Team | `/salesteam` | Team management |
| Tax | `/tax` | HST/GST codes, CRA compliance |
| Tech Admin | `/techadmin` | AI keys, platform config |
| Vendor | `/vendors` | Vendor master, rebates |
| Vision | `/vision` | Barcode scan, computer vision |

---

## 7. Frontend Architecture

**Total TypeScript files: 210** · **Pages: 68** · **Components: 42+**

### 7.1 App Layouts

| Layout | Component | Routes | Target User |
|---|---|---|---|
| ERP Desktop | `app-shell.ts` | `/erp/*` | Staff (Sandra, Ryan, yard team) |
| B2B Portal | `portal-layout.ts` | `/portal/*` | Contractor / dealer customers |
| Driver App | `driver-layout.ts` | `/driver/*` | Delivery drivers (mobile) |
| Yard App | `yard-layout.ts` | `/yard/*` | Yard/warehouse staff |
| POS Terminal | (inline) | `/pos` | Counter sales staff |

### 7.2 Page Inventory

| Section | Pages |
|---|---|
| **Accounting** | `BankReconciliation`, `ChartOfAccounts`, `JournalEntries`, `POMatching`, `TrialBalance` |
| **Accounts (CRM)** | `AccountDetailPage`, `AccountsPage`, `ActivityFeed`, `ContactList` |
| **Admin** | `Branches`, `BranchUsers`, `PendingApprovals`, `StaffManagement`, `TechAdminPage` |
| **Admin / Pricing** | `AccountRulesTable`, `CategoryTree`, `MatrixGrid`, `PricingMatrix`, `ResolutionPreview`, `RuleDrawer` |
| **Dashboard** | `Dashboard`, `DailyTill` |
| **Dispatch** | `DispatchBoard` |
| **Driver** | `DeliveryDetail`, `RouteList`, `StopList` |
| **Governance** | `NewRFC`, `RFCDashboard`, `RFCDetail` |
| **Inventory** | `ProductDetail`, `ProductList` *(+ yard: `InventoryLookup`, `CycleCount`)* |
| **Portal (B2B)** | `PortalHome`, `PortalMyAccount`, `PortalOrders`, `PortalProductDetail`, `PortalTeam` |
| **POS** | `POSTerminal` |
| **Projects** | `ProjectDashboard`, `ProjectList` |
| **Purchasing** | `NewPurchaseOrder`, `ProcurementDashboard`, `PurchaseOrderDetail`, `PurchaseOrderList`, `PurchasingRecommendations`, `RebatePrograms`, `RebateReport`, `ReplenishmentSettings`, `VendorDetail`, `VendorList` |
| **Quotes** | `QuoteBuilder`, `QuoteAnalytics`, `QuoteDetail`, `QuoteList` |
| **Reports** | `ARAgingReport`, `CustomerStatementPage`, `ReportBuilder`, `SavedReports` |
| **Yard** | `CycleCount`, `InventoryLookup`, `PickDetail`, `PickQueue`, `ReceivePO` |

### 7.3 Component Inventory

| Category | Components |
|---|---|
| **UI Atoms** | `Button`, `Card`, `Tooltip`, `LoadingScreen`, `brand-logo`, `toast-container`, `shortcuts-modal`, `omnibar` |
| **Layout** | `app-shell`, `branch-switcher`, `portal-layout`, `driver-layout`, `yard-layout` |
| **Dashboard Widgets** | `KPICard`, `RevenueTrendChart`, `OrderStatusChart`, `RecentOrdersFeed`, `TopCustomersTable`, `InventoryAlertsWidget` |
| **Inventory** | `InventoryTable`, `AddProductModal`, `StockAdjustmentModal`, `InventoryTransferModal`, `ProductMarginModal` |
| **Logistics** | `DeliveryList`, `RouteList`, `RouteMap`, `AssignOrderModal`, `CreateRouteModal` |
| **Quotes** | `LineItemEditor`, `EscalatorToggle`, `MaterialListUpload`, `ParsedResultsPanel` |
| **Portal** | `CartSidebar`, `ProductCard` |
| **POS** | `OfflineBanner` |
| **Payments** | `RunPaymentsForm`, `PaymentModal` |
| **CRM** | `LogActivityModal` |
| **Customers** | `CustomerSelect` |
| **Products** | `ProductSelect` |
| **Location** | `LocationManager` |
| **Common** | `PermissionFallbackModal`, `BarcodeScanner` |

---

## 8. General Ledger Module (A-01)

The GL module (`internal/gl`) is the financial core. It implements a native double-entry bookkeeping engine inside the ERP, eliminating the Dynamics GP integration.

| Capability | Status |
|---|---|
| Chart of Accounts (CRUD) | ✅ Schema seeded (16 LBM accounts); API gaps in progress |
| Journal Entries (DRAFT→POSTED→VOID) | ✅ Core engine; ~70% complete |
| Fiscal Period open/close | ✅ Exists |
| Period reopen + admin approval | 🔧 Gap — A-01 execution plan |
| Trial Balance compute + CSV export | ✅ Exists |
| Balance Sheet | 🔧 Gap — A-01 execution plan |
| P&L / Income Statement | 🔧 Gap — A-01 execution plan |
| ERP auto-posting (NATS consumer) | 🔧 Gap — A-01 execution plan |
| Dashboard financial widget | 🔧 Gap — A-01 execution plan |
| GL Audit Log | 🔧 Gap — A-01 execution plan |

Full spec: `.agents/handoff/features/A-01-native-general-ledger/` (21-artifact pipeline complete)

---

## 9. Active Feature Work (Current Branch)

**Branch:** `feature/a20-internal-permissions`

| Feature | Status | Key Files |
|---|---|---|
| A-20: Internal Permission Model | 🔧 In Progress | `internal/location/role_model.go`, `role_repository.go`, `role_service.go`, `role_test.go`; Migration `079_staff_roles_and_approval_requests.sql`; `app/src/pages/admin/PendingApprovals.ts`, `StaffManagement.ts`, `components/common/PermissionFallbackModal.ts` |

---

## 10. Design System Reference

**Theme:** Industrial Dark — High Contrast, Data Density, Zero Clutter

| Token | Value | Usage |
|---|---|---|
| Gable Green | `#00FFA3` | Primary actions, success states, active glow |
| Deep Space | `#0A0B10` | Global background |
| Slate Steel | `#161821` | Card backgrounds, sidebar, modals |
| Safety Red | `#F43F5E` | Errors, stockouts, credit hold |
| Blueprint Blue | `#38BDF8` | Technical data, dimensions, links |
| Glass overlay | `rgba(255,255,255,0.05)` + `backdrop-filter: blur(12px)` | Modal overlays, cards |
| UI Font | Inter (400/500/600) | Labels, headers, body |
| Data Font | JetBrains Mono | SKUs, prices, quantities, dimensions |

Full design system: `docs/design-system.md`

---

## 11. Legacy Interop & Migration Strategy

| Legacy System | Adapter Location | Protocol |
|---|---|---|
| Epicor BisTrack | `internal/integrations` | REST JSON |
| ECI Spruce | `internal/integrations` | SOAP XML |
| DMSi Agility | `internal/integrations` | REST JSON |
| Dynamics GP (GL) | Replacing with A-01 native GL | N/A — full replacement |

Sync engine: `pkg/sync` (bidirectional data flow during phase-in).
