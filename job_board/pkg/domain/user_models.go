package domain

import "github.com/jackc/pgx/v5/pgtype"

type User struct { // Includes all types of users
    ID           int32              `json:"id"`
    Username     string             `json:"username"`
    PasswordHash string             `json:"password_hash"` // Added PasswordHash
    Email        string             `json:"email"`           // Added Email
    UserType     int32              `json:"user_type"`       // Replaced Role with UserType
    CreatedAt    pgtype.Timestamptz `json:"created_at"`
    UpdatedAt    pgtype.Timestamptz `json:"updated_at"` // Added UpdatedAt
}