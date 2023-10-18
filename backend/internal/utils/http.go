package utils

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/image-randomizer/internal/errors"
)

// Response is standard api response model.
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

// ResponseWithJSON to write response with JSON format.
func ResponseWithJSON(w http.ResponseWriter, code int, data interface{}, err error) {
	r := Response{
		Status:  code,
		Message: strings.ToLower(http.StatusText(code)),
	}

	r.Data = data
	if err != nil {
		r.Message = err.Error()
	}

	rJSON, _ := json.Marshal(r)

	// Set response header.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(rJSON)))
	w.WriteHeader(code)

	_, _ = w.Write(rJSON)
}

// ResponseWithImage serve image as response.
func ResponseWithImage(ctx context.Context, w http.ResponseWriter, image string) {
	img, err := http.Get(image)
	if err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, nil, stack.Wrap(ctx, err, errors.ErrInternalServer))
		return
	}

	if img.StatusCode != http.StatusOK {
		ResponseWithJSON(w, img.StatusCode, nil, stack.Wrap(ctx, _errors.New(http.StatusText(img.StatusCode))))
		return
	}

	defer img.Body.Close()
	w.Header().Set("Content-Type", "image/jpeg")

	if _, err := io.Copy(w, img.Body); err != nil {
		ResponseWithJSON(w, http.StatusInternalServerError, nil, stack.Wrap(ctx, err, errors.ErrInternalServer))
	}
}

// Recoverer is custom recoverer middleware.
// Will return 500.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				ResponseWithJSON(
					w,
					http.StatusInternalServerError,
					nil,
					stack.Wrap(r.Context(),
						fmt.Errorf("%s", debug.Stack()),
						fmt.Errorf("%v", rvr),
						errors.ErrInternalServer))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
