package httperror

import (
	"external-metrics/pkg/tools/logging"
	"fmt"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Result bool        `json:"result"`
	Errors []error     `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func Response(errors []error, data interface{}) BaseResponse {
	return BaseResponse{
		Result: len(errors) == 0,
		Errors: errors,
		Data:   data,
	}
}

// ErrorView.
type ErrorView struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// ErrorView.
func (err ErrorView) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func ErrorWrapper(logger *logging.Logger, handler func(c *gin.Context) (resp interface{}, status int, err error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		errors := make([]error, 0)
		resp, status, err := handler(c)
		if err != nil {
			logger.Error(err)
			errors = append(errors, ErrorView{Message: err.Error()})
		}
		c.JSON(status, Response(errors, resp))
	}
}
