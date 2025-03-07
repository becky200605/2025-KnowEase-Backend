package routes

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserPageSvc,
	NewPostSvc,
	NewUserSvc,
	NewQASvc,
	NewSearchSvc,
	NewAIChatSvc,
)
