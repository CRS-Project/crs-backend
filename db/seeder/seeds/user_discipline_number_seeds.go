package seeds

import (
	"encoding/json"
	"os"

	"github.com/CRS-Project/crs-backend/internal/entity"
	mylog "github.com/CRS-Project/crs-backend/internal/pkg/logger"
	"gorm.io/gorm"
)

func SeederUserDisciplineNumber(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding user discipline numbers...")
	jsonFile, err := os.Open("./db/seeder/data/user_discipline_number_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.UserDisciplineNumber
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding user discipline numbers completed")
	return nil
}
