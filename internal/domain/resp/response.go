package resp

import (
	"e-klinik/infra/pg"
	"e-klinik/internal/domain/entity"
	"e-klinik/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struktur response sukses / error utama
type BaseHttpResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Code    string             `json:"code,omitempty"`
	Result  any                `json:"result,omitempty"`
	Error   *pkg.ErrorResponse `json:"error,omitempty"`
}

// Response sukses
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, BaseHttpResponse{
		Success: true,
		Message: message,
		Result:  data,
	})
}

// HandleErrorResponse kirim respons error standar
func HandleErrorResponse(c *gin.Context, message string, err error) {
	appErr, ok := pkg.AsAppError(err)
	if !ok {
		appErr = ToAppError(err, pkg.ErrorCodeInternal, message)
	}

	resp := BaseHttpResponse{
		Success: false,
		Message: message,
		Error: &pkg.ErrorResponse{
			Message:     appErr.Message,
			Details:     appErr.Orig,
			Validations: appErr.Validations,
		},
	}

	c.JSON(appErr.HTTPStatus(), resp)
}

// HandleSuccessResponse kirim respons sukses standar
func HandleSuccessResponse(c *gin.Context, message string, result any) {
	c.JSON(http.StatusOK, BaseHttpResponse{
		Success: true,
		Message: message,
		Result:  result,
	})
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

// ToAppError memastikan error biasa dikonversi jadi *AppError
func ToAppError(err error, code pkg.ErrorCode, msg string) *pkg.AppError {
	if ae, ok := pkg.AsAppError(err); ok {
		return ae
	}

	wrapped := pkg.WrapError(err, code, msg)
	if ae, ok := wrapped.(*pkg.AppError); ok {
		return ae
	}

	return &pkg.AppError{
		Code:    code,
		Message: msg,
		Orig:    err,
		Expose:  false,
	}
}

func BuildMenuTree(items []pg.GetR1ViewRecursiveRow, parentID *int32) []*entity.MenuNode {
	var nodes []*entity.MenuNode

	for _, item := range items {
		// Skip jika view == nil (tidak ditampilkan di frontend)
		if item.View == nil {
			continue
		}

		// Root (ParentID == nil) atau cocok dengan parent
		if (parentID == nil && item.ParentID == nil) ||
			(parentID != nil && item.ParentID != nil && *item.ParentID == *parentID) {

			node := &entity.MenuNode{
				ID:          item.ID,
				Label:       item.Label,
				ResourceKey: item.ResourceKey,
				Action:      item.Action,
				View:        item.View,
				Data:        item.Data,
				Level:       item.Level,
				Path:        item.Path,
			}

			children := BuildMenuTree(items, &item.ID)
			if len(children) > 0 {
				node.Children = children
			}

			nodes = append(nodes, node)
		}
	}

	return nodes
}
