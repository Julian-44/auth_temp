package sqlite

import (
	"auth-server/internal/domain"
	"auth-server/internal/repository"
	"context"
	"database/sql"
)

type UserRepositorySQLite struct {
	db *sql.DB
}

func NewUserRepositorySQLite(db *sql.DB) repository.UserRepository {
	return &UserRepositorySQLite{db: db}
}

func (r *UserRepositorySQLite) Create(ctx context.Context, user *domain.User) (int64, error) {
	query := `INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepositorySQLite) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositorySQLite) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
