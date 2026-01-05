INSERT INTO roles (id, name, description, created_at, created_by, updated_at, updated_by) VALUES
    (gen_random_uuid(), 'admin', 'Administrator with full access', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'user', 'Regular user with limited access', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

INSERT INTO permissions (id, name, version, resource, operation, description, created_at, created_by, updated_at, updated_by) VALUES
    (gen_random_uuid(), 'create', 'v1', 'users', 'post', 'Permission to create a new user', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'delete', 'v1', 'users', 'delete', 'Permission to delete a user', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'update', 'v1', 'users', 'patch', 'Permission to update user information', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'get', 'v1', 'users', 'get', 'Permission to view user information', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'list', 'v1', 'users', 'get', 'Permission to list users', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'list_followers', 'v1', 'users', 'get', 'Permission to list followers', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'list_followees', 'v1', 'users', 'get', 'Permission to list followees', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'follow', 'v1', 'users', 'put', 'Permission to follow another user', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'unfollow', 'v1', 'users', 'delete', 'Permission to unfollow another user', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'generate_presigned_url', 'v1', 'users', 'post', 'Permission to create presigned url for user profile picture upload', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),

    (gen_random_uuid(), 'create', 'v1', 'fighters', 'post', 'Permission to create a new fighter', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'delete', 'v1', 'fighters', 'delete', 'Permission to delete a fighter', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'update', 'v1', 'fighters', 'put', 'Permission to update fighter information', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    (gen_random_uuid(), 'get', 'v1', 'fighters', 'get', 'Permission to view fighter information', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- Assign all above permissions to admin and user role
INSERT INTO role_permissions (role_id, permission_id, created_at, created_by, updated_at, updated_by)
SELECT r.id, p.id, now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now(), 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
FROM roles r, permissions p
-- admin gets all permissions. user gets all user related permissions plus get and list operations on other resources
WHERE r.name = 'admin' OR (r.name = 'user' AND p.resource = 'users' OR p.operation IN ('get', 'list'));

-- assign admin role to user 'adminuser'
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by)
VALUES (
    (SELECT id FROM users WHERE username = 'adminuser'),
    (SELECT id FROM roles WHERE name = 'admin'),
    now(),
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    now(),
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
);

-- assign user role to user 'exampleuser'
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by)
VALUES (
    (SELECT id FROM users WHERE username = 'exampleuser'),
    (SELECT id FROM roles WHERE name = 'user'),
    now(),
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    now(),
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
);