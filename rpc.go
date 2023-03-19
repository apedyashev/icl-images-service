package main

import (
	"context"
	"errors"
	"fmt"
	"icl-images-service/data"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DB_NAME string = "logs"
var COLLECTION_NAME string = "logs"

type RPCServer struct{}

type SaveImagePayload struct {
	PostId    uint
	ImageBody string
}

func (srv *RPCServer) SaveImage(payload SaveImagePayload, resp *data.Image) error {
	fmt.Println("SaveImage is called", payload)
	collection := client.Database(DB_NAME).Collection(COLLECTION_NAME)
	result, err := collection.InsertOne(
		context.TODO(),
		data.Image{
			PostId:    payload.PostId,
			ImageBody: payload.ImageBody,
		},
	)
	if err != nil {
		log.Println("error writing to mongodb: ", err)
		return err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("unable to convert primitive.ObjectID to string")
	}
	fmt.Println("created entry", oid.String())
	*resp = data.Image{
		ID: oid.Hex(),
	}
	return nil
}

type GetPostImagesPayload struct {
	PostId uint
}
type GetPostImageResponse struct {
	Images []*data.Image
}

func (srv *RPCServer) GetPostImages(payload GetPostImagesPayload, response *GetPostImageResponse) error {
	fmt.Println("GetPostImages is called", payload)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database(DB_NAME).Collection(COLLECTION_NAME)

	// opts := options.Find()
	// opts.SetSort(bson.D{{"created_at", -1}})

	filter := bson.D{} // bson.D{{"postId", payload.PostId}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Failed to find all images", err)
		return err
	}

	defer cursor.Close(ctx)

	// slice of pointer to LogEntry
	var images []*data.Image
	for cursor.Next(ctx) {
		var item data.Image
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Failed to decode log into slice", err)
			return err
		}
		log.Println("item", item)
		images = append(images, &item)
	}
	fmt.Printf("images %+v\n", images)
	response.Images = images

	return nil
}
