package api

import (
	"fmt"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/render"
)

// Response renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type Response struct {
	Err            error `json:"-,omitempty"` // low-level runtime error
	HTTPStatusCode int   `json:"-"`           // http response status code

	Status    string `json:"status"`          // user-level status message
	AppCode   int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText string `json:"error,omitempty"` // application-level error message, for debugging
}

var ErrNotFound = &Response{
	HTTPStatusCode: 404,
	Status:         "Resource not found.",
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: 400,
		Status:         "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalError(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: 503,
		Status:         "Internal error.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: 422,
		Status:         "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrValidationError(errs validator.ValidationErrors) render.Renderer {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		//TODO: forbid past dates
		/* 		case "date":
		errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field())) */
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return &Response{
		HTTPStatusCode: 400,
		Status:         "Bad request",
		ErrorText:      strings.Join(errMsgs, ", "),
	}
}

func OK() *Response {
	return &Response{
		HTTPStatusCode: 200,
		Status:         "OK.",
	}
}
