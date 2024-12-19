package handler

import (
	"borrow_book/internal/domain/request"
	"borrow_book/internal/domain/response"
	"borrow_book/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthorHandler handles author-related HTTP requests.
type AuthorHandler struct {
	svc service.AuthorService
}

// NewAuthorHandler creates a new AuthorHandler.
func NewAuthorHandler(svc service.AuthorService) *AuthorHandler {
	return &AuthorHandler{svc: svc}
}

// ListAuthors godoc
// @Summary List authors
// @Description Get a list of authors with optional filters, sorts, and selected fields
// @Tags Authors
// @Accept json
// @Produce json
// @Param filter query []string false "Filter conditions"
// @Param sort query []string false "Sort conditions"
// @Param fields query string false "Fields to select"
// @Success 200 {array} response.AuthorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /authors [get]
func (h *AuthorHandler) ListAuthors(c *gin.Context) {
	filters := c.QueryArray("filter")
	sorts := c.QueryArray("sort")
	fields := c.Query("fields")

	authors, err := h.svc.ListAuthors(c.Request.Context(), filters, sorts, fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert model.Author to response.AuthorResponse
	var authorResponses []response.AuthorResponse
	for _, author := range authors {
		authorResponses = append(authorResponses, response.AuthorResponse{
			ID:   author.ID,
			Name: author.Name,
		})
	}

	c.JSON(http.StatusOK, authorResponses)
}

// GetAuthor godoc
// @Summary Get an author by ID
// @Description Retrieve a single author using their unique ID
// @Tags Authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} response.AuthorResponse
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse
// @Router /authors/{id} [get]
func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	author, err := h.svc.GetAuthor(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if author == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "not found"})
		return
	}

	// Convert model.Author to response.AuthorResponse
	authorResponse := response.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}

	c.JSON(http.StatusOK, authorResponse)
}

// CreateAuthor godoc
// @Summary Create a new author
// @Description Add a new author to the system
// @Tags Authors
// @Accept json
// @Produce json
// @Param author body request.CreateAuthorRequest true "Author to create"
// @Success 201 {object} response.AuthorResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse
// @Router /authors [post]
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req request.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	author, err := h.svc.CreateAuthor(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert model.Author to response.AuthorResponse
	authorResponse := response.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}

	c.JSON(http.StatusCreated, authorResponse)
}

// UpdateAuthor godoc
// @Summary Update an existing author
// @Description Modify the details of an existing author using their ID
// @Tags Authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param author body request.UpdateAuthorRequest true "Author data to update"
// @Success 200 {object} response.AuthorResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "Author not found"
// @Failure 500 {object} response.ErrorResponse
// @Router /authors/{id} [put]
func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	var req request.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	author, err := h.svc.UpdateAuthor(c.Request.Context(), id, req.Name)
	if err != nil {
		if err.Error() == "author not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert model.Author to response.AuthorResponse
	authorResponse := response.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}

	c.JSON(http.StatusOK, authorResponse)
}

// DeleteAuthor godoc
// @Summary Delete an author
// @Description Remove an author from the system using their ID
// @Tags Authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 404 {object} response.ErrorResponse "Author not found"
// @Failure 500 {object} response.ErrorResponse
// @Router /authors/{id} [delete]
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
		return
	}

	err = h.svc.DeleteAuthor(c.Request.Context(), id)
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
