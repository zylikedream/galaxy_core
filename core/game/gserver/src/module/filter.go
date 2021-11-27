package module

import (
	"github.com/zylikedream/galaxy/core/gcontext"
)

type ModuleFilter func(ctx gcontext.Context, msg interface{}) error
