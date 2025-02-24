package api

import (
	"context"
	"net/http"

	db "FangResv/db/sqlc"

	"github.com/gin-gonic/gin"
)

// 创建活动场地
func CreateVenue(c *gin.Context, queries *db.Queries) {
	var req struct {
		Name        string `json:"name"`
		Capacity    int32  `json:"capacity"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 插入场地信息
	_, err := queries.CreateVenue(context.Background(), db.CreateVenueParams{
		Name:        req.Name,
		Address:     req.Description,
		MaxCapacity: req.Capacity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create venue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venue created successfully"})
}
