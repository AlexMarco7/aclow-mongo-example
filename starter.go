package main

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/AlexMarco7/aclow"
	"github.com/AlexMarco7/aclow-mongo-example/actions"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {

	startOpt := aclow.StartOptions{
		Debug:         true,
		Host:          "localhost",
		Port:          4222,
		ClusterPort:   8222,
		ClusterRoutes: []*url.URL{
			//&url.URL{Host: fmt.Sprintf("localhost:%d", 8223)},
		},
	}

	var app = &aclow.App{}

	app.Start(startOpt)

	connectOnMongo(app)

	app.RegisterModule("mongo", []aclow.Node{
		&actions.Execute{},
	})

	time.Sleep(time.Second * 2)
	app.Publish("mongo@execute", aclow.Message{})

	app.Wait()
}

func connectOnMongo(app *aclow.App) {
	log.Println("connecting on mongo...")
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	app.Config["db_name"] = "db_name"

	if err != nil {
		time.Sleep(time.Second * 1)
		log.Println(err.Error())
		log.Println("trying again...")
		connectOnMongo(app)
	} else {
		app.Resources["db"] = client
	}
}
