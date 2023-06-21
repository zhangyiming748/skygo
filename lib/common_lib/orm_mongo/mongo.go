package orm_mongo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"skygo_detection/lib/common_lib/log"
	"skygo_detection/service"
)

type MongoConn struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
	collections   *mongo.Collection
}

var mongoConn *MongoConn

func InitMongoClient() error {
	mongoConn = new(MongoConn)

	config := service.LoadConfig().MongoDB
	user := config.Username
	password := config.Password
	url := config.Host + ":" + strconv.Itoa(int(config.Port))
	dbname := "admin"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	// construct url: mongodb://username:password@127.0.0.1:27017/dbname
	mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
	mongoConn.clientOptions = options.Client().ApplyURI(mongoUrl)
	var err error
	mongoConn.client, err = mongo.Connect(ctx, mongoConn.clientOptions)
	if err != nil {
		log.GetHttpLogLogger().Fatal(fmt.Sprintf("connect to mongodb error: %v", err))
	}

	// Check the connection
	err = mongoConn.client.Ping(context.TODO(), nil)
	if err != nil {
		log.GetHttpLogLogger().Fatal(fmt.Sprintf("check the connection to mongo error: %v", err))
	}
	return nil
}

func GetMongoClient() *mongo.Client {
	if mongoConn == nil {
		InitMongoClient()
	}
	return mongoConn.client
}

func GetDefaultMongoDatabase() *mongo.Database {
	if mongoConn == nil {
		InitMongoClient()
	}

	config := service.LoadMongoDBConfig()
	dbname := config.DBName
	return mongoConn.client.Database(dbname)
}
