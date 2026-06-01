package main

import (
	"log"
	"net/http"
	"os"

	"blog/go-backend/database"
	"blog/go-backend/handlers"
)

func main() {
    if err := database.InitDB(); err != nil {
        log.Printf("WARNING: Database unavailable (%v) — running without DB", err)
    } else {
        defer database.DB.Close()
    }

    mux := http.NewServeMux()
    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "Blog Backend (Go) is running! ✅", "status": "ok"}`))
    })
    mux.HandleFunc("GET /api/posts", handlers.GetPosts)
    mux.HandleFunc("GET /api/posts/{id}", handlers.GetPost)
    mux.HandleFunc("POST /api/posts", handlers.CreatePost)
    mux.HandleFunc("PUT /api/posts/{id}", handlers.UpdatePost)
    mux.HandleFunc("DELETE /api/posts/{id}", handlers.DeletePost)
    mux.HandleFunc("GET /api/md/{slug}", handlers.GetMDFromGitHub)

    handler := handlers.CORSMiddleware(mux)

    log.Println("Server listening on :8080")
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Fatal(http.ListenAndServe(":"+port, handler))
}