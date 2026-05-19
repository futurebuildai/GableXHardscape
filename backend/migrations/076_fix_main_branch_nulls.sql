-- 076: Backfill NOT-NULL-friendly address fields on the upstream-created
-- 'Main Branch' row.
--
-- Upstream migration 059 creates a default BRANCH row but doesn't set the
-- address/city/state/zip/phone columns (they're nullable, default NULL).
-- Upstream's Go branches handler scans into *string for those columns;
-- pgx v5 returns an error on NULL → *string scans, which 500s
-- /api/v1/branches and /api/v1/me/branches. That cascades to a blank
-- frontend (Lit shell waits on /me/branches to set the user's context).
--
-- Per migration 075 we now have real Dibbits BRANCH rows for Trenton +
-- Kingston. The upstream 'Main Branch' placeholder is no longer needed,
-- but a DELETE would orphan any pre-existing rows that backfilled to it
-- in migration 060. Safer to backfill its address columns.

UPDATE locations
   SET address               = COALESCE(address,               'n/a'),
       city                  = COALESCE(city,                  'n/a'),
       state                 = COALESCE(state,                 'ON'),
       zip                   = COALESCE(zip,                   'n/a'),
       phone                 = COALESCE(phone,                 'n/a'),
       tax_jurisdiction_code = COALESCE(tax_jurisdiction_code, 'CA-ON'),
       default_tax_rate      = COALESCE(default_tax_rate,      0.1300)
 WHERE type = 'BRANCH'
   AND (address IS NULL OR city IS NULL OR state IS NULL
        OR zip IS NULL OR phone IS NULL);

-- Also point the default_branch_id at the *Dibbits Trenton* branch
-- now that real Dibbits branches exist. Anything backfilled to the
-- placeholder via migration 060 stays linked, but new writes default
-- to Trenton.
UPDATE system_settings
   SET value = 'a0000001-0000-0000-0000-000000000001'
 WHERE key = 'default_branch_id';
