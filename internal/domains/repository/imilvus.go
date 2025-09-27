package repository

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

//go:generate mockgen -package=repository -destination=imilvus_mock.go -source=imilvus.go

type IMilvusRepositoryInterface interface {
	Close() error
	CreateCollection(schema *entity.Schema, shardsNum int32, opts ...client.CreateCollectionOption) error
	HasCollection(collectionName string) (bool, error)
	CreateCollectionCustom(schema *entity.Schema, shardNums int32, async bool,
		indexColumns []string) error
	LoadCollection(collectionName string, async bool, opts ...client.LoadCollectionOption) error
	Flush(collectionName string, async bool, opts ...client.FlushOption) error
	DropCollection(collectionName string, opts ...client.DropCollectionOption) error
	CreatePartition(collectionName, partitionName string, opts ...client.CreatePartitionOption) error
	HasPartition(collectionName, partitionName string) (bool, error)
	CreateIndex(collectionName string, fieldName string, index entity.Index, async bool, opts ...client.IndexOption) error
	DropIndex(collectionName, fieldName string, opts ...client.IndexOption) error
	Insert(collectionName, partitionName string, columns ...entity.Column) (entity.Column, error)
	Upsert(collectionName, partitionName string, columns ...entity.Column) (entity.Column, error)
	FindTopK(collectionName string, partitions []string, expr string, outputFields []string,
		vectors []entity.Vector, vectorField string, metricType entity.MetricType,
		topK int, sp entity.SearchParam, opts ...client.SearchQueryOptionFunc) ([]client.SearchResult, error)
	DeleteEntityByID(collectionName, partitionName string, id int64) error
	DeleteEntitiesByIDs(collectionName, partitionName string, ids []int64) error
	InitCollection(collectionName string) error
}
