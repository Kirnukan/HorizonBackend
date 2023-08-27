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
		log.Printf("Service error fetching images for family: %s and group: %s, Error: %v", family, group, err)
		return nil, err
	}

	// 3. Преобразование данных
	for i := range images {
		images[i].UsageCount++
	}

	return images, nil
}

func (s *ImageService) GetImageByIDAndIncreaseUsage(imageID int) (model.Image, error) {
	// Увеличиваем счетчик использования
	err := s.repo.IncreaseUsageCount(imageID)
	if err != nil {
		return model.Image{}, err
	}

	// Получаем изображение
	img, err := s.repo.GetImageByID(imageID)
	return img, err
}

func (s *ImageService) SearchImages(keyword, family string) ([]model.Image, error) {
	return s.repo.SearchImagesByKeywordAndFamily(keyword, family)
}

func (s *ImageService) GetImageByNumber(family, group, imageNumber string) (*model.Image, error) {
	image, err := s.repo.FindImageByNumber(family, group, imageNumber)
	if err != nil {
		log.Printf("Service error fetching image by number for family: %s, group: %s, number: %s, Error: %v", family, group, imageNumber, err)
		return nil, err
	}
	return image, nil
}
