package model

import (
	"borrow_book/internal/domain/response"
	"time"
)

type Borrow struct {
	ID         int    `db:"id" json:"id"`
	BookID     int    `db:"book_id" json:"book_id"`
	UserName   string `db:"user_name" json:"user_name"` // Name of the user borrow that book
	BorrowedAt int64  `db:"borrowed_at" json:"borrowed_at"`
}

func (b *Borrow) ConvertToResponse() response.BorrowResponse {
	tm := time.Unix(b.BorrowedAt, 0).UTC() // Convert timestamp to time.Time
	return response.BorrowResponse{
		ID:         b.ID,
		BookID:     b.BookID,
		UserName:   b.UserName,
		BorrowedAt: tm.Format("2006-01-02"), // Format as "YYYY-MM-DD"
	}
}
