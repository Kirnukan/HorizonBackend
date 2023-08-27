package postgres

import (
	"HorizonBackend/internal/model"
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) GetImagesByFamilyAndGroup(family, group string) ([]model.Image, error) {
	// Получение ID группы по имени семейства и группы
	var groupID int
	err := r.db.QueryRow(`SELECT g.id FROM "groups" g JOIN "families" f ON g.family_id = f.id WHERE f.name = $1 AND g.name = $2`, family, group).Scan(&groupID)
	if err != nil {
		return nil, err
	}

	// Получение изображений по ID группы
	rows, err := r.db.Query(`SELECT id, group_id, name, file_path, thumb_path, usage_count, meta_tags FROM "images" WHERE group_id = $1`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		err := rows.Scan(&img.ID, &img.GroupID, &img.Name, &img.FilePath, &img.ThumbPath, &img.UsageCount, pq.Array(&img.MetaTags))
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepository) IncreaseUsageCount(imageID int) error {
	_, err := r.db.Exec("UPDATE Images SET usage_count = usage_count + 1 WHERE id = $1", imageID)
	return err
}

func (r *ImageRepository) GetImageByID(imageID int) (model.Image, error) {
	var img model.Image
	err := r.db.QueryRow(`SELECT id, group_id, name, file_path, usage_count, meta_tags FROM "images" WHERE id = $1`, imageID).Scan(&img.ID, &img.GroupID, &img.Name, &img.FilePath, &img.UsageCount, pq.Array(&img.MetaTags))
	return img, err
}

func (r *ImageRepository) SearchImagesByKeywordAndFamily(keyword, family string) ([]model.Image, error) {
	query := `
        SELECT i.id, i.group_id, i.name, i.file_path, i.thumb_path, i.usage_count, i.meta_tags
        FROM images i
        JOIN groups g ON i.group_id = g.id
        JOIN families f ON g.family_id = f.id
        WHERE (i.name ILIKE $1 OR EXISTS (SELECT 1 FROM unnest(i.meta_tags) AS tag WHERE tag ILIKE $1))
           AND (f.name ILIKE $2)
    `

	log.Printf("Query: %s", query)
	log.Printf("Keyword: %s, Family: %s", keyword, family)

	rows, err := r.db.Query(query, "%"+keyword+"%", "%"+family+"%")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		var metaTags pq.StringArray // Create a pq.StringArray to scan the array

		err := rows.Scan(&img.ID, &img.GroupID, &img.Name, &img.FilePath, &img.ThumbPath, &img.UsageCount, &metaTags)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		img.MetaTags = []string(metaTags) // Convert pq.StringArray to regular []string
		images = append(images, img)
	}

	return images, nil
}

func (r *ImageRepository) getGroupIDByName(groupName string) (int, error) {
	query := `SELECT id FROM groups WHERE name = $1` // Уточните название таблицы и колонок, если оно отличается.
	var groupID int
	err := r.db.QueryRow(query, groupName).Scan(&groupID)
	if err != nil {
		return 0, err
	}
	return groupID, nil
}

func (r *ImageRepository) FindImageByNumber(family, group, imageNumber string) (*model.Image, error) {
	// Получаем ID группы на основе ее имени (группы).
	groupID, err := r.getGroupIDByName(group)
	if err != nil {
		return nil, err
	}

	imageNamePattern := family + "_" + group + "_%" + imageNumber

	query := `SELECT id, group_id, name, file_path, thumb_path, usage_count, meta_tags 
              FROM images 
              WHERE group_id = $1 AND name LIKE $2`

	row := r.db.QueryRow(query, groupID, imageNamePattern)

	image := &model.Image{}
	err = row.Scan(&image.ID, &image.GroupID, &image.Name, &image.FilePath, &image.ThumbPath, &image.UsageCount, pq.Array(&image.MetaTags))

	if err != nil {
		return nil, err
	}
	return image, nil
}
