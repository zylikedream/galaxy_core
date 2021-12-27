package gsconfig

import (
	"encoding/json"
	"io/ioutil"

	gameconfig "github.com/zylikedream/galaxy/core/game/gserver/gameconfig/src"
)

type GameConfig struct {
	*gameconfig.Tables
}

func NewGameConfig() (*GameConfig, error) {
	gc := &GameConfig{}
	if err := gc.initTables(); err != nil {
		return nil, err
	}
	return gc, nil
}

func (gc *GameConfig) initTables() error {
	tables, err := gameconfig.NewTables(loader)
	if err != nil {
		return err
	}
	gc.Tables = tables
	return nil
}

func loader(file string) ([]map[string]interface{}, error) {
	if bytes, err := ioutil.ReadFile("data/" + file + ".json"); err != nil {
		return nil, err
	} else {
		jsonData := make([]map[string]interface{}, 0)
		if err = json.Unmarshal(bytes, &jsonData); err != nil {
			return nil, err
		}
		return jsonData, nil
	}
}
