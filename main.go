package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"campus_connect_api/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)

	mux.HandleFunc("/discover", handlers.Discover)

	mux.HandleFunc("/opportunities/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/opportunities/" || r.URL.Path == "/opportunities" {
			if r.Method == http.MethodGet {
				handlers.OpportunitiesList(w, r)
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/opportunities/") && strings.HasSuffix(r.URL.Path, "/applications") {
			handlers.Apply(w, r)
			return
		}
		handlers.OpportunityByID(w, r)
	})
	mux.HandleFunc("/opportunities", handlers.OpportunitiesList)

	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/events/" || r.URL.Path == "/events" {
			handlers.EventsList(w, r)
			return
		}
		handlers.EventByID(w, r)
	})
	mux.HandleFunc("/events", handlers.EventsList)

	mux.HandleFunc("/groups/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimSuffix(r.URL.Path, "/")
		if strings.HasSuffix(p, "/join") {
			handlers.GroupJoin(w, r)
			return
		}
		handlers.GroupByID(w, r)
	})
	mux.HandleFunc("/groups", handlers.GroupsList)

	mux.HandleFunc("/me", handlers.Me)
	mux.HandleFunc("/users/me", handlers.Me)

	addr := resolveListenAddr()
	log.Printf("campus_connect_api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// resolveListenAddr: LISTEN_ADDR (ex. ":8080", "127.0.0.1:3000") tem prioridade;
// senão PORT (ex. "8080" → ":8080"); padrão ":8080".
func resolveListenAddr() string {
	if a := os.Getenv("LISTEN_ADDR"); a != "" {
		return a
	}
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}
