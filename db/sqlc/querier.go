// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CheckVenueAvailability(ctx context.Context, arg CheckVenueAvailabilityParams) (int64, error)
	CountEventParticipants(ctx context.Context, id int32) (pgtype.Int4, error)
	CountUpcomingEvents(ctx context.Context) (int32, error)
	CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error)
	CreateUser(ctx context.Context, arg CreateUserParams) error
	CreateVenue(ctx context.Context, arg CreateVenueParams) (Venue, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetVenueByID(ctx context.Context, id int32) (Venue, error)
	GetVenueCapacity(ctx context.Context, id int32) (int32, error)
	JoinEvent(ctx context.Context, arg JoinEventParams) (JoinEventRow, error)
	ListUpcomingEvents(ctx context.Context, arg ListUpcomingEventsParams) ([]ListUpcomingEventsRow, error)
}

var _ Querier = (*Queries)(nil)
