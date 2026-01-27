CREATE EXTENSION citext;

CREATE TYPE gender AS ENUM ('male', 'female', 'other');

-- users table to store user information and authentication details
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    username CITEXT NOT NULL CHECK (char_length(username) BETWEEN 3 AND 32) UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    dob DATE CHECK (dob <= CURRENT_DATE) NOT NULL,
    gender gender NOT NULL,
    location VARCHAR(255),
    bio TEXT,
    profile_picture VARCHAR(255),
    password_hash TEXT NOT NULL,

    -- Email verification fields
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token_hash BYTEA,
    email_verification_token_expires_at TIMESTAMPTZ,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    CHECK (
        (email_verification_token_hash IS NULL AND email_verification_token_expires_at IS NULL)
        OR
        (email_verification_token_hash IS NOT NULL AND email_verification_token_expires_at IS NOT NULL)
    ),
    CHECK (
        email_verification_token_expires_at IS NULL
        OR email_verification_token_expires_at > created_at
    ),
    CHECK (char_length(first_name) BETWEEN 2 AND 255),
    CHECK (char_length(last_name) BETWEEN 2 AND 255)
);

CREATE UNIQUE INDEX uniq_users_email_verification_token_hash
ON users (email_verification_token_hash)
WHERE email_verification_token_hash IS NOT NULL;

CREATE INDEX idx_users_email_verification_token_hash
ON users (email_verification_token_hash)
WHERE email_verified = FALSE;

-- roles table to define different user roles (e.g., admin, moderator, user)
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    CHECK (char_length(name) BETWEEN 3 AND 50)
);

CREATE TYPE resource_type AS ENUM ('users', 'fighters', 'cards', 'fights', 'picks');

-- permissions table to define specific actions that can be performed on resources
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(16) NOT NULL,
    resource resource_type NOT NULL,
    operation VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    UNIQUE (name, version, resource, operation)
);

-- role_permissions table to link roles to their permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

-- user_roles table to link users to their roles
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    PRIMARY KEY (user_id, role_id)
);

-- refresh_tokens table to store issued refresh tokens for users
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    jti UUID NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    replaced_by UUID REFERENCES refresh_tokens(id),
    ip_address TEXT NULL, -- IP address from which the token was issued
    user_agent TEXT NULL, -- User agent string of the client device, these are stored for forensics (security, auditing) purposes

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX ON refresh_tokens (user_id, revoked);
CREATE INDEX ON refresh_tokens (expires_at);

