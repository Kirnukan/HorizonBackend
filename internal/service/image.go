package service

import (
	"HorizonBackend/internal/model"
	"HorizonBackend/internal/repository/postgres"
	"errors"
	"log"
)

type ImageService struct {
	repo *postgres.ImageRepository
}

func NewImageService(repo *postgres.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) GetImagesByFamilyAndGroup(family, group string) ([]model.Image, error) {
	// 1. Валидация
	if family == "" || group == "" {
		log.Println("Invalid input: family or group is empty")
		return nil, errors.New("family and group cannot be empty")
	}

	// 2. Получение изображений
	images, err := s.repo.GetImagesByFamilyAndGroup(family, group)
	if err != nil {
		log.Printf("Error fetching images for family: %s and group: %s, Error: %v", family, group, err)
		return nil, err
	}

	// 3. Преобразование данных (простой пример: увеличение счетчика использования)
	// В реальной ситуации вы, возможно, захотите обновить это значение в базе данных
	for i := range images {
		images[i].UsageCount++
	}

	return images, nil
}
