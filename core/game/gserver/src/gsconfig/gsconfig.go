package gsconfig

import (
	tabtoy "github.com/davyxu/tabtoy/v3/api/golang"
	"github.com/zylikedream/galaxy/core/game/gserver/gameconfig"
)

type GameConfig struct {
	ItemTable *gameconfig.ItemTable
	BagTable  *gameconfig.BagTable
}

func NewGameConfig() (*GameConfig, error) {
	gc := &GameConfig{}
	if err := gc.initTables(); err != nil {
		return nil, err
	}
	return gc, nil
}

func (gc *GameConfig) initTables() error {
	gc.ItemTable = gameconfig.NewItemTable()
	if err := tabtoy.LoadFromFile(gc.ItemTable, "data/item.json"); err != nil {
		return err
	}
	return nil
}
