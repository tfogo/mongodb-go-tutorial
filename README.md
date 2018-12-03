# MongoDB Go Driver Tutorial

This tutorial will help you get started with the official [MongoDB Go driver](https://github.com/mongodb/mongo-go-driver/). We will be coding a simple program to demonstrate how to:

- Install the MongoDB Go Driver
- Connect to MongoDB using the Go Driver
- Use BSON objects in Go
- Send CRUD operations to MongoDB

You can view the complete code for this tutorial on [this GitHub repository](https://github.com/tfogo/go-tutorial). In order to follow along, you will need a MongoDB database which you can connect. You can use a MongoDB database running locally, or easily create a free 500 MB database using [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).

## Installing the MongoDB Go Driver

The MongoDB Go driver is made up of several packages. Run the following command to install all the Go driver packages using `go get`:

```
go get github.com/mongodb/mongo-go-driver
```

If you are using the `dep` package manager, install the main `mongo` package as well as the `bson` package using this command:

```
dep ensure --add github.com/mongodb/mongo-go-driver/mongo \
github.com/mongodb/mongo-go-driver/bson \
github.com/mongodb/mongo-go-driver/options
```

Create the file `main.go` and import the `bson`, `mongo`, and `options` packages:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/options"
)

// We will be using this Person type later in the program
type Trainer struct {
	Name string 
	Age  int    
	City string
}

func main() {
    // Rest of the code will go here
}
```

This code also imports some standard libraries and defines a `Trainer` type. We will be using these later in the tutorial.


## Connect to MongoDB using the Go driver

Once the MongoDB driver has been imported, you can connect to a MongoDB deployment using the `mongo.Connect()` function. You must pass a context and connection string to `mongo.Connect()`. Optionally, you can also pass in an `options.ClientOptions` object as a third argument to configure driver settings such as write concerns, socket timeouts, and more. [The mongo/options package documentation](https://godoc.org/github.com/mongodb/mongo-go-driver/mongo/options) has more information about what client options are available.

```go
client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

if err != nil {
    log.Fatal(err)
} else {
    fmt.Println("Connected to MongoDB!")
}
```

Get a handle for the `trainers` collection in the `test` database using the following line of code:

```go
collection := client.Database("test").Collection("trainers")
```

We will use this collection handle to query the `trainers` collection.

It is best practice to keep a connection open to MongoDB so that the application can make use of connection pooling - you don't want to open and close a connection for each query. However, if your application no longer requires a connection, the connection can be closed with `client.Disconnect()`:

```go
err := client.Disconnect(context.TODO())

if err != nil {
    log.Fatal(err)
} else {
    fmt.Println("Connection to MongoDB closed.")
}
```

Run the code (`go run main.go`) to test that your program can successfully connect to your MongoDB server.

## Using BSON Objects in Go

Documents in MongoDB are stored as a type of binary-encoded JSON called BSON. The Go driver has two families of types for representing BSON data: The `D` types and the `Raw` types.

The `D` family of types is used to concisely build BSON objects using native Go types. This can be particularly useful for constructing commands passed to MongoDB. The `D` family consists of four types:

- `D`: A BSON document. This type should be used in situations where order matters, such as MongoDB commands.
- `M`: An unordered map. It is the same as `D`, except it does not preserve order.
- `A`: A BSON array.
- `E`: A single element inside a `D`.

Here is an example of a filter document built using `D` types which may be used to find documents where the `name` field matches either Alice or Bob:

```go
bson.D{
    {
        "name", 
        bson.D{
            {
                "$in", 
                bson.A{"Alice", "Bob"}
            }
        }
    }
}
```

The `Raw` family of types is used for validating a slice of bytes. You can also retrieve single elements from Raw types using a [`Lookup()`](https://godoc.org/github.com/mongodb/mongo-go-driver/bson#Raw.Lookup). This is useful if you don't want the overhead of having to unmarshall the BSON into another type.

## CRUD Operations

Once you have connceted to the database, it's time to start adding and manipulating some data. The `Collection` type has several methods which allow you to send queries to the database.

### Inserting documents

Create some new `Trainer` structs to insert into the database:

```go
ash := Trainer{"Ash", 10, "Pallet Town"}
misty := Trainer{"Misty", 10, "Cerulean City"}
brock := Trainer{"Brock", 15, "Pewter City"}
```

To insert a single document, use the `collection.InsertOne()` method:

```go
res, err = collection.InsertOne(context.Background(), ash)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Inserted a single document: ", insertResult.InsertedID)
```

To insert a multiple documents at a time, the `collection.InsertMany()` method will take an array of objects:

```go
trainers := []interface{}{misty, brock}

insertManyResult, err := collection.InsertMany(context.Background(), trainers)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
```

### Updating documents

The `collection.UpdateOne()` method allows you to update a single document. It requires a filter document to match documents in the database and an update document to describe the update operation. We can build these using `bson.D` types:

```go
filter := bson.D{{"name", "Ash"}}

update := bson.D{
    {"$inc", bson.D{
        {"age", 1},
    }},
}
```

This code will then match the document where the name is Ash and will increment Ash's age by 1 - happy birthday Ash!

```go
updateResult, err := collection.UpdateOne(context.Background(), filter, update)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
```

### Finding Documents

To find a document, you will need a filter document as well as a pointer to a value into which the result can be decoded. To find a single document, use `collection.FindOne()`. This method returns a single result which can be decoded into a value. We'll use the same `filter` variable we used in the update query to find a document where the name is Ash.

```go
result := &Trainer{}

err = collection.FindOne(context.Background(), filter).Decode(result)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %+v\n", *result)
```

To find multiple documents, use `collection.Find()`. This method returns a `Cursor`. A `Cursor` provides a stream of documents through which you can iterate and decode one at a time. Once a `Cursor` has been exhausted, you should close the `Cursor`. 

```go
options := options.Find()
options.SetLimit(2)

var results []*Trainer

// Finding multiple documents returns a cursor
cur, err := collection.Find(context.TODO(), nil, options)
if err != nil {
    log.Fatal(err)
}

// Iterate through the cursor
for cur.Next(context.TODO()) {
    elem := &Trainer{}
    err := cur.Decode(elem)
    if err != nil {
        log.Fatal(err)
    }

    results = append(results, elem)
}

cur.Close(context.TODO())

fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
if err := cur.Err(); err != nil {
    log.Fatal(err)
}
```

### Deleting Documents

Finally, you can delete documents using `collection.DeleteOne()` or `collection.DeleteMany()`. Here we pass `nil` as the filter argument, which will match all documents in the collection. You could also use [`collection.Drop()`](https://godoc.org/github.com/mongodb/mongo-go-driver/mongo#Collection.Drop) to delete an entire collection.

```go
deleteResult, err := collection.DeleteMany(context.TODO(), nil)
if err := cur.Err(); err != nil {
    log.Fatal(err)
}
fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
```   

## Next steps

You can view the final code from this tuoprial in [this GitHub repository](). Documentation for the MongoDB Go Driver is available on [GoDoc](https://godoc.org/github.com/mongodb/mongo-go-driver). You may be particularly interested in the documentation about using [aggregations](https://godoc.org/github.com/mongodb/mongo-go-driver/mongo#Collection.Aggregate) or [transactions](https://godoc.org/github.com/mongodb/mongo-go-driver/mongo#Session). 

If you have any questions, please get in touch in the [mongo-go-driver Google Group](https://groups.google.com/forum/#!forum/mongodb-go-driver). Please file any bug reports on the Go project in the [MongoDB JIRA](https://www.google.com/url?q=https%3A%2F%2Fjira.mongodb.org%2Fprojects%2FGODRIVER&sa=D&sntz=1&usg=AFQjCNEOEt6d3ZNOMKzmT23RYOVYdjSD6g). We would love your feedback on the Go Driver, so please get in touch with us to let us know your thoughts.