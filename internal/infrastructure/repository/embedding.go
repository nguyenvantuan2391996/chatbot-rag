package repository

import (
	"context"

	"chatbot-rag/base_common/database/entities"
	"gorm.io/gorm"
)

type EmbeddingRepository struct {
	db *gorm.DB
}

func NewEmbeddingRepository(db *gorm.DB) *EmbeddingRepository {
	return &EmbeddingRepository{db: db}
}

func (er *EmbeddingRepository) Create(ctx context.Context, record *entities.Embedding) error {
	return er.db.WithContext(ctx).Create(&record).Error
}
