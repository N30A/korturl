package render

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func Component(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render component", "error", err)
	}
}
