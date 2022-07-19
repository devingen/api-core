package core

import "context"

type Controller func(ctx context.Context, req Request) (*Response, error)
