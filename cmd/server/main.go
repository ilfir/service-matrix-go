package main

import (
	"log"
	"net/http"
	"os"

	"service-matrix-go/internal/api/handlers"
	"service-matrix-go/internal/core/services"
	"service-matrix-go/internal/infrastructure/storage"
)

func main() {
	// Initialize dependencies
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileHelper := storage.NewFileHelper(cwd)
	wordService := services.NewWordService(fileHelper)
	httpHandlers := handlers.NewHTTPHandlers(wordService)

	// Router setup
	mux := http.NewServeMux()

	// Register endpoints (matching C# routes approx)
	// C# Route: [Route("[controller]")] -> /Words
	// Actions: [HttpPost("Search")] -> /Words/Search

	mux.HandleFunc("/Words/Search", httpHandlers.Search)
	mux.HandleFunc("/Words/Update", httpHandlers.Update)
	mux.HandleFunc("/Words/List", httpHandlers.GetList)
	mux.HandleFunc("/Words/Merge", httpHandlers.MergeWords)
	mux.HandleFunc("/Words/CleanMerge", httpHandlers.CleanMerge)
	mux.HandleFunc("/Words/LookupWord", httpHandlers.LookupWord)

	// Add CORS middleware if needed (found in C# Program.cs)
	handler := corsMiddleware(mux)

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
