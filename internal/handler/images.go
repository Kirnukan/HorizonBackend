package handler

import (
	"HorizonBackend/config"
	"HorizonBackend/internal/service"
	"encoding/json"
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
		subgroup := vars["subgroup"] // Добавлено

		images, err := s.GetImagesByFamilyGroupSubgroup(family, group, subgroup) // Обновлено
		if err != nil {
			log.Printf("Error fetching images by family, group and subgroup: %v", err) // Обновлено
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

func GetImageByID(s service.ImageService) http.HandlerFunc {
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

func GetImageByNumber(service service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := cfg.BaseURL

		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]
		subgroup := vars["subgroup"] // Добавлено
		number := vars["number"]

		image, err := service.GetImageByNumber(family, group, subgroup, number) // Обновлено
		if err != nil {
			log.Printf("Error fetching image by number: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = service.IncreaseUsageCount(image.ID)
		if err != nil {
			log.Printf("Error increasing usage count: %v", err)
		}

		image.FilePath = baseURL + image.FilePath
		image.ThumbPath = baseURL + image.ThumbPath

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(image)
		if err != nil {
			log.Printf("Failed to encode image to JSON: %v", err)
			http.Error(w, "Failed to encode image to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func GetLeastUsedImages(s service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлеките параметры family и count из строки запроса
		family := r.URL.Query().Get("family")
		if family == "" {
			http.Error(w, "Family parameter is missing", http.StatusBadRequest)
			return
		}

		// Извлечение параметра count из запроса и установка значения по умолчанию на 6
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
			count = 6 // значение по умолчанию
		}

		// Логирование входящих параметров
		log.Printf("Fetching least used images for family: %s and count: %d", family, count)

		// Оставшаяся логика остается без изменений
		images, err := s.GetLeastUsedImages(family, count)
		if err != nil {
			// Логирование ошибки
			log.Printf("Error fetching least used images: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Логирование количества извлеченных изображений
		log.Printf("Fetched %d images", len(images))

		// Отправка ответа в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}
