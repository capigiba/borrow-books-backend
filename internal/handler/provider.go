package handler

import "github.com/google/wire"

var ProviderSetHandler = wire.NewSet(
	NewBookHandler,
	NewAuthorHandler,
	NewBorrowHandler,
)
