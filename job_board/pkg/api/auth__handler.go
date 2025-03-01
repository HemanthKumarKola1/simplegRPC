package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"poultry-management.com/internal/auth"
	"poultry-management.com/pkg/repo"
)

type AuthHandler struct {
    config Config
    user   repo.UserRepo
    Router *gin.Engine
}

type Config struct {
    JWTSecret string
}

func NewAuthHandler(config Config, repo repo.UserRepo, e *gin.Engine) *AuthHandler {
    a := &AuthHandler{
        config: config,
        user:   repo,
        Router: e,
    }
    setupRoutes(a)
    return a
}

func setupRoutes(a *AuthHandler) {
    a.Router.POST("/login", a.loginHandler)
    a.Router.POST("/signup", a.signupHandler)

    protected := a.Router.Group("/verified")
    protected.Use(authMiddleware([]byte(a.config.JWTSecret)))
    protected.GET("/ping", a.pingHandler)
}

func (a *AuthHandler) signupHandler(c *gin.Context) {
    var req auth.SignupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := a.user.CreateUser(c, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"user-signed-up": true})
}

func (a *AuthHandler) loginHandler(c *gin.Context) {
    var req auth.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := a.user.GetUser(c, req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := auth.GenerateJWT(user,[]byte(a.config.JWTSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthHandler) pingHandler(c *gin.Context) {
    userID := c.GetInt("user_id")
    username := c.GetString("username")
    userType := c.GetInt("user_type") // Get UserType instead of Role
    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("pong user id: %d, username: %s, user type: %d", userID, username, userType),
    })
}

func authMiddleware(jwtKey []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
            return
        }

        claims, err := auth.ValidateJWT(tokenString, jwtKey)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("user_type", claims.UserType)
        c.Next()
    }
}