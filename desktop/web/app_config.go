package web

import (
	"desktop/storage"
	"net/http"
	"strings"
)

func handleAppConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	st := storage.DefaultStorage
	if st == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}
	settings, err := st.GetSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	language := "zh-CN"
	themeMode := "system"
	if settings != nil {
		if v := strings.TrimSpace(settings.System.Language); v != "" {
			language = v
		}
		if v := strings.TrimSpace(settings.System.ThemeMode); v != "" {
			themeMode = strings.ToLower(v)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"language":  language,
		"themeMode": themeMode,
	})
}
