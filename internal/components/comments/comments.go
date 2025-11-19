package comments

import (
	"checkpoint/internal/components/util"
	sqlc "checkpoint/internal/db"
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

const (
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
		switch r.Header.Get("Datastar-Request") {
		case "true":
			query_params := r.URL.Query()
			// Special case because this is the first load. We need to find the latest comment to view.
			next := query_params.Get("next")
			// If this is the first load, we will use the first functionality to get the first message
			if next == "" {
				h.first(w, r)
			} else {
				id, err := strconv.Atoi(next)
				if err != nil {
					slog.Error("messages list", "error", err)
				}
				h.load(w, r, int64(id))
			}
		default:
			templ.Handler(MessageBoardFull()).ServeHTTP(w, r)
		}
	case http.MethodPost:
		h.post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) first(w http.ResponseWriter, r *http.Request) {
	slog.Info("load")
	sse := datastar.NewSSE(w, r)
	offset := 0
	sqlcParams := sqlc.GetRecentCommentsParams{
		Limit:  1,
		Offset: int64(offset),
	}
	comments, err := h.queries.GetRecentComments(r.Context(), sqlcParams)
	slog.Info("load", "comments", comments)
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
	err = sse.PatchElementTempl(MessagePost(comments[0]), datastar.WithModeAppend(), datastar.WithSelectorID("messageboard"))
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
}

func (h *Handler) load(w http.ResponseWriter, r *http.Request, id int64) {
	slog.Info("load")
	sse := datastar.NewSSE(w, r)
	comment, err := h.queries.GetCommentByID(r.Context(), id)
	slog.Info("load", "comment", comment)
	if err != nil {
		if err == sql.ErrNoRows {
			sse.PatchElementTempl(EndOfMessages(), datastar.WithModeAppend(), datastar.WithSelectorID("messageboard"))
			return
		}
		util.InternalError(sse, w, err)
		return
	}
	err = sse.PatchElementTempl(MessagePost(comment), datastar.WithModeAppend(), datastar.WithSelectorID("messageboard"))
	if err != nil {
		util.InternalError(sse, w, err)
		return
	}
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
	err = sse.PatchSignals([]byte("{message_left: true, show_postmessage: false}"))
	if err != nil {
		util.InternalError(sse, w, err)
	}
	err = sse.PatchElementTempl(MessagePost(comment), datastar.WithModePrepend(), datastar.WithSelectorID("messageboard"))
	if err != nil {
		util.InternalError(sse, w, err)
	}
}
