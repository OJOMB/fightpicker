INSERT INTO roles (id, name, description, created_at, created_by, updated_at, updated_by) VALUES
    (gen_random_uuid(), 'admin', 'Administrator with full access', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'user', 'Regular user with limited access', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98');

INSERT INTO permissions (id, name, version, resource, operation, description, created_at, created_by, updated_at, updated_by) VALUES
    (gen_random_uuid(), 'create', 'v1', 'users', 'post', 'Permission to create a new user', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'delete', 'v1', 'users', 'delete', 'Permission to delete a user', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'update', 'v1', 'users', 'patch', 'Permission to update user information', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'get', 'v1', 'users', 'get', 'Permission to view user information', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'list', 'v1', 'users', 'get', 'Permission to list users', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'list_followers', 'v1', 'users', 'get', 'Permission to list followers', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'list_followees', 'v1', 'users', 'get', 'Permission to list followees', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'follow', 'v1', 'users', 'put', 'Permission to follow another user', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'unfollow', 'v1', 'users', 'delete', 'Permission to unfollow another user', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'generate_presigned_url', 'v1', 'users', 'post', 'Permission to create presigned url for user profile picture upload', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),

    (gen_random_uuid(), 'create', 'v1', 'fighters', 'post', 'Permission to create a new fighter', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'delete', 'v1', 'fighters', 'delete', 'Permission to delete a fighter', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'update', 'v1', 'fighters', 'put', 'Permission to update fighter information', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'get', 'v1', 'fighters', 'get', 'Permission to view fighter information', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'),
    (gen_random_uuid(), 'ingest', 'v1', 'fighters', 'post', 'Permission to list fighters', now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98');

-- Assign all above permissions to admin and user role
INSERT INTO role_permissions (role_id, permission_id, created_at, created_by, updated_at, updated_by)
SELECT r.id, p.id, now(), '019b938b-69b9-76c5-973b-547d442dbf98', now(), '019b938b-69b9-76c5-973b-547d442dbf98'
FROM roles r, permissions p
-- admin gets all permissions. user gets all user related permissions plus get and list operations on other resources
WHERE r.name = 'admin' OR (r.name = 'user' AND p.resource = 'users' OR p.operation IN ('get', 'list'));

-- assign admin role to user 'adminuser'
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by)
VALUES (
    (SELECT id FROM users WHERE username = 'adminuser'),
    (SELECT id FROM roles WHERE name = 'admin'),
    now(),
    '019b938b-69b9-76c5-973b-547d442dbf98',
    now(),
    '019b938b-69b9-76c5-973b-547d442dbf98'
);

-- assign user role to user 'exampleuser'
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by)
VALUES (
    (SELECT id FROM users WHERE username = 'exampleuser'),
    (SELECT id FROM roles WHERE name = 'user'),
    now(),
    '019b938b-69b9-76c5-973b-547d442dbf98',
    now(),
    '019b938b-69b9-76c5-973b-547d442dbf98'
);