package comments

import (
	sqlc "checkpoint/internal/db"
	"database/sql"
	"net/http"
)

type frontendSignals struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}
type Handler struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
