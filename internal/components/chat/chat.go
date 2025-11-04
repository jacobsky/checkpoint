package chat

import (
	components "checkpoint/internal/components/util"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

const (
	maincontent         = "main"
	chatselector        = "chatboard"
	channelBuffer       = 10
	maxRetainedMessages = 10
)

type Message struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

type handler struct {
	// Note probably need some kind of mutex lock or channel for modifying things internally
	message_store []Message
	connections   map[context.Context]*datastar.ServerSentEventGenerator
	tx            chan Message
}

func newHandler() *handler {
	tx := make(chan Message, channelBuffer)
	h := &handler{
		message_store: []Message{},
		connections:   make(map[context.Context]*datastar.ServerSentEventGenerator),
		tx:            tx,
	}
	go h.sendUpdates(tx)
	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.chat(w, r)
	case http.MethodPost:
		h.postMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) pushMessage(message Message) {
	if messageLength := len(h.message_store); messageLength > maxRetainedMessages {
		h.message_store = append(h.message_store[1:messageLength-2], message)
	} else {
		h.message_store = append(h.message_store, message)
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /join", join)
	mux.Handle("/chat", newHandler())
}

func join(w http.ResponseWriter, r *http.Request) {
	slog.Info("User joined")
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(ChatBox(), datastar.WithSelectorID(maincontent))
	if err != nil {
		components.InternalError(sse, w, err)
	}
}

// TODO: Find a way to use channels to better sequentially update the state
func (h *handler) chat(w http.ResponseWriter, r *http.Request) {
	var store Message
	err := datastar.ReadSignals(r, &store)
	if err != nil {
		slog.Error("datastar error occurred", "error", err)
	}

	slog.Info("Chat Connected", "user", store.Nickname)
	sse := datastar.NewSSE(w, r)
	err = sse.PatchElementTempl(ChatBoxMessages(h.message_store))
	h.connections[sse.Context()] = sse

	if err != nil {
		components.InternalError(sse, w, err)
	}
	// Keep the context open until the connection closes (detectable via the request context)
	<-sse.Context().Done()
	slog.Info("Chat disconnected", "user", store.Nickname)
	delete(h.connections, sse.Context())
}

// Post a new message that will be polled
func (h *handler) postMessage(w http.ResponseWriter, r *http.Request) {
	store := &Message{}
	err := datastar.ReadSignals(r, store)

	sse := datastar.NewSSE(w, r)
	if err != nil {
		components.InternalError(sse, w, err)
		return
	}

	message := Message{store.Nickname, store.Message}
	h.pushMessage(message)
	h.tx <- message
	// Do something to indicate that there is a new message
	store.Message = ""
	signals, err := json.Marshal(store)
	if err != nil {
		slog.Error("Json marshall error", "error", err)
	}

	err = sse.PatchSignals(signals)
	if err != nil {
		slog.Error("Patch Element Error", "message", err)
	}
}

func (h *handler) sendUpdates(messages chan Message) {
	slog.Info("Updater worker started")
	for {
		msg := <-messages
		slog.Info("Sending message", "message", msg, "active connections", len(h.connections))
		for i, rx := range h.connections {
			slog.Debug("Connection Update", "count", i)
			// skip if nil or closed
			if rx == nil {
				continue
			}
			if rx.IsClosed() {
				slog.Info("Rx is closed skipping")
				// TODO: prune the connection (or perhaps make a separate goroutine for pruning closed connections)
				continue
			}
			err := rx.PatchElementTempl(ChatMessage(msg.Nickname, msg.Message), datastar.WithSelectorID(chatselector), datastar.WithModeAppend())
			if err != nil {
				slog.Error("Error occurred when patching", "error", err)
			}
		}
	}
}
