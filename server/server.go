package server

import (
	"context"
	"log"

	"FangResv/api"
	db "FangResv/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

// NewServer 初始化服务器
func NewServer(dbURL string) *Server {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	queries := db.New(pool)

	return &Server{
		Pool:    pool,
		Queries: queries,
	}
}

// SetupRouter 设置 Gin 路由
func (s *Server) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/register", func(c *gin.Context) { api.RegisterUser(c, s.Queries) })
	router.POST("/login", func(c *gin.Context) { api.LoginUser(c, s.Queries) })

	auth := router.Group("/")
	auth.Use(api.AuthMiddleware())
	auth.POST("/venues", func(c *gin.Context) { api.CreateVenue(c, s.Queries) })
	auth.POST("/events", func(c *gin.Context) { api.CreateEvent(c, s.Queries) })
	auth.POST("/events/:id/join", func(c *gin.Context) { api.JoinEvent(c, s.Queries) })

	return router
}
