package models

import (
	"chatbot-rag/internal/domains/user/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FactRequest struct {
	Facts []string `json:"facts"`
}

type ChatRequest struct {
	Convention string `form:"convention"`
}

func (r *FactRequest) ToFactInput() *models.FactInput {
	out := &models.FactInput{}
	if r == nil {
		return out
	}

	out.Facts = r.Facts

	return out
}

func (r *ChatRequest) ToChatInput() *models.ChatInput {
	out := &models.ChatInput{}
	if r == nil {
		return out
	}

	out.Convention = r.Convention

	return out
}

func (r *FactRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Facts, validation.Required),
	)
}

func (r *ChatRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Convention, validation.Required),
	)
}
