package chat

import (
	components "checkpoint/internal/components/util"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

type Message struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /chat", chat)
	// mux.HandleFunc("POST /chat", postMessage)
}

func chat(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(ChatBox())
	if err != nil {
		components.InternalError(sse, w, err)
	}
}

// func postMessage(w http.ResponseWriter, r *http.Request) {
// 	store := &Message{}
// 	err := datastar.ReadSignals(r, store)
//
// 	sse := datastar.NewSSE(w, r)
//
// 	if err != nil {
// 		components.InternalError(sse, w, err)
// 		return
// 	}
//
// }
