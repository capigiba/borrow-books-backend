package request

type BaseBorrowRequest struct {
	BookID     int    `json:"book_id"`
	UserName   string `json:"user_name"`
	BorrowedAt string `json:"borrowed_at"`
}

type CreateBorrowRequest struct {
	BaseBorrowRequest
}

type UpdateBorrowRequest struct {
	BaseBorrowRequest
}
