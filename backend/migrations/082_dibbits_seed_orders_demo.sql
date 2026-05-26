-- 082: Dibbits Demo Orders and History
-- Requires 080_dibbits_seed_customers_vendors.sql and 081_hardscape_product_seed.sql

DO $$
DECLARE
    -- Customer UUIDs (from 080)
    quinte_ltd uuid := 'c0000001-0000-0000-0000-000000000001'::uuid;
    kingston_co uuid := 'c0000001-0000-0000-0000-000000000002'::uuid;
    pe_pools uuid := 'c0000001-0000-0000-0000-000000000005'::uuid;
    
    -- Product UUIDs (from 081)
    blu60_slate uuid := 'f0000001-0000-0000-0000-000000000001'::uuid;
    brussels uuid := 'f0000001-0000-0000-0000-000000000009'::uuid;
    poly_sand uuid := 'f0000001-0000-0000-0000-000000000032'::uuid;
    
    -- Branches
    trenton_branch uuid := 'a0000001-0000-0000-0000-000000000001'::uuid;
    kingston_branch uuid := 'a0000001-0000-0000-0000-000000000002'::uuid;
    
    -- Variables
    order_id1 uuid := gen_random_uuid();
    order_id2 uuid := gen_random_uuid();
    order_id3 uuid := gen_random_uuid();
    inv_id1 uuid := gen_random_uuid();
    inv_id2 uuid := gen_random_uuid();
    inv_id3 uuid := gen_random_uuid();
    sales_rep uuid;
BEGIN
    -- Get a default salesperson
    SELECT id INTO sales_rep FROM sales_team LIMIT 1;
    IF sales_rep IS NULL THEN
        -- Fallback if no sales team
        INSERT INTO sales_team (id, name, email, phone, role) 
        VALUES (gen_random_uuid(), 'System Rep', 'rep@dibbits.ca', '111-222-3333', 'Sales Rep') RETURNING id INTO sales_rep;
    END IF;

    -- Create Orders
    INSERT INTO orders (id, customer_id, branch_id, total_amount, status, salesperson_id, created_at) VALUES
    (order_id1, quinte_ltd, trenton_branch, 4250.00, 'FULFILLED', sales_rep, NOW() - INTERVAL '45 days'),
    (order_id2, kingston_co, kingston_branch, 1850.50, 'FULFILLED', sales_rep, NOW() - INTERVAL '15 days'),
    (order_id3, pe_pools, trenton_branch, 3100.00, 'CONFIRMED', sales_rep, NOW() - INTERVAL '2 days');

    -- Order Lines
    -- Order 1: 500 sqft Blu60 + 10 bags Poly Sand
    INSERT INTO order_lines (order_id, product_id, quantity, price_each) VALUES
    (order_id1, blu60_slate, 500, 8.50),
    (order_id1, poly_sand, 10, 32.00);

    -- Order 2: 200 sqft Brussels
    INSERT INTO order_lines (order_id, product_id, quantity, price_each) VALUES
    (order_id2, brussels, 246.7, 7.50);

    -- Order 3: 300 sqft Blu60
    INSERT INTO order_lines (order_id, product_id, quantity, price_each) VALUES
    (order_id3, blu60_slate, 364.7, 8.50);

    -- Create Invoices (for fulfilled orders)
    INSERT INTO invoices (id, order_id, customer_id, branch_id, status, total_amount, subtotal, tax_rate, tax_amount, due_date, payment_terms, created_at) VALUES
    (inv_id1, order_id1, quinte_ltd, trenton_branch, 'PAID', 4802.50, 4250.00, 0.13, 552.50, (NOW() - INTERVAL '45 days') + INTERVAL '30 days', 'NET30', NOW() - INTERVAL '45 days'),
    (inv_id2, order_id2, kingston_co, kingston_branch, 'UNPAID', 2091.06, 1850.50, 0.13, 240.56, (NOW() - INTERVAL '15 days') + INTERVAL '45 days', 'NET45', NOW() - INTERVAL '15 days');

    -- Create Payments
    INSERT INTO payments (invoice_id, amount, method, reference, notes) VALUES
    (inv_id1, 4802.50, 'CARD', 'TX-109841', 'Paid via online portal');

    -- Update balances
    UPDATE customers SET balance_due = balance_due + 2091.06 WHERE id = kingston_co;

END $$;
