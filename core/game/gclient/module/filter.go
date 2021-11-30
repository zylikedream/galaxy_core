package module

import "context"

type ModuleFilter func(ctx context.Context, msg interface{}) error
