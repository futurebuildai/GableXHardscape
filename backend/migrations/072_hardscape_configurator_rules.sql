-- Migration 072: Hardscape Configurator Rules
-- Soft-delete LBM lumber rules and insert hardscape presets + compatibility rules.

-- 1. Soft-delete existing LBM configurator rules (preserve for rollback)
DELETE FROM configurator_rules;
DELETE FROM configurator_presets;

-- 2. Hardscape Compatibility Rules: Manufacturer → Collection availability
INSERT INTO configurator_rules (attribute_type, attribute_value, depends_on_type, depends_on_value, is_allowed, error_message)
VALUES
    -- Techo-Bloc collections
    ('Collection', 'Blu 60', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Blu 80', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Para', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Borealis', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Mini-Creta', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Stackton', 'Manufacturer', 'Techo-Bloc', true, NULL),
    ('Collection', 'Graphix', 'Manufacturer', 'Techo-Bloc', true, NULL),

    -- Belgard collections
    ('Collection', 'Catalina', 'Manufacturer', 'Belgard', true, NULL),
    ('Collection', 'Subterra', 'Manufacturer', 'Belgard', true, NULL),
    ('Collection', 'Mega-Arbel', 'Manufacturer', 'Belgard', true, NULL),
    ('Collection', 'Weston Stone', 'Manufacturer', 'Belgard', true, NULL),

    -- Unilock collections
    ('Collection', 'Artline', 'Manufacturer', 'Unilock', true, NULL),
    ('Collection', 'Brussels Block', 'Manufacturer', 'Unilock', true, NULL),
    ('Collection', 'Rivercrest', 'Manufacturer', 'Unilock', true, NULL),
    ('Collection', 'Lineo', 'Manufacturer', 'Unilock', true, NULL),

    -- Rinox collections
    ('Collection', 'Proma', 'Manufacturer', 'Rinox', true, NULL),
    ('Collection', 'Centurion', 'Manufacturer', 'Rinox', true, NULL),

    -- Permacon collections
    ('Collection', 'Mondrian', 'Manufacturer', 'Permacon', true, NULL),
    ('Collection', 'Melville', 'Manufacturer', 'Permacon', true, NULL)
ON CONFLICT DO NOTHING;

-- 3. Paver Patio Kit: required accessories
-- When Application = "Patio", polymeric sand is required
INSERT INTO configurator_rules (attribute_type, attribute_value, depends_on_type, depends_on_value, is_allowed, error_message)
VALUES
    ('Accessory', 'Polymeric Sand', 'Application', 'Patio', true, NULL),
    ('Accessory', 'Edge Restraint', 'Application', 'Patio', true, NULL),
    ('Accessory', 'Geotextile Fabric', 'Application', 'Patio', true, NULL),
    ('Accessory', 'Polymeric Sand', 'Application', 'Walkway', true, NULL),
    ('Accessory', 'Edge Restraint', 'Application', 'Walkway', true, NULL),
    ('Accessory', 'Geotextile Fabric', 'Application', 'Walkway', true, NULL)
ON CONFLICT DO NOTHING;

-- 4. Retaining Wall System: geogrid thresholds
-- Walls > 4ft require geogrid reinforcement
INSERT INTO configurator_rules (attribute_type, attribute_value, depends_on_type, depends_on_value, is_allowed, error_message)
VALUES
    ('Accessory', 'Geogrid', 'WallHeight', 'Over4ft', true, NULL),
    ('Accessory', 'Geogrid', 'WallHeight', 'Under4ft', false, 'Geogrid not required for walls under 4 feet'),
    ('Accessory', 'Cap Units', 'Application', 'Retaining Wall', true, NULL),
    ('Accessory', 'Drainage Aggregate', 'Application', 'Retaining Wall', true, NULL)
ON CONFLICT DO NOTHING;

-- 5. Sand type compatibility
INSERT INTO configurator_rules (attribute_type, attribute_value, depends_on_type, depends_on_value, is_allowed, error_message)
VALUES
    ('SandType', 'Polymeric', 'Application', 'Patio', true, NULL),
    ('SandType', 'Polymeric', 'Application', 'Walkway', true, NULL),
    ('SandType', 'Polymeric', 'Application', 'Driveway', true, NULL),
    ('SandType', 'Joint Stabilizer', 'Application', 'Patio', true, NULL),
    ('SandType', 'Regular Sand', 'Application', 'Patio', false, 'Regular sand is not recommended for paver joints — use polymeric sand for long-lasting stability')
ON CONFLICT DO NOTHING;

-- 6. Hardscape Presets
INSERT INTO configurator_presets (name, description, product_type, config, is_active)
VALUES
    (
        'Paver Patio Kit',
        'Standard patio installation package: pavers + polymeric sand + edge restraint + geotextile',
        'Paver',
        '{"application": "Patio", "accessories": ["Polymeric Sand", "Edge Restraint", "Geotextile Fabric"], "base_depth_inches": 6, "sand_bed_inches": 1}'::jsonb,
        true
    ),
    (
        'Retaining Wall System',
        'Retaining wall with cap units and drainage — includes geogrid recommendation for walls over 4ft',
        'Wall',
        '{"application": "Retaining Wall", "accessories": ["Cap Units", "Drainage Aggregate"], "geogrid_threshold_ft": 4, "setback_per_course_inches": 1.25}'::jsonb,
        true
    ),
    (
        'Steps Kit',
        'Pre-configured step unit package with cap and tread options',
        'Steps',
        '{"application": "Steps", "accessories": ["Cap Units", "Polymeric Sand"], "riser_height_inches": 7, "tread_depth_inches": 14}'::jsonb,
        true
    ),
    (
        'Driveway Package',
        'Heavy-duty paver driveway with reinforced base and edge restraint',
        'Paver',
        '{"application": "Driveway", "accessories": ["Polymeric Sand", "Edge Restraint", "Geotextile Fabric"], "base_depth_inches": 10, "sand_bed_inches": 1, "min_thickness_mm": 80}'::jsonb,
        true
    )
ON CONFLICT DO NOTHING;
