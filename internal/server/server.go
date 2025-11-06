package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

type Server struct {
	port int
	db   *sql.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	conn, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_ADDRESS"))
	if err != nil {
		log.Panicf("Database could not be opened %v", err)
	}
	NewServer := &Server{
		port: port,
		db:   conn,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
		// IdleTimeout:  time.Minute,
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 30 * time.Second,
	}

	return server
}
