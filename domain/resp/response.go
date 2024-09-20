package resp

type BaseHttpResponse struct {
	Result     any        `json:"result,omitempty"`
	Success    bool       `json:"success"`
	ResultCode ResultCode `json:"rc"`
	Error      any        `json:"error,omitempty"`
}

type SuccessResponse struct {
	Data       any `json:"data,omitempty"`
	Pagination any `json:"page,omitempty"`
}

func GenerateBaseResponse(result any, success bool, resultCode ResultCode) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success:    success,
		ResultCode: resultCode,
		Result:     result,
	}
}

func GenerateBaseResponseWithError(result any, success bool, resultCode ResultCode, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success:    success,
		ResultCode: resultCode,
		Result:     result,
		Error:      err.Error(),
	}

}

func GenerateBaseResponseWithAnyError(result any, success bool, resultCode ResultCode, err any) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success:    success,
		ResultCode: resultCode,
		Result:     result,
		Error:      err,
	}
}

func GenerateBaseResponseWithValidationError(result any, success bool, resultCode ResultCode, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success:    success,
		ResultCode: resultCode,
		Result:     result,
		Error:      err.Error(),
	}
}
func WithPaginate(data any, pagination any) *SuccessResponse {
	return &SuccessResponse{
		Data:       data,
		Pagination: pagination,
	}
}

type ResultCode int

const (
	Success         ResultCode = 0
	ValidationError ResultCode = 40001
	AuthError       ResultCode = 40101
	ForbiddenError  ResultCode = 40301
	NotFoundError   ResultCode = 40401
	LimiterError    ResultCode = 42901
	OtpLimiterError ResultCode = 42902
	CustomRecovery  ResultCode = 50001
	InternalError   ResultCode = 50002
)
