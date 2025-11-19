package components

import (
	"checkpoint/internal/components/util"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /showchat", showchat)
	mux.HandleFunc("GET /showmessages", showmessages)
	mux.HandleFunc("GET /join", join)
	mux.Handle("/", templ.Handler(Loading()))
	mux.HandleFunc("GET /intro", intro)
}

var introduction_text = []string{
	"Greetings weary traveler.",
	"What brought you here?",
	"Well, that's not important.",
	"Regardless, please take a rest.",
	"From where ever you hail.",
	"Whatever your story",
	"Welcome to the checkpoint",
	"I hope you enjoy your stay",
}

func intro(w http.ResponseWriter, r *http.Request) {
	// Play introduction
	// 1. Show an empty page
	// 2. Fade in a description of the checkpoint
	// 3. Fade in the header
	// 4, Fade in the footer
	sse := datastar.NewSSE(w, r)

	err := sse.PatchElementTempl(IntroContainer())
	if err != nil {
		util.InternalError(sse, w, err)
	}

	for _, text := range introduction_text {
		err := sse.PatchElementTempl(Intro(text), datastar.WithSelectorID("introcontainer"), datastar.WithModeAppend())
		if err != nil {
			util.InternalError(sse, w, err)
		}
		time.Sleep(time.Second * 3)
	}
	// TODO: Implement a skip intro
	// Skipping should likely be a second endpoint getting called from the frontend
	err = sse.PatchSignals([]byte(`{introcomplete: 'true'}`))
	if err != nil {
		util.InternalError(sse, w, err)
	}

	time.Sleep(time.Second * 3)
	err = sse.PatchElementTempl(LandingFrag())
	if err != nil {
		util.InternalError(sse, w, err)
	}

}

func join(w http.ResponseWriter, r *http.Request) {
	slog.Info("User joined")
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(MainContentFrag(), datastar.WithSelectorID("main"), datastar.WithModeInner())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}

func showchat(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	// Transitino into the chat content
	err := sse.PatchElementTempl(ChatContent())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}

func showmessages(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	// Transitino into the chat content
	err := sse.PatchElementTempl(MessageBoardContent())
	if err != nil {
		util.InternalError(sse, w, err)
	}
}
