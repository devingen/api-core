package dvnruntime

import "context"

type ControllerFunc func(ctx context.Context, req Request) (interface{}, int, error)
