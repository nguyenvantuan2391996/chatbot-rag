package response

import (
	"context"
	"fmt"
	"time"

	"chatbot-rag/base_common/constants"
)

const (
	ErrorCodeOk = iota
	ErrorCodeAuth
	ErrorCodeInput
	ErrorCodeServer
	ErrorCodeDatabase
)

var (
	ErrorMsgOk     = "OK"
	ErrorMsgServer = "Internal Server Error"

	ErrorMsgDatabase = "Database Error"
	ErrorMsgInput    = "Invalid Input Error"
	ErrorMsgAuth     = "Authorization Error"
)

type Response struct {
	Data         interface{}             `json:"data,omitempty"`
	Agg          *map[string]interface{} `json:"agg,omitempty"`
	Meta         *map[string]interface{} `json:"meta,omitempty"`
	RequestID    string                  `json:"request_id"`
	ErrorMessage string                  `json:"message"`
	ErrorCode    int                     `json:"code"`
	ServerTime   int64                   `json:"server_time"`
	Count        int                     `json:"count,omitempty"`
}

func NewResponse(ctx context.Context) *Response {
	resp := new(Response)
	resp.RequestID = fmt.Sprintf("%v", ctx.Value(constants.RequestIDField))
	resp.ServerTime = time.Now().Unix()
	resp.NoError()
	resp.Data = map[string]interface{}{}
	return resp
}

func (resp *Response) ToResponse(code int, data interface{}, message string) *Response {
	resp.ErrorCode = code
	if len(message) > 0 {
		resp.ErrorMessage = message
	}

	if data != nil {
		resp.Data = data
	}

	return resp
}

func (resp *Response) Code(code int) *Response {
	resp.ErrorCode = code
	return resp
}

func (resp *Response) Msg(msg string) *Response {
	resp.ErrorMessage = msg
	return resp
}

func (resp *Response) ServerError() *Response {
	return resp.Code(ErrorCodeServer).Msg(ErrorMsgServer)
}

func (resp *Response) DatabaseError() *Response {
	return resp.Code(ErrorCodeDatabase).Msg(ErrorMsgDatabase)
}

func (resp *Response) InputError() *Response {
	return resp.Code(ErrorCodeInput).Msg(ErrorMsgInput)
}

func (resp *Response) AuthError() *Response {
	return resp.Code(ErrorCodeAuth).Msg(ErrorMsgAuth)
}

func (resp *Response) NoError() *Response {
	return resp.Code(ErrorCodeOk).Msg(ErrorMsgOk)
}
