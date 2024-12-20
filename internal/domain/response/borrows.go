package response

type BorrowResponse struct {
	ID         int    `json:"id"`
	BookID     int    `json:"book_id"`
	UserName   string `json:"user_name"`
	BorrowedAt string `json:"borrowed_at"` // Format: "YYYY-MM-DD"
}
