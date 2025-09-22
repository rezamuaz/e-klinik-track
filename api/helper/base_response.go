package helper

import "e-klinik/api/validation"

type BaseHttpResponse struct {
	Result           any                           `json:"result"`
	Success          bool                          `json:"success"`
	ResultCode       ResultCode                    `json:"rc"`
	ValidationErrors *[]validation.ValidationError `json:"validation_errors,omitempty"`
	Error            any                           `json:"error,omitempty"`
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
		Success:          success,
		ResultCode:       resultCode,
		Result:           result,
		ValidationErrors: validation.GetValidationErrors(err),
	}
}
func WithPaginate(data any, pagination any) *SuccessResponse {
	return &SuccessResponse{
		Data:       data,
		Pagination: pagination,
	}
}
