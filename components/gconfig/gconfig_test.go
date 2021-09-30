package gconfig

import (
	"fmt"
	"path"
	"testing"
)

func TestPath(t *testing.T) {
	fmt.Println(path.Ext("config"))
	fmt.Println(path.Ext("config.txt"))
}
