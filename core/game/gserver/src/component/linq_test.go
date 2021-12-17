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
	grp := []Item{}
	linq.From(itemList).GroupBy(
		func(it interface{}) interface{} { return it.(Item).PropID },
		func(it interface{}) interface{} { return it.(Item).Num },
	).Select(func(i interface{}) interface{} {
		return Item{
			PropID: i.(linq.Group).Key.(int),
			Num:    linq.From(i.(linq.Group).Group).SumUInts(),
		}
	}).ToSlice(&grp)
	t.Logf("resultAll:%v", grp)

}
