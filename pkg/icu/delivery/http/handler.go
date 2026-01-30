package http

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hamilton/icu-app/pkg/domain"
	httputil "github.com/hamilton/icu-app/pkg/util/http"
	"github.com/labstack/echo/v4"
)

type ICUHandler struct {
	Usecase domain.ICUUseCase
}

func NewICUHandler(e *echo.Echo, us domain.ICUUseCase) {
	handler := &ICUHandler{
		Usecase: us,
	}

	// Middleware Group
	g := e.Group("")
	g.Use(handler.AuthAndLogMiddleware)

	g.GET("/search", handler.SearchPasien)
	g.POST("/result", handler.SubmitResult)
	g.GET("/log", handler.GetLogs)
	g.GET("/report", handler.FetchReport)
	g.POST("/admission", handler.AdmitPatient)
	g.GET("/patients", handler.FetchPatients)
}

func (h *ICUHandler) AuthAndLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Log Request (Start)
		// We'll log AFTER processing to get the status code, or we can log "attempt"
		// The requirement says "data log with url /log", implies logging the traffic.

		// 2. Headers Check
		// Content-Type: application/json (for POST)
		if c.Request().Method == http.MethodPost {
			if c.Request().Header.Get("Content-Type") != "application/json" {
				return httputil.WriteErrorResponse(c, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			}
		}
		// Accept: application/json
		if c.Request().Header.Get("Accept") != "application/json" {
			// Some clients might not send it strictly, but let's enforce as per requirement
			return httputil.WriteErrorResponse(c, http.StatusNotAcceptable, "Accept header must be application/json")
		}

		// 3. Auth Check (Bearer Token)
		authHeader := c.Request().Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		expectedToken := os.Getenv("APP_API_TOKEN") // Static token from ENV

		if expectedToken == "" {
			// Fallback or error if env not set
			expectedToken = "STATIC_TOKEN_123"
		}

		if token != expectedToken {
			// Log Failed Auth
			_ = h.Usecase.SaveLog(context.Background(), domain.AppLog{
				Endpoint: c.Request().URL.Path,
				Method:   c.Request().Method,
				Status:   http.StatusUnauthorized,
				Message:  "Unauthorized access attempt",
			})
			return httputil.WriteErrorResponse(c, http.StatusUnauthorized, "Authentication failed")
		}

		// Execute next handler
		err := next(c)

		// 4. Log Result
		res := c.Response()
		_ = h.Usecase.SaveLog(context.Background(), domain.AppLog{
			Endpoint: c.Request().URL.Path,
			Method:   c.Request().Method,
			Status:   res.Status,
			Message:  "Request processed", // Or capture error message if we could
		})

		return err
	}
}

func (h *ICUHandler) SearchPasien(c echo.Context) error {
	ctx := c.Request().Context()
	kodeReg := c.QueryParam("KODE_REG")
	if kodeReg == "" {
		return httputil.WriteErrorResponse(c, http.StatusBadRequest, "KODE_REG is required")
	}

	result, err := h.Usecase.SearchPasien(ctx, kodeReg)
	if err != nil {
		// Assuming error means not found or db error
		// For simplicity, treating as not found or error
		return httputil.WriteErrorResponse(c, http.StatusNotFound, "Data not found or error: "+err.Error())
	}

	return httputil.WriteOkResponse(c, result, "Data found")
}

func (h *ICUHandler) SubmitResult(c echo.Context) error {
	ctx := c.Request().Context()
	var payload domain.ResultSubmission
	if err := c.Bind(&payload); err != nil {
		return httputil.WriteErrorResponse(c, http.StatusBadRequest, "Invalid JSON body")
	}

	if err := h.Usecase.SubmitResult(ctx, payload); err != nil {
		return httputil.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to save result: "+err.Error())
	}

	return httputil.WriteOkResponse(c, nil, "Result saved successfully")
}

func (h *ICUHandler) GetLogs(c echo.Context) error {
	ctx := c.Request().Context()
	date := c.QueryParam("date")
	limitStr := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 100
	}

	logs, err := h.Usecase.GetLogs(ctx, date, limit)
	if err != nil {
		return httputil.WriteErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return httputil.WriteOkResponse(c, logs, "Logs retrieved")
}

func (h *ICUHandler) FetchReport(c echo.Context) error {
	ctx := c.Request().Context()
	results, err := h.Usecase.GetReport(ctx)
	if err != nil {
		return httputil.WriteErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return httputil.WriteOkResponse(c, results, "Report data retrieved")
}

func (h *ICUHandler) AdmitPatient(c echo.Context) error {
	ctx := c.Request().Context()
	var req domain.AdmissionRequest
	if err := c.Bind(&req); err != nil {
		return httputil.WriteErrorResponse(c, http.StatusBadRequest, "Invalid JSON body")
	}

	kodeReg, err := h.Usecase.RegisterAdmission(ctx, req)
	if err != nil {
		return httputil.WriteErrorResponse(c, http.StatusInternalServerError, "Failed to admit patient: "+err.Error())
	}

	return httputil.WriteOkResponse(c, map[string]string{"kode_reg": kodeReg}, "Patient admitted successfully")
}

func (h *ICUHandler) FetchPatients(c echo.Context) error {
	ctx := c.Request().Context()
	limitStr := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10
	}

	patients, err := h.Usecase.GetPatients(ctx, limit)
	if err != nil {
		return httputil.WriteErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return httputil.WriteOkResponse(c, patients, "Patients retrieved")
}
