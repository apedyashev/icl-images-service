package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		Image: Image{},
	}
}

type Models struct {
	Image Image
}

type Image struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	PostId    uint      `bson:"postId,omitempty" json:"postId,omitempty"`
	ImageBody string    `bson:"imageBody,omitempty" json:"imageBody,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (l *Image) Insert(entry Image) error {
	collection := client.Database("images").Collection("images")

	_, err := collection.InsertOne(context.TODO(), Image{
		PostId:    entry.PostId,
		ImageBody: entry.ImageBody,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting image", err)
		return err
	}

	return nil
}

// []*LogEntry - slice of pointer to LogEntry
func (l *Image) All() ([]*Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("images").Collection("images")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Failed to find all images", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	// slice of pointer to LogEntry
	var logs []*Image
	for cursor.Next(ctx) {
		var item Image
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Failed to decode log into slice", err)
			return nil, err
		}
		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *Image) GetOne(id string) (*Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("images").Collection("images")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry Image
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (l *Image) FindByPostId(postId uint) ([]*Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("images").Collection("images")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	filter := bson.D{{"postId", postId}}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println("Failed to find all images", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	// slice of pointer to LogEntry
	var logs []*Image
	for cursor.Next(ctx) {
		var item Image
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Failed to decode log into slice", err)
			return nil, err
		}
		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *Image) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("images").Collection("images")
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (l *Image) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("images").Collection("images")
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"postId", l.PostId},
				{"imageBody", l.ImageBody},
				{"updated_at", time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
