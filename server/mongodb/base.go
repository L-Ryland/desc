package mongodb

import (
	"context"
	"flag"
	"server/util"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongodbAddr = flag.String("mongodb.addr", "localhost:27017", "mongodb addr, default localhost:27017")
var mongodbName = flag.String("mongodb.name", "webStorage", "mongodb name, default webStorage")
var db *mongo.Database

const dbTimeoutTime = 5 * time.Second

func NewDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+*mongodbAddr))
	if err != nil {
		panic(util.Errorf("init mongodb error").WithCause(err))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(util.Errorf("ping mongodb error").WithCause(err))
	}

	db = client.Database(*mongodbName)

	InitMongoDB()
}

var dbDatas []dbData

func registerDBData(t dbData) {
	dbDatas = append(dbDatas, t)
}

type dbData interface {
	initTable()
}

func InitMongoDB() {
	for _, t := range dbDatas {
		t.initTable()
	}
}
