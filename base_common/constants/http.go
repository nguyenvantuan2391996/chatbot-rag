package constants

const (
	RequestIDField = "request_id"
)

const (
	InternalServerError = -1
)

var ResponseMessage = map[int]string{
	InternalServerError: "The system has an error!",
}
