package repository

import (
	"database/sql"

	"github.com/vincentkho67/jwt/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Email, user.Password).Scan(&user.ID)
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateNote(note *domain.Note) error {
	query := `INSERT INTO notes (user_id, content) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, note.UserID, note.Content).Scan(&note.ID)
}

func (r *Repository) GetNotesByUserID(userID int) ([]*domain.Note, error) {
	query := `SELECT id, user_id, content FROM notes WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*domain.Note
	for rows.Next() {
		note := &domain.Note{}
		if err := rows.Scan(&note.ID, &note.UserID, &note.Content); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
