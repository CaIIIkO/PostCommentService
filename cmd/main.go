package main

import (
	"log"
	"net/http"
	"os"
	"postcommentservice/graph"
	"postcommentservice/internal/config"
	"postcommentservice/internal/db"
	"postcommentservice/internal/gateway"
	in_memory "postcommentservice/internal/gateway/inmemory"
	"postcommentservice/internal/gateway/postgres"
	"postcommentservice/internal/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

func main() {

	//reading the configuration file
	envFile := ".env"
	if len(os.Args) >= 2 {
		envFile = os.Args[1]
	}

	if err := config.InitConfig(envFile); err != nil {
		log.Fatal(err)
	}

	//get the settings db from the configuration file
	options := db.PostgresOptions{
		Name:     os.Getenv("POSTGRES_DBNAME"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	//connect db
	postgresDb, err := db.NewPostgresDB(options)

	if err != nil {
		log.Fatal(err)
	}

	var gateways *gateway.Gateway

	//choosing the data storage method
	if os.Getenv("USE_IN_MEMORY") == "true" {
		posts := in_memory.NewPostMemory()
		comments := in_memory.NewCommentsMemory()
		gateways = gateway.NewGateway(posts, comments)
	} else {
		posts := postgres.NewPostPostgres(postgresDb)
		comments := postgres.NewCommentsPostgres(postgresDb)
		gateways = gateway.NewGateway(posts, comments)
	}

	services := service.NewService(gateways)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		PostsService:        services.Post,
		CommentsService:     services.Comment,
		SubscriptionService: service.NewCommentSubscription(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
