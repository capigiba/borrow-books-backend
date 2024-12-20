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

// BookHandler handles book-related HTTP requests.
type BookHandler struct {
	svc service.BookService
}

// NewBookHandler creates a new BookHandler.
func NewBookHandler(svc service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// ListBooks godoc
// @Summary List books
// @Description Get a list of books with optional filters, sorts, and selected fields
// @Tags Books
// @Accept json
// @Produce json
// @Param filter query []string false "Filter conditions"
// @Param sort query []string false "Sort conditions"
// @Param fields query string false "Fields to select"
// @Success 200 {array} response.BookResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {
	filters := c.QueryArray("filter")
	sorts := c.QueryArray("sort")
	fields := c.Query("fields")

	books, err := h.svc.ListBooks(c.Request.Context(), filters, sorts, fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert each Book to BookResponse
	resp := make([]response.BookResponse, len(books))
	for i, b := range books {
		resp[i] = b.ConvertToResponse()
	}

	c.JSON(http.StatusOK, resp)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieve a single book using its unique ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} response.BookResponse
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	book, err := h.svc.GetBook(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if book == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "not found"})
		return
	}

	resp := book.ConvertToResponse()
	c.JSON(http.StatusOK, resp)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the system
// @Tags Books
// @Accept json
// @Produce json
// @Param book body request.CreateBookRequest true "Book to create"
// @Success 201 {object} response.BookResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req request.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Parse the date string to UNIX timestamp
	tm, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	book, err := h.svc.CreateBook(c.Request.Context(), req.Title, req.AuthorID, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := book.ConvertToResponse()
	c.JSON(http.StatusCreated, resp)
}

// UpdateBook godoc
// @Summary Update an existing book
// @Description Modify the details of an existing book using its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body request.UpdateBookRequest true "Book data to update"
// @Success 200 {object} response.BookResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "Book not found"
// @Failure 500 {object} response.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	var req request.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Parse the date string to UNIX timestamp
	tm, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid published_at format"})
		return
	}
	timestamp := tm.Unix()

	book, err := h.svc.UpdateBook(c.Request.Context(), id, req.Title, req.AuthorID, timestamp)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := book.ConvertToResponse()
	c.JSON(http.StatusOK, resp)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Remove a book from the system using its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 404 {object} response.ErrorResponse "Book not found"
// @Failure 500 {object} response.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	err = h.svc.DeleteBook(c.Request.Context(), id)
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
