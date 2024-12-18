package handler

import (
	"borrow_book/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExtraHandler struct {
	svc service.ExtraService
}

func NewExtraHandler(svc service.ExtraService) *ExtraHandler {
	return &ExtraHandler{svc: svc}
}

func (h *ExtraHandler) ExecuteRawQuery(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.svc.RunRawQuery(c.Request.Context(), req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
