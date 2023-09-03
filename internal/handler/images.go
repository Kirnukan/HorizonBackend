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

func GetImagesByFamilyAndGroup(s service.ImageService, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := cfg.BaseURL

		vars := mux.Vars(r)
		family := vars["family"]
		group := vars["group"]

		images, err := s.GetImagesByFamilyAndGroup(family, group)
		if err != nil {
			log.Printf("Error fetching images by family and group: %v", err)
			http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
			return
		}

		// Преобразование путей к изображениям
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
		number := vars["number"]

		image, err := service.GetImageByNumber(family, group, number)
		if err != nil {
			log.Printf("Error fetching image by number: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Увеличение счетчика использования
		err = service.IncreaseUsageCount(image.ID) // Предполагая, что у вашей структуры изображения есть поле ID
		if err != nil {
			log.Printf("Error increasing usage count: %v", err)
			// Ошибка при увеличении счетчика использования не должна мешать
			// отправке изображения пользователю, поэтому здесь мы просто логируем ошибку.
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

		// Оставшаяся логика остается без изменений
		images, err := s.GetLeastUsedImages(family, count)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Отправка ответа в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}
