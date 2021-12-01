package gmongo

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	"time"

// 	"github.com/zylikedream/galaxy/core/gconfig"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const (
// 	MONGO_MSG_TYPE_FINDONE = iota + 1
// 	MONGO_MSG_TYPE_FIND
// 	MONGO_MSG_TYPE_UPDATE_ONE
// 	MONGO_MSG_TYPE_UPDATE
// 	MONGO_MSG_TYPE_INSERT
// 	MONGO_MSG_TYPE_INSERT_ONE
// 	MONGO_MSG_TYPE_DELETE
// 	MONGO_MSG_TYPE_DELETE_ONE
// )

// type gmongoClient struct {
// 	cmdChan chan *mongoCmd
// 	config  *clientConfig
// 	client  *mongo.Client
// 	cmds    []*mongoCmd
// }

// type clientConfig struct {
// 	Database string `toml:"database"`
// 	BatchNum int    `toml:"batch_num"`
// 	Interval int    `toml:"interval"`
// }

// type CmdResult struct {
// 	data interface{}
// 	err  error
// }

// type mongoCmd struct {
// 	Type     int
// 	Database string
// 	Col      string
// 	filter   interface{}
// 	doc      interface{}
// 	opts     interface{}
// 	result   chan CmdResult
// }

// func newMongoClient(ctx context.Context, configure *gconfig.Configuration) (*gmongoClient, error) {
// 	conf := &clientConfig{}
// 	if err := configure.UnmarshalKey("client", conf); err != nil {
// 		return nil, err
// 	}
// 	return &gmongoClient{
// 		cmdChan: make(chan *mongoCmd, 1024),
// 		config:  conf,
// 		cmds:    make([]*mongoCmd, 1024),
// 	}, nil
// }

// func (c *gmongoClient) newMongoCmd(ctx context.Context, msgType int, col string, filter interface{}, doc interface{}, opts interface{}) *mongoCmd {
// 	return &mongoCmd{
// 		Type:     msgType,
// 		Database: c.GetDatabase(ctx),
// 		Col:      col,
// 		filter:   filter,
// 		doc:      doc,
// 		opts:     opts,
// 		result:   make(chan CmdResult, 1),
// 	}
// }

// func (c *gmongoClient) GetDatabase(ctx context.Context) string {
// 	return c.config.Database
// }

// func (c *gmongoClient) FindOne(ctx context.Context, result interface{}, col string, filter interface{}, opts ...*options.FindOneOptions) error {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_FINDONE, col, filter, nil, opts)
// 	c.cmdChan <- cmd
// 	res := <-cmd.result
// 	if res.err != nil {
// 		return res.err
// 	}
// 	data := res.data.([]byte)
// 	return bson.Unmarshal(data, result)

// }

// func (c *gmongoClient) Find(ctx context.Context, results interface{}, col string, filter interface{}, opts ...*options.FindOptions) error {
// 	resultsVal := reflect.ValueOf(results)
// 	if resultsVal.Kind() != reflect.Ptr {
// 		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
// 	}

// 	sliceVal := resultsVal.Elem()
// 	if sliceVal.Kind() == reflect.Interface {
// 		sliceVal = sliceVal.Elem()
// 	}

// 	if sliceVal.Kind() != reflect.Slice {
// 		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
// 	}

// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_FIND, col, filter, nil, opts)
// 	c.cmdChan <- cmd
// 	res := <-cmd.result
// 	if res.err != nil {
// 		return res.err
// 	}
// 	datas := res.data.([][]byte)
// 	elemType := sliceVal.Type().Elem()

// 	for i, data := range datas {
// 		if sliceVal.Len() == i {
// 			// slice is full
// 			newElem := reflect.New(elemType)
// 			sliceVal = reflect.Append(sliceVal, newElem.Elem())
// 			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
// 		}

// 		currElem := sliceVal.Index(i).Addr().Interface()
// 		if err := bson.Unmarshal(data, currElem); err != nil {
// 			return err
// 		}
// 	}
// 	resultsVal.Elem().Set(sliceVal.Slice(0, len(datas)))
// 	return nil
// }

// // 更新操作，默认不关心结果
// func (c *gmongoClient) UpdateSetOne(ctx context.Context, col string, filter interface{}, Set interface{}, opts ...*options.UpdateOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_UPDATE_ONE, col, filter, bson.M{"$set": Set}, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) UpdateOne(ctx context.Context, col string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_UPDATE_ONE, col, filter, update, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) UpdateMany(ctx context.Context, col string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_UPDATE, col, filter, update, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) InsertOne(ctx context.Context, col string, doc interface{}, opts ...*options.InsertOneOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_INSERT_ONE, col, nil, doc, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) InsertMany(ctx context.Context, col string, docs []interface{}, opts ...*options.InsertManyOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_INSERT, col, nil, docs, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) DeleteOne(ctx context.Context, col string, filter interface{}, opts ...*options.DeleteOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_DELETE_ONE, col, filter, nil, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) DeleteMany(ctx context.Context, col string, filter interface{}, opts ...*options.DeleteOptions) {
// 	cmd := c.newMongoCmd(ctx, MONGO_MSG_TYPE_DELETE, col, filter, nil, opts)
// 	c.cmdChan <- cmd
// }

// func (c *gmongoClient) worker(ctx context.Context) {
// 	tick := time.NewTicker(time.Second)
// 	defer tick.Stop()
// 	for {
// 		select {
// 		case cmd := <-c.cmdChan:
// 			c.cmds = append(c.cmds, cmd)
// 			c.handleCmds(c.cmds)
// 		}
// 		select {
// 		case <-tick.C:
// 			c.handleCmds(c.cmds)
// 			// 写入数据库
// 		}
// 	}
// }

// func (c *gmongoClient) handleCmds(cmds []*mongoCmd) {
// 	if len(cmds) < c.config.BatchNum {
// 		return
// 	}
// 	for i := len(c.cmds) - 1; i >= 0; i-- {
// 	}
// }
