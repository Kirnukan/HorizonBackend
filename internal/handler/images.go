package handler

import (
	"HorizonBackend/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetImagesByFamilyAndGroup(s *service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]

		images, err := s.GetImagesByFamilyAndGroup(family, group)
		if err != nil {
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(images)
		if err != nil {
			http.Error(w, "Failed to encode images to JSON", http.StatusInternalServerError)
		}
	}
}
