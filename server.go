package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ComputationTime/finesse-api/graph"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

func main() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"math.church", "http://math.church", "https://math.church", "http://localhost:*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)


	port := "8000"

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Check against your desired domains here
				var allowed map[string]bool
				allowed["http://math.church"] = true
				allowed["https://math.church"] = true
				allowed["math.church"] = true
                return allowed[r.Host]
            },
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        },
    })

	router.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// serve create-react-app build folder
	router.Handle("/*", http.FileServer(http.Dir("../finesse-frontend/build")))

	log.Printf("Listening on http://localhost:%s...", port)

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

