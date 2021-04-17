package core

import "context"

type Controller func(ctx context.Context, req Request) (interface{}, int, error)
