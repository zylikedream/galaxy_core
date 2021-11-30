package gmongo

import (
	"context"

	"github.com/zylikedream/galaxy/core/gconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type gmongo struct {
	config *mongoConfig
	client *mongo.Client
}

type mongoConfig struct {
	Addr        string `toml:"addr"`
	DataBase    string `toml:"db"`
	MaxPoolSize int    `toml:"pool_size.max"`
	MinPoolSize int    `toml:"pool_size.min"`
}

func NewMongo(ctx context.Context, configFile string) (*gmongo, error) {
	conf := &mongoConfig{}
	configure := gconfig.New(configFile)
	if err := configure.UnmarshalKey("mongo", conf); err != nil {
		return nil, err
	}
	opt := options.Client()
	opt.ApplyURI(conf.Addr)
	opt.SetMinPoolSize(uint64(conf.MinPoolSize))
	opt.SetMaxPoolSize(uint64(conf.MaxPoolSize))
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &gmongo{
		config: conf,
		client: client,
	}, nil
}

func (m *gmongo) GetDatabase(ctx context.Context) string {
	return m.config.DataBase
}

func (m *gmongo) FindOne(ctx context.Context, Col string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.FindOne(ctx, filter, opts...)
}

func (m *gmongo) Find(ctx context.Context, Col string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.Find(ctx, filter, opts...)
}

func (m *gmongo) UpdateSetOne(ctx context.Context, Col string, filter interface{}, Set interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.UpdateOne(ctx, filter, bson.M{"$set": Set}, opts...)
}

func (m *gmongo) UpdateOne(ctx context.Context, Col string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.UpdateOne(ctx, filter, update, opts...)
}

func (m *gmongo) UpdateMany(ctx context.Context, Col string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.UpdateMany(ctx, filter, update, opts...)
}

func (m *gmongo) InsertOne(ctx context.Context, Col string, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.InsertOne(ctx, doc, opts...)
}

func (m *gmongo) InsertMany(ctx context.Context, Col string, docs []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.InsertMany(ctx, docs, opts...)
}

func (m *gmongo) DeleteOne(ctx context.Context, Col string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.DeleteOne(ctx, filter, opts...)
}

func (m *gmongo) DeleteMany(ctx context.Context, Col string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.DeleteOne(ctx, filter, opts...)
}
