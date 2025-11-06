package comments

import (
	"checkpoint/internal/components/util"
	sqlc "checkpoint/internal/db"
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

const (
	scrollByLimit    = 10
	minMessageLength = 8
)

type frontendSignals struct {
	Nickname    string `json:"nickname"`
	WallMessage string `json:"wall_message"`
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
	switch r.Method {
	case http.MethodGet:
		query_params := r.URL.Query()
		// Special case because this is the first load.
		var offset = 0
		if p := query_params.Get("offset"); p != "" {
			val, err := strconv.Atoi(p)
			if err != nil {
				slog.Error("messages list", "error", err)
			}
			offset = val
		}
		if offset == 0 {
			h.list(w, r)
		} else {
			h.load(w, r, offset)
		}
	case http.MethodPost:
		h.post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	slog.Info("load")
	sse := datastar.NewSSE(w, r)
	offset := 0
	sqlcParams := sqlc.GetRecentCommentsParams{
		Limit:  scrollByLimit,
		Offset: int64(offset),
	}
	comments, err := h.queries.GetRecentComments(r.Context(), sqlcParams)
	slog.Info("load", "comments", comments)
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
	err = sse.PatchElementTempl(MessagePostGroup(comments), datastar.WithModeAppend(), datastar.WithSelectorID("messageboard"))
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
}

func (h *Handler) load(w http.ResponseWriter, r *http.Request, offset int) {
	slog.Error("Load handler not implemented yet. This is a bug")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var store frontendSignals
	err := datastar.ReadSignals(r, &store)
	sse := datastar.NewSSE(w, r)
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}

	params := sqlc.AddCommentParams{
		Postdate: time.Now(),
		Poster:   store.Nickname,
		Message:  store.WallMessage,
	}
	if len(params.Message) < minMessageLength {
		params.Message = "Was here"
	}
	comment, err := h.queries.AddComment(r.Context(), params)
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
	err = sse.PatchElementTempl(MessagePost(comment), datastar.WithModePrepend(), datastar.WithSelectorID("messageboard"))
	if err != nil {
		util.InternalError(sse, w, err)
	}
}
