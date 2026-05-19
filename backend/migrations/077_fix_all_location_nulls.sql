-- 077: Fix NULL address columns on ALL locations (not just BRANCH-type).
--
-- Migration 076 only touched type='BRANCH' rows, but the upstream
-- branches handler returns multiple rows of varying type and 500s on
-- the first NULL address it scans. Backfill every location with NULL
-- address-shaped columns to harmless defaults so the handler can scan
-- cleanly.

UPDATE locations
   SET address               = COALESCE(address,               'n/a'),
       city                  = COALESCE(city,                  'n/a'),
       state                 = COALESCE(state,                 'ON'),
       zip                   = COALESCE(zip,                   'n/a'),
       phone                 = COALESCE(phone,                 'n/a'),
       tax_jurisdiction_code = COALESCE(tax_jurisdiction_code, 'CA-ON'),
       default_tax_rate      = COALESCE(default_tax_rate,      0.1300)
 WHERE address IS NULL OR city IS NULL OR state IS NULL
    OR zip IS NULL OR phone IS NULL
    OR tax_jurisdiction_code IS NULL OR default_tax_rate IS NULL;

-- Also fix any remaining NULL name on locations (handler may scan it
-- as *string too). The `name` column was added in 057 nullable.
UPDATE locations
   SET name = COALESCE(name, description, code, 'unnamed-location')
 WHERE name IS NULL;
