package handler

import (
	"borrow_book/internal/domain/model"
	"borrow_book/pkg/response"
	"time"
)

func convertToResponse(b model.Book) response.BookResponse {
	tm := time.Unix(b.PublishedAt, 0).UTC() // Convert timestamp to time.Time
	return response.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		AuthorID:    b.AuthorID,
		PublishedAt: tm.Format("2006-01-02"), // Format as "YYYY-MM-DD"
	}
}
