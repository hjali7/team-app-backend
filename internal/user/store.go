package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	GoogleID  *string   `json:"google_id,omitempty"`
	GithubID  *string   `json:"github_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(email, hashedPassword string) (*User, error) {
	query := `
		INSERT INTO users (email, hashed_password) 
		VALUES ($1, $2) 
		RETURNING id, email, google_id, github_id, created_at
	`
	
	u := new(User)

	err := s.db.QueryRow(context.Background(), query, email, hashedPassword).Scan(
		&u.ID,
		&u.Email,
		&u.GoogleID,
		&u.GithubID,
		&u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}