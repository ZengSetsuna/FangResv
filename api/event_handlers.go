package api

import (
	db "FangResv/db/sqlc"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// 创建活动
func (s *Server) CreateEvent(c *gin.Context) {
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt, ok := userID.(int32)
	if !ok || userIDInt != req.HostID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// 检查场地是否可用
	startTime := pgtype.Timestamp{Time: req.StartTime, Valid: true}
	endTime := pgtype.Timestamp{Time: req.EndTime, Valid: true}
	availableCount, err := s.Queries.CheckVenueAvailability(context.Background(), db.CheckVenueAvailabilityParams{
		VenueID: pgtype.Int4{Int32: int32(req.VenueID), Valid: true},
		Column2: startTime,
		Column3: endTime,
	})
	if err != nil || availableCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Venue is not available at this time"})
		return
	}
	// 创建活动
	res, err := s.Queries.CreateEvent(context.Background(), db.CreateEventParams{
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
	_, err = s.Queries.JoinEvent(context.Background(), db.JoinEventParams{
		EventID: res.ID,
		UserID:  req.HostID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully"})
}

func (s *Server) JoinEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt, ok := userID.(int32)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		EventID int32 `json:"event_id"`
		// UserID  int32 `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用 SQL 查询，直接检查活动是否已满，并加入活动
	_, err := s.Queries.JoinEvent(context.Background(), db.JoinEventParams{
		EventID: req.EventID,
		UserID:  userIDInt,
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
	user, err := s.Queries.GetUserByID(context.Background(), userIDInt)
	event, err := s.Queries.GetEventByID(context.Background(), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user or event"})
		log.Println(err)
		return
	}
	eventString := event.StartTime.Time.GoString() + "的" + event.Name
	s.Mailer.SendEmail(user.Email, "成功加入活动："+eventString, "您已成功加入"+eventString+"，请准时参加！")

}

// GetUpcomingEvents 处理获取未来活动的 API
func (s *Server) GetUpcomingEvents(c *gin.Context) {
	// 从 URL 查询参数获取分页信息，提供默认值
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	// 计算 OFFSET 和 LIMIT
	offset := (page - 1) * pageSize
	limit := pageSize

	// 获取未来活动的总数
	count, err := s.Queries.CountUpcomingEvents(c)
	if err != nil {
		log.Println("Error getting total count of upcoming events:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get event count"})
		return
	}

	// 获取分页数据
	events, err := s.Queries.ListUpcomingEvents(c, db.ListUpcomingEventsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		log.Println("Error getting upcoming events:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get events"})
		return
	}

	// 计算总页数
	totalPages := (count + int32(pageSize) - 1) / int32(pageSize)

	// 返回 JSON 响应
	c.JSON(http.StatusOK, gin.H{
		"total_count": count,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"events":      events,
	})
}

func (s *Server) CancelEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt, ok := userID.(int32)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		EventID int32 `json:"event_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	event, err := s.Queries.GetEventByID(context.Background(), req.EventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.CreatorID.Int32 != userIDInt {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	attendees, err := s.Queries.ListEventAttendees(context.Background(), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendees"})
		log.Println(err)
		return
	}

	for _, attendeeID := range attendees {
		user, err := s.Queries.GetUserByID(context.Background(), attendeeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			log.Println(err)
			return
		}
		eventString := event.StartTime.Time.GoString() + "的" + event.Name
		s.Mailer.SendEmail(user.Email, "活动取消："+eventString, "很抱歉，"+eventString+"已取消。")
	}

	// 取消活动
	err = s.Queries.DeleteEvent(context.Background(), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel event"})
		log.Println(err)
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{"message": "Event cancelled successfully"})
}
