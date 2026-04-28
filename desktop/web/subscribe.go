package web

import (
	"desktop/build"
	"desktop/internal"
	"desktop/storage"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

func handleSubscribeServers(w http.ResponseWriter, r *http.Request) {
	st := storage.DefaultStorage
	if st == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		servers, err := st.GetServers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, servers)
	case http.MethodPost, http.MethodPut:
		var payload []storage.ServerEntry
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := st.SetServers(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		emitServersChanged()
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSubscribeImportServers(w http.ResponseWriter, r *http.Request) {
	st := storage.DefaultStorage
	if st == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var imported []storage.ServerEntry
	if err := json.NewDecoder(r.Body).Decode(&imported); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := st.SetServers(imported); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	emitServersChanged()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "count": len(imported)})
}

func handleSubscribeExportServers(w http.ResponseWriter, r *http.Request) {
	st := storage.DefaultStorage
	if st == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	servers, err := st.GetServers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="subscribe_servers.json"`)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(servers)
}

type subscriptionConvertRequest struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func handleSubscribeSubscriptionConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	subURL := strings.TrimSpace(req.URL)
	if subURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	subType := strings.ToLower(strings.TrimSpace(req.Type))
	if subType == "" {
		subType = "prism"
	}
	servers, err := convertSubscriptionURL(subURL, subType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"count":   len(servers),
		"servers": servers,
	})
}

func emitServersChanged() {
	if internal.App == nil {
		return
	}
	internal.App.Event.Emit(internal.AppConfig.EventNameServersChanged, time.Now().UnixMilli())
}

func handleSubscribeLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(build.Icon) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(build.Icon)
}

//go:embed all:static/subscribe
var staticSubscribeFS embed.FS

func handleSubscribeFileServer() http.Handler {
	sub, err := fs.Sub(staticSubscribeFS, "static/subscribe")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}
