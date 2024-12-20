package model

import (
	"borrow_book/internal/domain/response"
	"time"
)

type Book struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	AuthorID    int    `db:"author_id" json:"author_id"`
	PublishedAt int64  `db:"published_at" json:"published_at"`
}

func (b *Book) ConvertToResponse() response.BookResponse {
	tm := time.Unix(b.PublishedAt, 0).UTC() // Convert timestamp to time.Time
	return response.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		AuthorID:    b.AuthorID,
		PublishedAt: tm.Format("2006-01-02"), // Format as "YYYY-MM-DD"
	}
}