-- fighters table to represent individual fighters
CREATE TABLE IF NOT EXISTS fighters (
    id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    gender gender NOT NULL,
    dob DATE CHECK (dob <= CURRENT_DATE) NOT NULL,
    height NUMERIC(5,2) NOT NULL,
    weight NUMERIC(5,2) NOT NULL,
    reach NUMERIC(5,2) NOT NULL,
    stance VARCHAR(50) NOT NULL,
    country VARCHAR(100) NOT NULL,
    fighting_out_of VARCHAR(100) NOT NULL,
    profile_picture VARCHAR(255),
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,
    disqualifications INTEGER NOT NULL DEFAULT 0,
    no_contests INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TYPE source AS ENUM ('ufcstats', 'tapology', 'sherdog');

-- fighter_external_ids table to store external IDs for fighters from various sources
-- e.g. Sherdog ID, Tapology ID, etc.
CREATE TABLE fighter_external_ids (
    fighter_id UUID NOT NULL REFERENCES fighters(id) ON DELETE CASCADE ON UPDATE CASCADE,
    source source NOT NULL,
    external_id TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    PRIMARY KEY (source, external_id),
    UNIQUE (fighter_id, source)
);

CREATE TYPE fight_status AS ENUM ('scheduled', 'in_progress', 'completed', 'cancelled');

-- fights table to represent individual fights between fighters
CREATE TABLE IF NOT EXISTS fights (
    id UUID PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL,
    no_of_rounds INTEGER NOT NULL CHECK (no_of_rounds > 0) DEFAULT 3,
    weight_class VARCHAR(50) NOT NULL,
    is_championship BOOLEAN NOT NULL DEFAULT FALSE,
    status fight_status NOT NULL DEFAULT 'scheduled',
    referee VARCHAR(255),
    location VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
);

-- cards table to represent fight cards/events
CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    promotion VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TYPE card_stage AS ENUM ('early_prelims', 'prelims', 'main');

CREATE TABLE IF NOT EXISTS card_fights (
    id UUID PRIMARY KEY,
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE ON UPDATE CASCADE,
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stage card_stage NOT NULL,
    order_in_card INTEGER NOT NULL, -- 1 is main event, 2 is co-main, etc.

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    UNIQUE (card_id, fight_id),
    UNIQUE (card_id, order_in_card)
);

-- ko = knockout: fighter is rendered unconscious
-- tko = technical knockout: referee stops the fight due to one fighter being unable to defend themselves
-- dq = disqualification: fighter is disqualified for illegal actions
-- sub = submission: fighter taps out or verbally submits
-- ud = unanimous decision: all three judges score for one fighter
-- sd = split decision: two judges score for one fighter, one for the other
-- md = majority decision: two judges score for one fighter, one scores a draw
CREATE TYPE win_method AS ENUM (
    'ko', 'tko', 'dq', 'sub', 'ud', 'sd', 'md'
);

CREATE TYPE corner AS ENUM ('red', 'blue');

-- fight_participants table to link fighters to fights and their corner assignments
CREATE TABLE fight_participants (
    id UUID PRIMARY KEY,
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE ON UPDATE CASCADE,
    fighter_id UUID NOT NULL REFERENCES fighters(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    corner corner NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    UNIQUE (fight_id, fighter_id), -- a fighter can only participate once in a fight
    UNIQUE (fight_id, corner), -- only one fighter per corner
    UNIQUE (fight_id, id) -- fight_participant id must be unique per fight for referencing in picks table
);

CREATE INDEX ON fight_participants (fight_id);
CREATE INDEX ON fight_participants (fighter_id);

-- fight_results table to store the recorded outcome of fights to compare against user picks
CREATE TABLE IF NOT EXISTS fight_results (
    id UUID PRIMARY KEY,
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE ON UPDATE CASCADE UNIQUE,
    is_draw BOOLEAN NOT NULL DEFAULT FALSE,
    is_no_contest BOOLEAN NOT NULL DEFAULT FALSE,
    winner_id UUID REFERENCES fight_participants(id),
    method win_method,
    round INTEGER,
    fight_time INTERVAL, -- Time elapsed in the fight when it ended
    judges_scorecard JSONB,
    referee VARCHAR(255),
    fight_of_the_night BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    CHECK (
        (
            is_draw = FALSE AND
            is_no_contest = FALSE AND
            winner_id IS NOT NULL AND
            method IS NOT NULL
        )
    OR
        (
            is_draw = TRUE AND
            is_no_contest = FALSE AND
            winner_id IS NULL AND
            method IS NULL
        )
    OR
        (
            is_no_contest = TRUE AND
            winner_id IS NULL AND
            method IS NULL
        )
    )
);

-- picks table to store user predictions for fight outcomes
CREATE TABLE picks (
    id UUID PRIMARY KEY,
    fight_id UUID NOT NULL REFERENCES fights(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    fight_participant_id UUID NOT NULL REFERENCES fight_participants(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    method win_method,
    round INTEGER CHECK (round > 0),

    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    FOREIGN KEY (fight_id, fight_participant_id) REFERENCES fight_participants (fight_id, id), -- ensure the picked participant is in the correct fight
    UNIQUE (user_id, fight_id) -- a user can only make one pick per fight
);

CREATE INDEX ON picks (fight_id);
CREATE INDEX ON picks (user_id);

-- followers table to represent "user A follows user B" relationships
-- the id field is not strictly necessary, we could us a composite primary key of (follower_id, followee_id)
-- but because we use UUID7 it means we can have consistent and simple pagination based on the id field alone
-- I've chosen to go this way for simplicity and consistency with pagination across the api
CREATE TABLE IF NOT EXISTS followers (
    id UUID PRIMARY KEY,
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    followee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,

    UNIQUE (follower_id, followee_id),
    CHECK (follower_id <> followee_id)
);

-- We return followers/followees in descending order of id (newest first).
-- Postgres can scan a Btree index backwards but explicitly matching direction improves planner confidence.
CREATE INDEX followers_followee_id_id_desc
  ON followers (followee_id, id DESC);

CREATE INDEX followers_follower_id_id_desc
  ON followers (follower_id, id DESC);

-- REVOKE UPDATE (created_at) ON ALL TABLES IN SCHEMA fightpicker FROM postgres;
-- REVOKE UPDATE (created_by) ON ALL TABLES IN SCHEMA fightpicker FROM postgres;