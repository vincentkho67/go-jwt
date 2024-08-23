package usecase

import (
	"github.com/vincentkho67/jwt/internal/domain"
)

type NoteRepository interface {
	CreateNote(note *domain.Note) error
	GetNotesByUserID(userID int) ([]*domain.Note, error)
}

type NoteUseCase struct {
	repo NoteRepository
}

func NewNoteUseCase(repo NoteRepository) *NoteUseCase {
	return &NoteUseCase{repo: repo}
}

func (uc *NoteUseCase) CreateNote(userID int, content string) error {
	note := &domain.Note{
		UserID:  userID,
		Content: content,
	}
	return uc.repo.CreateNote(note)
}

func (uc *NoteUseCase) GetNotesByUserID(userID int) ([]*domain.Note, error) {
	return uc.repo.GetNotesByUserID(userID)
}
