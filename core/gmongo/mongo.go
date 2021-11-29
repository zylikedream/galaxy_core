package gmongo

import (
	"context"

	"github.com/zylikedream/galaxy/core/gconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type gmongo struct {
	config   *mongoConfig
	client   *mongo.Client
	database string
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
		config:   conf,
		database: conf.DataBase,
		client:   client,
	}, nil
}

func (m *gmongo) GetDatabase(ctx context.Context) string {
	return m.database
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
	return col.UpdateOne(ctx, filter, bson.M{"$set", Set}, opts...)
}

func (m *gmongo) UpdateOne(ctx context.Context, Col string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col := m.client.Database(m.GetDatabase(ctx)).Collection(Col)
	return col.UpdateOne(ctx, filter, update, opts...)
}
