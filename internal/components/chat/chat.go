package chat

import (
	components "checkpoint/internal/components/util"
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
	tx            chan Message
	rx            []chan Message
	addRx         chan chan Message
	delRx         chan (<-chan Message)
}

func newHandler() *handler {
	h := &handler{
		message_store: []Message{},
		tx:            make(chan Message, channelBuffer),
		rx:            make([]chan Message, 0),
		addRx:         make(chan chan Message, channelBuffer),
		delRx:         make(chan (<-chan Message), channelBuffer),
	}
	go h.serve()
	return h
}

func (h *handler) serve() {
	slog.Info("Updater worker started")
	for {
		select {
		case msg := <-h.tx:
			slog.Info("Sending message", "message", msg, "active connections", len(h.rx))
			for i, rx := range h.rx {
				slog.Debug("broadcasting", "rxid", i)
				rx <- msg
			}
		case channel := <-h.addRx:
			slog.Debug("Opening channel", "channel", channel)
			h.rx = append(h.rx, channel)
		case channel := <-h.delRx:
			slog.Debug("Closing channel", "channel", channel)
			for i, ch := range h.rx {
				if ch == channel {
					h.rx[i] = h.rx[len(h.rx)-1]
					h.rx = h.rx[:len(h.rx)-1]
					close(ch)
					break

				}
			}
		}
	}
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
		h.message_store = append(h.message_store[1:], message)
	} else {
		h.message_store = append(h.message_store, message)
	}
	h.tx <- message
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

// Note: Due to the way that the handlers work, when there is a disconnection, the entire state of the chat history is sent to the chat.
// This site is intended to be an ephemeral chat with minimal history and only holds about 50 most recent messages.
// It's a checkpoint, not a hangout spot.
// There is a board for leaving a messages in internal/components/comments
func (h *handler) chat(w http.ResponseWriter, r *http.Request) {
	var store Message
	err := datastar.ReadSignals(r, &store)
	if err != nil {
		slog.Error("datastar error occurred", "error", err)
	}

	slog.Info("Chat Connected", "user", store.Nickname)
	sse := datastar.NewSSE(w, r)

	err = sse.PatchElementTempl(ChatBoxMessages(h.message_store))

	if err != nil {
		components.InternalError(sse, w, err)
	}
	listener := make(chan Message)
	h.addRx <- listener
	// Keep the context open until the connection closes (detectable via the request context)
	for {
		select {
		case <-sse.Context().Done():
			slog.Info("Chat disconnected", "user", store.Nickname)
			h.delRx <- listener
			return
		case msg := <-listener:
			err := sse.PatchElementTempl(ChatMessage(msg.Nickname, msg.Message), datastar.WithSelectorID(chatselector), datastar.WithModeAppend())
			if err != nil {
				slog.Error("Error occurred when patching", "error", err)
			}
		}
	}
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
	// Do something to indicate that there is a new message
	store.Message = ""

	err = sse.MarshalAndPatchSignals(&store)
	if err != nil {
		slog.Error("Patch Element Error", "message", err)
	}
}
