package repository

import "github.com/google/wire"

var ProviderSetRepository = wire.NewSet(
	NewBookRepository,
	NewExtraRepository,
	NewAuthorRepository,
)
