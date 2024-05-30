package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
)

func (h *Handler) handle(rw http.ResponseWriter, r *http.Request) {
	image := chi.URLParam(r, "image")

	if image == "" {
		http.Error(rw, "image is required", http.StatusBadRequest)
		return
	}

	if cached, ok := h.cache.Get(image); ok {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(cached)
		return
	}

	tags, err := h.service.GetTags(image)
	if err != nil {
		http.Error(rw, fmt.Sprintf("GetTags: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	res := make(map[string]map[string]int)

	for _, tag := range tags {
		manifest, err := h.service.GetManifest(image, tag)
		if err != nil {
			http.Error(rw, fmt.Sprintf("GetManifest: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		var size int
		for _, layer := range manifest.Layers {
			size += layer.Size
		}

		key := fmt.Sprintf("%s:%s", image, tag)
		res[key] = map[string]int{
			"layers_count": len(manifest.Layers),
			"total_size":   size,
		}
	}

	h.cache.Set(image, res, cache.DefaultExpiration)

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
