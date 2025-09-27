package repository

import (
	"context"
	"fmt"
	"time"

	"chatbot-rag/base_common/constants"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	// Collection
	idColumnName        = "id"
	productIDColumnName = "product_id"
	vectorColumnName    = "vector"
	dimension           = 512
	consistencyLevel    = entity.DefaultConsistencyLevel

	// Index
	metricType = entity.IP
	nList      = 2048
	nProbe     = 64
)

type MilvusRepository struct {
	client  client.Client
	timeout time.Duration
}

func NewMilvusRepository(client client.Client, timeout time.Duration) *MilvusRepository {
	return &MilvusRepository{
		client:  client,
		timeout: timeout,
	}
}

func (m *MilvusRepository) Close() error {
	err := m.client.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) CreateCollection(schema *entity.Schema, shardsNum int32, opts ...client.CreateCollectionOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.CreateCollection(ctx, schema, shardsNum, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) HasCollection(collectionName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	existed, err := m.client.HasCollection(ctx, collectionName)
	if err != nil {
		return false, err
	}

	return existed, nil
}

func (m *MilvusRepository) CreateCollectionCustom(schema *entity.Schema, shardNums int32, async bool,
	indexColumns []string) error {
	metricTypeOpt := client.WithMetricsType(metricType)
	consistencyLevelOpt := client.WithConsistencyLevel(consistencyLevel)

	err := m.CreateCollection(schema, shardNums, metricTypeOpt, consistencyLevelOpt)
	if err != nil {
		return err
	}

	index, errIndex := entity.NewIndexIvfFlat(metricType, nList)
	if errIndex != nil {
		return errIndex
	}

	for _, col := range indexColumns {
		errIndex = m.CreateIndex(schema.CollectionName, col, index, async)
		if errIndex != nil {
			return errIndex
		}
	}

	err = m.LoadCollection(schema.CollectionName, async)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) LoadCollection(collectionName string, async bool, opts ...client.LoadCollectionOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.LoadCollection(ctx, collectionName, async, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) Flush(collectionName string, async bool, opts ...client.FlushOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.Flush(ctx, collectionName, async, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) DropCollection(collectionName string, opts ...client.DropCollectionOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.DropCollection(ctx, collectionName, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) CreatePartition(collectionName, partitionName string, opts ...client.CreatePartitionOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if partitionName == "" {
		return fmt.Errorf("partitionName must not be empty")
	}

	err := m.client.CreatePartition(ctx, collectionName, partitionName, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) HasPartition(collectionName, partitionName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if partitionName == "" {
		return false, fmt.Errorf("partitionName must not be empty")
	}

	existed, err := m.client.HasPartition(ctx, collectionName, partitionName)
	if err != nil {
		return false, err
	}

	return existed, nil
}

func (m *MilvusRepository) CreateIndex(collectionName string, fieldName string, index entity.Index, async bool, opts ...client.IndexOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.CreateIndex(ctx, collectionName, fieldName, index, async, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) DropIndex(collectionName, fieldName string, opts ...client.IndexOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.client.DropIndex(ctx, collectionName, fieldName, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) Insert(collectionName, partitionName string, columns ...entity.Column) (entity.Column, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	// insert data
	column, err := m.client.Insert(ctx, collectionName, partitionName, columns...)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (m *MilvusRepository) Upsert(collectionName, partitionName string, columns ...entity.Column) (entity.Column, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	column, err := m.client.Upsert(ctx, collectionName, partitionName, columns...)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (m *MilvusRepository) FindTopK(collectionName string, partitions []string, expr string, outputFields []string, vectors []entity.Vector, vectorField string, metricType entity.MetricType, topK int, sp entity.SearchParam, opts ...client.SearchQueryOptionFunc) ([]client.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	results, err := m.client.Search(
		ctx,
		collectionName,
		partitions,
		expr,
		outputFields,
		vectors,
		vectorField,
		metricType,
		topK,
		sp,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MilvusRepository) DeleteEntityByID(collectionName, partitionName string, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	idColumn := entity.NewColumnInt64(idColumnName, []int64{id})

	err := m.client.DeleteByPks(ctx, collectionName, partitionName, idColumn)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) DeleteEntitiesByIDs(collectionName, partitionName string, ids []int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	idColumn := entity.NewColumnInt64(idColumnName, ids)

	err := m.client.DeleteByPks(ctx, collectionName, partitionName, idColumn)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilvusRepository) InitCollection(collectionName string) error {
	existed, err := m.HasCollection(collectionName)
	if err != nil {
		return err
	}

	if !existed {
		schema := &entity.Schema{
			CollectionName: collectionName,
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     false,
				},
				{
					Name:       "is_visible",
					DataType:   entity.FieldTypeBool,
					PrimaryKey: false,
					AutoID:     false,
				},
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": fmt.Sprintf("%d", constants.SmallEmbedding3Dims),
					},
				},
			},
		}

		err = m.CreateCollectionCustom(schema, entity.DefaultShardNumber,
			true, []string{"vector"})
		if err != nil {
			return err
		}
	}

	return nil
}
