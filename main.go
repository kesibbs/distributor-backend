package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Distributor struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Capacity int    `json:"capacity"`
	Active   bool   `json:"active"`
}

var distributors = []Distributor{
	{1, "Northline Logistics", "north", 1200, true},
	{2, "Harbor Freight Co", "coastal", 800, true},
	{3, "Inland Cargo Partners", "central", 650, false},
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	port := env("SERVER_PORT", "8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"service":   "distributor-backend",
			"version":   env("DIST_VERSION", "dev"),
			"port":      port,
			"endpoints": []string{"/health", "/api/distributors"},
		})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/distributors", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, distributors)
	})

	log.Printf("distributor-backend listening on :%s (version=%s)", port, env("DIST_VERSION", "dev"))
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
