package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(r *http.Request, err error) render.Renderer {
	hlog.FromRequest(r).Error().Err(err).Send()
	reqID, _ := hlog.IDFromRequest(r)
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		Msg:            err.Error(),
		RequestID:      reqID.String(),
	}
}

func ErrNotFound(r *http.Request, err error) render.Renderer {
	hlog.FromRequest(r).Error().Err(err).Send()
	reqID, _ := hlog.IDFromRequest(r)
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		Msg:            err.Error(),
		RequestID:      reqID.String(),
	}
}

func ErrUnauthorized(r *http.Request, err error) render.Renderer {
	hlog.FromRequest(r).Error().Err(err).Send()
	reqID, _ := hlog.IDFromRequest(r)
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		Msg:            err.Error(),
		RequestID:      reqID.String(),
	}
}

func ErrServerError(r *http.Request, err error) render.Renderer {
	hlog.FromRequest(r).Error().Err(err).Send()
	reqID, _ := hlog.IDFromRequest(r)
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		Msg:            err.Error(),
		RequestID:      reqID.String(),
	}
}
func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	log.Ctx(ctx).Error().Err(err).Send()
	reqID, _ := hlog.IDFromCtx(ctx)
	w.Header().Set("Content-Type", "application/json")
	if headerer, ok := err.(openapi3filter.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(openapi3filter.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	body, _ := json.Marshal(&ErrResponse{
		HTTPStatusCode: code,
		Msg:            err.Error(),
		RequestID:      reqID.String(),
	})

	w.Write(body)
}
