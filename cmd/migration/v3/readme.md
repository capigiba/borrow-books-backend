# This migration v3 supports to add data to the database

Use --authors and --books flags for specifying the paths to the data files.

```base
go run cmd/migration/v3/main.go --authors=cmd/migration/v3/data/author_x.txt --books=cmd/migration/v3/data/book_x.txt
```