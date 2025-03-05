-- name: CreateVenue :one
INSERT INTO venues (name, address, max_capacity) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVenueByID :one
SELECT * FROM venues WHERE id = $1;

-- name: CreateEvent :one
INSERT INTO events (creator_id, venue_id, name, description, start_time, end_time, location, max_participants)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CheckVenueAvailability :one
SELECT COUNT(*) 
FROM events
WHERE venue_id = $1
AND EXISTS (
    SELECT 1
    FROM (SELECT $2::TIMESTAMP AS start_time, $3::TIMESTAMP AS end_time) AS t
    WHERE (t.start_time, t.end_time) OVERLAPS (events.start_time, events.end_time)
);

-- name: CountEventParticipants :one
SELECT current_participants FROM events WHERE id = $1;

-- name: GetVenueCapacity :one
SELECT max_capacity FROM venues WHERE id = $1;

-- name: JoinEvent :one
WITH updated AS (
    UPDATE events 
    SET current_participants = current_participants + 1
    WHERE id = $1 AND current_participants < max_participants
    RETURNING id, current_participants
)
INSERT INTO event_attendees (event_id, user_id)
SELECT $1, $2
FROM updated
RETURNING event_id, user_id;

-- name: CreateUser :exec
INSERT INTO users (username, password) VALUES ($1, $2);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CountUpcomingEvents :one
SELECT CAST(COUNT(*) AS INT4) 
FROM events
WHERE start_time > NOW();

-- name: ListUpcomingEvents :many
SELECT
    id,
    name,
    start_time,
    end_time,
    location,
    max_participants,
    current_participants
FROM events
WHERE start_time > NOW()
ORDER BY start_time ASC
LIMIT $1 OFFSET $2;
