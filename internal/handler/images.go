package handler

import (
	"HorizonBackend/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetImagesByFamilyAndGroup(s *service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]

		images, err := s.GetImagesByFamilyAndGroup(family, group)
		if err != nil {
			log.Printf("Error fetching images by family and group: %v", err)
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(images)
		if err != nil {
			log.Printf("Failed to encode images to JSON: %v", err)
			http.Error(w, "Failed to encode images to JSON", http.StatusInternalServerError)
		}
	}
}

func GetImageByID(s *service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		imageIDStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Image ID is missing", http.StatusBadRequest)
			return
		}

		imageID, err := strconv.Atoi(imageIDStr)
		if err != nil {
			http.Error(w, "Invalid Image ID", http.StatusBadRequest)
			return
		}

		img, err := s.GetImageByIDAndIncreaseUsage(imageID)
		if err != nil {
			http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(img)
		if err != nil {
			http.Error(w, "Failed to encode image to JSON", http.StatusInternalServerError)
		}
	}
}

func SearchImages(s *service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		family := r.URL.Query().Get("family")

		images, err := s.SearchImages(keyword, family)
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

func GetImageByNumber(service *service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]
		number := vars["number"]

		image, err := service.GetImageByNumber(family, group, number)
		if err != nil {
			log.Printf("Error fetching image by number: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(image)
		if err != nil {
			log.Printf("Failed to encode image to JSON: %v", err)
			http.Error(w, "Failed to encode image to JSON", http.StatusInternalServerError)
			return
		}
	}
}
