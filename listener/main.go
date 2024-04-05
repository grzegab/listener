package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListResponse struct {
	ID     interface{} `json:"id"`
	Method interface{} `json:"method"`
	Date   interface{} `json:"date"`
}

var mongoUrl, databaseName, collectionName string

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/listen", saveRequest)
	mux.HandleFunc("/list", readAll)
	mux.HandleFunc("/read/{id}", read)
	mux.HandleFunc("/remove", deleteAll)

	collectionName = "requests"
	databaseName = os.Getenv("MONGO_DATABASE")
	mongoUrl = os.Getenv("MONGO_URL")

	err := http.ListenAndServe(":80", mux)
	fmt.Printf("Listening...\n")

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func saveRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("got new request\n")
	header := bson.D{}
	for i, h := range r.Header {
		for _, v := range h {
			{
				header = append(header, bson.E{i, v})
			}
		}
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	request := bson.D{
		{"method", r.Method},
		{"url", r.URL.Host + " " + r.URL.Path},
		{"uri", r.RequestURI},
		{"host", r.Host},
		{"header", header},
		{"body", string(bodyBytes)},
		{"created", time.Now().Format("2006-01-02 15:04:05")},
	}

	ctx := context.TODO()
	c := openConnection(ctx)
	defer closeConnection(ctx, c)

	requestCollection := c.Database(databaseName).Collection(collectionName)
	result, err := requestCollection.InsertOne(ctx, request)
	if err != nil {
		panic(err)
	}

	log.Printf("inset into database: %s\n", result.InsertedID)
	io.WriteString(w, "OK\n")
}

func read(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	if idString == "" {
		io.WriteString(w, "{}")

		return
	}

	ctx := context.Background()
	c := openConnection(ctx)
	defer closeConnection(ctx, c)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	requestCollection := c.Database(databaseName).Collection(collectionName)

	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		panic(err)
	}

	filter := bson.M{"_id": objectId}
	cursor, err := requestCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
}

func readAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	c := openConnection(ctx)
	defer closeConnection(ctx, c)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	requestCollection := c.Database(databaseName).Collection(collectionName)

	filter := bson.D{{}}
	cursor, err := requestCollection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	var lp []ListResponse
	for _, result := range results {
		lr := ListResponse{
			ID:     result["_id"],
			Method: result["method"],
			Date:   result["created"],
		}

		lp = append(lp, lr)
	}

	err = json.NewEncoder(w).Encode(lp)
	if err != nil {
		panic(err)
	}
}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	c := openConnection(ctx)
	defer closeConnection(ctx, c)

	requestCollection := c.Database(databaseName).Collection(collectionName)
	result, err := requestCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		log.Printf("error deleting documents: %v\n", err)
	}

	deleteString := fmt.Sprintf("Database cleared! Deleted count: %d\n", result.DeletedCount)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, deleteString)
}

func openConnection(ctx context.Context) *mongo.Client {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}

	return c
}

func closeConnection(ctx context.Context, c *mongo.Client) {
	if err := c.Disconnect(ctx); err != nil {
		panic(err)
	}
}
