package component

import (
	"github.com/ahmetb/go-linq"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
)

type bagItem struct {
	PropID int    `json:"prop_id"`
	Num    uint64 `json:"num"`
	Grid   int    `json:"grid"` // 占用的格子数
}

type RoleBag struct {
	Items   map[uint64]bagItem `json:"items"`
	GridUse int                `json:"grid_use"`
}

type Item struct {
	PropID int    `json:"prop_id"`
	Num    uint64 `json:"num"`
}

func (r *RoleBag) AddItem(ctx *gscontext.Context, itemList []Item) error {
	// gameconf := ctx.GetGameConfig()
	linq.From(itemList).GroupByT(
		func(it Item) int { return it.PropID },
		func(it Item) uint64 { return it.Num },
	)
	return nil
}
