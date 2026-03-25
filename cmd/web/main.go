package main

import (
	"log"
	"net/http"
	"os"

	"single-page-developer-portfolio/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	submissionsPath := os.Getenv("SUBMISSIONS_PATH")
	if submissionsPath == "" {
		submissionsPath = "data/submissions.jsonl"
	}

	app, err := handlers.NewApp("web/templates/index.html", "web/static", submissionsPath)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	log.Printf("listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
