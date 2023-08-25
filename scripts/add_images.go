package scripts

import (
	"database/sql"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func compressImage(inputPath string, outputPath string) error {
	fmt.Println("Processing:", inputPath) // <-- Добавим эту строку для вывода имени обрабатываемого файла
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported format for file: %s", inputPath)
	}

	m := resize.Resize(100, 100, img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if filepath.Ext(inputPath) == ".jpg" || filepath.Ext(inputPath) == ".jpeg" {
		err = jpeg.Encode(out, m, nil)
	} else if filepath.Ext(inputPath) == ".png" {
		err = png.Encode(out, m)
	}

	return err
}

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
				thumbPath := ""

				if familyName == "Forms" {
					thumbPath = imagePath
				} else {
					thumbPath = filepath.Join("static", "images", familyName, groupName, imageName+"_thumb"+filepath.Ext(imageFile.Name()))

					err = compressImage(filepath.Join(baseFolder, familyName, groupName, imageFile.Name()), filepath.Join(baseFolder, familyName, groupName, imageName+"_thumb"+filepath.Ext(imageFile.Name())))
					if err != nil {
						panic(err)
					}
				}

				_, err := db.Exec(`INSERT INTO Images (name, file_path, thumb_path, group_id) VALUES ($1, $2, $3, (SELECT id FROM Groups WHERE name = $4))`, imageName, imagePath, thumbPath, groupName)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
