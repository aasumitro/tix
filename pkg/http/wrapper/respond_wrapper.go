package wrapper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonRespond struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type SuccessWithPaginationRespond struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	Count    int    `json:"total_data"`
	Total    int    `json:"total_page"`
	Current  int    `json:"current_page"`
	Next     Paging `json:"next"`
	Previous Paging `json:"previous"`
	Data     any    `json:"data"`
}

type Paging struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type ErrorRespond struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

func NewHTTPRespondWrapper(context *gin.Context, code int, data interface{}, args ...any) {
	if code == http.StatusOK || code == http.StatusCreated {
		if len(args) > 0 {
			context.JSON(code, SuccessWithPaginationRespond{
				Code:     code,
				Total:    args[0].(int),
				Current:  args[1].(int),
				Next:     args[2].(Paging),
				Previous: args[3].(Paging),
				Count:    args[4].(int),
				Status:   http.StatusText(code),
				Data:     data,
			})

			return
		}

		context.JSON(code, CommonRespond{
			Code:   code,
			Status: http.StatusText(code),
			Data:   data,
		})

		return
	}

	if code == http.StatusUnprocessableEntity {
		context.JSON(code, CommonRespond{
			Code:   code,
			Status: http.StatusText(code),
			Data:   data,
		})

		return
	}

	msg := func() string {
		switch {
		case data != nil:
			return data.(string)
		case code == http.StatusBadRequest:
			return "something went wrong with the request"
		default:
			return "something went wrong with the server"
		}
	}()

	context.JSON(code, ErrorRespond{
		Code:   code,
		Status: http.StatusText(code),
		Data:   msg,
	})
}
