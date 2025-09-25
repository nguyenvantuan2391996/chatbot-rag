package repository

import (
	"context"

	"chatbot-rag/base_common/database/entities"
)

type IEmbeddingRepositoryInterface interface {
	Create(ctx context.Context, record *entities.Embedding) error
	GetListFacts(ctx context.Context, ids []int64) ([]string, error)
}
