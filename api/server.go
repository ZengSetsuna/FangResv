package api

import (
	"context"
	"log"
	"time"

	db "FangResv/db/sqlc"
	"FangResv/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
	Mailer  *util.Mailer
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 预检缓存 12 小时
	}))
	router.POST("/register", s.RegisterUser)
	router.POST("/login", s.LoginUser)

	auth := router.Group("/")
	auth.Use(AuthMiddleware())
	auth.POST("/venues", s.CreateVenue)
	auth.GET("/events", s.GetUpcomingEvents)
	auth.POST("/events", s.CreateEvent)
	auth.GET("/events/:id", s.GetEventDetails)
	auth.POST("/events/:id/join", s.JoinEvent)

	return router
}
