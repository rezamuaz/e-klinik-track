package api

import (
	"e-klinik/pkg"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	// "go.opentelemetry.io/otel/trace"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
}

func RenderErrorResponse(c *gin.Context, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var ierr *pkg.Error
	if !errors.As(err, &ierr) {
		resp.Error = "internal error"
	} else {
		switch ierr.Code() {
		case pkg.ErrorCodeNotFound:
			status = http.StatusNotFound
		case pkg.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest

			var verrors validation.Errors
			if errors.As(ierr, &verrors) {
				resp.Validations = verrors
			}
		case pkg.ErrorCodeUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
		}
	}

	// if err != nil {
	// 	_, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "rest.renderErrorResponse")
	// 	defer span.End()

	// 	span.RecordError(err)
	// }

	// XXX fmt.Printf("Error: %v\n", err)

	renderResponse(c, resp, status)
}

func renderResponse(c *gin.Context, res interface{}, status int) {
	c.JSON(status, res)
}
