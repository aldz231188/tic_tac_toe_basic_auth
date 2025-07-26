package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
	"t03/internal/di"
)

func main() {
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	app := fx.New(di.Module)

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
