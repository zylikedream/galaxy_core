package cookie

import "github.com/zylikedream/galaxy/core/game/gserver/src/entity"

type Cookie struct {
	Role *entity.RoleEntity
}

func NewCookie() *Cookie {
	return &Cookie{}
}
