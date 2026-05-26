-- 081: Dibbits Demo Hardscape Products Seed
-- Requires 070_hardscape_product_model.sql and 080_dibbits_seed_customers_vendors.sql (for vendors)

INSERT INTO products (
    id, sku, description, category, uom_primary, base_price, average_unit_cost, vendor, weight_lbs,
    manufacturer, collection, color, finish, pieces_per_sf, pallet_count, weight_per_unit
) VALUES
-- Techo-Bloc Pavers
('f0000001-0000-0000-0000-000000000001'::uuid, 'TB-BLU60-SLATE', 'Blu 60 mm Slate', 'Pavers', 'SF', 8.50, 5.25, 'Techo-Bloc', 28.0, 'Techo-Bloc', 'Blu 60', 'Slate', 'Smooth', 1.0, 116, 28.0),
('f0000001-0000-0000-0000-000000000002'::uuid, 'TB-BLU60-SMOOTH', 'Blu 60 mm Smooth', 'Pavers', 'SF', 8.75, 5.50, 'Techo-Bloc', 28.0, 'Techo-Bloc', 'Blu 60', 'Grey', 'Smooth', 1.0, 116, 28.0),
('f0000001-0000-0000-0000-000000000003'::uuid, 'TB-VILLAGIO-ONYX', 'Villagio Onyx Black', 'Pavers', 'SF', 10.20, 6.80, 'Techo-Bloc', 35.0, 'Techo-Bloc', 'Villagio', 'Onyx Black', 'Textured', 5.0, 117, 7.0),
('f0000001-0000-0000-0000-000000000004'::uuid, 'TB-VALET-SAND', 'Valet Sandlewood', 'Pavers', 'SF', 9.50, 6.00, 'Techo-Bloc', 33.0, 'Techo-Bloc', 'Valet', 'Sandlewood', 'Textured', 2.5, 93, 13.0),
('f0000001-0000-0000-0000-000000000005'::uuid, 'TB-LINEA-GREY', 'Linea Greyed Nickel', 'Pavers', 'SF', 11.00, 7.50, 'Techo-Bloc', 38.0, 'Techo-Bloc', 'Linea', 'Greyed Nickel', 'Smooth', 3.0, 77, 12.5),

-- Techo-Bloc Retaining Walls
('f0000001-0000-0000-0000-000000000006'::uuid, 'TB-MINICRETA-ARCH', 'Mini-Creta Architectural', 'Retaining Walls', 'SF', 15.50, 10.20, 'Techo-Bloc', 95.0, 'Techo-Bloc', 'Mini-Creta', 'Architectural', 'Split Face', 0.5, 40, 47.5),
('f0000001-0000-0000-0000-000000000007'::uuid, 'TB-GRAPHIX-ONYX', 'Graphix Wall Onyx Black', 'Retaining Walls', 'SF', 18.00, 12.50, 'Techo-Bloc', 85.0, 'Techo-Bloc', 'Graphix', 'Onyx Black', 'Smooth/Split', 0.8, 30, 42.5),
('f0000001-0000-0000-0000-000000000008'::uuid, 'TB-GSS-CAP', 'Graphix Double Sided Cap', 'Retaining Walls', 'LF', 22.00, 15.00, 'Techo-Bloc', 45.0, 'Techo-Bloc', 'Graphix Cap', 'Grey', 'Smooth', 1.0, 48, 45.0),

