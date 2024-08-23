package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincentkho67/jwt/internal/usecase"
)

type Handler struct {
	userUseCase *usecase.UserUseCase
	noteUseCase *usecase.NoteUseCase
}

func NewHandler(userUC *usecase.UserUseCase, noteUC *usecase.NoteUseCase) *Handler {
	return &Handler{
		userUseCase: userUC,
		noteUseCase: noteUC,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUseCase.Register(input.Email, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userUseCase.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) CreateNote(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.noteUseCase.CreateNote(userID.(int), input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully"})
}

func (h *Handler) GetNotes(c *gin.Context) {
	userID, _ := c.Get("user_id")

	notes, err := h.noteUseCase.GetNotesByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}
