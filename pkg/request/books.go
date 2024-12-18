package request

type BaseBookRequest struct {
	Title       string `json:"title"`
	AuthorID    int    `json:"author_id"`
	PublishedAt string `json:"published_at"` // Expected format: "YYYY-MM-DD"
}

type CreateBookRequest struct {
	BaseBookRequest
}

type UpdateBookRequest struct {
	BaseBookRequest
}
