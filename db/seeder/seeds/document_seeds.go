package seeds

import (
	"encoding/json"
	"os"

	"github.com/CRS-Project/crs-backend/internal/entity"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"gorm.io/gorm"
)

func SeederDocument(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding documents...")
	jsonFile, err := os.Open("./db/seeder/data/document_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.Document
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding documents completed")
	return nil
}
