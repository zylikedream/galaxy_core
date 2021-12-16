package component

import "context"

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

func (r *RoleBag) AddItem(ctx context.Context, itemList []Item) error {
	return nil
}
