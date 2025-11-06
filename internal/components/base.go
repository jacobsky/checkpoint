package components

import (
	"checkpoint/internal/components/util"
	"log/slog"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

const maincontent = "main"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /join", join)
}

func join(w http.ResponseWriter, r *http.Request) {
	slog.Info("User joined")
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(SplitLayout(), datastar.WithSelectorID(maincontent))
	if err != nil {
		util.InternalError(sse, w, err)
	}
}
