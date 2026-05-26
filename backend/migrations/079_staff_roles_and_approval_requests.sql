-- 079: Database migration for internal staff roles and override approval requests.
-- Strict 3-role staff taxonomy: Super Admin, Scoped Manager, Scoped Core/Hourly.
-- Interactive UI fallbacks approval queue.

-- 1. Create staff_roles table
CREATE TABLE IF NOT EXISTS staff_roles (
    user_sub TEXT PRIMARY KEY,
    role VARCHAR(50) NOT NULL CHECK (role IN (
        'General Manager',
        'Branch Manager',
        'Procurement Manager',
        'Sales Manager',
        'Inside Sales',
        'Outside Sales',
        'Yard Manager',
        'Yard Team',
        'Logistics Manager',
        'Drivers',
        'HR',
        'Financial Controller',
        'Payables/Receivables'
    )),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_staff_roles_user ON staff_roles(user_sub);

-- 2. Create permission_approval_requests table for override approval hooks
CREATE TABLE IF NOT EXISTS permission_approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_sub TEXT NOT NULL,
    branch_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    policy_type VARCHAR(50) NOT NULL CHECK (policy_type IN ('MIN_MARGIN', 'CREDIT_LIMIT', 'COD_CONSTRAINT', 'BRANCH_SCOPING', 'PRICE_RULES')),
    details JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    decided_at TIMESTAMPTZ,
    decided_by TEXT,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_perm_approvals_status ON permission_approval_requests(status);
CREATE INDEX IF NOT EXISTS idx_perm_approvals_branch ON permission_approval_requests(branch_id);
CREATE INDEX IF NOT EXISTS idx_perm_approvals_user ON permission_approval_requests(user_sub);
