package api

import (
	"context"
	"log"

	db "FangResv/db/sqlc"
	"FangResv/util"

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

	router.POST("/register", s.RegisterUser)
	router.POST("/login", s.LoginUser)

	auth := router.Group("/")
	auth.Use(AuthMiddleware())
	auth.POST("/venues", s.CreateVenue)
	auth.POST("/events", s.CreateEvent)
	auth.POST("/events/:id/join", s.JoinEvent)

	return router
}
