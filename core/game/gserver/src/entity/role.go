package entity

import "github.com/zylikedream/galaxy/core/game/gserver/src/component"

type RoleEntity struct {
	roleInfo component.RoleInfo
}

func NewRoleEntity() *RoleEntity {
	return &RoleEntity{}
}
