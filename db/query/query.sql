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
INSERT INTO users (username, password, email) VALUES ($1, $2, $3);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

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

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: ListEventAttendees :many
SELECT user_id FROM event_attendees WHERE event_id = $1;

-- name: LeaveEvent :one
WITH updated AS (
    UPDATE events 
    SET current_participants = current_participants - 1
    WHERE id = $1 AND current_participants > 0
    RETURNING id, current_participants
)
DELETE FROM event_attendees
WHERE event_id = $1 AND user_id = $2
RETURNING event_id, user_id;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: ListEventsByCreatorID :many
SELECT
    id,
    name,
    start_time,
    end_time,
    location,
    max_participants,
    current_participants
FROM events
WHERE creator_id = $1
ORDER BY start_time ASC
LIMIT $2 OFFSET $3;

-- name: ListEventsByAttendeeID :many
SELECT
    events.id,
    events.name,
    events.start_time,
    events.end_time,
    events.location,
    events.max_participants,
    events.current_participants
FROM events
JOIN event_attendees ON events.id = event_attendees.event_id
WHERE event_attendees.user_id = $1
ORDER BY events.start_time ASC
LIMIT $2 OFFSET $3;

-- name: GetEventDetails :one
WITH event_info AS (
    SELECT 
        e.id AS event_id,
        e.name AS event_name,
        e.description,
        e.start_time,
        e.end_time,
        v.name AS venue_name,
        e.max_participants,
        e.current_participants,
        u.username AS organizer
    FROM events e
    JOIN venues v ON e.venue_id = v.id
    JOIN users u ON e.creator_id = u.id
    WHERE e.id = $1
),
attendees AS (
    SELECT 
        ea.event_id,
        json_agg(u.username) AS participant_usernames
    FROM event_attendees ea
    JOIN users u ON ea.user_id = u.id
    WHERE ea.event_id = $1
    GROUP BY ea.event_id
)
SELECT 
    ei.event_id,
    ei.event_name,
    ei.description,
    ei.start_time,
    ei.end_time,
    ei.venue_name AS location,
    ei.max_participants,
    ei.current_participants,
    ei.organizer,
    COALESCE(a.participant_usernames, '[]'::json) AS participants,
    (ei.current_participants < ei.max_participants) AS can_join
FROM event_info ei
LEFT JOIN attendees a ON ei.event_id = a.event_id;

-- name: NewPendingUser :exec
INSERT INTO pending_users (email, username, code, expires_at) VALUES ($1, $2, $3, $4);

-- name: GetPendingUserByEmail :one
SELECT * FROM pending_users WHERE email = $1;

-- name: DeletePendingUserByEmail :exec
DELETE FROM pending_users WHERE email = $1;