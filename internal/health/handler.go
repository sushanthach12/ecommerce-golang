package health

import (
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/httpx"
)

type handler struct{}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) Health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, "All Good")
}
