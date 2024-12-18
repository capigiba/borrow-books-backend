package repository

type TableName string

const (
	TableBooks   TableName = "books"
	TableAuthors TableName = "author"
	TableBorrows TableName = "borrows"
)
