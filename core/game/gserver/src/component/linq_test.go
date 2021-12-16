package component

import (
	"testing"

	"github.com/ahmetb/go-linq"
)

func TestLinq(t *testing.T) {
	itemList := []Item{
		{1, 100},
		{1, 101},
		{2, 200},
		{3, 300},
	}
	uniqItemList := []Item{}
	linq.From(itemList).GroupByT(
		func(it Item) int { return it.PropID },
		func(it Item) uint64 { return it.Num },
	).ToSlice(uniqItemList)
}
