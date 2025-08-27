package httpapi

import (
	"net/http"
	"time"

	"github.com/itdyaingenieria/otp-service/internal/domain"
	"github.com/itdyaingenieria/otp-service/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	gen usecase.GenerateOTP
	val usecase.ValidateOTP
}

func NewHandler(gen usecase.GenerateOTP, val usecase.ValidateOTP) Handler {
	return Handler{gen: gen, val: val}
}

func (h Handler) Register(e *echo.Echo) {
	g := e.Group("/api/v1")
	g.POST("/otp", h.generate)
	g.POST("/otp/validate", h.validate)
}

type generateReq struct {
	TenantID    string         `json:"tenant_id"`
	Channel     domain.Channel `json:"channel"` // "sms" or "email"
	Destination string         `json:"destination"`
}

type generateResp struct {
	ID        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}

func (h Handler) generate(c echo.Context) error {
	var req generateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	out, err := h.gen.Execute(c.Request().Context(), usecase.GenerateInput{TenantID: req.TenantID, Channel: req.Channel, Destination: req.Destination})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, generateResp{ID: out.ID.String(), ExpiresAt: out.ExpiresAt.Format(time.RFC3339)})
}

type validateReq struct {
	TenantID string `json:"tenant_id"`
	ID       string `json:"id"`
	Code     string `json:"code"`
}

type validateResp struct {
	Valid bool `json:"valid"`
}

func (h Handler) validate(c echo.Context) error {
	var req validateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	out, err := h.val.Execute(c.Request().Context(), usecase.ValidateInput{TenantID: req.TenantID, ID: id, Code: req.Code})
	if err != nil {
		switch err {
		case domain.ErrExpired, domain.ErrMaxAttempts, domain.ErrAlreadyUsed, domain.ErrInvalidCode, domain.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, validateResp{Valid: out.Valid})
}
