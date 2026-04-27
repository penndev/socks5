package web

import (
	"desktop/build"
	"desktop/storage"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

func handleServers(w http.ResponseWriter, r *http.Request) {
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
		for i := range payload {
			if payload[i].ID == "" {
				payload[i].ID = payload[i].Host
			}
		}
		if err := st.SetServers(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleImportServers(w http.ResponseWriter, r *http.Request) {
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
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(8 << 20); err != nil {
			http.Error(w, "parse form failed", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "missing file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&imported); err != nil {
			http.Error(w, "invalid file json", http.StatusBadRequest)
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&imported); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
	}
	for i := range imported {
		if imported[i].ID == "" {
			imported[i].ID = imported[i].Host
		}
	}
	if err := st.SetServers(imported); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "count": len(imported)})
}

func handleExportServers(w http.ResponseWriter, r *http.Request) {
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
	URL string `json:"url"`
}

func handleSubscriptionConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req subscriptionConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	subURL := strings.TrimSpace(req.URL)
	if subURL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	servers, err := convertSubscriptionURL(subURL)
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

func convertSubscriptionURL(subURL string) ([]storage.ServerEntry, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(subURL)
	if err != nil {
		return nil, fmt.Errorf("fetch subscription failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("subscription response status: %d", resp.StatusCode)
	}
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read subscription failed: %w", err)
	}
	text := strings.TrimSpace(string(rawBody))
	if text == "" {
		return nil, fmt.Errorf("empty subscription content")
	}

	decodedText, ok := decodeBase64Text(text)
	if ok {
		text = decodedText
	}

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	servers := make([]storage.ServerEntry, 0, len(lines))
	idSeen := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		entry, ok := parseSubscriptionLine(line)
		if !ok {
			continue
		}
		baseID := entry.ID
		if baseID == "" {
			baseID = entry.Host
		}
		n := idSeen[baseID]
		idSeen[baseID] = n + 1
		if n > 0 {
			entry.ID = fmt.Sprintf("%s-%d", baseID, n+1)
		} else {
			entry.ID = baseID
		}
		servers = append(servers, entry)
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no https/socks5 nodes found in subscription")
	}
	return servers, nil
}

func parseSubscriptionLine(line string) (storage.ServerEntry, bool) {
	raw := strings.TrimSpace(line)
	u, err := url.Parse(raw)
	if err != nil || u == nil {
		if decoded, ok := decodeBase64Text(raw); ok {
			u, err = url.Parse(strings.TrimSpace(decoded))
			if err != nil || u == nil {
				return storage.ServerEntry{}, false
			}
		} else {
			return storage.ServerEntry{}, false
		}
	}

	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	protocol := ""
	switch scheme {
	case "https":
		protocol = "httpovertls"
	case "socks5":
		protocol = "socks5overtls"
	default:
		return storage.ServerEntry{}, false
	}

	host := strings.TrimSpace(u.Host)
	if host == "" {
		return storage.ServerEntry{}, false
	}
	if _, _, err := net.SplitHostPort(host); err != nil {
		// 要求 host:port；无端口的订阅节点跳过。
		return storage.ServerEntry{}, false
	}

	username := ""
	password := ""
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	remark := decodeRemark(u.Fragment)
	if remark == "" {
		remark = host
	}
	idSeed := fmt.Sprintf("%s|%s|%s", host, protocol, remark)
	return storage.ServerEntry{
		ID:       idSeed,
		Host:     host,
		Remark:   remark,
		Username: username,
		Password: password,
		Protocol: protocol,
	}, true
}

func decodeRemark(fragment string) string {
	s := strings.TrimSpace(fragment)
	if s == "" {
		return ""
	}
	if unescaped, err := url.QueryUnescape(s); err == nil {
		s = strings.TrimSpace(unescaped)
	}
	if decoded, ok := decodeBase64Text(s); ok {
		return strings.TrimSpace(decoded)
	}
	return s
}

func decodeBase64Text(s string) (string, bool) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", false
	}
	compact := strings.NewReplacer("\n", "", "\r", "", "\t", "", " ", "").Replace(trimmed)
	try := []func(string) ([]byte, error){
		base64.StdEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawURLEncoding.DecodeString,
	}
	for _, fn := range try {
		b, err := fn(compact)
		if err != nil {
			continue
		}
		if len(b) == 0 || !utf8.Valid(b) {
			continue
		}
		out := strings.TrimSpace(string(b))
		if out == "" {
			continue
		}
		return out, true
	}
	return "", false
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

func subscribeFileServer() http.Handler {
	sub, err := fs.Sub(staticSubscribeFS, "static/subscribe")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}

func registerSubscribe(mux *http.ServeMux) {
	mux.HandleFunc("/subscribe/logo.png", handleSubscribeLogo)
	mux.HandleFunc("/subscribe/api/servers", handleServers)
	mux.HandleFunc("/subscribe/api/servers/import", handleImportServers)
	mux.HandleFunc("/subscribe/api/servers/export", handleExportServers)
	mux.HandleFunc("/subscribe/api/subscription/convert", handleSubscriptionConvert)
	mux.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/subscribe/", http.StatusFound)
	})
	mux.Handle("/subscribe/", http.StripPrefix("/subscribe/", subscribeFileServer()))
}
