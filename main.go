package main

import (
	"context"
	"fmt"
	"icl-images-service/data"
	"log"
	"net"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// as usual, when working with a context, we defer cancel()
	defer cancel()

	// close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	app.rpcListen()
	fmt.Println("After listen")
	// log.Println("Starting server on port", webPort)
	// srv := &http.Server{
	// 	Addr:    ":" + webPort,
	// 	Handler: app.routes(),
	// }

	// err = srv.ListenAndServe()
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}
	clientOptions := options.Client().ApplyURI(mongoUrl).SetMonitor(cmdMonitor)
	clientOptions.SetAuth(options.Credential{
		// TODO: env vars
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo: ", err)
		return nil, err
	}

	log.Println("Connected to mongo!")
	return client, nil
}
