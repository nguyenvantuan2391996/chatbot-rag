package handler

import (
	"fmt"
	"net/http"

	"chatbot-rag/base_common/constants"
	"chatbot-rag/base_common/response"
	"chatbot-rag/handler/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Index(ctx *gin.Context) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginAPI, "Index"))
	request := models.FactRequest{}
	responseAPI := response.NewResponse(ctx)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Warnf(constants.FormatTaskErr, "ShouldBindJSON", err)
		ctx.JSON(http.StatusBadRequest, responseAPI.InputError().Msg(response.ErrorMsgInput))
		return
	}

	if err = request.Validate(); err != nil {
		logrus.Errorf(constants.FormatTaskErr, "Validate", err)
		ctx.JSON(http.StatusBadRequest, responseAPI.InputError().Msg(err.Error()))
		return
	}

	result, err := h.userService.Index(ctx, request.ToFactInput())
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "Index", err)
		ctx.JSON(http.StatusInternalServerError, responseAPI.ToResponse(constants.InternalServerError,
			nil, constants.ResponseMessage[constants.InternalServerError]))
		return
	}

	responseAPI.ToResponse(result.Status, result.Data, result.Message)
	ctx.JSON(result.Status, responseAPI)
}

func (h *Handler) Chat(ctx *gin.Context) {
	logrus.Info(fmt.Sprintf(constants.FormatBeginAPI, "Chat"))
	request := models.ChatRequest{}
	responseAPI := response.NewResponse(ctx)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Warnf(constants.FormatTaskErr, "ShouldBindJSON", err)
		ctx.JSON(http.StatusBadRequest, responseAPI.InputError().Msg(response.ErrorMsgInput))
		return
	}

	if err = request.Validate(); err != nil {
		logrus.Errorf(constants.FormatTaskErr, "Validate", err)
		ctx.JSON(http.StatusBadRequest, responseAPI.InputError().Msg(err.Error()))
		return
	}

	result, err := h.userService.Chat(ctx, request.ToChatInput())
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "Chat", err)
		ctx.JSON(http.StatusInternalServerError, responseAPI.ToResponse(constants.InternalServerError,
			nil, constants.ResponseMessage[constants.InternalServerError]))
		return
	}

	responseAPI.ToResponse(result.Status, result.Data, result.Message)
	ctx.JSON(result.Status, responseAPI)
}
