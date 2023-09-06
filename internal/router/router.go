package router

import (
	"HorizonBackend/config"
	"HorizonBackend/internal/handler"
	"HorizonBackend/internal/repository/postgres"
	"HorizonBackend/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	// Инициализация репозитория и службы
	imageRepo := postgres.NewImageRepository(db)
	imageService := service.NewImageService(imageRepo)

	// Регистрация обработчика
	r.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("./static/images/"))))
	r.HandleFunc("/{family}/{group}/{subgroup}/{number:[0-9]+}", handler.GetImageByNumber(imageService, cfg)).Methods("GET") // Обновлено
	r.HandleFunc("/{family}/{group}/{subgroup}/", handler.GetImagesByFamilyGroupSubgroup(imageService, cfg)).Methods("GET")  // Обновлено
	//можете использовать этот роут для поиска изображений,
	//передавая в запросе параметры keyword и family.
	//Например: /search?keyword=Form&family=Forms.
	r.HandleFunc("/least-used", handler.GetLeastUsedImages(imageService, cfg)).Methods("GET")
	r.HandleFunc("/search", handler.SearchImages(imageService, cfg)).Methods("GET")

	return r
}
