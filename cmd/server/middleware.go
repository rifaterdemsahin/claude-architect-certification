package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func isAdminRequest(r *http.Request) bool {
	if c, err := r.Cookie("admin_logged_in"); err == nil && c.Value == "true" {
		return true
	}
	return false
}

func observe(cfg config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		setCORS(sw, r)
		if r.Method == http.MethodOptions {
			sw.WriteHeader(http.StatusNoContent)
			return
		}

		defer func() {
			if rec := recover(); rec != nil {
				http.Error(sw, "internal server error", http.StatusInternalServerError)
				go shipToAxiom(cfg, []map[string]any{{
					"_time":  time.Now().UTC().Format(time.RFC3339),
					"level":  "error",
					"method": r.Method,
					"path":   r.URL.Path,
					"panic":  fmt.Sprintf("%v", rec),
				}})
				log.Printf("PANIC %s %s: %v", r.Method, r.URL.Path, rec)
			}
		}()

		next.ServeHTTP(sw, r)

		dur := time.Since(start)
		go shipToAxiom(cfg, []map[string]any{{
			"_time":       time.Now().UTC().Format(time.RFC3339),
			"level":       "info",
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      sw.status,
			"duration_ms": dur.Milliseconds(),
		}})
		log.Printf("%s %s %d %dms", r.Method, r.URL.Path, sw.status, dur.Milliseconds())
	})
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.HasSuffix(origin, ".github.io") || strings.HasPrefix(origin, "http://localhost") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, apikey")
}
