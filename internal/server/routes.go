package server

import (
	"net/http"

	"checkpoint/internal/components"
	"checkpoint/internal/components/chat"
	"checkpoint/internal/components/comments"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	// TODO add in health check
	// mux.HandleFunc("/health", s.healthHandler)
	components.RegisterRoutes(mux)
	chat.RegisterRoutes(mux)
	mux.Handle("/messages", comments.NewHandler(s.db))

	fileServer := http.FileServer(http.FS(Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/", templ.Handler(components.Layout()))
	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
