BEGIN;

TRUNCATE TABLE
    password_reset_tokens,
    oauth_providers,
    refresh_tokens,
    user_phone_numbers,
    user_addresses,
    users,
    roles
RESTART IDENTITY CASCADE;

INSERT INTO roles (id, name, description, is_active, is_deleted, created_at, updated_at) VALUES
    (1, 'user', 'Basic authenticated user with limited access', TRUE, FALSE, '2026-01-01T09:00:00Z', '2026-01-01T09:00:00Z'),
    (2, 'editor', 'Content editor with create/edit/moderate permissions', TRUE, FALSE, '2026-01-01T09:00:00Z', '2026-01-01T09:00:00Z'),
    (3, 'admin', 'Full system administrator with complete access', TRUE, FALSE, '2026-01-01T09:00:00Z', '2026-01-01T09:00:00Z');

INSERT INTO users (
    id, username, password, email, first_name, middle_name, last_name, avatar, gender,
    date_of_birth, phone_number, role_id, is_active, is_deleted, created_at, updated_at
) VALUES
    (
        1, 'ava_admin', '$2b$12$xAKs8dILkj0o3aft7lQE4eKIWgzOVlVy8r6bifsGiYCEOmx3m301G',
        'ava.admin@gearbox.local', 'Ava', NULL, 'Admin',
        'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop',
        'female', '1991-04-12', '+1-202-555-0101', 3, TRUE, FALSE,
        '2026-01-03T10:00:00Z', '2026-03-25T09:00:00Z'
    ),
    (
        2, 'eli_editor', '$2b$12$xAKs8dILkj0o3aft7lQE4eKIWgzOVlVy8r6bifsGiYCEOmx3m301G',
        'eli.editor@gearbox.local', 'Eli', NULL, 'Editor',
        'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=400&h=400&fit=crop',
        'male', '1993-08-27', '+1-202-555-0102', 2, TRUE, FALSE,
        '2026-01-06T11:30:00Z', '2026-03-12T13:10:00Z'
    ),
    (
        3, 'nina_ops', '$2b$12$xAKs8dILkj0o3aft7lQE4eKIWgzOVlVy8r6bifsGiYCEOmx3m301G',
        'nina.ops@gearbox.local', 'Nina', 'R.', 'Patel',
        'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=400&h=400&fit=crop',
        'female', '1995-02-18', '+1-202-555-0103', 1, TRUE, FALSE,
        '2026-01-10T08:45:00Z', '2026-03-18T16:00:00Z'
    ),
    (
        4, 'mason_shop', '$2b$12$xAKs8dILkj0o3aft7lQE4eKIWgzOVlVy8r6bifsGiYCEOmx3m301G',
        'mason.shop@gearbox.local', 'Mason', NULL, 'Lee',
        'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=400&h=400&fit=crop',
        'male', '1990-11-03', '+1-202-555-0104', 1, TRUE, FALSE,
        '2026-02-02T14:20:00Z', '2026-03-29T18:20:00Z'
    ),
    (
        5, 'zoe_oauth', NULL,
        'zoe.oauth@gearbox.local', 'Zoe', NULL, 'Kim',
        'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=400&h=400&fit=crop',
        'prefer_not_to_say', '1997-06-21', '+1-202-555-0105', 1, TRUE, FALSE,
        '2026-02-14T17:05:00Z', '2026-03-30T11:50:00Z'
    );

INSERT INTO user_addresses (
    id, user_id, address_type, first_name, last_name, company, address_line_1, address_line_2,
    city, state, country, postal_code, phone, email, is_verified, is_default, is_deleted,
    created_at, updated_at
) VALUES
    (
        1, 1, 'billing', 'Ava', 'Admin', 'Gearbox HQ', '120 Market Street', 'Suite 400',
        'San Francisco', 'California', 'United States', '94105', '+1-202-555-0101',
        'ava.admin@gearbox.local', TRUE, TRUE, FALSE, '2026-01-03T10:15:00Z', '2026-03-25T09:00:00Z'
    ),
    (
        2, 3, 'shipping', 'Nina', 'Patel', NULL, '42 River Walk', 'Apt 9B',
        'Austin', 'Texas', 'United States', '78701', '+1-202-555-0103',
        'nina.ops@gearbox.local', TRUE, TRUE, FALSE, '2026-01-11T09:15:00Z', '2026-03-20T10:00:00Z'
    ),
    (
        3, 4, 'shipping', 'Mason', 'Lee', NULL, '88 Cedar Avenue', NULL,
        'Seattle', 'Washington', 'United States', '98101', '+1-202-555-0104',
        'mason.shop@gearbox.local', TRUE, TRUE, FALSE, '2026-02-02T14:30:00Z', '2026-03-29T18:25:00Z'
    ),
    (
        4, 4, 'billing', 'Mason', 'Lee', 'Northwind Labs', '200 Industrial Road', 'Building C',
        'Seattle', 'Washington', 'United States', '98109', '+1-202-555-0174',
        'billing@northwind.example', FALSE, FALSE, FALSE, '2026-02-03T10:30:00Z', '2026-03-15T12:00:00Z'
    ),
    (
        5, 5, 'home', 'Zoe', 'Kim', NULL, '17 Lake Shore Drive', NULL,
        'Chicago', 'Illinois', 'United States', '60601', '+1-202-555-0105',
        'zoe.oauth@gearbox.local', TRUE, TRUE, FALSE, '2026-02-14T17:15:00Z', '2026-03-30T11:50:00Z'
    );

