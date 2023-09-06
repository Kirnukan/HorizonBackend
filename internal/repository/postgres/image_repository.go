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

func (r *ImageRepository) GetImagesByFamilyGroupSubgroup(family, group, subgroup string) ([]model.Image, error) {
	// Получение ID подгруппы по имени семейства, группы и подгруппы
	var subgroupID int
	err := r.db.QueryRow(`
        SELECT s.id 
        FROM "subgroups" s 
        JOIN "groups" g ON s.group_id = g.id 
        JOIN "families" f ON g.family_id = f.id 
        WHERE f.name = $1 AND g.name = $2 AND s.name = $3`, family, group, subgroup).Scan(&subgroupID)
	if err != nil {
		return nil, err
	}

	// Получение изображений по ID подгруппы
	rows, err := r.db.Query(`
		SELECT i.id, i.subgroup_id, i.name, i.file_path, i.thumb_path, i.usage_count, i.meta_tags 
		FROM "images" i
		JOIN "subgroups" sg ON i.subgroup_id = sg.id 
		WHERE sg.id = $1`, subgroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		err := rows.Scan(&img.ID, &img.SubgroupID, &img.Name, &img.FilePath, &img.ThumbPath, &img.UsageCount, pq.Array(&img.MetaTags))
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
	err := r.db.QueryRow(`SELECT id, subgroup_id, name, file_path, usage_count, meta_tags FROM "images" WHERE id = $1`, imageID).Scan(&img.ID, &img.SubgroupID, &img.Name, &img.FilePath, &img.UsageCount, pq.Array(&img.MetaTags))
	return img, err
}

func (r *ImageRepository) SearchImagesByKeywordAndFamily(keyword, family string) ([]model.Image, error) {
	query := `
        SELECT i.id, i.subgroup_id, i.name, i.file_path, i.thumb_path, i.usage_count, i.meta_tags
        FROM images i
        JOIN subgroups s ON i.subgroup_id = s.id
        JOIN groups g ON s.group_id = g.id
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

		err := rows.Scan(&img.ID, &img.SubgroupID, &img.Name, &img.FilePath, &img.ThumbPath, &img.UsageCount, &metaTags)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		img.MetaTags = []string(metaTags) // Convert pq.StringArray to regular []string
		images = append(images, img)
	}

	return images, nil
}

func (r *ImageRepository) getSubgroupIDByName(subgroupName string) (int, error) {
	query := `SELECT id FROM subgroups WHERE name = $1`
	var subgroupID int
	err := r.db.QueryRow(query, subgroupName).Scan(&subgroupID)
	if err != nil {
		return 0, err
	}
	return subgroupID, nil
}

func (r *ImageRepository) FindImageByNumber(family, group, subgroup, imageNumber string) (*model.Image, error) {
	// Получаем ID подгруппы на основе ее имени.
	subgroupID, err := r.getSubgroupIDByName(subgroup)
	if err != nil {
		return nil, err
	}

	imageNamePattern := family + "_" + group + "_" + subgroup + "_%" + imageNumber

	query := `SELECT id, subgroup_id, name, file_path, thumb_path, usage_count, meta_tags 
              FROM images 
              WHERE subgroup_id = $1 AND name LIKE $2`

	row := r.db.QueryRow(query, subgroupID, imageNamePattern)

	image := &model.Image{}
	err = row.Scan(&image.ID, &image.SubgroupID, &image.Name, &image.FilePath, &image.ThumbPath, &image.UsageCount, pq.Array(&image.MetaTags))

	if err != nil {
		return nil, err
	}
	return image, nil
}

func (r *ImageRepository) GetLeastUsedImages(family string, limit int) ([]model.Image, error) {
	// Запрос к базе данных для получения 6 изображений с наименьшим счетчиком использования для определенного семейства
	const query = `
    SELECT i.id, i.subgroup_id, i.name, i.file_path, i.thumb_path, i.usage_count, i.meta_tags 
		FROM "images" i
		JOIN "subgroups" sg ON i.subgroup_id = sg.id 
		JOIN "groups" g ON sg.group_id = g.id 
		JOIN "families" f ON g.family_id = f.id 
		WHERE f.name = $1 
		ORDER BY i.usage_count ASC 
		LIMIT $2;
    `
	rows, err := r.db.Query(query, family, limit)
	if err != nil {
		log.Printf("Error querying the database: %v", err)
		return nil, err
	}

	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		err := rows.Scan(&img.ID, &img.SubgroupID, &img.Name, &img.FilePath, &img.ThumbPath, &img.UsageCount, pq.Array(&img.MetaTags))
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
