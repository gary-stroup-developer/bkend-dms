package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gary-stroup-developer/bkend-dms/driver"
	"github.com/gary-stroup-developer/bkend-dms/handlers"
	"github.com/gary-stroup-developer/bkend-dms/routes"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	//load the environment variables
	godotenv.Load(".env")
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")
	host := os.Getenv("MONGO_HOST")
	database := os.Getenv("MONGO_DB")

	//create a new server
	srv := http.NewServeMux()
	handler := cors.Default().Handler(srv)

	//connect to the database
	client, err := driver.ConnectDB(user, pass, host)
	if err != nil {
		log.Fatalln(err)
	}

	//use this context to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close the connection to Mongo when application exits
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	//create a variable to hold the mongo DB *mongo.Database
	mongodb := client.Database(database)

	//create a new repo that holds the database
	repo := handlers.CreateRepo(mongodb)
	handlers.NewRepo(repo)

	//pass the server to the routes func to handle the routes
	routes.Routes(srv)
	log.Fatalln(http.ListenAndServe(":8080", handler))
}
