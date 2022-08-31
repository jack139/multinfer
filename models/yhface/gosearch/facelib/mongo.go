package facelib

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/event"
)

var (
	MongoURL, MongoUser, MongoPwd string // db 参数
	Client *mongo.Client
)

func ConnectDb() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.PoolMonitor{
		Event: HandlePoolMonitor,
	}

	credential := options.Credential{
		Username: MongoUser,
		Password: MongoPwd,
		AuthSource: "face_db",
	}

	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(MongoURL).
			SetAuth(credential).
			SetMinPoolSize(0).
			SetMaxPoolSize(100).
			SetPoolMonitor(monitor))
	if err != nil {
		log.Println("Connect DB fail: ", err.Error())
		return false
	}
	Client = client
	return Ping()
}

func reconnect(client *mongo.Client) {
	for {
		if ConnectDb() {
			client = Client
			break
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func HandlePoolMonitor(evt *event.PoolEvent) {
	switch evt.Type {
	case event.PoolClosedEvent:
		log.Println("DB connection closed.")
		reconnect(Client)
	}
}

func Ping() bool {
	if err := Client.Ping(context.TODO(), nil); err != nil {
		return false
	}
	return true
}
