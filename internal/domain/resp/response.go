package resp

import (
	"e-klinik/pkg"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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

type Pagination struct {
	Limit       int64 `json:"limit"`
	TotalPage   int64 `json:"total_page"`
	TotalRows   int64 `json:"total_rows"`
	CurrentPage int32 `json:"current_page"`
	NextPage    bool  `json:"next"`
	PrevPage    bool  `json:"prev"`
}

func GenerateBaseResponse(result any, success bool, resultCode ResultCode) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success:    success,
		ResultCode: resultCode,
		Result:     result,
	}
}

func GenerateBaseResponseWithError(c *gin.Context, msg string, err error) {
	resp := pkg.ErrorResponse{Error: pkg.ErrorDetails{
		Message: msg,
	}}
	status := http.StatusInternalServerError
	var resultCode ResultCode
	var ierr *pkg.Error
	if !errors.As(err, &ierr) {
		resp.Error.Message = "internal error"
		resp.Error.Details = ierr.Unwrap().Error()
	} else {
		switch ierr.Code() {
		case pkg.ErrorCodeNotFound:
			status = http.StatusNotFound
			resultCode = NotFoundError
			resp.Error.Details = ierr.Error()
		case pkg.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest
			resp.Error.Details = ierr.Error()
			resultCode = ValidationError
			var verrors validation.Errors
			if errors.As(ierr, &verrors) {
				resp.Validations = verrors
			}
		case pkg.ErrorCodeUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
			resultCode = InternalError
			resp.Error.Details = ierr.Unwrap().Error()
		}
	}

	bodyReponse := BaseHttpResponse{
		Success:    false,
		ResultCode: resultCode,
		Result:     nil,
		Error:      resp}
	c.JSON(status, bodyReponse)
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

func CalculatePagination(page int32, limit int32, totalRows int64) *Pagination {
	// Pagination parameters
	if limit == 0 {
		limit = 10 // Number of records per page
	}
	// Calculate total pages
	totalPages := totalRows / int64(limit)
	if totalRows%int64(limit) != 0 {
		totalPages++
	}

	return &Pagination{
		Limit:       int64(limit),
		TotalRows:   totalRows,
		TotalPage:   totalPages,
		CurrentPage: page,
		NextPage:    page < int32(totalPages),
		PrevPage:    page > 1,
	}
}

type ResultCode int

const (
	Success         ResultCode = 0
	ValidationError ResultCode = 4000
	AuthError       ResultCode = 4001
	ForbiddenError  ResultCode = 4003
	NotFoundError   ResultCode = 4004
	LimiterError    ResultCode = 4291
	OtpLimiterError ResultCode = 4292
	CustomRecovery  ResultCode = 5001
	InternalError   ResultCode = 5002
)
