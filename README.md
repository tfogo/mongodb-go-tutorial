# MongoDB Go Driver Tutorial Part 1: Connecting, Using BSON, and CRUD Operations

The official MongoDB Go Driver [recently moved to GA](#) with the release of version 1.0.0. It's now regarded as feature complete and ready for production use. This tutorial will help you get started with the [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver/). You will create a simple program and learn how to:

- Install the MongoDB Go Driver
- Connect to MongoDB using the Go Driver
- Use BSON objects in Go
- Send CRUD operations to MongoDB

You can view the complete code for this tutorial on [this GitHub repository](https://github.com/tfogo/mongodb-go-tutorial). In order to follow along, you will need a MongoDB database to which you can connect. You can use a MongoDB database running locally, or easily create a free 500 MB database using [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).

## Install the MongoDB Go Driver

The MongoDB Go Driver is made up of several packages. If you are just using `go get`, you can install the driver using:

```
go get go.mongodb.org/mongo-driver
```

The output of this may look like a warning stating something like `package go.mongodb.org/mongo-driver: no Go files in (...)`. This is expected output from `go get`.

If you are using the [`dep`](https://golang.github.io/dep/docs/introduction.html) package manager, you can install the main `mongo` package as well as the `bson` and `mongo/options` package using this command:

```
dep ensure --add go.mongodb.org/mongo-driver/mongo \
go.mongodb.org/mongo-driver/bson \
go.mongodb.org/mongo-driver/mongo/options
```

If you are using [`go mod`](https://github.com/golang/go/wiki/Modules), the correct packages should be retrieved at build time.

## Create the wireframe

Create the file `main.go` and import the `bson`, `mongo`, and `mongo/options` packages:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
    Name string
    Age  int
    City string
}

func main() {
    // Rest of the code will go here
}
```

This code also imports some standard libraries and defines a `Trainer` type. You will be using these later in the tutorial.


## Connect to MongoDB using the Go Driver

Once the MongoDB Go Driver has been imported, you can connect to a MongoDB deployment using the `mongo.Connect()` function. You must pass a context and a `options.ClientOptions` object to `mongo.Connect()`. The client options are used to set the connection string. It can also be used to configure driver settings such as write concerns, socket timeouts, and more. [The options package documentation](https://godoc.org/go.mongodb.org/mongo-driver/mongo/options) has more information about what client options are available.

Add this code in the main function:

```go
// Set client options
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// Connect to MongoDB
client, err := mongo.Connect(context.TODO(), clientOptions)

if err != nil {
    log.Fatal(err)
}

// Check the connection
err = client.Ping(context.TODO(), nil)

if err != nil {
    log.Fatal(err)
}

fmt.Println("Connected to MongoDB!")
```

Once you have connected, you can now get a handle for the `trainers` collection in the `test` database by adding the following line of code at the end of the main function:

```go
collection := client.Database("test").Collection("trainers")
```

The following code will use this collection handle to query the `trainers` collection.

It is best practice to keep a client that is connected to MongoDB around so that the application can make use of connection pooling - you don't want to open and close a connection for each query. However, if your application no longer requires a connection, the connection can be closed with `client.Disconnect()` like so:

```go
err = client.Disconnect(context.TODO())

if err != nil {
    log.Fatal(err)
}
fmt.Println("Connection to MongoDB closed.")
```

Run the code (`go run main.go`) to test that your program can successfully connect to your MongoDB server. Go will complain about the unused `bson` and `mongo/options` packages and the unused `collection` variable, since we haven't done anything with them yet. You have to comment these out until they are used to make your program run and test the connection. 

## Use BSON Objects in Go

Before we start sending queries to the database, it's important to understand how the Go Driver works with BSON objects. JSON documents in MongoDB are stored in a binary representation called BSON (Binary-encoded JSON). Unlike other databases that store JSON data as simple strings and numbers, the BSON encoding extends the JSON representation to include additional types such as int, long, date, floating point, and decimal128. This makes it much easier for applications to reliably process, sort, and compare data. The Go Driver has two families of types for representing BSON data: The `D` types and the `Raw` types.

The `D` family of types is used to concisely build BSON objects using native Go types. This can be particularly useful for constructing commands passed to MongoDB. The `D` family consists of four types:

- `D`: A BSON document. This type should be used in situations where order matters, such as MongoDB commands.
- `M`: An unordered map. It is the same as `D`, except it does not preserve order.
- `A`: A BSON array.
- `E`: A single element inside a `D`.

Here is an example of a filter document built using `D` types which may be used to find documents where the `name` field matches either Alice or Bob:

```go
bson.D{{
    "name", 
    bson.D{{
        "$in", 
        bson.A{"Alice", "Bob"}
    }}
}}
```

The `Raw` family of types is used for validating a slice of bytes. You can also retrieve single elements from Raw types using a [`Lookup()`](https://godoc.org/go.mongodb.org/mongo-driver/bson#Raw.Lookup). This is useful if you don't want the overhead of having to unmarshall the BSON into another type. This tutorial will just use the `D` family of types.

## CRUD Operations

Once you have connected to the database, it's time to start adding and manipulating some data. The `Collection` type has several methods which allow you to send queries to the database.

### Insert documents

First, create some new `Trainer` structs to insert into the database:

```go
ash := Trainer{"Ash", 10, "Pallet Town"}
misty := Trainer{"Misty", 10, "Cerulean City"}
brock := Trainer{"Brock", 15, "Pewter City"}
```

To insert a single document, use the `collection.InsertOne()` method:

```go
insertResult, err := collection.InsertOne(context.TODO(), ash)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Inserted a single document: ", insertResult.InsertedID)
```

To insert multiple documents at a time, the `collection.InsertMany()` method will take a slice of objects:

```go
trainers := []interface{}{misty, brock}

insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
```

### Update documents

The `collection.UpdateOne()` method allows you to update a single document. It requires a filter document to match documents in the database and an update document to describe the update operation. You can build these using `bson.D` types:

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
updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
```

### Find documents

To find a document, you will need a filter document as well as a pointer to a value into which the result can be decoded. To find a single document, use `collection.FindOne()`. This method returns a single result which can be decoded into a value. You'll use the same `filter` variable you used in the update query to match a document where the name is Ash.

```go
// create a value into which the result can be decoded
var result Trainer

err = collection.FindOne(context.TODO(), filter).Decode(&result)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found a single document: %+v\n", result)
```

To find multiple documents, use `collection.Find()`. This method returns a `Cursor`. A `Cursor` provides a stream of documents through which you can iterate and decode one at a time. Once a `Cursor` has been exhausted, you should close the `Cursor`. Here you'll also set some options on the operation using the `options` package. Specifically, you'll set a limit so only 2 documents are returned.

```go
// Pass these options to the Find method
findOptions := options.Find()
findOptions.SetLimit(2)

// Here's an array in which you can store the decoded documents
var results []*Trainer

// Passing bson.D{{}} as the filter matches all documents in the collection
cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
if err != nil {
    log.Fatal(err)
}

// Finding multiple documents returns a cursor
// Iterating through the cursor allows us to decode documents one at a time
for cur.Next(context.TODO()) {
    
    // create a value into which the single document can be decoded
    var elem Trainer
    err := cur.Decode(&elem)
    if err != nil {
        log.Fatal(err)
    }

    results = append(results, &elem)
}

if err := cur.Err(); err != nil {
    log.Fatal(err)
}

// Close the cursor once finished
cur.Close(context.TODO())

fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
```

### Delete Documents

Finally, you can delete documents using `collection.DeleteOne()` or `collection.DeleteMany()`. Here you pass `bson.D{{}}` as the filter argument, which will match all documents in the collection. You could also use [`collection.Drop()`](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Drop) to delete an entire collection.

```go
deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
```   

## Next steps

You can view the final code from this tutorial in [this GitHub repository](https://github.com/tfogo/mongodb-go-tutorial). Documentation for the MongoDB Go Driver is available on [GoDoc](https://godoc.org/go.mongodb.org/mongo-driver). You may be particularly interested in the documentation about using [aggregations](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Aggregate) or [transactions](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Session). 

If you have any questions, please get in touch in the [mongo-go-driver Google Group](https://groups.google.com/forum/#!forum/mongodb-go-driver). Please file any bug reports on the Go project in the [MongoDB JIRA](https://www.google.com/url?q=https%3A%2F%2Fjira.mongodb.org%2Fprojects%2FGODRIVER&sa=D&sntz=1&usg=AFQjCNEOEt6d3ZNOMKzmT23RYOVYdjSD6g). We would love your feedback on the Go Driver, so please get in touch with us to let us know your thoughts.

This is part one of several tutorials on the Go Driver. The next part will be:

- _MongoDB Go Driver Tutorial Part 2: BSON Tags and Primitive Types_

So stay tuned for more parts in the coming weeks.