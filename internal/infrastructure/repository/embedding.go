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

func (er *EmbeddingRepository) GetListFacts(ctx context.Context, ids []int64) ([]string, error) {
	var facts []string

	err := er.db.WithContext(ctx).Model(&entities.Embedding{}).Where("id IN ?", ids).Pluck("fact", &facts).Error
	if err != nil {
		return nil, err
	}

	return facts, nil
}
