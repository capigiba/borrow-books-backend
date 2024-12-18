package model

type Borrow struct {
	ID         int   `db:"id" json:"id"`
	BookID     int   `db:"book_id" json:"book_id"`
	BorrowedAt int64 `db:"borrowed_at" json:"borrowed_at"`
}
