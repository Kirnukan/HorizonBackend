package handler

import (
	"HorizonBackend/config"
	"HorizonBackend/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetImagesByFamilyGroupSubgroup(s service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := cfg.BaseURL

		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]
		subgroup := vars["subgroup"] 

		images, err := s.GetImagesByFamilyGroupSubgroup(family, group, subgroup) 
		if err != nil {
			log.Printf("Error fetching images by family, group and subgroup: %v", err) 
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}

		for i := range images {
			images[i].FilePath = baseURL + images[i].FilePath
			images[i].ThumbPath = baseURL + images[i].ThumbPath
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(images)
		if err != nil {
			log.Printf("Failed to encode images to JSON: %v", err)
			http.Error(w, "Failed to encode images to JSON", http.StatusInternalServerError)
		}
	}
}


func SearchImages(s service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := cfg.BaseURL

		keyword := r.URL.Query().Get("keyword")
		family := r.URL.Query().Get("family")

		images, err := s.SearchImages(keyword, family)
		if err != nil {
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}

		for i := range images {
			images[i].FilePath = baseURL + images[i].FilePath
			images[i].ThumbPath = baseURL + images[i].ThumbPath
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(images)
		if err != nil {
			log.Printf("Failed to encode images to JSON: %v", err)
			http.Error(w, "Failed to encode images to JSON", http.StatusInternalServerError)
		}
	}
}

type ImageResponse struct {
	FilePath string `json:"file_path"`
}

func IncreaseImageUsage(service service.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("IncreaseImageUsage handler called")
		vars := mux.Vars(r)
		thumbPath, ok := vars["thumbPath"]
		fmt.Printf("Extracted thumbPath: %s, success: %v\n", thumbPath, ok)

		if !ok {
			http.Error(w, "Thumb path is required", http.StatusBadRequest)
			return
		}

		err := service.IncreaseUsageCount(thumbPath)
		if err != nil {
			fmt.Printf("Error increasing usage count: %v\n", err)
			http.Error(w, fmt.Sprintf("Error increasing usage count: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Usage count increased"))
	}
}

func GetImageByNumber(service service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := cfg.BaseURL

		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]
		subgroup := vars["subgroup"] 
		number := vars["number"]

		image, err := service.GetImageByNumber(family, group, subgroup, number) 
		if err != nil {
			log.Printf("Error fetching image by number: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//err = service.IncreaseUsageCount(image.ID)
		//if err != nil {
		//	log.Printf("Error increasing usage count: %v", err)
		//}

		response := ImageResponse{
			FilePath: baseURL + image.FilePath,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Failed to encode image to JSON: %v", err)
			http.Error(w, "Failed to encode image to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func GetLeastUsedImages(s service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		family := r.URL.Query().Get("family")
		if family == "" {
			http.Error(w, "Family parameter is missing", http.StatusBadRequest)
			return
		}

		count := 6
		countStr := r.URL.Query().Get("count")
		if countStr != "" {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil {
				http.Error(w, "Invalid count parameter", http.StatusBadRequest)
				return
			}
		} else {
			count = 6
		}

		log.Printf("Fetching least used images for family: %s and count: %d", family, count)

		images, err := s.GetLeastUsedImages(family, count)
		if err != nil {
			log.Printf("Error fetching least used images: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Fetched %d images", len(images))

		for i := range images {
			images[i].FilePath = cfg.BaseURL + images[i].FilePath
			images[i].ThumbPath = cfg.BaseURL + images[i].ThumbPath
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}
