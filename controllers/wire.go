package controllers

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewLikeControllers,
	NewPostControllers,
	NewUserControllers,
	NewQAControllers,
	NewSearchController,
	NewAIChatControllers,
)
