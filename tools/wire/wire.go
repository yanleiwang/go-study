//go:build wireinject

package main

import (
	"github.com/google/wire"
)

func InitServer() *Server {
	wire.Build(
		InitRedis, InitDB,
		NewServer,
	)
	return new(Server)
}
