package handlers

import (
	"net/http"
	"strings"

	"campus_connect_api/internal/httpx"
	"campus_connect_api/internal/models"
)

// Health liveness/readiness simples.
func Health(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Discover GET /discover?filter=all|internships|events|groups|projects
func Discover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}
	_ = filter // TODO: aplicar filtro quando houver fonte de dados
	httpx.JSON(w, http.StatusOK, []models.DiscoverItem{})
}

// OpportunitiesList GET /opportunities?q=...
func OpportunitiesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	_ = r.URL.Query().Get("q")
	httpx.JSON(w, http.StatusOK, []models.Opportunity{})
}

// OpportunityByID GET /opportunities/{id}
func OpportunityByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/opportunities/")
	if id == "" || id == "opportunities" {
		httpx.Error(w, http.StatusBadRequest, "invalid_id", "missing opportunity id")
		return
	}
	httpx.Error(w, http.StatusNotFound, "not_found", "opportunity not found")
}

// EventsList GET /events
func EventsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	httpx.JSON(w, http.StatusOK, []models.CampusEvent{})
}

// EventByID GET /events/{id}
func EventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/events/")
	if id == "" || id == "events" {
		httpx.Error(w, http.StatusBadRequest, "invalid_id", "missing event id")
		return
	}
	httpx.Error(w, http.StatusNotFound, "not_found", "event not found")
}

// GroupsList GET /groups
func GroupsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	httpx.JSON(w, http.StatusOK, []models.StudyGroup{})
}

// GroupByID GET /groups/{id} (opcional no app; stub 404 até o domínio existir).
func GroupByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/groups/")
	id = strings.Trim(id, "/")
	if id == "" {
		httpx.Error(w, http.StatusBadRequest, "invalid_id", "missing group id")
		return
	}
	httpx.Error(w, http.StatusNotFound, "not_found", "group not found")
}

// GroupJoin POST /groups/{id}/join
func GroupJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/groups/")
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) < 2 || parts[1] != "join" {
		httpx.Error(w, http.StatusBadRequest, "invalid_path", "expected /groups/{id}/join")
		return
	}
	_ = parts[0]
	httpx.Error(w, http.StatusNotImplemented, "not_implemented", "join requires auth and persistence")
}

// Me GET /me ou /users/me
func Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	httpx.Error(w, http.StatusUnauthorized, "unauthorized", "authentication required")
}

// Apply POST /opportunities/{id}/applications
func Apply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	const prefix = "/opportunities/"
	if !strings.HasPrefix(r.URL.Path, prefix) || !strings.HasSuffix(r.URL.Path, "/applications") {
		httpx.Error(w, http.StatusBadRequest, "invalid_path", "expected /opportunities/{id}/applications")
		return
	}
	mid := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, prefix), "/applications")
	mid = strings.Trim(mid, "/")
	if mid == "" {
		httpx.Error(w, http.StatusBadRequest, "invalid_id", "missing opportunity id")
		return
	}
	httpx.Error(w, http.StatusNotImplemented, "not_implemented", "application flow not wired yet")
}
