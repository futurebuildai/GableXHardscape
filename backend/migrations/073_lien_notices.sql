-- Migration 073: Ontario Construction Act — Lien Notices
-- Tracks preservation deadlines and holdback amounts per the Ontario Construction Act.

CREATE TABLE IF NOT EXISTS lien_notices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES customers(id),
    project_name VARCHAR(500),
    supply_date DATE NOT NULL,
    preservation_deadline DATE NOT NULL,    -- supply_date + 60 days
    holdback_amount DECIMAL(19,4) NOT NULL, -- invoice_total × 0.10
    invoice_id UUID REFERENCES invoices(id),
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, preserved, expired, released
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lien_notices_account ON lien_notices(account_id);
CREATE INDEX IF NOT EXISTS idx_lien_notices_deadline ON lien_notices(preservation_deadline);
CREATE INDEX IF NOT EXISTS idx_lien_notices_status ON lien_notices(status);