INSERT INTO user_phone_numbers (
    id, user_id, phone_type, phone_number, country_code, is_verified, is_primary, is_deleted, created_at, updated_at
) VALUES
    (1, 1, 'mobile', '2025550101', '+1', TRUE, TRUE, FALSE, '2026-01-03T10:10:00Z', '2026-03-25T09:00:00Z'),
    (2, 2, 'work', '2025550102', '+1', TRUE, TRUE, FALSE, '2026-01-06T11:40:00Z', '2026-03-12T13:10:00Z'),
    (3, 3, 'mobile', '2025550103', '+1', TRUE, TRUE, FALSE, '2026-01-10T08:50:00Z', '2026-03-18T16:00:00Z'),
    (4, 4, 'mobile', '2025550104', '+1', TRUE, TRUE, FALSE, '2026-02-02T14:25:00Z', '2026-03-29T18:20:00Z'),
    (5, 5, 'mobile', '2025550105', '+1', TRUE, TRUE, FALSE, '2026-02-14T17:10:00Z', '2026-03-30T11:50:00Z');

INSERT INTO oauth_providers (
    id, user_id, provider, provider_user_id, email, access_token, refresh_token, expires_at, created_at, updated_at
) VALUES
    (
        1, 5, 'github', 'github-zoe-kim-001', 'zoe.oauth@gearbox.local',
        'gho_demo_access_token', 'ghr_demo_refresh_token', '2026-06-30T00:00:00Z',
        '2026-02-14T17:20:00Z', '2026-03-30T11:50:00Z'
    );

INSERT INTO refresh_tokens (
    id, user_id, refresh_token, user_agent, ip_address, expires_at, last_used_at, is_revoked, created_at, updated_at
) VALUES
    (
        1, 1, 'demo-refresh-admin-token', 'Mozilla/5.0 (Macintosh; Intel Mac OS X)', '127.0.0.1',
        '2026-04-30T00:00:00Z', '2026-03-31T06:30:00Z', FALSE, '2026-03-01T08:00:00Z', '2026-03-31T06:30:00Z'
    ),
    (
        2, 4, 'demo-refresh-customer-token', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', '127.0.0.1',
        '2026-04-20T00:00:00Z', '2026-03-30T19:00:00Z', FALSE, '2026-03-10T12:00:00Z', '2026-03-30T19:00:00Z'
    );

INSERT INTO password_reset_tokens (
    id, user_id, token, expires_at, used, created_at, updated_at
) VALUES
    (
        1, 4, 'demo-password-reset-token-mason', '2026-04-01T12:00:00Z', FALSE,
        '2026-03-31T08:00:00Z', '2026-03-31T08:00:00Z'
    );

SELECT setval('roles_id_seq', COALESCE((SELECT MAX(id) FROM roles), 1), TRUE);
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1), TRUE);
SELECT setval('user_addresses_id_seq', COALESCE((SELECT MAX(id) FROM user_addresses), 1), TRUE);
SELECT setval('user_phone_numbers_id_seq', COALESCE((SELECT MAX(id) FROM user_phone_numbers), 1), TRUE);
SELECT setval('oauth_providers_id_seq', COALESCE((SELECT MAX(id) FROM oauth_providers), 1), TRUE);
SELECT setval('refresh_tokens_id_seq', COALESCE((SELECT MAX(id) FROM refresh_tokens), 1), TRUE);
SELECT setval('password_reset_tokens_id_seq', COALESCE((SELECT MAX(id) FROM password_reset_tokens), 1), TRUE);

COMMIT;
