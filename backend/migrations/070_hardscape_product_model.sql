-- Migration 070: Hardscape Product Model Transformation
-- Renames LBM-specific columns (species, grade) to hardscape equivalents (manufacturer, collection)
-- and adds hardscape-specific product attributes.

-- 1. Rename LBM columns to hardscape equivalents
ALTER TABLE products RENAME COLUMN species TO manufacturer;
ALTER TABLE products ALTER COLUMN manufacturer TYPE VARCHAR(255);
ALTER TABLE products ALTER COLUMN manufacturer SET DEFAULT '';

ALTER TABLE products RENAME COLUMN grade TO collection;
ALTER TABLE products ALTER COLUMN collection TYPE VARCHAR(255);
ALTER TABLE products ALTER COLUMN collection SET DEFAULT '';

-- 2. Add hardscape-specific columns
ALTER TABLE products ADD COLUMN IF NOT EXISTS coverage_sf_per_unit DECIMAL(19,4);
ALTER TABLE products ADD COLUMN IF NOT EXISTS color VARCHAR(255) DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS finish VARCHAR(255) DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS application VARCHAR(255) DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS dimensions_lwh VARCHAR(100) DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS pallet_count INTEGER;
ALTER TABLE products ADD COLUMN IF NOT EXISTS weight_per_unit DECIMAL(19,4);
ALTER TABLE products ADD COLUMN IF NOT EXISTS pieces_per_sf DECIMAL(19,4);

-- 3. Drop old LBM indexes and create new hardscape indexes
DROP INDEX IF EXISTS idx_products_species;
DROP INDEX IF EXISTS idx_products_grade;

CREATE INDEX IF NOT EXISTS idx_products_manufacturer ON products(manufacturer);
CREATE INDEX IF NOT EXISTS idx_products_collection ON products(collection);
CREATE INDEX IF NOT EXISTS idx_products_color ON products(color);
CREATE INDEX IF NOT EXISTS idx_products_application ON products(application);
CREATE INDEX IF NOT EXISTS idx_products_finish ON products(finish);
CREATE INDEX IF NOT EXISTS idx_products_manufacturer_collection ON products(manufacturer, collection);
