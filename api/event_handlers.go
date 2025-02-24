package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	db "FangResv/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// 创建活动
func CreateEvent(c *gin.Context, queries *db.Queries) {
	var req struct {
		VenueID         int32     `json:"venue_id"`
		HostID          int32     `json:"host_id"`
		Name            string    `json:"name"`
		StartTime       time.Time `json:"start_time"`
		EndTime         time.Time `json:"end_time"`
		MaxParticipants int32     `json:"max_participants"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	// userIDInt, ok := userID.(int32)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	// req.HostID = userIDInt
	// 检查场地是否可用
	startTime := pgtype.Timestamp{Time: req.StartTime, Valid: true}
	endTime := pgtype.Timestamp{Time: req.EndTime, Valid: true}
	availableCount, err := queries.CheckVenueAvailability(context.Background(), db.CheckVenueAvailabilityParams{
		VenueID: pgtype.Int4{Int32: int32(req.VenueID), Valid: true},
		Column2: startTime,
		Column3: endTime,
	})
	if err != nil || availableCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Venue is not available at this time"})
		return
	}
	// 创建活动
	res, err := queries.CreateEvent(context.Background(), db.CreateEventParams{
		VenueID:         pgtype.Int4{Int32: int32(req.VenueID), Valid: true},
		CreatorID:       pgtype.Int4{Int32: req.HostID, Valid: true},
		Name:            req.Name,
		StartTime:       startTime,
		EndTime:         endTime,
		MaxParticipants: req.MaxParticipants,
	})
	log.Println(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create event"})
		log.Println(err)
		return
	}
	_, err = queries.JoinEvent(context.Background(), db.JoinEventParams{
		EventID: res.ID,
		UserID:  req.HostID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully"})
}

func JoinEvent(c *gin.Context, queries *db.Queries) {
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	// userIDInt, ok := userID.(int32)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }
	var req struct {
		EventID int32 `json:"event_id"`
		UserID  int32 `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用 SQL 查询，直接检查活动是否已满，并加入活动
	_, err := queries.JoinEvent(context.Background(), db.JoinEventParams{
		EventID: req.EventID,
		UserID:  req.UserID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "User already joined the event"})
			return
		}
		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join event"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Joined event successfully"})
}
