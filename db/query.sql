-------------------------------------
-- SQL queries for user management --
-------------------------------------

-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    profile_picture,
    location,
    password_hash,
    dob,
    gender,
    created_at,
    updated_at,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14
);

-- name: GetUserByID :one
SELECT
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    profile_picture,
    location,
    password_hash,
    dob,
    gender,
    created_at,
    updated_at,
    updated_by
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    profile_picture,
    location,
    password_hash,
    dob,
    gender,
    created_at,
    updated_at,
    updated_by
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    profile_picture,
    location,
    password_hash,
    dob,
    gender,
    created_at,
    updated_at,
    updated_by
FROM users
WHERE username = $1;

-- name: UpdateUserByID :exec
UPDATE users SET
    email = $2,
    first_name = $3,
    last_name = $4,
    username = $5,
    bio = $6,
    profile_picture = $7,
    location = $8,
    password_hash = $9,
    dob = $10,
    gender = $11,
    updated_at = $12,
    updated_by = $13
WHERE id = $1;

-- name: UpdateUserProfilePictureByID :exec
UPDATE users SET
    profile_picture = $2,
    updated_at = $3,
    updated_by = $1
WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    profile_picture,
    location,
    dob,
    gender,
    created_at,
    updated_at,
    updated_by
FROM users
WHERE id > $1
ORDER BY id
LIMIT $2;

-- name: FollowUser :exec
INSERT INTO followers (id, follower_id, followee_id, created_at, created_by)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (follower_id, followee_id) DO NOTHING;

-- name: UnfollowUser :exec
DELETE FROM followers
WHERE follower_id = $1 AND followee_id = $2;

-- name: IsFollowing :one
SELECT EXISTS (
  SELECT 1
  FROM followers
  WHERE follower_id = $1 AND followee_id = $2
) AS is_following;


-- for first page (i.e. no last_seen_id) the app passes sentinel max value of all Fs for UUID7
-- name: ListFollowers :many
SELECT
  u.id, u.email, u.first_name, u.last_name, u.username,
  u.bio, u.profile_picture, u.location, u.dob, u.gender,
  u.created_at, u.updated_at, u.updated_by
FROM followers f
JOIN users u ON u.id = f.follower_id
WHERE f.followee_id = $1
  AND f.id < $2
ORDER BY f.id DESC
LIMIT $3;

-- for first page (i.e. no last_seen_id) the app passes sentinel max value of all Fs for UUID7
-- name: ListFollowees :many
SELECT
  u.id, u.email, u.first_name, u.last_name, u.username,
  u.bio, u.profile_picture, u.location, u.dob, u.gender,
  u.created_at, u.updated_at, u.updated_by
FROM followers f
JOIN users u ON u.id = f.followee_id
WHERE f.follower_id = $1
  AND f.id < $2
ORDER BY f.id DESC
LIMIT $3;

-- name: CountFollowers :one
SELECT COUNT(*) AS follower_count
FROM followers
WHERE followee_id = $1;

-- name: CountFollowing :one
SELECT COUNT(*) AS following_count
FROM followers
WHERE follower_id = $1;

-- name: IsUserEmailVerifiedByID :one
SELECT email_verified
FROM users
WHERE id = $1;

-- name: UpdateEmailVerificationTokenHashByUserID :exec
UPDATE users SET
    email_verification_token_hash = $2,
    email_verification_token_expires_at = $3,
    updated_at = $4
WHERE id = $1;

-- name: VerifyUserEmailByTokenHash :one
UPDATE users
SET
    email_verified = TRUE,
    email_verification_token_hash = NULL,
    email_verification_token_expires_at = NULL,
    updated_at = $2
WHERE
    email_verification_token_hash = $1
    AND email_verification_token_expires_at > NOW()
    AND email_verified = FALSE
RETURNING id;

-------------------------------------
-- SQL queries for auth management --
-------------------------------------

-- name: StoreRefreshToken :exec
INSERT INTO refresh_tokens (
    id,
    user_id,
    token_hash,
    jti,
    expires_at,
    ip_address,
    user_agent,
    created_at,
    updated_at
)
VALUES ($1,$2,$3,$4,$5,$6, $7, $8, $9);

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, jti, expires_at, revoked, replaced_by, ip_address, user_agent, created_at, updated_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeAndRotateRefreshTokenByHash :exec
UPDATE refresh_tokens SET
    revoked = TRUE,
    replaced_by = $2,
    updated_at = $3
WHERE token_hash = $1;

-- name: RevokeRefreshTokenByHash :exec
UPDATE refresh_tokens SET
    revoked = TRUE,
    updated_at = $2
WHERE token_hash = $1;

-- name: GetUserPermissionsByID :many
SELECT p.name, p.resource, p.version, p.operation
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE u.id = $1
ORDER BY p.resource, p.name;

-- name: AssignRoleToUserByRoleName :exec
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by)
VALUES (
    $1,
    (SELECT id FROM roles WHERE name = $2),
    $3,
    $4,
    $5,
    $6
);

-- name: GetUserRolesByID :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1
ORDER BY r.name;

-- name: GetUserPermissionsAndRolesByID :many
SELECT
    r.name AS role_name,
    p.name AS permission_name,
    p.resource,
    p.version,
    p.operation
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE u.id = $1
ORDER BY r.name, p.resource, p.name;

-- name: GetUserRBAC :one
SELECT
    u.id AS user_id,
    json_agg(DISTINCT r.name) AS roles,
    json_agg(
        DISTINCT jsonb_build_object(
            'id', p.id,
            'name', p.name,
            'resource', p.resource,
            'version', p.version,
            'operation', p.operation
        )
    ) AS permissions
FROM users u
JOIN user_roles ur            ON ur.user_id = u.id
JOIN roles r                  ON r.id = ur.role_id
JOIN role_permissions rp      ON rp.role_id = r.id
JOIN permissions p            ON p.id = rp.permission_id
WHERE u.id = $1
GROUP BY u.id;

----------------------------------------
-- SQL queries for fighter management --
----------------------------------------

-- name: GetFighterByID :one
SELECT
    id,
    first_name,
    last_name,
    nickname,
    gender,
    dob,
    height,
    weight,
    reach,
    stance,
    country,
    fighting_out_of,
    bio,
    profile_picture,
    wins,
    losses,
    draws,
    disqualifications,
    no_contests,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM fighters
WHERE id = $1;

-- name: CreateFighter :exec
INSERT INTO fighters (
    id,
    first_name,
    last_name,
    nickname,
    gender,
    dob,
    height,
    weight,
    reach,
    stance,
    country,
    fighting_out_of,
    bio,
    profile_picture,
    wins,
    losses,
    draws,
    disqualifications,
    no_contests,
    created_at,
    created_by,
    updated_at,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20,
    $21, $22, $23
);