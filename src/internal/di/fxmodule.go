package di

import (
	"go.uber.org/fx"
	"os"
	handler "t03/internal/api/http"
	"t03/internal/app"
	"t03/internal/infra/memory"
)

var Module = fx.Options(

	fx.Provide(memory.NewPGConfig),
	fx.Provide(memory.NewStorage),
	fx.Provide(memory.NewGameRepository),
	fx.Provide(app.NewGameService),
	fx.Provide(app.NewUserService),
	fx.Provide(handler.NewGameHandler),

	fx.Invoke(func(g fx.DotGraph) {
		err := os.WriteFile("graph.dot", []byte(g), 0644)
		if err != nil {
			panic(err)
		}
	}),

	fx.Invoke(handler.RegisterRoutes),
)
