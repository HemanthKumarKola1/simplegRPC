package repo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"poultry-management.com/internal/auth"
	db "poultry-management.com/internal/db/sqlc"
	"poultry-management.com/pkg/domain"
)

type Repository struct {
    *db.Queries
    db *pgxpool.Pool
}

type UserRepo interface {
    CreateUser(c *gin.Context, req auth.SignupRequest) (error)
    GetUser(c *gin.Context, req auth.LoginRequest) (domain.User, error)
}

func NewRepository(dbConn *pgxpool.Pool) *Repository {
    repo := &Repository{
        Queries: db.New(dbConn),
        db:      dbConn,
    }
    return repo
}

func (r *Repository) CreateUser(c *gin.Context, req auth.SignupRequest) ( error) {
    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    err = r.Queries.CreateUser(c.Request.Context(), db.CreateUserParams{
        Username:     req.Username,
        Passwordhash: hashedPassword,
        Email:        req.Email,
        Usertype:     int32(req.UserType),
    })

    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

func (r *Repository) GetUser(c *gin.Context, req auth.LoginRequest) (domain.User, error) {
    user, err := r.Queries.GetUserByUserName(c.Request.Context(), req.Username)
    if err != nil {
        if err == pgx.ErrNoRows {
            return domain.User{}, fmt.Errorf("user not found")
        }
        return domain.User{}, fmt.Errorf("failed to get user: %w", err)
    }

    err = auth.CheckPasswordHash(req.Password, user.Passwordhash)
    if err != nil {
        return domain.User{}, fmt.Errorf("invalid password")
    }

    return domain.User{
        ID:           user.Userid,
        Username:     user.Username,
        PasswordHash: user.Passwordhash,
        Email:        user.Email,
        UserType:     user.Usertype,
        CreatedAt:    user.CreatedAt,
        UpdatedAt:    user.CreatedAt,
    }, nil
}