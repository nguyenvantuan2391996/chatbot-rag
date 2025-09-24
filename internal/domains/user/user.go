package user

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"chatbot-rag/base_common/comoutput"
	"chatbot-rag/base_common/constants"
	"chatbot-rag/base_common/database/entities"
	"chatbot-rag/base_common/milvus"
	"chatbot-rag/base_common/utils"
	"chatbot-rag/internal/domains/repository"
	"chatbot-rag/internal/domains/user/models"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type User struct {
	milvusClient  milvus.IMilvusClientInterface
	openAI        *openai.Client
	embeddingRepo repository.IEmbeddingRepositoryInterface
}

func NewUserService(milvusClient milvus.IMilvusClientInterface, openAI *openai.Client,
	embeddingRepo repository.IEmbeddingRepositoryInterface) *User {
	return &User{
		milvusClient:  milvusClient,
		openAI:        openAI,
		embeddingRepo: embeddingRepo,
	}
}

func (us *User) Index(ctx context.Context, input *models.IndexInput) (*comoutput.BaseOutput, error) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginTask, "Index", input))

	resp, err := us.openAI.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: input.Text,
			Model: openai.SmallEmbedding3,
		},
	)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "CreateEmbeddings", err)
		return nil, err
	}

	for _, emb := range resp.Data {
		record := &entities.Embedding{
			Vector: base64.StdEncoding.EncodeToString([]byte(utils.Float32SliceToBase64(emb.Embedding))),
		}

		err = us.embeddingRepo.Create(ctx, record)
		if err != nil {
			logrus.Errorf(constants.FormatCreateEntityErr, "embedding", err)
			continue
		}

		err = us.milvusClient.Insert(emb.Embedding, viper.GetString("MILVUS_COLLECTION"), "", record.ID)
		if err != nil {
			logrus.Errorf(constants.FormatTaskErr, "Insert", err)
		}
	}

	return &comoutput.BaseOutput{
		Status: http.StatusOK,
	}, nil
}

func (us *User) Chat(ctx context.Context, input *models.ChatInput) (*comoutput.BaseOutput, error) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginTask, "Chat", input))

	return &comoutput.BaseOutput{
		Status: http.StatusOK,
	}, nil
}
