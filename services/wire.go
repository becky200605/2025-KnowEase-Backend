package services

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewEmailService,
	NewLikeService,
	NewPostService,
	NewTokenService,
	NewUserService,
	NewQAService,
	NewSearchService,
	NewAIChatService,
)
