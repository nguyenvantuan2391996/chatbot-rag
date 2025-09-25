package milvus

import "github.com/milvus-io/milvus-sdk-go/milvus"

type IMilvusClientInterface interface {
	Insert(vector []float32, collectionName, partitionTag string, id int64) error
	Search(vector []float32, collectionName string, partitionTag []string, topK int64) (milvus.TopkQueryResult, error)
	Delete(collectionName, partitionTag string, id int64) error
	DropCollection(collectionName string) error
	HasCollection(collectionName string) (bool, error)
	CreateCollection(collectionName string, dimension, indexSize int64, metric milvus.MetricType) error
	CreateIndex(collectionName string, nList int64, indexType milvus.IndexType) error
	DropIndex(collectionName string) error
	CreatePartition(collectionName, partitionTag string) error
	DropPartition(collectionName, partitionTag string) error
}
