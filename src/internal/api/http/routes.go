package http

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"
	"t03/internal/domain"
)

func RegisterRoutes(lc fx.Lifecycle, gameHandler *GameHandler, authService domain.UserService) {
	mux := http.NewServeMux()

	authenticator := NewUserAuthenticator(authService)

	mux.HandleFunc("/signup", gameHandler.HandleSignUpRequest)
	mux.HandleFunc("/signin", gameHandler.HandleSignInRequest)

	mux.HandleFunc("/new-game", authenticator.Protect(gameHandler.HandleNewGame))
	mux.HandleFunc("/game/", authenticator.Protect(gameHandler.HandleGame))
	mux.HandleFunc("/games", authenticator.Protect(gameHandler.HandleGamesList))
	mux.HandleFunc("/stats/", authenticator.Protect(gameHandler.HandlePlayerStats))

	mux.Handle("/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP Server on :8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP Server")
			return server.Shutdown(ctx)
		},
	})
}
