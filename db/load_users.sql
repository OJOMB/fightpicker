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
    created_at,
    updated_at,
    updated_by
) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'admin@fightpicker.com', 'Powerful', 'Admin', 'adminuser', 'Administrator account', 'https://example.com/images/admin.jpg', 'Headquarters', '$2a$05$nTpkCGxah9xFgVRFG89jh.3RqX6dbgs04E/p9QXLFExL1q6IpdjtK', '1990-01-01', '2024-01-01T12:00:00Z', '2024-01-01T12:00:00Z', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'user@fightpicker.com', 'Lowly', 'User', 'exampleuser', 'Example user account', 'https://example.com/images/user.jpg', 'Remote', '$2a$05$nTpkCGxah9xFgVRFG89jh.3RqX6dbgs04E/p9QXLFExL1q6IpdjtK', '1995-05-15', '2024-01-01T12:00:00Z', '2024-01-01T12:00:00Z', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb');