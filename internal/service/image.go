package service

import (
	"HorizonBackend/internal/model"
	"HorizonBackend/internal/repository/postgres"
	"errors"
	"log"
)

type ImageService interface {
	GetImagesByFamilyAndGroup(family, group string) ([]model.Image, error)
	GetImageByIDAndIncreaseUsage(imageID int) (model.Image, error)
	SearchImages(keyword, family string) ([]model.Image, error)
	GetImageByNumber(family, group, imageNumber string) (*model.Image, error)
	IncreaseUsageCount(imageID int) error
	GetLeastUsedImages(family string, limit int) ([]model.Image, error)
}

type imageServiceImpl struct {
	repo *postgres.ImageRepository
}

func NewImageService(repo *postgres.ImageRepository) ImageService {
	return &imageServiceImpl{repo: repo}
}

func (s *imageServiceImpl) GetImagesByFamilyAndGroup(family, group string) ([]model.Image, error) {
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

func (s *imageServiceImpl) GetImageByIDAndIncreaseUsage(imageID int) (model.Image, error) {
	// Увеличиваем счетчик использования
	err := s.repo.IncreaseUsageCount(imageID)
	if err != nil {
		return model.Image{}, err
	}

	// Получаем изображение
	img, err := s.repo.GetImageByID(imageID)
	return img, err
}

func (s *imageServiceImpl) SearchImages(keyword, family string) ([]model.Image, error) {
	return s.repo.SearchImagesByKeywordAndFamily(keyword, family)
}

func (s *imageServiceImpl) GetImageByNumber(family, group, imageNumber string) (*model.Image, error) {
	image, err := s.repo.FindImageByNumber(family, group, imageNumber)
	if err != nil {
		log.Printf("Service error fetching image by number for family: %s, group: %s, number: %s, Error: %v", family, group, imageNumber, err)
		return nil, err
	}
	return image, nil
}

func (s *imageServiceImpl) IncreaseUsageCount(imageID int) error {
	return s.repo.IncreaseUsageCount(imageID)
}

func (s *imageServiceImpl) GetLeastUsedImages(family string, limit int) ([]model.Image, error) {
	return s.repo.GetLeastUsedImages(family, limit)
}
