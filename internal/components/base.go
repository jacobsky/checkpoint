package components

import (
	"checkpoint/internal/components/util"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

const maincontent = "main"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /showchat", showchat)
	mux.HandleFunc("GET /showmessages", showmessages)
	mux.HandleFunc("GET /join", join)
	mux.Handle("/", templ.Handler(Landing()))
}

func join(w http.ResponseWriter, r *http.Request) {
	slog.Info("User joined")
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(MainContent())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}

func showchat(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(ChatContent())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}

func showmessages(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(MessageBoardContent())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}
