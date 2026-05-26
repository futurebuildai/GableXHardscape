-- Feedback table for testing feedback from Dibbits team.
-- Submissions come from both the ERP UI and the Partner Portal;
-- the `source` column distinguishes origin.
CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source VARCHAR(20) NOT NULL,           -- 'ERP' or 'PORTAL'
    category VARCHAR(50) NOT NULL,         -- Bug, UI/UX, Feature Request, Data Issue, Question, Other
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    page_url TEXT,                          -- URL where feedback was submitted
    submitted_by_name VARCHAR(255),
    submitted_by_email VARCHAR(255),
    user_id UUID,                          -- nullable; ERP users have a mapped UUID
    status VARCHAR(30) NOT NULL DEFAULT 'NEW',
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    admin_notes TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedback_status ON feedback(status);
CREATE INDEX idx_feedback_source ON feedback(source);
CREATE INDEX idx_feedback_created ON feedback(created_at DESC);
