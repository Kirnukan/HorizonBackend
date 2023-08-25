package scripts

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

func AddImagesFromFolder(db *sql.DB, baseFolder string) {
	familyDirs, err := os.ReadDir(baseFolder)
	if err != nil {
		panic(err)
	}

	for _, familyDir := range familyDirs {
		if !familyDir.IsDir() {
			continue
		}
		familyName := familyDir.Name()

		_, err := db.Exec(`INSERT INTO Families (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, familyName)
		if err != nil {
			panic(err)
		}

		groupDirs, err := os.ReadDir(filepath.Join(baseFolder, familyName))
		if err != nil {
			panic(err)
		}

		for _, groupDir := range groupDirs {
			if !groupDir.IsDir() {
				continue
			}
			groupName := groupDir.Name()

			_, err := db.Exec(`INSERT INTO Groups (name, family_id) VALUES ($1, (SELECT id FROM Families WHERE name = $2)) ON CONFLICT (name) DO NOTHING`, groupName, familyName)
			if err != nil {
				panic(err)
			}

			imageFiles, err := os.ReadDir(filepath.Join(baseFolder, familyName, groupName))
			if err != nil {
				panic(err)
			}

			for _, imageFile := range imageFiles {
				if imageFile.IsDir() {
					continue
				}
				imageName := strings.TrimSuffix(imageFile.Name(), filepath.Ext(imageFile.Name()))
				imagePath := filepath.Join("static", "images", familyName, groupName, imageFile.Name())

				_, err := db.Exec(`INSERT INTO Images (name, file_path, group_id) VALUES ($1, $2, (SELECT id FROM Groups WHERE name = $3))`, imageName, imagePath, groupName)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
