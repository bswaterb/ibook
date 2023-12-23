//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"ibook/internal/conf"
	"ibook/internal/data"
	"ibook/internal/service"
	"ibook/internal/web"
)

// wireApp init gin application.
func wireApp(*conf.Secret, *conf.Data, *conf.Server) (*gin.Engine, func(), error) {
	panic(wire.Build(data.DataProviderSet, web.WebProviderSet, service.ServiceProviderSet, newApp))
}
