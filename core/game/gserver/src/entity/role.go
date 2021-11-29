package entity

import (
	"github.com/zylikedream/galaxy/core/game/gserver/src/component"
)

type RoleEntity struct {
	EntityID   uint64
	roleInfo   component.RoleInfo
	components []component.Component
}

func NewRoleEntity(enityID uint64) *RoleEntity {
	return &RoleEntity{
		EntityID: enityID,
	}
}

type Persit struct {
	ID      uint64
	persits []component.IPersit
}

func (p *Persit) Load() {

}
