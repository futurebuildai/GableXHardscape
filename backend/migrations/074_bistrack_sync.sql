-- Migration 074: Bistrack Sync Tables
-- Tracks sync jobs and discrepancies for bi-directional Bistrack WebTrack integration.

CREATE TABLE IF NOT EXISTS sync_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type VARCHAR(50) NOT NULL,             -- 'initial', 'incremental'
    direction VARCHAR(50) NOT NULL,            -- 'bistrack_to_native', 'native_to_bistrack', 'bidirectional'
    status VARCHAR(50) NOT NULL DEFAULT 'running', -- running, completed, failed, cancelled
    entity_type VARCHAR(50),                   -- 'customer', 'product', 'inventory', 'order', 'all'
    records_total INTEGER DEFAULT 0,
    records_synced INTEGER DEFAULT 0,
    records_failed INTEGER DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sync_jobs_status ON sync_jobs(status);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_started ON sync_jobs(started_at DESC);

CREATE TABLE IF NOT EXISTS sync_discrepancies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sync_job_id UUID REFERENCES sync_jobs(id),
    entity_type VARCHAR(50) NOT NULL,          -- 'customer', 'product', 'inventory', 'order'
    entity_id UUID NOT NULL,
    bistrack_id VARCHAR(255),
    field_name VARCHAR(255) NOT NULL,
    bistrack_value TEXT,
    native_value TEXT,
    resolution VARCHAR(50),                    -- 'use_bistrack', 'use_native', 'manual', NULL (unresolved)
    resolved_by UUID,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sync_discrepancies_unresolved ON sync_discrepancies(resolution) WHERE resolution IS NULL;
CREATE INDEX IF NOT EXISTS idx_sync_discrepancies_job ON sync_discrepancies(sync_job_id);
CREATE INDEX IF NOT EXISTS idx_sync_discrepancies_entity ON sync_discrepancies(entity_type, entity_id);
