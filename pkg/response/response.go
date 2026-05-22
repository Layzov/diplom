package response

const (
	CodeBadRequest   = "BAD_REQUEST"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
)

type Response struct {
	Error ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(code, msg string) Response {
	return Response{
		Error: ResponseError{
			Code:    code,
			Message: msg,
		},
	}
}
