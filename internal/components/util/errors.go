package util

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

type dsSignals struct {
	ShowError    bool   `json:"show_error"`
	Error        string `json:"error"`
	ErrorDetails string `json:"error_details"`
}

func InternalError(sse *datastar.ServerSentEventGenerator, w http.ResponseWriter, err error) {
	slog.Error("Internal Server Error", "error", err)
	SignalErrorToPatch(sse, w, http.StatusInternalServerError, fmt.Errorf("an internal error has occurred"))
}

func BadRequest(sse *datastar.ServerSentEventGenerator, w http.ResponseWriter, err error) {
	slog.Info("Bad Request", "error", err)
	SignalErrorToPatch(sse, w, http.StatusBadRequest, err)
}

func SignalErrorToPatch(sse *datastar.ServerSentEventGenerator, w http.ResponseWriter, code int, err error) {
	// No need to do this here if nil
	if err == nil {
		slog.Error("Incorrect usage of this function")
		return
	}

	err = sse.ConsoleError(err)
	if err != nil {
		slog.Error("datastar error", "error", err)
		http.Error(w, "an internal framework error has occured, please refer to the server logs to troubleshoot", http.StatusInternalServerError)
	}

	err = sse.MarshalAndPatchSignals(dsSignals{
		ShowError:    true,
		Error:        http.StatusText(code),
		ErrorDetails: err.Error(),
	})
	if err != nil {
		slog.Error("datastar error", "error", err)
		http.Error(w, "an internal framework error has occured, please refer to the server logs to troubleshoot", http.StatusInternalServerError)
	}
}
