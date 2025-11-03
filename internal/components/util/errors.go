package components

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

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
		slog.Warn("Incorrect usage of this function")
		return
	}

	err = sse.PatchElementTempl(ErrorCard(code, err.Error()))
	if err != nil {
		slog.Error("datastar error", "error", err)
		http.Error(w, "an internal error has occured, please refer to the server logs to troubleshoot", http.StatusInternalServerError)
	}
}
