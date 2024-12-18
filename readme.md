```.
├── Makefile
├── cmd
│   ├── migration
│   │   ├── v1
│   │   │   ├── main.go
│   │   │   └── readme.md
│   │   └── v2
│   │       ├── main.go
│   │       └── readme.md
│   ├── server
│   │   └── main.go
│   └── setup
│       └── main.go
├── config.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── i18n
├── internal
│   ├── config
│   │   └── config.go
│   ├── domain
│   │   └── model
│   │       ├── author.go
│   │       ├── book.go
│   │       └── borrow.go
│   ├── handler
│   │   ├── author_handler.go
│   │   ├── book_handler.go
│   │   ├── provider.go
│   │   └── raw_book_handler.go
│   ├── infra
│   │   ├── database
│   │   │   ├── postgres
│   │   │   │   └── postgres.go
│   │   │   └── query
│   │   │       ├── builder.go
│   │   │       ├── constants.go
│   │   │       └── query_parser.go
│   │   └── server
│   │       ├── http
│   │       │   ├── http.go
│   │       │   └── http_test.go
│   │       └── server.go
│   ├── initialize
│   │   ├── database.go
│   │   ├── router.go
│   │   ├── run.go
│   │   ├── server.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── repository
│   │   ├── author_repo.go
│   │   ├── book_repo.go
│   │   ├── constants.go
│   │   ├── provider.go
│   │   └── raw_book_repo.go
│   ├── router
│   │   ├── provider.go
│   │   ├── router.go
│   │   └── swagger_router.go
│   └── service
│       ├── author_service.go
│       ├── book_service.go
│       ├── provider.go
│       └── raw_book_service.go
├── pkg
│   ├── localization
│   │   └── localizer.go
│   ├── logger
│   │   └── logger.go
│   └── reason
│       └── reason.go
└── scripts```