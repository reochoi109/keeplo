package app

import (
	"context"
	"keeplo/internal/adapter/rest/router"
)

func Run() {
	ctx := context.Background()
	router.Run(ctx)
}
