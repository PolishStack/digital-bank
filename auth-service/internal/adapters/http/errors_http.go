package http

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	errs "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/errors"
	"github.com/rs/zerolog"
)

type APIErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func mapErrToHTTP(err error) (int, APIErrorResp) {
	// If it's a coded error, map by code
	var ce *errs.CodedError
	if errors.As(err, &ce) {
		switch ce.Code() {
		case errs.CodeBadRequest:
			return http.StatusBadRequest, APIErrorResp{Code: string(ce.Code()), Message: "invalid request"}
		case errs.CodeNotFound:
			return http.StatusNotFound, APIErrorResp{Code: string(ce.Code()), Message: "resource not found"}
		case errs.CodeConflict:
			return http.StatusConflict, APIErrorResp{Code: string(ce.Code()), Message: "conflict"}
		case errs.CodeUnauth:
			return http.StatusUnauthorized, APIErrorResp{Code: string(ce.Code()), Message: "unauthenticated"}
		default:
			return http.StatusInternalServerError, APIErrorResp{Code: string(errs.CodeInternal), Message: "internal server error"}
		}
	}

	// If not coded, treat as internal
	return http.StatusInternalServerError, APIErrorResp{Code: string(errs.CodeInternal), Message: "internal server error"}
}

func ErrorHandler(logger *zerolog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Map
		status, api := mapErrToHTTP(err)

		// obtain request id / user id if present
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = "none"
		}

		// log full error chain
		logger.Error().
			Err(err).
			Str("req_id", reqID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Msg("request error")

		// attach request id
		c.Set("X-Request-ID", reqID)
		return c.Status(status).JSON(api)
	}
}
