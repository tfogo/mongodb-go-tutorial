# MongoDB Go Driver Tutorial

This tutorial will help you get started with the official [MongoDB Go driver](). We will be coding a simple program to demonstrate how to:

- Install the MongoDB Go Driver
- Connect to MongoDB using the Go Driver
- Use BSON objects in Go
- Send CRUD operations to MongoDB

You can view the complete code for this tutorial on [this GitHub repository](). In order to follow along, you will need a MongoDB database which you can connect. You can use a MongoDB database running locally, or easily create a free 500 MB database using [MongoDB Atlas]().

## Installing the MongoDB Go Driver

Run the following command to install the Go driver using `go get`:

```
go get github.com/mongodb/mongo-go-driver
```

If you are using the `dep` package manager, install the Go driver using:

```
dep ensure --add github.com/mongodb/mongo-go-driver
```

The MongoDB Go driver is made up of several packages which will need to be imported. Create the file `main.go` and import the `bson` and `mongo` packages:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// We will be using this Person type later in the program
type Person struct {
	Name string 
	Age  int    
}

func main() {
    // Code goes here
}
```

This code also imports some standard libraries and defines a `Person` type. We will be using these later in the tutorial.


## Connect to MongoDB using the Go driver

Once the MongoDB driver has been imported, you can connect to a MongoDB deployment using the `mongo.Connect` function. You must pass a context and connection string to `mongo.Connect`. Optionally, you can also pass in an `options.ClientOptions` object as a third argument to configure driver settings such as write concerns, socket timeouts, and more. [The mongo/options package documentation](https://godoc.org/github.com/mongodb/mongo-go-driver/mongo/options) has more information about what client options are available.

```go
client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

if err != nil {
    log.Fatal(err)
} else {
    fmt.Println("Connected to MongoDB!")
}
```

Get a handle for the `people` collection in the `test` database using the following line of code:

```go
collection := client.Database("test").Collection("people")
```

We will use this collection handle to query the `people` collection.

It is best practice to keep a connection open to MongoDB so that the application can make use of connection pooling - you don't want to open and close a connection for each query. However, if your application no longer requires a connection, the connection can be closed with `client.Disconnect`:

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

The `Raw` family of types is used for validating a slice of bytes. You can also retreive single elements from Raw types using a [`Lookup()`](https://godoc.org/github.com/mongodb/mongo-go-driver/bson#Raw.Lookup). This is useful if you don't want the overhead of having to unmarshall the BSON into another type.


## CRUD Operations

### Inserting documents

Create a new `Person` struct to insert into the database:

```go
person := Person{"Ash Ketchum", 10}
```

Then use the `collection.InsertOne()` method to insert the `person` struct:

```go
res, err = collection.InsertOne(context.Background(), person)
if err != nil {
    log.Fatal(err)
}

id := res.InsertedID
fmt.Println("Inserted ID: ", id)
```

### Finding Documents

To find a document, you will need a filter document as well as a value into which the result can be decoded. To find a single document, use `collection.FindOne()`. This method returns a single result which can be decoded into a value.

```go
result := &Person{}
filter := bson.D{{"name", "Tim"}}

err = collection.FindOne(context.Background(), filter).Decode(result)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %+v\n", *result)
```

To find multiple documents, use `collection.Find()`. This method returns a `Cursor`. A `Cursor` provides a stream of documents through which you can iterate and decode one at a time. Once a `Cursor` has been exhausted, you should close the `Cursor`. 

```go
result := &Person{}
filter := bson.D{{"name", "Tim"}}

cursor, err = collection.Find(context.Background(), filter)
if err != nil {
    log.Fatal(err)
}


ctx := context.Background()
defer cur.Close(ctx)

for cur.Next(ctx) {
	elem := bson.NewDocument()
	if err := cur.Decode(elem); err != nil {
		log.Fatal(err)
	}
}

if err := cur.Err(); err != nil {
	log.Fatal(err)
}


fmt.Printf("Result: %+v\n", *result)
```


### Updating Documents

### Deleting Documents

## Next steps

You can. MongoDB Go Driver is available on [GoDoc. Questions can be asked through the mongo-go-driver Google Group and bug reports should be filed against the Go project in the MongoDB JIRA. Your feedback on the Go driver is greatly appreciated.