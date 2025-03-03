// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkVenueAvailability = `-- name: CheckVenueAvailability :one
SELECT COUNT(*) 
FROM events
WHERE venue_id = $1
AND EXISTS (
    SELECT 1
    FROM (SELECT $2::TIMESTAMP AS start_time, $3::TIMESTAMP AS end_time) AS t
    WHERE (t.start_time, t.end_time) OVERLAPS (events.start_time, events.end_time)
)
`

type CheckVenueAvailabilityParams struct {
	VenueID pgtype.Int4      `json:"venue_id"`
	Column2 pgtype.Timestamp `json:"column_2"`
	Column3 pgtype.Timestamp `json:"column_3"`
}

func (q *Queries) CheckVenueAvailability(ctx context.Context, arg CheckVenueAvailabilityParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkVenueAvailability, arg.VenueID, arg.Column2, arg.Column3)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countEventParticipants = `-- name: CountEventParticipants :one
SELECT current_participants FROM events WHERE id = $1
`

func (q *Queries) CountEventParticipants(ctx context.Context, id int32) (pgtype.Int4, error) {
	row := q.db.QueryRow(ctx, countEventParticipants, id)
	var current_participants pgtype.Int4
	err := row.Scan(&current_participants)
	return current_participants, err
}

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (creator_id, venue_id, name, description, start_time, end_time, location, max_participants)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, creator_id, venue_id, name, description, start_time, end_time, location, max_participants, created_at, current_participants
`

type CreateEventParams struct {
	CreatorID       pgtype.Int4      `json:"creator_id"`
	VenueID         pgtype.Int4      `json:"venue_id"`
	Name            string           `json:"name"`
	Description     pgtype.Text      `json:"description"`
	StartTime       pgtype.Timestamp `json:"start_time"`
	EndTime         pgtype.Timestamp `json:"end_time"`
	Location        string           `json:"location"`
	MaxParticipants int32            `json:"max_participants"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createEvent,
		arg.CreatorID,
		arg.VenueID,
		arg.Name,
		arg.Description,
		arg.StartTime,
		arg.EndTime,
		arg.Location,
		arg.MaxParticipants,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.VenueID,
		&i.Name,
		&i.Description,
		&i.StartTime,
		&i.EndTime,
		&i.Location,
		&i.MaxParticipants,
		&i.CreatedAt,
		&i.CurrentParticipants,
	)
	return i, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (username, password) VALUES ($1, $2)
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser, arg.Username, arg.Password)
	return err
}

const createVenue = `-- name: CreateVenue :one
INSERT INTO venues (name, address, max_capacity) VALUES ($1, $2, $3) RETURNING id, name, address, max_capacity
`

type CreateVenueParams struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	MaxCapacity int32  `json:"max_capacity"`
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (Venue, error) {
	row := q.db.QueryRow(ctx, createVenue, arg.Name, arg.Address, arg.MaxCapacity)
	var i Venue
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.MaxCapacity,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password, created_at FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getVenueByID = `-- name: GetVenueByID :one
SELECT id, name, address, max_capacity FROM venues WHERE id = $1
`

func (q *Queries) GetVenueByID(ctx context.Context, id int32) (Venue, error) {
	row := q.db.QueryRow(ctx, getVenueByID, id)
	var i Venue
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.MaxCapacity,
	)
	return i, err
}

const getVenueCapacity = `-- name: GetVenueCapacity :one
SELECT max_capacity FROM venues WHERE id = $1
`

func (q *Queries) GetVenueCapacity(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, getVenueCapacity, id)
	var max_capacity int32
	err := row.Scan(&max_capacity)
	return max_capacity, err
}

const joinEvent = `-- name: JoinEvent :one
WITH updated AS (
    UPDATE events 
    SET current_participants = current_participants + 1
    WHERE id = $1 AND current_participants < max_participants
    RETURNING id, current_participants
)
INSERT INTO event_attendees (event_id, user_id)
SELECT $1, $2
FROM updated
RETURNING event_id, user_id
`

type JoinEventParams struct {
	EventID int32 `json:"event_id"`
	UserID  int32 `json:"user_id"`
}

type JoinEventRow struct {
	EventID int32 `json:"event_id"`
	UserID  int32 `json:"user_id"`
}

func (q *Queries) JoinEvent(ctx context.Context, arg JoinEventParams) (JoinEventRow, error) {
	row := q.db.QueryRow(ctx, joinEvent, arg.EventID, arg.UserID)
	var i JoinEventRow
	err := row.Scan(&i.EventID, &i.UserID)
	return i, err
}
