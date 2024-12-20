package handler

import (
	"borrow_book/internal/domain/request"
	"borrow_book/internal/domain/response"
	"borrow_book/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// BorrowHandler handles borrow-related HTTP requests.
type BorrowHandler struct {
	svc service.BorrowService
}

// NewBorrowHandler creates a new BorrowHandler.
func NewBorrowHandler(svc service.BorrowService) *BorrowHandler {
	return &BorrowHandler{svc: svc}
}

// ListBorrows godoc
// @Summary List borrows
// @Description Get a list of borrows with optional filters, sorts, and selected fields
// @Tags Borrows
// @Accept json
// @Produce json
// @Param filter query []string false "Filter conditions"
// @Param sort query []string false "Sort conditions"
// @Param fields query string false "Fields to select"
// @Success 200 {array} response.BorrowResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /borrows [get]
func (h *BorrowHandler) ListBorrows(c *gin.Context) {
	filters := c.QueryArray("filter")
	sorts := c.QueryArray("sort")
	fields := c.Query("fields")

	borrows, err := h.svc.ListBorrowLists(c.Request.Context(), filters, sorts, fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := make([]response.BorrowResponse, len(borrows))
	for i, b := range borrows {
		resp[i] = b.ConvertToResponse()
	}

	c.JSON(http.StatusOK, resp)
}

// GetBorrow godoc
// @Summary Get a borrow by ID
// @Description Retrieve a single borrow using its unique ID
// @Tags Borrows
// @Accept json
// @Produce json
// @Param id path int true "Borrow ID"
// @Success 200 {object} response.BorrowResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /borrows/{id} [get]
func (h *BorrowHandler) GetBorrow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	borrow, err := h.svc.GetBorrow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if borrow == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "not found"})
		return
	}

	resp := borrow.ConvertToResponse()
	c.JSON(http.StatusOK, resp)
}

// CreateBorrow godoc
// @Summary Create a new borrow
// @Description Add a new borrow to the system
// @Tags Borrows
// @Accept json
// @Produce json
// @Param borrow body request.CreateBorrowRequest true "Borrow to create"
// @Success 200 {object} response.BorrowResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /borrows [post]
func (h *BorrowHandler) CreateBorrow(c *gin.Context) {
	var req request.CreateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	tm, err := time.Parse("2006-01-02", req.BorrowedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	borrow, err := h.svc.CreateBorrow(c.Request.Context(), req.BookID, req.UserName, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := borrow.ConvertToResponse()
	c.JSON(http.StatusCreated, resp)
}

// UpdateBorrow godoc
// @Summary Update an existing borrow
// @Description Modify the details of an existing borrow using its ID
// @Tags Borrows
// @Accept json
// @Produce json
// @Param id path int true "Borrow ID"
// @Param borrow body request.UpdateBorrowRequest true "Borrow data to update"
// @Success 200 {object} response.BorrowResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /borrows/{id} [put]
func (h *BorrowHandler) UpdateBorrow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	var req request.UpdateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	tm, err := time.Parse("2006-01-02", req.BorrowedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	borrow, err := h.svc.UpdateBorrow(c.Request.Context(), id, req.BookID, req.UserName, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := borrow.ConvertToResponse()
	c.JSON(http.StatusOK, resp)
}

// DeleteBorrow godoc
// @Summary Delete a borrow
// @Description Remove a borrow from the system using its ID
// @Tags Borrows
// @Accept json
// @Produce json
// @Param id path int true "Borrow ID"
// @Success 200 {object} response.BorrowResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /borrows/{id} [delete]
func (h *BorrowHandler) DeleteBorrow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	err = h.svc.DeleteBorrow(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "no rows deleted" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