-- Unilock Pavers
('f0000001-0000-0000-0000-000000000009'::uuid, 'UL-BRUSSELS-SIERRA', 'Brussels Block Sierra', 'Pavers', 'SF', 7.50, 4.80, 'Unilock', 32.0, 'Unilock', 'Brussels Block', 'Sierra', 'Tumbled', 3.4, 96, 9.4),
('f0000001-0000-0000-0000-000000000010'::uuid, 'UL-BEACON-BAV', 'Beacon Hill Flagstone Bavarian', 'Pavers', 'SF', 9.25, 6.10, 'Unilock', 29.0, 'Unilock', 'Beacon Hill', 'Bavarian', 'Flagstone', 0.7, 104, 41.5),
('f0000001-0000-0000-0000-000000000011'::uuid, 'UL-ARTLINE-WINT', 'Artline Winter Marvel', 'Pavers', 'SF', 12.50, 8.50, 'Unilock', 36.0, 'Unilock', 'Artline', 'Winter Marvel', 'Smooth', 1.2, 85, 30.0),
('f0000001-0000-0000-0000-000000000012'::uuid, 'UL-COPB-TOWN', 'Copthorne Town Hall', 'Pavers', 'SF', 14.00, 9.80, 'Unilock', 25.0, 'Unilock', 'Copthorne', 'Town Hall', 'Distressed', 4.5, 90, 5.5),

-- Unilock Retaining Walls
('f0000001-0000-0000-0000-000000000013'::uuid, 'UL-U-CARA-PIT', 'U-Cara Pitched Face', 'Retaining Walls', 'SF', 19.50, 13.00, 'Unilock', 60.0, 'Unilock', 'U-Cara', 'Pitched Face', 'Textured', NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000014'::uuid, 'UL-LINEO-ALM', 'Lineo Dimensional Stone Almond', 'Retaining Walls', 'LF', 16.00, 11.20, 'Unilock', 55.0, 'Unilock', 'Lineo', 'Almond', 'Smooth', NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000015'::uuid, 'UL-BRUSS-DIM', 'Brussels Dimensional Stone', 'Retaining Walls', 'LF', 12.00, 8.00, 'Unilock', 40.0, 'Unilock', 'Brussels Dimensional', 'Classic', 'Tumbled', NULL, NULL, NULL),

-- Permacon Pavers
('f0000001-0000-0000-0000-000000000016'::uuid, 'PC-BOREALIS-SMOKE', 'Borealis Smoked Pine', 'Pavers', 'SF', 13.50, 9.00, 'Permacon', 30.0, 'Permacon', 'Borealis', 'Smoked Pine', 'Wood Grain', 0.8, 48, 37.5),
('f0000001-0000-0000-0000-000000000017'::uuid, 'PC-MELVILLE-ROCK', 'Melville 60 Rockland Black', 'Pavers', 'SF', 8.25, 5.50, 'Permacon', 27.0, 'Permacon', 'Melville 60', 'Rockland Black', 'Smooth', 2.0, 100, 13.5),
('f0000001-0000-0000-0000-000000000018'::uuid, 'PC-CASSARA-MARG', 'Cassara Margaux Beige', 'Pavers', 'SF', 9.80, 6.50, 'Permacon', 31.0, 'Permacon', 'Cassara', 'Margaux Beige', 'Textured', 1.5, 88, 20.6),

-- Permacon Retaining Walls & Steps
('f0000001-0000-0000-0000-000000000019'::uuid, 'PC-CELTIK-WALL', 'Celtik Wall 90', 'Retaining Walls', 'SF', 14.50, 9.80, 'Permacon', 85.0, 'Permacon', 'Celtik Wall 90', 'Grey', 'Rustic', 1.0, 40, 85.0),
('f0000001-0000-0000-0000-000000000020'::uuid, 'PC-MELVILLE-STEP', 'Melville Step Rockland Black', 'Steps', 'EA', 125.00, 85.00, 'Permacon', 200.0, 'Permacon', 'Melville', 'Rockland Black', 'Smooth', NULL, NULL, NULL),

