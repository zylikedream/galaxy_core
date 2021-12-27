package component

import (
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var (
	ErrItemAddNoGrid        = errors.New("bag full no grid to add item")
	ErrItemDecItemNotEnough = errors.New("dec item not enough")
)

type bagItem struct {
	PropID     int       `json:"prop_id"`
	Num        uint64    `json:"num"`
	Grid       uint64    `json:"grid"` // 占用的格子数
	UpdateTime time.Time `json:"update_time"`
}

type RoleBag struct {
	RoleID  primitive.ObjectID `json:"_id"`
	Items   map[int]bagItem    `json:"items"`
	GridUse int                `json:"grid_use"`

	logger *glog.GalaxyLog `json:"-"`
}

type Item struct {
	PropID int    `json:"prop_id"`
	Num    uint64 `json:"num"`
}

type itemChange struct {
	PropID  int
	PreNum  uint64
	Num     uint64
	PreGrid uint64
	Grid    uint64
}

func (it *bagItem) update(propID int, Num uint64, Grid uint64) *itemChange {
	chg := &itemChange{
		PropID:  propID,
		PreNum:  it.Num,
		Num:     Num,
		PreGrid: it.Grid,
		Grid:    Grid,
	}

	it.PropID = propID
	it.Num = Num
	it.Grid = Grid
	it.UpdateTime = time.Now()

	return chg
}

func (r *RoleBag) CreateByID(ctx *gscontext.Context, ID primitive.ObjectID) {
	r.RoleID = ID
	r.Items = make(map[int]bagItem)
}

func (r *RoleBag) Init(ctx *gscontext.Context) error {
	r.logger = ctx.GetLogger().With(zap.Namespace("role_bag"))
	return nil
}

func (r *RoleBag) AddItem(ctx *gscontext.Context, itemList []Item) error {
	var chgs []*itemChange
	itemList = ClassifyItemList(itemList)
	for _, item := range itemList {
		if chg, err := r.addSingleItem(ctx, item); err != nil {
			// todo 格子满了的处理
			return err
		} else {
			chgs = append(chgs, chg)
		}
	}
	r.notifyItemUpdate(ctx, chgs)
	return nil
}

func (r *RoleBag) addSingleItem(ctx *gscontext.Context, item Item) (*itemChange, error) {
	itemTable := ctx.GetGameConfig().TbItem
	itemconf := itemTable.Get(int32(item.PropID))
	have := r.Items[item.PropID]
	newGrid := (have.Num + item.Num - 1) / uint64(itemconf.MaxOverlap)
	gridAdd := int(newGrid - have.Grid)
	if newGrid > have.Grid && r.IsGridFull(ctx, gridAdd) {
		return nil, errors.Wrapf(ErrItemAddNoGrid, "item_num:%d grid_use:%d ", item.Num, r.GridUse)
	}
	if itemconf.AutoUse {
		// todo 自动使用物品
		return nil, nil
	}
	chg := have.update(item.PropID, have.Num+item.Num, newGrid)
	r.GridUse += gridAdd
	r.logger.Debug("add item success", zap.Any("item", item), zap.Int("gridadd", gridAdd))
	return chg, nil
}

func (r *RoleBag) GetItem(propid int) Item {
	bagItem := r.Items[propid]
	return Item{
		PropID: propid,
		Num:    bagItem.Num,
	}
}

func (r *RoleBag) IsGridFull(ctx *gscontext.Context, add int) bool {
	bagMaxGrid := ctx.GetGameConfig().TbBag.Get().MaxGrid
	return int32(r.GridUse+add) > bagMaxGrid
}

func (r *RoleBag) CheckItem(ctx *gscontext.Context, itemList []Item) bool {
	itemList = ClassifyItemList(itemList)
	return linq.From(itemList).All(func(i interface{}) bool {
		item := i.(Item)
		have := r.GetItem(item.PropID)
		return have.Num >= item.Num
	})
}

func (r *RoleBag) DecItem(ctx *gscontext.Context, itemList []Item) error {
	itemList = ClassifyItemList(itemList)
	var chgs []*itemChange
	for _, item := range itemList {
		if chg, err := r.decSingleItem(ctx, item); err != nil {
			return err
		} else {
			chgs = append(chgs, chg)
		}
	}
	r.notifyItemUpdate(ctx, chgs)
	return nil
}

func (r *RoleBag) decSingleItem(ctx *gscontext.Context, item Item) (*itemChange, error) {
	have := r.Items[item.PropID]
	if item.Num > have.Num {
		return nil, errors.Wrapf(ErrItemDecItemNotEnough, "have:%v, need:%v", have, item)
	}
	itemTable := ctx.GetGameConfig().TbItem
	itemconf := itemTable.Get(int32(item.PropID))
	newGrid := (have.Num - item.Num - 1) / uint64(itemconf.MaxOverlap)
	gridDec := int(newGrid - have.Grid)

	chg := have.update(item.PropID, have.Num-item.Num, newGrid)

	if newGrid == 0 {
		delete(r.Items, item.PropID)
	}
	r.GridUse -= gridDec
	r.logger.Debug("dec item success", zap.Any("item", item), zap.Int("griddec ", gridDec))
	return chg, nil
}

func (r *RoleBag) notifyItemUpdate(ctx *gscontext.Context, chgs []*itemChange) {
	sess := ctx.GetSession()
	msg := proto.NtfItemUpdate{
		Items: []proto.PItemInfo{},
	}
	linq.From(chgs).SelectT(func(i interface{}) interface{} {
		return proto.PItemInfo{
			PropID: i.(*itemChange).PropID,
			Num:    i.(*itemChange).Num,
		}
	}).ToSlice(&msg.Items)
	if err := sess.Send(msg); err != nil {
		r.logger.Error("notify item update failed", zap.Any("msg", msg))
	}
}

func ClassifyItemList(itemList []Item) []Item {
	classifyItemList := []Item{}
	linq.From(itemList).GroupBy(
		func(it interface{}) interface{} { return it.(Item).PropID },
		func(it interface{}) interface{} { return it.(Item).Num },
	).Select(func(i interface{}) interface{} {
		return Item{
			PropID: i.(linq.Group).Key.(int),
			Num:    linq.From(i.(linq.Group).Group).SumUInts(),
		}
	}).ToSlice(&classifyItemList)
	return classifyItemList
}
