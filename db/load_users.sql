INSERT INTO users (
    id,
    email,
    first_name,
    last_name,
    username,
    bio,
    gender,
    profile_picture,
    location,
    password_hash,
    dob,
    created_at,
    updated_at,
    updated_by
) VALUES
    ('019b938b-69b9-76c5-973b-547d442dbf98', 'admin@fightpicker.com', 'Powerful', 'Admin', 'adminuser', 'Administrator account', 'other', 'http://localhost:8080/static/images/default-profile-picture-admin.webp', 'Headquarters', '$2a$05$nTpkCGxah9xFgVRFG89jh.3RqX6dbgs04E/p9QXLFExL1q6IpdjtK', '1990-01-01', '2024-01-01T12:00:00Z', '2024-01-01T12:00:00Z', '019b938b-69b9-76c5-973b-547d442dbf98'),
    ('019b938c-227a-7c55-b48f-0d3fde1c1131', 'user@fightpicker.com', 'Lowly', 'User', 'exampleuser', 'Example user account', 'other', '', 'Remote', '$2a$05$nTpkCGxah9xFgVRFG89jh.3RqX6dbgs04E/p9QXLFExL1q6IpdjtK', '1995-05-15', '2024-01-01T12:00:00Z', '2024-01-01T12:00:00Z', '019b938c-227a-7c55-b48f-0d3fde1c1131');