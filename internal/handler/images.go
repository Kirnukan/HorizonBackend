package handler

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func GetImageTest(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, err := rdb.HGetAll(r.Context(), "image:1").Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encodeErr := json.NewEncoder(w).Encode(image)
		if encodeErr != nil {
			return
		}
	}
}
