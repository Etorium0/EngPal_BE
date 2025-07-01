package repository

import "context"

type User struct {
    ID           int
    Username     string
    PasswordHash string
}

type UserRepository interface {
    GetUserByUsername(ctx context.Context, username string) (*User, error)
    CreateUser(ctx context.Context, username string, passwordHash string) error
    UpdatePassword(ctx context.Context, username string, newPasswordHash string) error
}