-- Belgard Pavers
('f0000001-0000-0000-0000-000000000021'::uuid, 'BG-MEGA-ARBEL', 'Mega-Arbel Patio', 'Pavers', 'SF', 10.50, 7.00, 'Belgard', 35.0, 'Belgard', 'Mega-Arbel', 'Patio', 'Textured', 0.9, 50, 38.8),
('f0000001-0000-0000-0000-000000000022'::uuid, 'BG-DUBLIN-COB', 'Dublin Cobble', 'Pavers', 'SF', 9.00, 6.00, 'Belgard', 30.0, 'Belgard', 'Dublin Cobble', 'Classic', 'Tumbled', 3.0, 120, 10.0),
('f0000001-0000-0000-0000-000000000023'::uuid, 'BG-LAFITT-RUSTIC', 'Lafitt Rustic Slab', 'Pavers', 'SF', 11.20, 7.80, 'Belgard', 34.0, 'Belgard', 'Lafitt Rustic', 'Slab', 'Textured', 1.1, 70, 30.9),

-- Belgard Retaining Walls
('f0000001-0000-0000-0000-000000000024'::uuid, 'BG-WESTON-STONE', 'Weston Stone Universal Kit', 'Retaining Walls', 'LF', 15.00, 10.00, 'Belgard', 45.0, 'Belgard', 'Weston Stone', 'Universal Kit', 'Textured', NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000025'::uuid, 'BG-TANDEM-WALL', 'Tandem Modular Grid', 'Retaining Walls', 'SF', 21.00, 14.50, 'Belgard', 50.0, 'Belgard', 'Tandem', 'Modular Grid', 'Textured', NULL, NULL, NULL),

-- Aggregates (Lakefront Aggregate)
('f0000001-0000-0000-0000-000000000026'::uuid, 'AGG-A-GRAVEL', 'A-Gravel (Ton)', 'Aggregates', 'TON', 22.00, 14.00, 'Lakefront Aggregate', 2000.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000027'::uuid, 'AGG-HPB', 'High Performance Bedding (Ton)', 'Aggregates', 'TON', 28.00, 18.00, 'Lakefront Aggregate', 2000.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000028'::uuid, 'AGG-LIMESTONE-SCR', 'Limestone Screenings (Ton)', 'Aggregates', 'TON', 20.00, 12.50, 'Lakefront Aggregate', 2000.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000029'::uuid, 'AGG-RIVER-ROCK', '3/4" River Rock (Ton)', 'Aggregates', 'TON', 45.00, 28.00, 'Lakefront Aggregate', 2000.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),

-- Accessories & Sealers (Quinte Logistics / Generic)
('f0000001-0000-0000-0000-000000000030'::uuid, 'ACC-EDGE-ALUM', 'Aluminum Edge Restraint 8ft', 'Accessories', 'EA', 18.50, 12.00, 'Quinte Logistics', 2.5, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000031'::uuid, 'ACC-SPIKES-10', 'Landscape Spikes 10" (Box of 50)', 'Accessories', 'EA', 45.00, 28.00, 'Quinte Logistics', 15.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000032'::uuid, 'ACC-POLY-SAND-G', 'Polymeric Sand Grey 50lb', 'Accessories', 'EA', 32.00, 21.00, 'Quinte Logistics', 50.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000033'::uuid, 'ACC-POLY-SAND-T', 'Polymeric Sand Tan 50lb', 'Accessories', 'EA', 32.00, 21.00, 'Quinte Logistics', 50.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000034'::uuid, 'ACC-SEALER-NAT', 'Natural Look Paver Sealer 1Gal', 'Accessories', 'EA', 55.00, 35.00, 'Lakefront Aggregate', 8.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
('f0000001-0000-0000-0000-000000000035'::uuid, 'ACC-SEALER-WET', 'Wet Look Paver Sealer 1Gal', 'Accessories', 'EA', 65.00, 42.00, 'Lakefront Aggregate', 8.0, NULL, NULL, NULL, NULL, NULL, NULL, NULL)

ON CONFLICT (sku) DO UPDATE SET 
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    uom_primary = EXCLUDED.uom_primary,
    base_price = EXCLUDED.base_price,
    average_unit_cost = EXCLUDED.average_unit_cost,
    vendor = EXCLUDED.vendor,
    weight_lbs = EXCLUDED.weight_lbs;


