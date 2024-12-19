package response

type BookResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AuthorID    int    `json:"author_id"`
	PublishedAt string `json:"published_at"` // Format: "YYYY-MM-DD"
}
