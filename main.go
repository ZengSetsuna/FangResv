package main

import (
	"FangResv/util"
	"context"
	"log"

	"FangResv/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 读取数据库 URL
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	databaseURL := config.DBSource
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 连接到 PostgreSQL 数据库（pgx 连接池）
	dbPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 创建 Gin 服务器
	s := api.NewServer(databaseURL)

	// 设置路由
	router := s.SetupRouter()

	s.Mailer = &util.Mailer{
		SmtpHost: config.SMTPHost,
		SmtpPort: config.SMTPPort,
		SmtpUser: config.SMTPUsername,
		SmtpPass: config.SMTPPassword,
	}
	// 启动服务器
	log.Println("Server is running on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
