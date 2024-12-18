package handler

import (
	"borrow_book/internal/service"
	"borrow_book/pkg/request"
	"borrow_book/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	svc service.BookService
}

func NewBookHandler(svc service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	filters := c.QueryArray("filter")
	sorts := c.QueryArray("sort")
	fields := c.Query("fields")

	books, err := h.svc.ListBooks(c.Request.Context(), filters, sorts, fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert each Book to BookResponse
	resp := make([]response.BookResponse, len(books))
	for i, b := range books {
		resp[i] = convertToResponse(b)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.svc.GetBook(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	resp := convertToResponse(*book)
	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req request.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the date string to UNIX timestamp
	tm, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	book, err := h.svc.CreateBook(c.Request.Context(), req.Title, req.AuthorID, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := convertToResponse(*book)
	c.JSON(http.StatusCreated, resp)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the date string to UNIX timestamp
	tm, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	book, err := h.svc.UpdateBook(c.Request.Context(), id, req.Title, req.AuthorID, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := convertToResponse(*book)
	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.DeleteBook(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
