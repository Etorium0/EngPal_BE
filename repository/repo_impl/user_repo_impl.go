package repo_impl

import (
    "context"
    "database/sql"
    "EngPal/repository"
)

type userRepoImpl struct {
    db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepository {
    return &userRepoImpl{db: db}
}

func (r *userRepoImpl) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
    user := &repository.User{}
    err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username=$1", username).
        Scan(&user.ID, &user.Username, &user.PasswordHash)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepoImpl) CreateUser(ctx context.Context, username string, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}

func (r *userRepoImpl) UpdatePassword(ctx context.Context, username string, newPasswordHash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET password_hash=$1 WHERE username=$2", newPasswordHash, username)
	return err
}