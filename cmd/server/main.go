package main

import (
	"HorizonBackend/config"
	"HorizonBackend/scripts"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := config.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}()

	scripts.AddImagesFromFolder(db, "./static/images")
}
