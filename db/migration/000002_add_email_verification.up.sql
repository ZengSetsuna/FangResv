CREATE TABLE pending_users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    code TEXT NOT NULL,  -- 6 位验证码
    expires_at TIMESTAMP NOT NULL,  -- 验证码过期时间
    created_at TIMESTAMP DEFAULT now()
);
