-- 080: Dibbits Demo Customers and Vendors Seed
-- Requires 075_dibbits_branches_seed.sql to have run so Trenton and Kingston branches exist.

-- 1. Vendors
INSERT INTO vendors (id, name, contact_email, phone, address_line1, city, state, zip, payment_terms, average_lead_time_days, fill_rate, total_spend_ytd)
VALUES
    ('d0000001-0000-0000-0000-000000000001'::uuid, 'Techo-Bloc', 'orders@techo-bloc.com', '1-800-463-0450', '5255 Albert-Millichamp', 'St-Hubert', 'QC', 'J3Y 8Z8', 'Net 30', 5, 95.0, 150000),
    ('d0000001-0000-0000-0000-000000000002'::uuid, 'Unilock', 'sales@unilock.com', '1-800-UNILOCK', '287 Armstrong Ave', 'Georgetown', 'ON', 'L7G 4X6', 'Net 30', 4, 96.5, 120000),
    ('d0000001-0000-0000-0000-000000000003'::uuid, 'Permacon', 'orderdesk@permacon.ca', '1-888-PERMACON', '8145 Bombardier', 'Anjou', 'QC', 'H1J 1A5', 'Net 45', 7, 92.0, 95000),
    ('d0000001-0000-0000-0000-000000000004'::uuid, 'Belgard', 'info@belgard.com', '1-877-BELGARD', '900 Ashwood Pkwy', 'Atlanta', 'GA', '30338', 'Net 30', 14, 88.0, 45000),
    ('d0000001-0000-0000-0000-000000000005'::uuid, 'Quinte Logistics', 'dispatch@quintelogistics.ca', '613-555-9988', '15 Transport Way', 'Trenton', 'ON', 'K8V 5P8', 'Net 15', 2, 99.0, 32000),
    ('d0000001-0000-0000-0000-000000000006'::uuid, 'Lakefront Aggregate', 'sales@lakefrontagg.ca', '613-555-8877', '100 Quarry Rd', 'Kingston', 'ON', 'K7K 7J5', 'Net 30', 3, 98.0, 28000)
ON CONFLICT (name) DO UPDATE SET 
    contact_email = EXCLUDED.contact_email,
    phone = EXCLUDED.phone,
    address_line1 = EXCLUDED.address_line1,
    city = EXCLUDED.city,
    state = EXCLUDED.state,
    zip = EXCLUDED.zip;

-- 2. Customers
-- We'll use DO block to fetch branch IDs dynamically and insert customers.
DO $$
DECLARE
    trenton_branch uuid := 'a0000001-0000-0000-0000-000000000001'::uuid;
    kingston_branch uuid := 'a0000001-0000-0000-0000-000000000002'::uuid;
    contractor_pl uuid;
    retail_pl uuid;
BEGIN
    -- Ensure Price Levels exist
    INSERT INTO price_levels (name, multiplier) VALUES ('Retail', 1.0000) ON CONFLICT DO NOTHING;
    INSERT INTO price_levels (name, multiplier) VALUES ('Contractor', 0.8500) ON CONFLICT DO NOTHING;
    INSERT INTO price_levels (name, multiplier) VALUES ('VIP Builder', 0.7500) ON CONFLICT DO NOTHING;
    INSERT INTO price_levels (name, multiplier) VALUES ('Municipal', 0.8000) ON CONFLICT DO NOTHING;

    SELECT id INTO contractor_pl FROM price_levels WHERE name = 'Contractor' LIMIT 1;
    SELECT id INTO retail_pl FROM price_levels WHERE name = 'Retail' LIMIT 1;

    INSERT INTO customers (
        id, name, account_number, email, phone, address, 
        credit_limit, balance_due, tier, payment_terms, price_level_id, primary_branch_id
    ) VALUES
    -- Commercial Landscapers
    (
        'c0000001-0000-0000-0000-000000000001'::uuid, 'Quinte Hardscapes Ltd', 'QH-001', 'ap@quintehardscapes.ca', '613-555-1234', '123 Main St, Trenton ON K8V 1A1', 
        50000, 15000, 'GOLD', 'NET30', contractor_pl, trenton_branch
    ),
    (
        'c0000001-0000-0000-0000-000000000002'::uuid, 'Kingston Landscape Co', 'KLC-001', 'billing@kingstonlandscape.ca', '613-555-5678', '456 Front Rd, Kingston ON K7M 4L7', 
        75000, 0, 'PLATINUM', 'NET45', contractor_pl, kingston_branch
    ),
    -- Residential Contractors
    (
        'c0000001-0000-0000-0000-000000000003'::uuid, 'Belleville Backyards', 'BB-001', 'info@bellevillebackyards.ca', '613-555-8765', '789 Bell Blvd, Belleville ON K8P 5H9', 
        20000, 5000, 'SILVER', 'NET30', contractor_pl, trenton_branch
    ),
    (
        'c0000001-0000-0000-0000-000000000004'::uuid, 'Limestone Patios', 'LP-001', 'sales@limestonepatios.ca', '613-555-4321', '321 Princess St, Kingston ON K7L 1B3', 
        25000, 2000, 'SILVER', 'NET30', contractor_pl, kingston_branch
    ),
    (
        'c0000001-0000-0000-0000-000000000005'::uuid, 'Prince Edward Pools', 'PEP-001', 'office@pepools.ca', '613-555-1122', '100 Picton Main, Picton ON K0K 2T0', 
        30000, 12000, 'GOLD', 'NET30', contractor_pl, trenton_branch
    ),
    -- Retail / DIY
    (
        'c0000001-0000-0000-0000-000000000006'::uuid, 'John Doe (DIY)', 'CASH-001', 'johndoe@email.com', '613-555-9999', '10 Applewood Ct, Trenton ON K8V 6R4', 
        0, 0, 'RETAIL', 'COD', retail_pl, trenton_branch
    ),
    (
        'c0000001-0000-0000-0000-000000000007'::uuid, 'Jane Smith (DIY)', 'CASH-002', 'janesmith@email.com', '613-555-8888', '20 Oak Ln, Kingston ON K7K 1A2', 
        0, 0, 'RETAIL', 'COD', retail_pl, kingston_branch
    ),
    -- Municipal
    (
        'c0000001-0000-0000-0000-000000000008'::uuid, 'City of Quinte West', 'CQW-001', 'purchasing@quintewest.ca', '613-392-2841', '7 Creswell Dr, Trenton ON K8V 5R6', 
        100000, 0, 'PLATINUM', 'NET60', contractor_pl, trenton_branch
    )
    ON CONFLICT (account_number) DO UPDATE SET 
        name = EXCLUDED.name,
        email = EXCLUDED.email,
        phone = EXCLUDED.phone,
        address = EXCLUDED.address,
        primary_branch_id = EXCLUDED.primary_branch_id;
END $$;
