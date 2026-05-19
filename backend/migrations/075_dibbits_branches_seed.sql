-- 075: Dibbits-specific BRANCH seed for the Trenton + Kingston yards.
-- Replaces the deferred-at-re-fork M7 seed (the old fork's 054 inserted
-- locations with type='ZONE'; upstream's multi-branch foundation now
-- requires type='BRANCH' with non-NULL address columns).
--
-- Why this exists in the *ERP* repo (not just the seed cmd):
-- Without these BRANCH rows the upstream /api/v1/me/branches endpoint
-- 500s on a NULL address column scan, the Lit shell can't initialize,
-- and the frontend gets stuck on its initial loader. This migration
-- is the minimum data needed for the ERP to render against an empty DB.

-- 1. Two BRANCH-type locations with FULL non-NULL metadata so the Go
-- scan into *string columns succeeds (pgx v5 doesn't accept NULL into
-- a non-Null *string).
INSERT INTO locations (
    id, parent_id, path, type, code, description,
    name, address, city, state, zip, phone,
    tax_jurisdiction_code, default_tax_rate, timezone, active
) VALUES
    (
        'a0000001-0000-0000-0000-000000000001'::uuid,
        NULL, '', 'BRANCH', 'TRN', 'Dibbits Trenton',
        'Dibbits Trenton', '275 Glen Miller Rd', 'Trenton', 'ON', 'K8V 5P8', '613-555-0001',
        'CA-ON', 0.1300, 'America/Toronto', TRUE
    ),
    (
        'a0000001-0000-0000-0000-000000000002'::uuid,
        NULL, '', 'BRANCH', 'KGS', 'Dibbits Kingston',
        'Dibbits Kingston', '1200 Midland Ave', 'Kingston', 'ON', 'K7P 2X8', '613-555-0002',
        'CA-ON', 0.1300, 'America/Toronto', TRUE
    )
ON CONFLICT (id) DO NOTHING;

-- 2. branch_id self-reference (upstream's denorm trigger from migration
-- 058 should also fire on BRANCH rows; the trigger sets branch_id = id
-- for BRANCH type, but a fallback explicit set guards against trigger
-- ordering edge-cases on bootstrap).
UPDATE locations
SET branch_id = id
WHERE type = 'BRANCH' AND branch_id IS NULL;

-- 3. Set the default_branch_id system_setting that upstream's portal /
-- branch context helpers read on every request.
INSERT INTO system_settings (key, value)
VALUES ('default_branch_id', 'a0000001-0000-0000-0000-000000000001')
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;

-- 4. Note: customers/orders/quotes/invoices backfill (the 062-067 chain)
-- already ran on a previously-empty DB, so existing data — none —
-- doesn't need re-touching. New writes flow through the default_branch_id
-- path correctly.
