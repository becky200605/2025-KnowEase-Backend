//go:build wireinject
// +build wireinject

package main

import (
	rag "KnowEase/RAG/rag_go/client"
	"KnowEase/controllers"
	"KnowEase/dao"
	"KnowEase/middleware"
	"KnowEase/routes"
	"KnowEase/services"

	"github.com/google/wire"
)

func InitializeApp() *routes.APP {
	wire.Build(
		dao.ProviderSet,
		rag.NewRAGService,
		services.ProviderSet,
		controllers.ProviderSet,
		routes.ProviderSet,
		middleware.NewMiddleWare,
		routes.NewApp,
		ProvideDBConnectionString,
	)
	return nil
}

