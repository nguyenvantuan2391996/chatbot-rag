package handler

import (
	"chatbot-rag/internal/domains/user"
)

type Handler struct {
	userService *user.User
}

func NewHandler(userService *user.User) *Handler {
	return &Handler{
		userService: userService,
	}
}
