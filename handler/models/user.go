package models

import (
	"chatbot-rag/internal/domains/user/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IndexRequest struct {
	TextIndex string `form:"text_index"`
}

type ChatRequest struct {
	Convention string `form:"convention"`
}

func (r *IndexRequest) ToIndexInput() *models.IndexInput {
	out := &models.IndexInput{}
	if r == nil {
		return out
	}

	return out
}

func (r *ChatRequest) ToChatInput() *models.ChatInput {
	out := &models.ChatInput{}
	if r == nil {
		return out
	}

	return out
}

func (r *IndexRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TextIndex, validation.Required),
	)
}

func (r *ChatRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Convention, validation.Required),
	)
}
