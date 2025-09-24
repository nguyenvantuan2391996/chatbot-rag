package user

import (
	"context"
	"fmt"
	"net/http"

	"chatbot-rag/base_common/comoutput"
	"chatbot-rag/base_common/constants"
	"chatbot-rag/internal/domains/user/models"
	"github.com/sirupsen/logrus"
)

type User struct {
}

func NewUserService() *User {
	return &User{}
}

func (us *User) Index(ctx context.Context, input *models.IndexInput) (*comoutput.BaseOutput, error) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginTask, "Index", input))

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
