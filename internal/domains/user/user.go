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

func (us *User) Index(ctx context.Context, input *models.FactInput) (*comoutput.BaseOutput, error) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginTask, "Index", input))

	resp, err := us.openAI.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: input.Facts,
			Model: openai.SmallEmbedding3,
		},
	)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "CreateEmbeddings", err)
		return nil, err
	}

	for idx, emb := range resp.Data {
		record := &entities.Embedding{
			Vector: base64.StdEncoding.EncodeToString([]byte(utils.Float32SliceToBase64(emb.Embedding))),
			Fact:   input.Facts[idx],
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

	resp, err := us.openAI.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{input.Convention},
			Model: openai.SmallEmbedding3,
		},
	)
	if err != nil || len(resp.Data) == 0 {
		logrus.Errorf(constants.FormatTaskErr, "CreateEmbeddings", err)
		return nil, err
	}

	search, err := us.milvusClient.Search(resp.Data[0].Embedding, viper.GetString("MILVUS_COLLECTION"), nil, 5)
	if err != nil || len(search.QueryResultList) == 0 {
		logrus.Errorf(constants.FormatTaskErr, "Search", err)
		return nil, err
	}

	ids := make([]int64, 0)
	for _, id := range search.QueryResultList[0].Ids {
		ids = append(ids, id)
	}

	facts, err := us.embeddingRepo.GetListFacts(ctx, ids)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "GetListFacts", err)
		return nil, err
	}

	factsPrompt := "Thông tin:\n"
	for i, f := range facts {
		factsPrompt += fmt.Sprintf("%d. %s\n", i+1, f)
	}

	prompt := fmt.Sprintf(`
Bạn là một trợ lý phim chiếu rạp.
Hãy trả lời câu hỏi dựa trên những thông tin được cung cấp dưới đây.
Nếu thông tin không đủ, hãy nói "Tôi không biết dựa trên những thông tin đó."
%s
Câu hỏi: %s`, factsPrompt, input.Convention)

	completion, err := us.openAI.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "CreateChatCompletion", err)
		return nil, err
	}

	return &comoutput.BaseOutput{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"answer": completion.Choices[0].Message.Content,
		},
	}, nil
}
