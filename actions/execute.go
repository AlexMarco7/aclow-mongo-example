package actions

import (
	"context"
	"log"

	"github.com/AlexMarco7/aclow"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type Execute struct {
	app *aclow.App
}

func (t *Execute) Address() []string { return []string{"execute"} }

func (t *Execute) Start(app *aclow.App) {
	t.app = app
}

func (t *Execute) Execute(msg aclow.Message, call aclow.Caller) (aclow.Message, error) {
	client := t.app.Resources["db"].(*mongo.Client)
	db := client.Database(t.app.Config["db_name"].(string))
	cmdResult := db.RunCommand(context.TODO(), command(), &options.RunCmdOptions{})

	var result interface{}

	if cmdResult.Err() != nil {
		log.Println(cmdResult.Err())
		return aclow.Message{}, cmdResult.Err()
	}

	cmdResult.Decode(&result)

	log.Println(result)

	return aclow.Message{Body: result}, nil
}

func command() bson.D {
	return bson.D{
		{"aggregate", "my_coll"},
		{"cursor", bson.M{"batchSize": 1000}},
		{"pipeline", bson.A{
			bson.M{"$match": bson.M{"a": 1}},
		}},
	}
}
