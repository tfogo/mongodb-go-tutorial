package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Person struct {
	Name string `bson:personName`
	Age  int    `bson:personAge`
}

func main() {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}

	collection := client.Database("baz").Collection("qux")

	res, err := collection.InsertOne(context.Background(), map[string]string{"hello": "world"})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID

	fmt.Println("ID ", id)

	person1 := Person{"Tim", 25}

	fmt.Printf("%+v\n", person1)

	//person1bson, err := bson.Marshal(person1)

	res, err = collection.InsertOne(context.Background(), person1)
	if err != nil {
		log.Fatal(err)
	}
	id2 := res.InsertedID
	fmt.Println("ID", id2)

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		elem := &bson.D{}
		err := cur.Decode(elem)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", elem)

	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	result := &Person{}
	filter := map[string]string{"name": "Tim"}

	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result\n")
	fmt.Printf("%+v\n", result)

}
