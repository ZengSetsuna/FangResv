package api

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	"FangResv/auth"
	db "FangResv/db/sqlc"

	"FangResv/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) PreregisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	allowedDomain := "sjtu.edu.cn"
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@`+allowedDomain+`$`, req.Email)
	if !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email domain"})
		return
	}
	verificationCode := util.GenerateVerificationCode()
	expiry := time.Now().Add(10 * time.Minute)
	err := s.Queries.NewPendingUser(context.Background(), db.NewPendingUserParams{
		Username:  req.Username,
		Email:     req.Email,
		Code:      verificationCode,
		ExpiresAt: pgtype.Timestamp{Time: expiry, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pending user"})
		log.Println(err)
		return
	}
	err = s.Mailer.SendEmail(req.Email, "您的注册验证码", "您的注册验证码是："+verificationCode+"，有效期 10 分钟")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// 注册用户
func (s *Server) RegisterUser(c *gin.Context) {
	var req struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		Email            string `json:"email"`
		VerificationCode string `json:"verification_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 检查验证码
	pendingUser, err := s.Queries.GetPendingUserByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or verification code"})
		return
	}
	if pendingUser.Code != req.VerificationCode || pendingUser.Username != req.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification failed"})
		return
	}
	if pendingUser.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification failed"})
		return
	}

	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 插入用户
	err = s.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		log.Fatal(err)
		return
	}
	err = s.Queries.DeletePendingUserByEmail(context.Background(), req.Email)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// 用户登录
func (s *Server) LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 查询用户
	user, err := s.Queries.GetUserByUsername(context.Background(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
