CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE venues (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    address TEXT NOT NULL,
    max_capacity INT NOT NULL CHECK (max_capacity > 0)
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    creator_id INT REFERENCES users(id) ON DELETE CASCADE,
    venue_id INT REFERENCES venues(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location TEXT NOT NULL,
    max_participants INT NOT NULL CHECK (max_participants > 0),
    created_at TIMESTAMP DEFAULT now(),
    current_participants INT DEFAULT 0
    --CONSTRAINT no_overlap UNIQUE (venue_id, start_time, end_time)  -- 确保同一场地活动时间不重叠
);


CREATE TABLE event_attendees (
    id SERIAL PRIMARY KEY,        -- 记录 ID，自增
    event_id INT NOT NULL,        -- 活动 ID（外键）
    user_id INT NOT NULL,         -- 用户 ID（外键）
    created_at TIMESTAMP DEFAULT now(), -- 记录报名时间
    UNIQUE (event_id, user_id),   -- 确保同一个用户不会重复报名同一个活动
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'confirmed'
);