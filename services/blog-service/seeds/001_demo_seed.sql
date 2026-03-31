BEGIN;

TRUNCATE TABLE blog_posts, categories RESTART IDENTITY CASCADE;

INSERT INTO categories (id, name, slug, description, color, created_at, updated_at) VALUES
    ('3cbf3ab3-1f9e-4282-8e20-c6f2fd56a001', 'Training', 'training', 'Workout strategy, routines, and recovery notes.', '#2563EB', '2026-03-01T08:00:00Z', '2026-03-25T08:00:00Z'),
    ('3cbf3ab3-1f9e-4282-8e20-c6f2fd56a002', 'Workspace', 'workspace', 'Tools, setups, and habits for focused work.', '#059669', '2026-03-01T08:05:00Z', '2026-03-25T08:05:00Z'),
    ('3cbf3ab3-1f9e-4282-8e20-c6f2fd56a003', 'Travel', 'travel', 'Packing, commuting, and gear-use stories.', '#D97706', '2026-03-01T08:10:00Z', '2026-03-25T08:10:00Z');

INSERT INTO blog_posts (
    id, title, slug, content, excerpt, author_id, author_name, author_email, status,
    featured_image, tags, category_id, published_at, created_at, updated_at, view_count, read_time
) VALUES
    (
        '8d3f8a09-6337-45ec-b2f4-98b1d2cb1001',
        'How to Build a 30-Minute Strength Routine That Actually Sticks',
        '30-minute-strength-routine',
        E'# How to Build a 30-Minute Strength Routine That Actually Sticks\n\nConsistency beats complexity. Start with three full-body sessions, keep the movement patterns simple, and track one small win every week.\n\n## A practical split\n\n- Squat or lunge\n- Push movement\n- Pull movement\n- Hinge\n- Carries or core finisher\n\nThe routine works because it is light on setup friction and easy to repeat on busy weeks.',
        'A simple full-body training structure for busy schedules.',
        '1', 'Ava Admin', 'ava.admin@gearbox.local', 'published',
        'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=1200&h=630&fit=crop',
        ARRAY['training','fitness','beginner'],
        '3cbf3ab3-1f9e-4282-8e20-c6f2fd56a001',
        '2026-03-12T09:00:00Z', '2026-03-11T18:00:00Z', '2026-03-12T09:00:00Z', 482, 6
    ),
    (
        '8d3f8a09-6337-45ec-b2f4-98b1d2cb1002',
        'The Five Desk Upgrades That Changed Our Team''s Focus',
        'desk-upgrades-that-improved-focus',
        E'# The Five Desk Upgrades That Changed Our Team''s Focus\n\nSmall environment changes compound quickly. Better light, cleaner cable runs, and one predictable place for everything can reduce the daily drag on attention.\n\nWe tested a handful of cheap upgrades and found that lighting and laptop elevation had the fastest payoff.',
        'A practical look at the workspace tweaks with the best daily payoff.',
        '2', 'Eli Editor', 'eli.editor@gearbox.local', 'published',
        'https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=1200&h=630&fit=crop',
        ARRAY['workspace','productivity','home-office'],
        '3cbf3ab3-1f9e-4282-8e20-c6f2fd56a002',
        '2026-03-18T10:30:00Z', '2026-03-17T17:00:00Z', '2026-03-18T10:30:00Z', 356, 5
    ),
    (
        '8d3f8a09-6337-45ec-b2f4-98b1d2cb1003',
        'Packing Light for a Three-Day City Trip',
        'packing-light-city-trip',
        E'# Packing Light for a Three-Day City Trip\n\nA small backpack beats an overpacked roller in most short-trip cases. The trick is picking one outfit system and being honest about what you will actually use.\n\nThis checklist prioritizes movement, flexibility, and one spare layer.',
        'A carry-on-first packing list for short trips.',
        '3', 'Nina Patel', 'nina.ops@gearbox.local', 'published',
        'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=1200&h=630&fit=crop',
        ARRAY['travel','packing','gear'],
        '3cbf3ab3-1f9e-4282-8e20-c6f2fd56a003',
        '2026-03-22T08:45:00Z', '2026-03-21T16:30:00Z', '2026-03-22T08:45:00Z', 291, 4
    ),
    (
        '8d3f8a09-6337-45ec-b2f4-98b1d2cb1004',
        'Why We Started Treating Inventory Alerts as Product Feedback',
        'inventory-alerts-as-product-feedback',
        E'# Why We Started Treating Inventory Alerts as Product Feedback\n\nLow-stock alerts are not just operations noise. They tell you where demand, merchandising, and forecasting are drifting apart.\n\nWhen the same SKUs keep surfacing, the right move is usually not a spreadsheet fix. It is a product decision.',
        'Inventory signals can reveal customer demand and catalog issues.',
        '1', 'Ava Admin', 'ava.admin@gearbox.local', 'draft',
        'https://images.unsplash.com/photo-1556740749-887f6717d7e4?w=1200&h=630&fit=crop',
        ARRAY['inventory','operations','product-strategy'],
        '3cbf3ab3-1f9e-4282-8e20-c6f2fd56a002',
        NULL, '2026-03-28T13:15:00Z', '2026-03-28T13:15:00Z', 0, 7
    ),
    (
        '8d3f8a09-6337-45ec-b2f4-98b1d2cb1005',
        'Choosing Between Daily Sneakers and Travel Sneakers',
        'daily-vs-travel-sneakers',
        E'# Choosing Between Daily Sneakers and Travel Sneakers\n\nThe best daily sneaker is not always the best travel sneaker. One needs quick comfort for repeat use. The other needs versatility, packability, and easy cleaning.\n\nHere is the framework we use when we compare both roles.',
        'A simple framework for picking sneakers by use case.',
        '4', 'Mason Lee', 'mason.shop@gearbox.local', 'published',
        'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=1200&h=630&fit=crop',
        ARRAY['footwear','travel','buying-guide'],
        '3cbf3ab3-1f9e-4282-8e20-c6f2fd56a003',
        '2026-03-29T14:00:00Z', '2026-03-29T09:00:00Z', '2026-03-29T14:00:00Z', 164, 5
    );

COMMIT;
