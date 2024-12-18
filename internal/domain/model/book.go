package model

type Book struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	AuthorID    int    `db:"author_id" json:"author_id"`
	PublishedAt int64  `db:"published_at" json:"published_at"`
}
