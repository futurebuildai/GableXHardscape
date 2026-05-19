-- Migration 071: UOM Expansion for Hardscape
-- Adds hardscape-specific units of measure to the uom_type enum.

ALTER TYPE uom_type ADD VALUE IF NOT EXISTS 'PLT';   -- Pallet
ALTER TYPE uom_type ADD VALUE IF NOT EXISTS 'TON';   -- Ton
ALTER TYPE uom_type ADD VALUE IF NOT EXISTS 'LYR';   -- Layer
ALTER TYPE uom_type ADD VALUE IF NOT EXISTS 'PC';    -- Piece (explicit single)
ALTER TYPE uom_type ADD VALUE IF NOT EXISTS 'CYD';   -- Cubic Yard (aggregates)
