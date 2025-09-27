package user

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"chatbot-rag/base_common/comoutput"
	"chatbot-rag/base_common/constants"
	"chatbot-rag/base_common/database/entities"
	"chatbot-rag/base_common/utils"
	"chatbot-rag/internal/domains/repository"
	"chatbot-rag/internal/domains/user/models"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type User struct {
	openAI        *openai.Client
	milvusRepo    repository.IMilvusRepositoryInterface
	embeddingRepo repository.IEmbeddingRepositoryInterface
}

func NewUserService(openAI *openai.Client, milvusRepo repository.IMilvusRepositoryInterface,
	embeddingRepo repository.IEmbeddingRepositoryInterface) *User {
	return &User{
		openAI:        openAI,
		milvusRepo:    milvusRepo,
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

		_, err = us.milvusRepo.Insert(
			viper.GetString(constants.MilvusCollection),
			"",
			entity.NewColumnInt64("id", []int64{record.ID}),
			entity.NewColumnBool("is_visible", []bool{true}),
			entity.NewColumnFloatVector("vector", 512, [][]float32{emb.Embedding}))
		if err != nil {
			logrus.Errorf(constants.FormatTaskErr, "milvusRepo.Insert", err)
			return nil, err
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

	searchParamCPU, err := entity.NewIndexIvfFlatSearchParam(64)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "NewIndexIvfFlatSearchParam", err)
		return nil, err
	}

	result, err := us.milvusRepo.FindTopK(
		viper.GetString(constants.MilvusCollection),
		[]string{},
		"",
		[]string{"id", "product_id", "vector"},
		[]entity.Vector{entity.FloatVector(resp.Data[0].Embedding)},
		"vector",
		entity.IP,
		viper.GetInt(constants.MilvusTopK),
		searchParamCPU,
	)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "FindTopK", err)
		return nil, err
	}

	if len(result) == 0 || result[0].ResultCount == 0 {
		logrus.Error("milvus doesnt have any vectors")
		return nil, fmt.Errorf("milvus doesnt have any vectors")
	}

	ids := make([]int64, 0)
	for i := 0; i < result[0].Fields.GetColumn("id").Len(); i++ {
		id, errGetColumn := result[0].Fields.GetColumn("id").GetAsInt64(i)
		if errGetColumn != nil {
			logrus.Errorf(constants.FormatTaskErr, "Fields.GetColumn", err)
			continue
		}

		if float64(result[0].Scores[i]) < viper.GetFloat64(constants.MilvusScore) {
			logrus.Errorf("score is small than threshold")
			continue
		}

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
Trả lời trả về là dạng html trong thẻ div.
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
		Data:   completion.Choices[0].Message.Content,
	}, nil
}
